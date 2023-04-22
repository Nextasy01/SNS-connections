package repository

import (
	"errors"

	"github.com/Nextasy01/SNS-connections/entity"
	"github.com/Nextasy01/SNS-connections/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	SaveUser(user entity.User) error
	DeleteUser(user entity.User)
	GetUserByID(uid string) (*entity.User, error)
	LoginCheck(username, password string) (string, error)
}

func NewUserRepository(db *Database) UserRepository {
	return db
}

func (db *Database) SaveUser(user entity.User) error {
	err := db.connection.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) DeleteUser(user entity.User) {
	db.connection.Delete(&user)
}

func (db *Database) LoginCheck(username, password string) (string, error) {
	u := entity.User{}

	err := db.connection.Model(entity.User{}).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return "", err
	}

	if err := VerifyPassword(password, u.Password); err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := utils.GenerateToken(u.ID)

	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (db *Database) GetUserByID(uid string) (*entity.User, error) {

	u := entity.NewUser()

	if err := db.connection.First(&u, "id=?", uid).Error; err != nil {
		return nil, errors.New("user not found")
	}

	u.PrepareGive()

	return u, nil

}
