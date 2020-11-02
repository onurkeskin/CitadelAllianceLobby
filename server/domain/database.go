package domain

import ()

type Query string

// Database interface
type IDatabase interface {
	Insert(obj interface{}) error
	Update(query Query, result interface{}) error
	UpdateAll(query Query) (int, error)
	FindOne(query Query, result interface{}) error
	FindAll(query Query, result interface{}, limit int, sort string) error
	Count(query Query) (int, error)
	RemoveOne(query Query) error
	RemoveAll(query Query) error
	Exists(query Query) bool
	DropCollection() error
	DropDatabase() error
}
