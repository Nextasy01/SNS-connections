package repository

import "github.com/Nextasy01/SNS-connections/entity"

type PostRepository interface {
	CreatePost(post entity.Post)
	DeletePost(post entity.Post)
	UpdatePost(post entity.Post)
	GetPost(videoId string) (*entity.Post, error)
}

func NewPostRepository(db *Database) PostRepository {
	return db
}

func (db *Database) CreatePost(post entity.Post) {
	db.connection.Create(&post)
}
func (db *Database) DeletePost(post entity.Post) {
	db.connection.Delete(&post)
}

func (db *Database) UpdatePost(post entity.Post) {
	db.connection.Save(&post)
}

func (db *Database) GetPost(videoId string) (*entity.Post, error) {
	post := entity.NewPost()
	err := db.connection.Where("video_id = ?", videoId).Find(&post).Error
	if err != nil {
		return nil, err
	}
	return post, nil
}
