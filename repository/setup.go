package repository

import (
	"github.com/Nextasy01/SNS-connections/entity"
	//"github.com/jinzhu/gorm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	connection *gorm.DB
}

func NewDatabase() Database {
	return Database{
		connection: ConnectToDB(),
	}
}

func ConnectToDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("faild to connect to database")
	}
	db.AutoMigrate(&entity.User{}, &entity.GoogleAccount{}, &entity.InstagramAccount{}, &entity.YoutubeCandidate{}, &entity.InstagramCandidate{}, &entity.Post{})
	db.Exec("PRAGMA foreign_keys = ON")
	return db
}
