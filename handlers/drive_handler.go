package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Nextasy01/SNS-connections/handlers/google"
	"github.com/Nextasy01/SNS-connections/handlers/instagram"
	"github.com/gin-gonic/gin"
	"github.com/kkdai/youtube/v2"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type ResponseObject struct {
	Url     []string `json:"url"`
	Title   []string `json:"title"`
	VideoId []string `json:"videoId"`
	Sns     string   `json:"sns"`
}

type UploadObject struct {
	//Response *http.Response
	Title string
	Body  io.ReadCloser
}

type DriveHandler struct {
	gh *google.GoogleAuthHandler
	ih *instagram.InstagramAuthHandler
}

func New(gh *google.GoogleAuthHandler, ih *instagram.InstagramAuthHandler) *DriveHandler {
	return &DriveHandler{gh, ih}
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
	log.Println("test")
	srv, err := drive.NewService(c, option.WithCredentialsFile("key.json"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
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
	log.Println(respObj.Sns)

	urlChan := make(chan UploadObject, len(respObj.VideoId))
	closeChan := make(chan struct{}, len(respObj.VideoId))
	log.Printf("Got %d links from client\n", len(respObj.VideoId))

	switch respObj.Sns {
	case "toggleYouTube":
		for i, code := range respObj.VideoId {
			fmt.Printf("%d channel downloading file\n", i)
			go downloadFromYT(code, c, urlChan)
		}
	case "toggleInsta":
		for i, url := range respObj.Url {
			fmt.Printf("%d channel downloading file\n", i)
			go downloadFile(url, respObj.Title[i], urlChan)
		}
	default:
		c.JSON(http.StatusBadRequest, "Invalid request: Wrong SNS, pick valid one!")
		return
	}

	// if respObj.Sns == "toggleYouTube" {
	// 	log.Println("Proceeding to download YouTube video")
	// 	if err := downloadFromYT(respObj.VideoId[0], c, urlChan); err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// 	return
	// }

	go func() {
		for i := 0; cap(closeChan) > i; i++ {
			log.Printf("closing channel #%d\n", i)
			<-closeChan
		}
		close(urlChan)
		close(closeChan)
		fmt.Println("All channels are closed!")
	}()

	fmt.Println("Proceeding to upload files to drive")
	for i := 0; cap(urlChan) > i; i++ {
		fmt.Printf("%d channel uploading file\n", i)
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
			defer video.Close()

			file := &drive.File{Name: uploadObj.Title, Parents: []string{"1Ytu2RGbIaQKA2f-ETWext45doJPrYae7"}}

			log.Printf("Started uploading: \"%s.mp4\"\n", uploadObj.Title)

			_, err = srv.Files.Create(file).Media(video).Do()
			if err != nil {
				log.Printf("[ERROR] Couldn't upload: %v\n", err)
				//c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			if err := dh.ih.UpdateVideos(respObj.VideoId[currentI]); err != nil {
				log.Printf("[ERROR] Couldn't update database: %v\n", err)
			}
			log.Printf("Video \"%s.mp4\" uploaded successfully\n", uploadObj.Title)
			closeChan <- struct{}{}
		}(urlChan, i)
	}
	c.JSON(http.StatusOK, "Successfully uploaded the videos")

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

func downloadFromYT(code string, ctx *gin.Context, urlCh chan<- UploadObject) error {
	client := youtube.Client{Debug: true, HTTPClient: http.DefaultClient}

	log.Println("Getting the video metadata")
	video, err := client.GetVideo("https://www.youtube.com/watch?v=" + code)
	if err != nil {
		log.Println("[ERROR] Couldn't get video metadata: ", err)
	}

	log.Println("Downloading the video, please wait..")

	formats := video.Formats.WithAudioChannels()
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		log.Println("[ERROR] Couldn't download YouTube video: ", err)
	}

	urlCh <- UploadObject{video.Title, stream}

	// file, err := os.Create(fmt.Sprintf("%s.mp4", video.Title))
	// if err != nil {
	// 	log.Println("[ERROR] Couldn't create file: ", err)
	// }
	// defer file.Close()

	// _, err = io.Copy(file, stream)
	// if err != nil {
	// 	log.Println("[ERROR] Couldn't copy file contents: ", err)
	// }

	log.Println("Video was downloaded successfully!")
	return nil
}
