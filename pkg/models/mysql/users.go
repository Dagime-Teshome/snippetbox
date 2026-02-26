package mysql

import (
	"database/sql"
	"strings"

	"github.com/Dagime-Teshome/snippetbox/pkg/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	Db *sql.DB
}

func (us *UserModel) Insert(name, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	sql := `insert into users (name,email, hashed_password,created) values (?,?,?,UTC_TIMESTAMP())`
	_, err = us.Db.Exec(sql, name, email, string(hash))
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 && strings.Contains(err.Error(), "email") {
				return models.ErrDuplicateEmail
			}
		}
	}
	return err
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashed_password []byte
	query := `select id,hashed_password from users where email = ?`
	row := m.Db.QueryRow(query, email)
	err := row.Scan(&id, &hashed_password)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword(hashed_password, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	usr := &models.User{}
	query := `SELECT id,name,email,created FROM users WHERE id = ?`
	err := m.Db.QueryRow(query, id).Scan(&usr.ID, &usr.Name, &usr.Email, &usr.Created)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return usr, nil

}
