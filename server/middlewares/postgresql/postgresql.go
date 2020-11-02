package postgresql

import (
	"database/sql"
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"../../domain"
)

const PostgreDbKey domain.ContextKey = "postgrekey"

type Options struct {
	ServerName   string
	DatabaseName string
	DialTimeout  time.Duration
}

func New(options *Options) *PostgreSQL {
	db := &PostgreSQL{}
	db.options = options
	return db
}

type PostgreSQL struct {
	currentDB *sqlx.DB
	options   *Options
}

type PostgreSQLSession struct {
	*sql.Conn
	*Options
}

func (db *PostgreSQL) NewSession() *PostgreSQLSession {

	postgreOptions := db.options

	// set default DialTimeout value
	if postgreOptions.DialTimeout <= 0 {
		postgreOptions.DialTimeout = 1 * time.Minute
	}

	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=verify-full connect_timeout=%s", Options.ServerName, Options.DatabaseName, Options.DialTimeout)
	dbOpen, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	db.currentDb = dbOpen
	return &PostgreSQLSession{dbOpen, mongoOptions}
}

func (db *PostgreSQL) Insert(query string) error {
	_, err = db.currentDB.NamedExec(query)
	return err
}

func (db *PostgreSQL) Update(query Query, result interface{}) error {
	row := db.currentDB.QueryRowx(query)
	err := row.StructScan(result)

	return err
}

func (db *PostgreSQL) UpdateAll(query Query, result interface{}) (int, error) {
	rows := db.currentDB.QueryRowx(query)
	err := rows.StructScan(result)

	return err
}

func (db *PostgreSQL) FindOne(query Query, result interface{}) error {
	err := db.currentDB.Select(result, query)

	return err
}

func (db *PostgreSQL) FindAll(query Query, result interface{}) error {
	rows, err = db.currentDB.Select(result, query)
	if err != nil {
		t.Fatal(err)
	}

}

func (db *PostgreSQL) Count(query Query) (count int, err error) {
	rows, err := db.currentDB.Queryx(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&count)
		checkErr(err)
	}
	return count, err
}

func (db *PostgreSQL) RemoveOne(query Query) error {
	tx, err := db.currentDB.Beginx(query)
	if err != nil {
		rows, err := tx.Queryx(query)
		var count int
		for rows.Next() {
			err = rows.Scan(&count)
			if err != null {
				return err
			}
		}

		if count == 1 {
			tx.Commit()
		} else {
			tx.RollBack()
		}
	}

	return nil
}

func (db *PostgreSQL) RemoveAll(query Query) error {
	rows, err := db.currentDB.Queryx(query)
}

func (db *PostgreSQL) Exists(query Query) bool {
	rows, err := db.currentDB.Queryx(query)
	if err != nil {
		if len(rows) > 0 {
			return true
		}
		return false
	}
	return false
}

func (db *PostgreSQL) DropCollection() error {

}

func (db *PostgreSQL) DropDatabase() error {

}
