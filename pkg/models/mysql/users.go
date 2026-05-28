package mysql

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"guthub.com/eduartepaiva/snippetbox/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, hashed_password string) (int, error) {
	stmt := `INSERT INTO users (name,email,hashed_password,created_at) VALUES (?, ?, ?, UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, name, email, hashed_password)

	if err != nil {
		// check if the error code is 1062 for unique constraint, if it's then it's a duplicate email
		mysql_err, ok := err.(*mysql.MySQLError)
		if ok && mysql_err.Number == 1062 {
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
	return 0, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
