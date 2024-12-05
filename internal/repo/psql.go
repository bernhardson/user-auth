package repo

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bernhardson/stub/internal/models"
)

type UserRepoPostgresImpl struct {
	db *sql.DB
}

// Insert creates a new user in the database.
func (p *UserRepoPostgresImpl) Insert(name, email, password string) (int, error) {
	query := `INSERT INTO users (username, email, password, create) VALUES ($1, $2, $3, $4)`
	res, err := p.db.Exec(query, name, email, password, time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// Get retrieves a user by ID.
func (p *UserRepoPostgresImpl) Get(id int64) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password, create FROM users WHERE id = $1`
	err := p.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Created)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with id %d", id)
		}
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}
	return &user, nil
}

// GetAll retrieves all users.
func (p *UserRepoPostgresImpl) GetAll() (*[]models.User, error) {
	query := `SELECT id, username, email, password, create FROM users`
	rows, err := p.db.Query(query)
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

// DeleteUser removes a user by ID.
func (p *UserRepoPostgresImpl) DeleteUser(id int) (int, error) {
	query := `DELETE FROM users WHERE id = $1`
	result, err := p.db.Exec(query, id)
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

func (p *UserRepoPostgresImpl) Authenticate(email, password string) (int, error) { return 0, nil }

// We'll use the Exists method to check if a user exists with a specific ID.

func (p *UserRepoPostgresImpl) Exists(id int) (bool, error) { return false, nil }

// ClearTable truncates a table.
func (p *UserRepoPostgresImpl) ClearTable(table string) error {
	query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)
	_, err := p.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to clear table %s: %v", table, err)
	}
	return nil
}
