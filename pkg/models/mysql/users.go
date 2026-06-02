package mysql

import (
	"database/sql"
	"strings"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"guthub.com/eduartepaiva/snippetbox/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) (int, error) {
	stmt := `INSERT INTO users (name,email,hashed_password,created_at) VALUES (?, ?, ?, UTC_TIMESTAMP())`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err

	}

	result, err := m.DB.Exec(stmt, name, email, hashedPassword)

	if err != nil {
		// check if the error code is 1062 for unique constraint, if it's then it's a duplicate email
		mysqlErr, ok := err.(*mysql.MySQLError)
		if ok && mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "users_uc_email") {
			return 0, models.ErrDuplicateEmail
		}

		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	stmt := "SELECT id, hashed_password FROM users WHERE users.email=?"

	var row struct {
		ID             int
		HashedPassword []byte
	}

	err := m.DB.QueryRow(stmt, email).Scan(&row.ID, &row.HashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(row.HashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	return row.ID, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	stmt := "SELECT id, name, email, created_at FROM users WHERE users.id = ?"

	user := new(models.User)

	err := m.DB.QueryRow(stmt, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return user, nil
}
