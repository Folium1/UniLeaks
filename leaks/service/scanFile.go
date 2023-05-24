package service

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"uniLeaks/leaks"
)

type scanResult struct {
	Result bool `json:"CleanResult"`
}

// scanFile checks file for viruses, returns false, if virus has been detected
func scanFile(file []byte) (bool, error) {
	url := "https://api.cloudmersive.com/virus/scan/file"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	part1, err := writer.CreateFormFile("inputFile", "file")
	_, err = io.Copy(part1, bytes.NewReader(file))
	if err != nil {
		log.Println(err)
		return true, leaks.FileCheckErr
	}
	err = writer.Close()
	if err != nil {
		log.Println(err)
		return true, leaks.FileCheckErr
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err)
		return true, leaks.FileCheckErr
	}
	req.Header.Add("Content-Type", "multipart/form-data")
	req.Header.Add("Apikey", os.Getenv("CLOUD_MERSIVE_API"))

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return true, leaks.FileCheckErr
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return true, leaks.FileCheckErr
	}
	var result scanResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println(err)
		return true, err
	}
	log.Println("File scan result:", result.Result)
	return result.Result, nil
}
