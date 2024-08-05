package services

import (
	"github.com/tsironis/syntropic-studio/db"
	"golang.org/x/crypto/bcrypt"
)

func NewUserServices(u User, uStore db.Store) *UserServices {

	return &UserServices{
		User:      u,
		UserStore: uStore,
	}
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type UserServices struct {
	User      User
	UserStore db.Store
}

func (us *UserServices) CreateUser(u User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		return err
	}

	user := User{Email: u.Email, Password: string(hashedPassword), Username: u.Username}

	return us.UserStore.Db.Create(&user).Error
}

func (us *UserServices) CheckEmail(email string) (User, error) {
	user := User{Email: email}
	result := us.UserStore.Db.Where(&user, "email").Find(&user)
	if result.Error != nil {
		return User{}, result.Error
	}

	us.User = user

	return us.User, nil
}

/* func (us *UserServices) GetUserById(id int) (User, error) {

	query := `SELECT id, email, password, username FROM users
		WHERE id = ?`

	stmt, err := us.UserStore.Db.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	us.User.ID = id
	err = stmt.QueryRow(
		us.User.ID,
	).Scan(
		&us.User.ID,
		&us.User.Email,
		&us.User.Password,
		&us.User.Username,
	)
	if err != nil {
		return User{}, err
	}

	return us.User, nil
} */
