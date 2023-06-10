package repository

import "github.com/Nextasy01/SNS-connections/entity"

type YouTubeRepository interface {
	SaveVideo(vid entity.YoutubeCandidate)
	SaveVideos(vid *[]entity.YoutubeCandidate)
	DeleteVideo(vid entity.YoutubeCandidate)
	UpdateVideo(vid entity.YoutubeCandidate)
	UpdateByYouTubeVideoId(videoId string) error
	GetVideosByAcc(uid string) (*[]entity.YoutubeCandidate, error)
}

func NewYouTubeRepository(db *Database) YouTubeRepository {
	return db
}

func (db *Database) SaveVideo(vid entity.YoutubeCandidate) {
	db.connection.Create(&vid)
}
func (db *Database) SaveVideos(vid *[]entity.YoutubeCandidate) {
	db.connection.Create(vid)
}

func (db *Database) DeleteVideo(vid entity.YoutubeCandidate) {
	db.connection.Delete(&vid)
}

func (db *Database) UpdateVideo(vid entity.YoutubeCandidate) {
	db.connection.Save(&vid)
}

func (db *Database) UpdateByYouTubeVideoId(videoId string) error {
	if err := db.connection.Model(&entity.YoutubeCandidate{}).Where("video_id = ?", videoId).Updates(entity.YoutubeCandidate{IsImported: true, IsImportedToInstagram: true}).Error; err != nil {
		return err
	}
	return nil
}

func (db *Database) GetVideosByAcc(uid string) (*[]entity.YoutubeCandidate, error) {
	videos := []entity.YoutubeCandidate{}
	if err := db.connection.Where("creator_id = ?", uid).Find(&videos).Error; err != nil {
		return nil, err
	}

	return &videos, nil
}
