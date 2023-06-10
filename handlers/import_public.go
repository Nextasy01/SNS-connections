package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
)

type responseObject struct {
	Status bool `json:"status"`
	Data   struct {
		File struct {
			Metadata struct {
				ID string `json:"id"`
			} `json:"metadata"`
		} `json:"file"`
	} `json:"data"`
}

func ImportToPublic(filename string, video *os.File) (string, error) {
	// params := url.Values{}
	envFile, err := godotenv.Read(".env")
	if err != nil {
		return "", err
	}

	// defer video.Close()

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(fw, video)
	if err != nil {
		return "", err
	}

	if err := w.Close(); err != nil {
		return "", err
	}

	// defer w.Close()
	// params.Set("token", envFile["Anonfiles_API_key"])

	// location := url.URL{Path: "https://api.anonfiles.com/upload", RawQuery: params.Encode()}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.anonfiles.com/upload?token=%s", envFile["Anonfiles_API_key"]), &buf)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	log.Println("Starting to upload file to storage")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		fmt.Println("Upload failed:", res.Status)
		return "", fmt.Errorf("bad request (status not ok)")
	}

	fmt.Println("Upload successful!")

	log.Println("Reading response body")
	respBodyData, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	log.Println(string(respBodyData))

	var responseObj responseObject
	log.Println("Unmarshalling JSON output")
	err = json.Unmarshal(respBodyData, &responseObj)
	if err != nil {
		return "", err
	}
	log.Println("Checking for status")
	if !responseObj.Status {
		return "", fmt.Errorf("error from 3rd party server side")
	}

	defer res.Body.Close()

	return responseObj.Data.File.Metadata.ID, nil
}

func GetTheDownloadLink(id string) (string, error) {
	downloadLink := ""

	res, err := http.Get("https://anonfiles.com/" + id)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, ok := s.Attr("href")
		if ok && i == 2 {
			downloadLink = link
		}
	})

	return downloadLink, nil
}
