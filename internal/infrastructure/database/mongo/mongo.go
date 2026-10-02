package mongo

import (
	"context"
	"time"

	"github.com/lion/go-framework/internal/config"
	domainuser "github.com/lion/go-framework/internal/domain/user"
	apperrors "github.com/lion/go-framework/internal/shared/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	client *mongo.Client
}

func Open(cfg config.MongoConfig) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, err
	}
	return &Client{client: client}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx, nil)
}

func (c *Client) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.client.Disconnect(ctx)
}

type userDoc struct {
	ID        string    `bson:"_id"`
	Name      string    `bson:"name"`
	Email     string    `bson:"email"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type UserRepository struct {
	col *mongo.Collection
}

func NewUserRepository(client *Client, database string) *UserRepository {
	col := client.client.Database(database).Collection("users")
	_, _ = col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return &UserRepository{col: col}
}

func (r *UserRepository) Create(ctx context.Context, user *domainuser.User) error {
	_, err := r.col.InsertOne(ctx, toDoc(user))
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domainuser.User, error) {
	var doc userDoc
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return fromDoc(doc), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domainuser.User, error) {
	var doc userDoc
	err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return fromDoc(doc), nil
}

func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*domainuser.User, int64, error) {
	total, err := r.col.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}
	opts := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.col.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var items []*domainuser.User
	for cursor.Next(ctx) {
		var doc userDoc
		if err := cursor.Decode(&doc); err != nil {
			return nil, 0, err
		}
		items = append(items, fromDoc(doc))
	}
	return items, total, cursor.Err()
}

func (r *UserRepository) Update(ctx context.Context, user *domainuser.User) error {
	res, err := r.col.ReplaceOne(ctx, bson.M{"_id": user.ID}, toDoc(user))
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func toDoc(u *domainuser.User) userDoc {
	return userDoc{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func fromDoc(d userDoc) *domainuser.User {
	return &domainuser.User{
		ID:        d.ID,
		Name:      d.Name,
		Email:     d.Email,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
