package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/bernhardson/stub/internal/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserRepoMySqlImpl struct {
	db *sql.DB
}

func (m *UserRepoMySqlImpl) Insert(name, email, password string) (int, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return -1, err
	}

	query := `INSERT INTO users (username, email, hashed_password, created)
    VALUES(?, ?, ?, UTC_TIMESTAMP())`

	res, err := m.db.Exec(query, name, email, string(hashedPassword))
	if err != nil {

		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return -1, models.ErrDuplicateEmail
			}
		}
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// Get retrieves a user by ID from the MySQL database.
func (r *UserRepoMySqlImpl) Get(id int64) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, email, hashed_password, created FROM users WHERE id = ?"
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Created)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with id %d", id)
		}
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}
	return &user, nil
}

// GetAll retrieves all users from the MySQL database.
func (r *UserRepoMySqlImpl) GetAll() (*[]models.User, error) {
	query := "SELECT id, username, email, hashed_password, create FROM users"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Created)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}
	return &users, nil
}

// DeleteUser removes a user by ID from the MySQL database.
func (r *UserRepoMySqlImpl) DeleteUser(id int) (int, error) {
	query := "DELETE FROM users WHERE id = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return id, fmt.Errorf("failed to delete user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return id, fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return id, fmt.Errorf("no user found with id %d", id)
	}
	return id, nil
}

func (m *UserRepoMySqlImpl) Authenticate(email, password string) (int, error) {

	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"

	err := m.db.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

// We'll use the Exists method to check if a user exists with a specific ID.

func (m *UserRepoMySqlImpl) Exists(id int) (bool, error) {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	err := m.db.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

// ClearTable truncates a specified table in the MySQL database.
func (r *UserRepoMySqlImpl) ClearTable(table string) error {
	// Using fmt.Sprintf here for dynamic table name injection (safe for admin-controlled tables)
	query := fmt.Sprintf("TRUNCATE TABLE %s", table)
	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to clear table %s: %v", table, err)
	}
	return nil
}
