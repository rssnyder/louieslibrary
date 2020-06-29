package main

import (
	"bytes"
	"os"
	"io/ioutil"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (app *App) UploadBytes(bucket, key string, data []byte) error {
	storage_connection := s3.New(app.Storage)

	// Upload a new object
	_, err := storage_connection.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(data),
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}

func (app *App) UploadFile(bucket, key, filename string) error {
	storage_connection := s3.New(app.Storage)

	// Read file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// Upload a new object
	_, err = storage_connection.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(data),
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}

func (app *App) DownloadObject(bucket, key, destination string) error {
	// Retrieve object
	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(app.Storage)
	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	
	return nil
}

func (app *App) DownloadBytes(bucket, key string) ([]byte, error) {
	var output []byte
	storage_connection := s3.New(app.Storage)

	result, err := storage_connection.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return output, err
	}

	output, err = ioutil.ReadAll(result.Body)
	if err != nil {
		return output, err
	}

	return output, nil
}

func (app *App) ServeFile(w http.ResponseWriter, bucket, key, name string) {
	data, err := app.DownloadBytes(bucket, key)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", name))
	w.Write(data)
}