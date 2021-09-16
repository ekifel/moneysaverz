package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ekifel/moneysaverz/internal/models"
	"github.com/ekifel/moneysaverz/pkg/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersRepo struct {
	db *mongo.Collection
}

func NewUsersRepo(db *mongo.Database) *UsersRepo {
	return &UsersRepo{
		db: db.Collection(usersCollection),
	}
}

func (r *UsersRepo) Create(ctx context.Context, user models.User) error {
	_, err := r.db.InsertOne(ctx, user)
	if mongodb.IsDuplicate(err) {
		return models.ErrUserAlreadyExists
	}

	return err
}

func (r *UsersRepo) GetByCreds(ctx context.Context, email, password string) (models.User, error) {
	var user models.User
	if err := r.db.FindOne(ctx, bson.M{"email": email, "password": password}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, models.ErrUserNotFound
		}

		return models.User{}, err
	}

	return user, nil
}

func (r *UsersRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (models.User, error) {
	var user models.User
	if err := r.db.FindOne(ctx, bson.M{
		"session.refreshToken": refreshToken,
		"session.expiresAt":    bson.M{"$gt": time.Now()},
	}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, models.ErrUserNotFound
		}

		return models.User{}, err
	}

	return user, nil
}

func (r *UsersRepo) SetSession(ctx context.Context, userID primitive.ObjectID, session models.Session) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"session": session, "lastVisitAt": time.Now()}})

	return err
}
