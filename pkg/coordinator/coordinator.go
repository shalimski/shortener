package coordinator

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Scalingo/go-etcd-lock/v5/lock"

	etcd "go.etcd.io/etcd/client/v3"
)

const (
	locker  = "/locker"
	counter = "/counter"
)

type Coordinator struct {
	cli    *etcd.Client
	locker lock.Locker
}

func NewCoordinator(endpoints []string) (*Coordinator, error) {
	cli, err := etcd.New(etcd.Config{
		Endpoints:   endpoints,
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	locker := lock.NewEtcdLocker(
		cli,
		lock.WithTryLockTimeout(500*time.Millisecond),
	)

	return &Coordinator{
		cli:    cli,
		locker: locker,
	}, nil
}

func (c *Coordinator) NextCounter(ctx context.Context) (int, error) {
	lock, err := c.locker.WaitAcquire(locker, 2)
	if err != nil {
		return 0, fmt.Errorf("failed to lock: %w", err)
	}
	defer lock.Release()

	resp, err := c.cli.Get(ctx, counter)
	if err != nil {
		return 0, fmt.Errorf("failed to get counter: %w", err)
	}

	var current int
	if resp.Count != 0 {
		current, err = strconv.Atoi(string(resp.Kvs[0].Value))
		if err != nil {
			return 0, fmt.Errorf("failed to convert counter: %w", err)
		}
	}

	next := current + 1

	_, err = c.cli.Put(ctx, counter, strconv.Itoa(next))

	if err != nil {
		return 0, fmt.Errorf("failed to put counter: %w", err)
	}

	return next, nil
}

func (c *Coordinator) Shutdown() {
	if c != nil && c.cli != nil {
		c.cli.Close()
	}
}
