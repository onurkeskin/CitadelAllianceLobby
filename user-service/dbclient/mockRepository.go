package database

import (
	"context"

	"github.com/stretchr/testify/mock"
	"keon.com/CitadelAllianceLobbyServer/user-service/domain"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(ctx context.Context, user domain.IUser) error {
	args := m.Mock.Called(user)
	return args.Error(0)
}

func (m *MockRepository) GetUsers(ctx context.Context) (domain.IUsers, error) {
	args := m.Mock.Called()
	return args.Get(0).(domain.IUsers), args.Error(1)
}

func (m *MockRepository) FilterUsers(ctx context.Context, field []string, query string, queryParams []interface{}, lastID string, limit int, sort string) domain.IUsers {
	args := m.Mock.Called()
	return args.Get(0).(domain.IUsers)
}

func (m *MockRepository) CountUsers(ctx context.Context, field []string, query string, queryParams []interface{}) int64 {
	args := m.Mock.Called()
	return args.Get(0).(int64)
}

func (m *MockRepository) DeleteUsers(ctx context.Context, ids []string) error {
	args := m.Mock.Called(ids)
	return args.Error(0)
}
func (m *MockRepository) DeleteAllUsers(ctx context.Context) error {
	args := m.Mock.Called()
	return args.Error(0)
}

func (m *MockRepository) GetUserById(ctx context.Context, id string) (domain.IUser, error) {
	args := m.Mock.Called(id)
	var r1 domain.IUser
	var r2 error
	if args.Get(0) == nil {
		r1 = nil
	} else {
		r1 = args.Get(0).(domain.IUser)
	}
	if args.Error(1) == nil {
		r2 = nil
	} else {
		r2 = args.Error(1)
	}
	return r1, r2
}
func (m *MockRepository) GetUserByUsername(ctx context.Context, username string) (domain.IUser, error) {
	args := m.Mock.Called(username)
	var r1 domain.IUser
	var r2 error
	if args.Get(0) == nil {
		r1 = nil
	} else {
		r1 = args.Get(0).(domain.IUser)
	}
	if args.Error(1) == nil {
		r2 = nil
	} else {
		r2 = args.Error(1)
	}
	return r1, r2
}
func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (domain.IUser, error) {
	args := m.Mock.Called(email)

	var r1 domain.IUser
	var r2 error
	if args.Get(0) == nil {
		r1 = nil
	} else {
		r1 = args.Get(0).(domain.IUser)
	}
	if args.Error(1) == nil {
		r2 = nil
	} else {
		r2 = args.Error(1)
	}
	return r1, r2
}

func (m *MockRepository) UserExistsByUsername(ctx context.Context, username string) bool {
	args := m.Mock.Called(username)
	return args.Get(0).(bool)
}

func (m *MockRepository) UserExistsByEmail(ctx context.Context, email string) bool {
	args := m.Mock.Called(email)
	return args.Get(0).(bool)
}

func (m *MockRepository) UpdateUser(ctx context.Context, id string, inUser domain.IUser) (domain.IUser, error) {
	args := m.Mock.Called(inUser)
	return args.Get(0).(domain.IUser), args.Error(1)
}

func (m *MockRepository) DeleteUser(ctx context.Context, id string) error {
	args := m.Mock.Called(id)
	return args.Error(0)
}
