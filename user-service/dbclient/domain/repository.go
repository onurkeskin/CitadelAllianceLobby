package users

import (
	"context"

	"github.com/jmoiron/sqlx"

	"keon.com/CitadelAllianceLobbyServer/user-service/domain"
)

type IUserRepositoryFactory interface {
	New(db sqlx.DB) IUserRepository
}

type IUserRepository interface {
	CreateUser(ctx context.Context, user domain.IUser) error
	GetUsers(ctx context.Context) (domain.IUsers, error)
	FilterUsers(ctx context.Context, field []string, query string, queryParams []interface{}, lastID string, limit int, sort string) domain.IUsers
	CountUsers(ctx context.Context, field []string, query string, queryParams []interface{}) int64
	DeleteUsers(ctx context.Context, ids []string) error
	DeleteAllUsers(ctx context.Context) error
	GetUserById(ctx context.Context, id string) (domain.IUser, error)
	GetUserByUsername(ctx context.Context, username string) (domain.IUser, error)
	GetUserByEmail(ctx context.Context, email string) (domain.IUser, error)
	UserExistsByUsername(ctx context.Context, username string) bool
	UserExistsByEmail(ctx context.Context, email string) bool
	UpdateUser(ctx context.Context, id string, inUser domain.IUser) (domain.IUser, error)
	DeleteUser(ctx context.Context, id string) error
}
