package repo

import (
	"database/sql"
	"errors"

	"github.com/bernhardson/stub/internal/models"
)

type UserRepository interface {
	Get(int64) (*models.User, error)
	GetAll() (*[]models.User, error)
	Insert(string, string, string) (int, error)
	DeleteUser(int) (int, error)
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	ClearTable(string) error
}

// create user repository implementation. so far we support either postgres or mysql.
func UserRepoFactory(datasource string, db *sql.DB) (UserRepository, error) {
	var repo UserRepository
	switch datasource {
	case "mysql":
		repo = &UserRepoMySqlImpl{db: db}
	case "postgres":
		repo = &UserRepoPostgresImpl{db: db}
	default:
		return nil, errors.New("unsupported database")
	}
	return repo, nil
}
