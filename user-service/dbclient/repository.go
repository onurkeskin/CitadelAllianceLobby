package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	. "github.com/volatiletech/sqlboiler/queries/qm"

	userServiceDomain "keon.com/CitadelAllianceLobbyServer/user-service/dbclient/domain"
	caModels "keon.com/CitadelAllianceLobbyServer/user-service/dbmodels"
	"keon.com/CitadelAllianceLobbyServer/user-service/domain"
	userModel "keon.com/CitadelAllianceLobbyServer/user-service/model"
)

// UserRepositoryFactory factory
type UserRepositoryFactory struct{}

// New creates a new UserRepository
func (factory *UserRepositoryFactory) New(db *sqlx.DB) userServiceDomain.IUserRepository {
	return &UserRepository{db}
}

// UserRepository repository for user-service
type UserRepository struct {
	DB *sqlx.DB
}

// GetUsers returns Users In database
func (r *UserRepository) GetUsers(ctx context.Context) (domain.IUsers, error) {
	users, err := caModels.Users(Load("UserPassword")).All(ctx, r.DB)
	var modelUsers []*userModel.User
	if err != nil {
		return nil, err
	}

	for _, us := range users {
		modelUser := userModel.FromDBModel(us)
		modelUsers = append(modelUsers, modelUser)
	}

	return modelUsers, nil
}

// CreateUser creates user
func (r *UserRepository) CreateUser(ctx context.Context, us domain.IUser) error {
	serviceUser := us.(*userModel.User)
	serviceUser.ID = uuid.New().String()

	var toCreateUser = userModel.ToDBModel(serviceUser)

	err := toCreateUser.Insert(ctx, r.DB, boil.Infer())
	if err != nil {
		return err
	}

	var toCreatePassword = &caModels.UserPassword{
		HashedPassword: null.StringFrom(serviceUser.HashedPassword),
	}
	toCreatePassword.ID = uuid.New().String()
	err = toCreateUser.SetUserPassword(ctx, r.DB, true, toCreatePassword)

	return err
}

//FilterUsers filter user by selecting @field, where @query and @queryParams,
func (r *UserRepository) FilterUsers(ctx context.Context, field []string, query string, queryParams []interface{}, lastID string, limit int, sort string) domain.IUsers {
	users, err := caModels.Users(
		Select(field...),
		Where(query, queryParams...),
		Limit(limit),
	).All(ctx, r.DB)

	var modelUsers []*userModel.User
	if err != nil {
		return nil
	}

	for _, us := range users {
		modelUser := userModel.FromDBModel(us)
		modelUsers = append(modelUsers, modelUser)
	}

	return modelUsers
}

//CountUsers returns count of users in db according to query
func (r *UserRepository) CountUsers(ctx context.Context, field []string, query string, queryParams []interface{}) int64 {
	usersCount, err := caModels.Users(
		Select(field...),
		Where(query, queryParams...),
	).Count(ctx, r.DB)

	if err != nil {
		return 0
	}
	return usersCount
}

// DeleteUsers delete users with @ids
func (r *UserRepository) DeleteUsers(ctx context.Context, ids []string) error {
	tx, err := boil.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	//sqlIds = strings.Join(ids,",")
	_, err = caModels.Users(WhereIn(caModels.UserColumns.ID+" in ?", ids)).DeleteAll(ctx, tx)
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
	return err
}

//DeleteAllUsers deletes all users from database
func (r *UserRepository) DeleteAllUsers(ctx context.Context) error {
	_, err := caModels.Users().DeleteAll(ctx, r.DB)
	return err
}

//GetUserById returns user with @id
func (r *UserRepository) GetUserById(ctx context.Context, id string) (domain.IUser, error) {
	foundUser, err := caModels.FindUser(ctx, r.DB, id)
	if err != nil {
		return nil, err
	}

	err = foundUser.L.LoadUserPassword(ctx, r.DB, true, foundUser, nil)
	if err != nil {
		return nil, err
	}

	returnUser := userModel.FromDBModel(foundUser)

	return returnUser, nil
}

// GetUserByUsername finds user by username
func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (domain.IUser, error) {
	foundUser, err := caModels.Users(Where(caModels.UserColumns.Username+"=?", username)).One(ctx, r.DB)
	if err != nil {
		return nil, err
	}
	err = foundUser.L.LoadUserPassword(ctx, r.DB, true, foundUser, nil)
	if err != nil {
		return nil, err
	}

	returnUser := userModel.FromDBModel(foundUser)

	return returnUser, nil
}

// GetUserByEmail finds user by email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (domain.IUser, error) {
	foundUser, err := caModels.Users(Where(caModels.UserColumns.Email+"=?", email)).One(ctx, r.DB)
	if err != nil {
		fmt.Println("email is " + email)
		fmt.Println(err)
		return nil, err
	}
	err = foundUser.L.LoadUserPassword(ctx, r.DB, true, foundUser, nil)
	if err != nil {
		return nil, err
	}

	returnUser := userModel.FromDBModel(foundUser)

	return returnUser, nil
}

// UserExistsByUsername checks if user exists by username
func (r *UserRepository) UserExistsByUsername(ctx context.Context, username string) bool {
	foundUser, err := caModels.Users(Where(caModels.UserColumns.Username+"=?", username)).Exists(ctx, r.DB)
	if err != nil {
		return false
	}

	return foundUser
}

// UserExistsByEmail checks if user exists by email
func (r *UserRepository) UserExistsByEmail(ctx context.Context, email string) bool {
	foundUser, err := caModels.Users(Where(caModels.UserColumns.Email+"=?", email)).Exists(ctx, r.DB)
	if err != nil {
		return false
	}
	return foundUser
}

// UpdateUser updates user
func (r *UserRepository) UpdateUser(ctx context.Context, id string, _inUser domain.IUser) (domain.IUser, error) {
	inUser := _inUser.(*userModel.User)

	foundUser, err := caModels.FindUser(ctx, r.DB, id)
	if err != nil {
		return nil, err
	}

	if inUser.ID != "" {
		foundUser.ID = inUser.ID
	}
	if inUser.Username != "" {
		foundUser.Username.String = inUser.Username
	}
	if inUser.Email != "" {
		foundUser.Email.String = inUser.Email
	}
	if inUser.Status != "" {
		foundUser.Status.String = inUser.Status
	}
	if inUser.ConfirmationCode != "" {
		foundUser.ConfirmationCode.String = inUser.ConfirmationCode
	}
	if inUser.Roles != nil {
		foundUser.Roles = make([]string, len(inUser.Roles))
		for _, r := range foundUser.Roles {
			foundUser.Roles = append(foundUser.Roles, r)
		}
	}

	if inUser.HashedPassword != "" {
		err = foundUser.L.LoadUserPassword(ctx, r.DB, true, foundUser, nil)
		if err != nil {
			return nil, err
		}
		foundUser.R.UserPassword.HashedPassword = null.StringFrom(inUser.HashedPassword)
		_, err := foundUser.R.UserPassword.Update(ctx, r.DB, boil.Infer())
		if err != nil {
			return nil, err
		}
	}

	_, err = foundUser.Update(ctx, r.DB, boil.Infer())
	foundUser.Reload(ctx, r.DB)

	returnUser := userModel.FromDBModel(foundUser)
	return returnUser, nil
}

// DeleteUser deletes user
func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	foundUser, err := caModels.FindUser(ctx, r.DB, id)
	if err != nil {
		return err
	}
	_, err = foundUser.Delete(ctx, r.DB)
	return err
}
