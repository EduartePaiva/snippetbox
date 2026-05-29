package mysql

import (
	"database/sql"
	"fmt"
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
	stmt := "SELECT id, hashed_password FROM users WHERE users.email='?'"

	var row any

	err := m.DB.QueryRow(stmt, email).Scan(&row)
	fmt.Println(row)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
