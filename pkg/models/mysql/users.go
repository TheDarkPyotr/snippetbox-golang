package mysql

import (
	"database/sql"
	"snippetbox/pkg/models"
	"strings"

	"github.com/go-sql-driver/mysql" // New import
	"golang.org/x/crypto/bcrypt"
	// New import
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	statement := `INSERT INTO users (name, email, hashed_password, created)
				  VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(statement, name, email, string(hashedPassword))
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
	}
	return err
}

func (m *UserModel) Authenticate(email, password string) (int, error) {

	var id int
	var hashedPassowrd []byte
	row := m.DB.QueryRow("SELECT id, hashed_password FROM users WHERE email = ?", email)
	err := row.Scan(&id, &hashedPassowrd)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassowrd, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials

	} else if err != nil {
		return 0, err
	}

	return id, nil

}

func (m *UserModel) Get(id int) (*models.User, error) {
	s := &models.User{}
	statemnt := `SELECT id, name, email, created FROM users WHERE id = ?`
	err := m.DB.QueryRow(statemnt, id).Scan(&s.ID, &s.Name, &s.Email, &s.Created)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return s, nil
}
