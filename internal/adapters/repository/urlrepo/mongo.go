package urlrepo

import (
	"context"
	"errors"

	"github.com/shalimski/shortener/internal/domain"
	"github.com/shalimski/shortener/internal/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const urlCollection = "links"

var _ ports.Repository = (*urlRepo)(nil)

// repository to save links
type urlRepo struct {
	collection *mongo.Collection
}

// NewURLRepo create instance of urlRepo
func NewURLRepo(db *mongo.Database) ports.Repository {
	return &urlRepo{
		collection: db.Collection(urlCollection),
	}
}

// Create add new value to DB
func (r *urlRepo) Create(ctx context.Context, url domain.URL) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	_, err := r.collection.InsertOne(ctx, url)

	return err
}

// Find first value by shortURL
func (r *urlRepo) Find(ctx context.Context, shortURL string) (domain.URL, error) {
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

// Delete value by short url
func (r *urlRepo) Delete(ctx context.Context, shortURL string) error {
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
