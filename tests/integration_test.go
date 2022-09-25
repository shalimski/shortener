//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/shalimski/shortener/config"
	"github.com/shalimski/shortener/internal/app"
	"github.com/shalimski/shortener/internal/web"
)

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(ShortenerSuit))
}

type ShortenerSuit struct {
	suite.Suite
	mongoContainer testcontainers.Container
	etcdContainer  testcontainers.Container
	redisContainer testcontainers.Container
	port           string
}

func (s *ShortenerSuit) SetupSuite() {
	t := s.T()

	// ctx := context.Background()

	// starting containers
	port, err := s.startMongoDB()
	s.NoError(err)
	s.NotEmpty(port)
	t.Setenv("MONGO_PORT", port)

	port, err = s.startEtcd()
	s.NoError(err)
	s.NotEmpty(port)
	t.Setenv("ETCD_ENDPOINTS", "http://127.0.0.1:"+port)

	port, err = s.startRedis()
	s.NoError(err)
	s.NotEmpty(port)
	t.Setenv("REDIS_DSN", "127.0.0.1:"+port)

	cfg, err := config.New()
	if err != nil {
		t.Logf("fail to parse config: %s", err.Error())
	}

	s.port = cfg.HTTP.Port

	go app.Run(cfg)

	time.Sleep(2 * time.Second)

	t.Log("Suite setup is done")
}

func (s *ShortenerSuit) TestFull() {
	jsonOk := []byte(`{ 
		"long_url":"https://commandcenter.blogspot.com/"
	}`)

	jsonBad := []byte(`{ 
		"url":"123"
	}`)

	jsonNotValid := []byte(`{ 
		"long_url":"123"
	}`)

	api := fmt.Sprintf("http://localhost:%s/api/v1", s.port)

	c := http.Client{}

	s.Run("create ok", func() {
		r, err := c.Post(api+"/shorten", "application/json", bytes.NewBuffer(jsonOk))
		s.NoError(err)
		defer r.Body.Close()

		s.Equal(http.StatusOK, r.StatusCode)

		var dto web.ResponseCreateDTO

		json.NewDecoder(r.Body).Decode(&dto)

		s.Equal("b", dto.ShortURL)
	})

	s.Run("create bad", func() {
		r, err := c.Post(api+"/shorten", "application/json", bytes.NewBuffer(jsonBad))
		s.NoError(err)
		defer r.Body.Close()

		s.Equal(http.StatusBadRequest, r.StatusCode)

		var dto web.ResponseMessage

		err = json.NewDecoder(r.Body).Decode(&dto)

		s.NoError(err)
	})

	s.Run("create not valid", func() {
		r, err := c.Post(api+"/shorten", "application/json", bytes.NewBuffer(jsonNotValid))
		s.NoError(err)
		defer r.Body.Close()

		s.Equal(http.StatusBadRequest, r.StatusCode)

		var dto web.ResponseMessage

		err = json.NewDecoder(r.Body).Decode(&dto)

		s.NoError(err)
		s.Equal("invalid long url", dto.Message)
	})

	s.Run("redirect ok", func() {
		r, err := c.Get(api + "/b")
		s.NoError(err)
		defer r.Body.Close()

		s.Equal(http.StatusOK, r.StatusCode)
	})

	s.Run("redirect bad", func() {
		r, err := c.Get(api + "/bbbbb")
		s.NoError(err)
		defer r.Body.Close()

		s.Equal(http.StatusNotFound, r.StatusCode)

		var dto web.ResponseMessage

		err = json.NewDecoder(r.Body).Decode(&dto)

		s.NoError(err)
		s.Equal("short url not found", dto.Message)
	})

	s.Run("redirect not valid", func() {
		r, err := c.Get(api + "/b*2")
		s.NoError(err)
		defer r.Body.Close()

		s.Equal(http.StatusBadRequest, r.StatusCode)

		var dto web.ResponseMessage

		err = json.NewDecoder(r.Body).Decode(&dto)

		s.NoError(err)
		s.Equal("invalid short url", dto.Message)
	})

	s.Run("delete ok", func() {
		req, err := http.NewRequest(http.MethodDelete, api+"/b", nil)
		s.NoError(err)

		r, err := c.Do(req)
		s.NoError(err)
		defer r.Body.Close()

		s.Equal(http.StatusOK, r.StatusCode)

		var dto web.ResponseMessage

		err = json.NewDecoder(r.Body).Decode(&dto)

		s.NoError(err)
		s.Equal("url deleted", dto.Message)
	})

	s.Run("delete bad", func() {
		req, err := http.NewRequest(http.MethodDelete, api+"/bbbbbbb", nil)
		s.NoError(err)

		r, err := c.Do(req)
		s.NoError(err)
		defer r.Body.Close()

		s.Equal(http.StatusNotFound, r.StatusCode)

		var dto web.ResponseMessage

		err = json.NewDecoder(r.Body).Decode(&dto)

		s.NoError(err)
		s.Equal("short url not found", dto.Message)
	})

	s.Run("delete not valid", func() {
		req, err := http.NewRequest(http.MethodDelete, api+"/bb-bbbbb", nil)
		s.NoError(err)

		r, err := c.Do(req)
		s.NoError(err)
		defer r.Body.Close()

		s.Equal(http.StatusBadRequest, r.StatusCode)

		var dto web.ResponseMessage

		err = json.NewDecoder(r.Body).Decode(&dto)

		s.NoError(err)
		s.Equal("invalid short url", dto.Message)
	})
}

