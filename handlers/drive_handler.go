package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Nextasy01/SNS-connections/entity"
	"github.com/Nextasy01/SNS-connections/handlers/google"
	"github.com/Nextasy01/SNS-connections/handlers/instagram"
	"github.com/Nextasy01/SNS-connections/repository"
	"github.com/Nextasy01/SNS-connections/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kkdai/youtube/v2"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type ResponseObject struct {
	Url     []string `json:"url"`
	Title   []string `json:"title"`
	VideoId []string `json:"videoId"`
	FromSns string   `json:"from_sns"`
	ToSns   string   `json:"to_sns"`
}

type UploadObject struct {
	//Response *http.Response
	Title string
	Body  io.ReadCloser
}

type DriveHandler struct {
	yt *google.YouTubeHandler
	ih *instagram.InstagramAuthHandler
	pr repository.PostRepository
}

func New(yt *google.YouTubeHandler, ih *instagram.InstagramAuthHandler, pr repository.PostRepository) *DriveHandler {
	return &DriveHandler{yt, ih, pr}
}

func (dh *DriveHandler) DriveAPICall(c *gin.Context) {
	srv, err := drive.NewService(c, option.WithCredentialsFile("key.json"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	res, err := srv.Files.List().Do()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var files []string

	for _, file := range res.Files {
		files = append(files, file.Name)
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}

func (dh *DriveHandler) Upload(c *gin.Context) {
	user_id, err := utils.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// srv, err := drive.NewService(c, option.WithCredentialsFile("key.json"))
	// if err != nil {
	// 	log.Println(err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	// 	return
	// }

	bodyData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Couldn't read body data")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var respObj ResponseObject
	if err = json.Unmarshal(bodyData, &respObj); err != nil {
		log.Println("Couldn't unmarshal")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	log.Printf("From %s To %s\n", respObj.FromSns, respObj.ToSns)

	/*
		Check if video was already imported, if yes, then just directly upload to SNS
		It works only for Instagram, because those geniuses insist
		that your videos should be at some public and accessible server, where they can get the video.
	*/
	if respObj.ToSns == "Insta" {
		temp := make([]string, len(respObj.VideoId))
		copy(temp, respObj.VideoId)
		for i, id := range temp {
			post, err := dh.pr.GetPost(id)
			if err != nil || post.VideoId == "" {
				continue
			}
			// videoLink := "https://drive.google.com/uc?export=view&id=" + post.DriveId

			videoLink, err := GetTheDownloadLink(post.DriveId)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			}

			dh.ih.Import(c, videoLink, post.Title)
			if i == 0 {
				respObj.VideoId = respObj.VideoId[1:]
				respObj.Title = respObj.Title[1:]
				continue
			}
			respObj.VideoId = append(respObj.VideoId[:i-1], respObj.VideoId[i:]...)
			respObj.Title = append(respObj.Title[:i-1], respObj.Title[i:]...)
		}
	}

	urlChan := make(chan UploadObject, len(respObj.VideoId))
	closeChan := make(chan struct{}, len(respObj.VideoId))
	log.Printf("Got %d links from client\n", len(respObj.VideoId))

	switch respObj.FromSns {
	case "YouTube":
		for i, code := range respObj.VideoId {
			fmt.Printf("%d channel downloading file\n", i)
			go downloadFromYT(code, c, urlChan)
		}
	case "Instagram":
		for i, url := range respObj.Url {
			fmt.Printf("%d channel downloading file\n", i)
			go downloadFile(url, respObj.Title[i], urlChan)
		}
	default:
		c.JSON(http.StatusBadRequest, "Invalid request: Wrong SNS, pick valid one!")
		return
	}

	go func() {
		for i := 0; cap(closeChan) > i; i++ {
			log.Printf("closing channel #%d\n", i)
			<-closeChan
		}
		close(urlChan)
		close(closeChan)
		fmt.Println("All channels are closed!")
	}()

	log.Println("Proceeding to upload files to storage")
	for i := 0; cap(urlChan) > i; i++ {
		go func(urlChan <-chan UploadObject, currentI int) {
			var uploadObj = <-urlChan

			defer uploadObj.Body.Close()

			tmpfile, err := os.CreateTemp("", "video-*.mp4")
			if err != nil {
				log.Println("[ERROR] Couldn't create temporary file: ", err)
				//c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			defer os.Remove(tmpfile.Name())

			_, err = io.Copy(tmpfile, uploadObj.Body)
			if err != nil {
				log.Printf("[ERROR] Failed to write video to temporary file: %v\n", err)
				//c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			tmpfile.Close()

			video, err := os.Open(tmpfile.Name())
			if err != nil {
				log.Printf("[ERROR] Failed to open temporary file: %v", err)
				//c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			// defer video.Close()

			post_id := ""

			if respObj.ToSns == "Insta" {
				id, err := dh.importToInstagram(c, video, uploadObj.Title)
				if err != nil {
					log.Printf("[ERROR] Failed to import to Instagram: %v", err)
					return
				}
				if err := dh.ih.UpdateVideos(respObj.VideoId[currentI]); err != nil {
					log.Printf("[ERROR] Couldn't update database: %v\n", err)
				}
				post_id = id
			} else {
				err := dh.importToYouTube(c, video, uploadObj.Title)
				if err != nil {
					log.Printf("[ERROR] Failed to import to YouTube: %v", err)
					return
				}
				if err := dh.yt.UpdateVideos(respObj.VideoId[currentI]); err != nil {
					log.Printf("[ERROR] Couldn't update database: %v\n", err)
				}
			}

			log.Printf("Video \"%s.mp4\" uploaded successfully\n", uploadObj.Title)
			dh.createPost(post_id, uploadObj.Title, respObj.FromSns, respObj.ToSns, respObj.VideoId[currentI], user_id)
			closeChan <- struct{}{}
			video.Close()
		}(urlChan, i)
	}

	// The old variant when videos should be uploaded to Google Drive first
	// Was not a good idea tbh, but still left as it is just in case

	// fmt.Println("Proceeding to upload files to drive")
	// for i := 0; cap(urlChan) > i; i++ {
	// 	fmt.Printf("%d channel uploading file\n", i)
	// 	go func(urlChan <-chan UploadObject, currentI int) {

	// 		var uploadObj = <-urlChan
	// 		defer uploadObj.Body.Close()

	// 		tmpfile, err := os.CreateTemp("", "video-*.mp4")
	// 		if err != nil {
	// 			log.Println("[ERROR] Couldn't create temporary file: ", err)
	// 			//c.AbortWithError(http.StatusInternalServerError, err)
	// 			return
	// 		}
	// 		defer os.Remove(tmpfile.Name())

	// 		_, err = io.Copy(tmpfile, uploadObj.Body)
	// 		if err != nil {
	// 			log.Printf("[ERROR] Failed to write video to temporary file: %v\n", err)
	// 			//c.AbortWithError(http.StatusInternalServerError, err)
	// 			return
	// 		}
	// 		tmpfile.Close()

	// 		video, err := os.Open(tmpfile.Name())
	// 		if err != nil {
	// 			log.Printf("[ERROR] Failed to open temporary file: %v", err)
	// 			//c.AbortWithError(http.StatusInternalServerError, err)
	// 			return
	// 		}
	// 		defer video.Close()

	// 		file := &drive.File{Name: fmt.Sprintf("%s.mp4", uploadObj.Title), Parents: []string{"1Ytu2RGbIaQKA2f-ETWext45doJPrYae7"}}

	// 		log.Printf("Started uploading: \"%s.mp4\"\n", uploadObj.Title)

	// 		driveFile, err := srv.Files.Create(file).Media(video).Do()
	// 		if err != nil {
	// 			log.Printf("[ERROR] Couldn't upload: %v\n", err)
	// 			//c.AbortWithError(http.StatusInternalServerError, err)
	// 			return
	// 		}
	// 		if respObj.FromSns == "YouTube" {
	// 			if err := dh.yt.UpdateVideos(respObj.VideoId[currentI]); err != nil {
	// 				log.Printf("[ERROR] Couldn't update database: %v\n", err)
	// 			}
	// 		} else {
	// 			if err := dh.ih.UpdateVideos(respObj.VideoId[currentI]); err != nil {
	// 				log.Printf("[ERROR] Couldn't update database: %v\n", err)
	// 			}
	// 		}

	// 		log.Printf("Video \"%s.mp4\" uploaded successfully\n", uploadObj.Title)
	// 		dh.createPost(driveFile, respObj.FromSns, respObj.ToSns, respObj.VideoId[currentI], user_id)
	// 		closeChan <- struct{}{}
	// 	}(urlChan, i)
	// }
	c.JSON(http.StatusOK, "Successfully uploaded the videos")

}

func (dh *DriveHandler) createPost(upload_id, title, fromSns, toSns, video_id, user_id string) {
	newPost := entity.NewPost()
	newPost.ID, _ = uuid.NewRandom()
	newPost.CreatorId, _ = uuid.Parse(user_id)
	newPost.DriveId = upload_id
	newPost.VideoId = video_id
	newPost.Title = title
	newPost.PlatformTo = toSns
	newPost.PlatformFrom = fromSns
	dh.pr.CreatePost(*newPost)
}

func (dh *DriveHandler) importToYouTube(ctx *gin.Context, tmpfile *os.File, title string) error {
	err := dh.yt.Import(ctx, title, tmpfile)
	if err != nil {
		return err
	}
	return nil
}

func (dh *DriveHandler) importToInstagram(ctx *gin.Context, tmpfile *os.File, title string) (string, error) {
	ffmpeg := new(utils.FFmpegConverter)
	err := ffmpeg.ConvertToReelsFormat(tmpfile.Name())
	if err != nil {
		log.Printf("[ERROR] Failed to convert a file to REEL format: %v", err)
		return "", err
	}

	defer os.Remove(strings.Replace(tmpfile.Name(), ".", "-converted.", 1))

	convertedVideo, err := os.Open(strings.Replace(tmpfile.Name(), ".", "-converted.", 1))
	if err != nil {
		log.Printf("[ERROR] Failed to open temporary file: %v", err)
		//c.AbortWithError(http.StatusInternalServerError, err)
		return "", err
	}
	defer convertedVideo.Close()

	id, err := ImportToPublic(fmt.Sprintf("%s.mp4", title), convertedVideo)
	if err != nil {
		log.Printf("[ERROR] Failed to upload a file to storage: %v", err)
		return "", err
	}

	videoLink, err := GetTheDownloadLink(id)
	if err != nil {
		log.Printf("[ERROR] Failed to get download link: %v", err)
		return "", err
	}

	dh.ih.Import(ctx, videoLink, title)

	return id, nil
}

func downloadFile(url string, title string, urlCh chan<- UploadObject) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("[ERROR] Couldn't download from this link: ", url)
	}
	if title == "" {
		title = "Untitled"
	}

	urlCh <- UploadObject{title, resp.Body}
	fmt.Println("Dowloaded file: ", title)
	return nil
}

func downloadFromYT(code string, ctx *gin.Context, urlCh chan<- UploadObject) {
	client := youtube.Client{Debug: true, HTTPClient: http.DefaultClient}

	log.Println("Getting the video metadata")
	video, err := client.GetVideo("https://www.youtube.com/watch?v=" + code)
	if err != nil {
		log.Println("[ERROR] Couldn't get video metadata: ", err)
		return

	}

	log.Println("Downloading the video, please wait..")

	formats := video.Formats.WithAudioChannels()
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		log.Println("[ERROR] Couldn't download YouTube video: ", err)
		return
	}

	urlCh <- UploadObject{video.Title, stream}
	log.Println("Video was downloaded successfully!")
}
