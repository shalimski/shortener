package urlrepo

import (
	"context"
	"errors"

	"github.com/shalimski/shortener/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const urlCollection = "links"

type URLRepo struct {
	collection *mongo.Collection
}

func NewURLRepo(db *mongo.Database) *URLRepo {
	return &URLRepo{
		collection: db.Collection(urlCollection),
	}
}

func (r *URLRepo) Create(ctx context.Context, url domain.URL) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	_, err := r.collection.InsertOne(ctx, url)

	return err
}

func (r *URLRepo) Find(ctx context.Context, shortURL string) (domain.URL, error) {
	select {
	case <-ctx.Done():
		return domain.URL{}, ctx.Err()
	default:
	}

	var url domain.URL

	if err := r.collection.FindOne(ctx, bson.M{"shorturl": shortURL}).Decode(&url); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.URL{}, domain.ErrNotFound
		}

		return domain.URL{}, err
	}

	return url, nil
}

func (r *URLRepo) Delete(ctx context.Context, shortURL string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	dresult, err := r.collection.DeleteOne(ctx, bson.M{"shorturl": shortURL})
	if err != nil {
		return err
	}

	if dresult.DeletedCount != 1 {
		return domain.ErrNotFound
	}

	return err
}
