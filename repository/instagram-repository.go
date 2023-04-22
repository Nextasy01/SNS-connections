package repository

import "github.com/Nextasy01/SNS-connections/entity"

type InstagramRepository interface {
	SaveInstaAcc(acc entity.InstagramAccount)
	DeleteInstaAcc(acc entity.InstagramAccount)
	UpdateInstaAcc(acc entity.InstagramAccount)
	GetInstaAccById(uid string) (*entity.InstagramAccount, error)
	GetInstaAccByUserId(uid string) (*entity.InstagramAccount, error)

	SaveInstaVideo(vid entity.InstagramCandidate)
	SaveInstaVideos(vid *[]entity.InstagramCandidate)
	DeleteInstaVideo(vid entity.InstagramCandidate)
	UpdateInstaVideo(vid entity.InstagramCandidate)
	GetInstaVideosByAcc(uid string) (*[]entity.InstagramCandidate, error)
}

func NewInstagramRepository(db *Database) InstagramRepository {
	return db
}

func (db *Database) SaveInstaAcc(acc entity.InstagramAccount) {
	db.connection.Create(&acc)
}

func (db *Database) DeleteInstaAcc(acc entity.InstagramAccount) {
	db.connection.Delete(&acc)
}

func (db *Database) UpdateInstaAcc(acc entity.InstagramAccount) {
	db.connection.Save(&acc)
}

func (db *Database) GetInstaAccByUserId(uid string) (*entity.InstagramAccount, error) {
	acc := entity.NewInstagramAccount()

	if err := db.connection.First(&acc, "user_id=?", uid).Error; err != nil {
		return nil, err
	}

	return acc, nil

}

func (db *Database) GetInstaAccById(uid string) (*entity.InstagramAccount, error) {
	acc := entity.NewInstagramAccount()

	if err := db.connection.First(&acc, "id=?", uid).Error; err != nil {
		return nil, err
	}

	return acc, nil

}

func (db *Database) SaveInstaVideo(vid entity.InstagramCandidate) {
	db.connection.Create(&vid)
}
func (db *Database) SaveInstaVideos(vid *[]entity.InstagramCandidate) {
	db.connection.Create(vid)
}
func (db *Database) DeleteInstaVideo(vid entity.InstagramCandidate) {
	db.connection.Delete(&vid)
}
func (db *Database) UpdateInstaVideo(vid entity.InstagramCandidate) {
	db.connection.Save(&vid)
}
func (db *Database) GetInstaVideosByAcc(uid string) (*[]entity.InstagramCandidate, error) {
	videos := []entity.InstagramCandidate{}
	if err := db.connection.Where("creator_id = ?", uid).Find(&videos).Error; err != nil {
		return &videos, err
	}
	return &videos, nil
}
