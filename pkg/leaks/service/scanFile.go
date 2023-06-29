package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	errHandler "leaks/pkg/err"
)

// scanFile checks file for viruses, returns false, if virus has been detected
func scanFile(file []byte) (bool, error) {
	url := "https://api.cloudmersive.com/virus/scan/file"
	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	part1, err := writer.CreateFormFile("inputFile", "file")
	_, err = io.Copy(part1, bytes.NewReader(file))
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't copy file to payload, err:", err))
		return true, errHandler.FileCheckErr
	}
	err = writer.Close()
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't close writer, err:", err))
		return true, errHandler.FileCheckErr
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't create a new request, err:", err))
		return true, errHandler.FileCheckErr
	}
	req.Header.Add("Content-Type", "multipart/form-data")
	req.Header.Add("Apikey", os.Getenv("CLOUD_MERSIVE_API"))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't send request, err:", err))
		return true, errHandler.FileCheckErr
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't read response body, err:", err))
		return true, errHandler.FileCheckErr
	}
	var result struct {
		Res bool `json:"CleanResult"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		logg.Error(fmt.Sprint("Couldn't unmarshal response body, err:", err))
		return true, err
	}
	return result.Res, nil
}
