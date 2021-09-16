package repository

import (
	"context"

	"github.com/ekifel/moneysaverz/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Users interface {
	Create(ctx context.Context, user models.User) error
	GetByCreds(ctx context.Context, email, password string) (models.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (models.User, error)
	SetSession(ctx context.Context, userID primitive.ObjectID, session models.Session) error
}

type Repositories struct {
	Users Users
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		Users: NewUsersRepo(db),
	}
}
