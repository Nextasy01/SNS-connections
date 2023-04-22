package repository

import (
	"github.com/Nextasy01/SNS-connections/entity"
)

type GoogleRepository interface {
	SaveAcc(acc entity.GoogleAccount)
	DeleteAcc(acc entity.GoogleAccount)
	UpdateAcc(acc entity.GoogleAccount)
	GetAccById(uid string) (*entity.GoogleAccount, error)
	GetAccByUserId(uid string) (*entity.GoogleAccount, error)
}

func NewGoogleRepository(db *Database) GoogleRepository {
	return db
}

func (db *Database) SaveAcc(acc entity.GoogleAccount) {
	db.connection.Create(&acc)
}

func (db *Database) DeleteAcc(acc entity.GoogleAccount) {
	db.connection.Delete(&acc)
}

func (db *Database) UpdateAcc(acc entity.GoogleAccount) {
	db.connection.Save(&acc)
}

func (db *Database) GetAccByUserId(uid string) (*entity.GoogleAccount, error) {
	acc := entity.NewGoogleAccount()

	if err := db.connection.First(&acc, "user_id=?", uid).Error; err != nil {
		return nil, err
	}

	return acc, nil

}

func (db *Database) GetAccById(uid string) (*entity.GoogleAccount, error) {
	acc := entity.NewGoogleAccount()

	if err := db.connection.First(&acc, "id=?", uid).Error; err != nil {
		return nil, err
	}

	return acc, nil

}