func (s *ShortenerSuit) TestPing() {
	c := http.Client{}
	s.Run("ping", func() {
		r, err := c.Get(fmt.Sprintf("http://localhost:%s/ping", s.port))
		s.NoErrorf(err, "ping request failed: %s", err)
		s.Equal(http.StatusOK, r.StatusCode)
	})
}

func (s *ShortenerSuit) TearDownSuite() {
	ctx := context.Background()
	err := s.mongoContainer.Terminate(ctx)
	s.NoError(err)

	err = s.etcdContainer.Terminate(ctx)
	s.NoError(err)

	err = s.redisContainer.Terminate(ctx)
	s.NoError(err)
}

func (s *ShortenerSuit) startMongoDB() (string, error) {
	ctx := context.Background()

	env := make(map[string]string)
	env["MONGO_INITDB_ROOT_USERNAME"] = "admin"
	env["MONGO_INITDB_ROOT_PASSWORD"] = "admin"

	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:5.0.9",
			Env:          env,
			ExposedPorts: []string{string("27017")},
			WaitingFor:   wait.ForLog("Waiting for connections").WithStartupTimeout(15 * time.Second),
			AutoRemove:   true,
		},
		Started: true,
	})
	s.NoError(err)

	natPort, err := mongoC.MappedPort(ctx, nat.Port("27017"))
	s.NoError(err)

	s.mongoContainer = mongoC

	return natPort.Port(), nil
}

func (s *ShortenerSuit) startEtcd() (string, error) {
	ctx := context.Background()

	env := make(map[string]string)
	env["ALLOW_NONE_AUTHENTICATION"] = "yes"
	env["ETCD_ADVERTISE_CLIENT_URLS"] = "http://127.0.0.1:2379"

	etcdC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "bitnami/etcd:3.5.5",
			Env:          env,
			ExposedPorts: []string{"2379"},
			WaitingFor:   wait.ForLog("ready to serve client requests").WithStartupTimeout(15 * time.Second),
			AutoRemove:   true,
		},
		Started: true,
	})
	s.NoError(err)

	natPort, err := etcdC.MappedPort(ctx, nat.Port("2379"))
	s.NoError(err)

	s.etcdContainer = etcdC

	return natPort.Port(), nil
}

func (s *ShortenerSuit) startRedis() (string, error) {
	ctx := context.Background()

	env := make(map[string]string)
	env["REDIS_PASSWORD"] = "admin"

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:7.0.5-alpine",
			Env:          env,
			ExposedPorts: []string{"6379"},
			WaitingFor:   wait.ForLog("Ready to accept connections").WithStartupTimeout(15 * time.Second),
			AutoRemove:   true,
		},
		Started: true,
	})
	s.NoError(err)

	natPort, err := redisC.MappedPort(ctx, nat.Port("6379"))
	s.NoError(err)

	s.redisContainer = redisC

	return natPort.Port(), nil
}
