package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// UploadBytes
// Take bytes and put into a bucket object
func (app *App) UploadBytes(bucket, key string, data []byte) error {

	// Connection to s3 server
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

	// Return no error
	return nil
}

// UploadFile
// Take local file and put into a bucket object
func (app *App) UploadFile(bucket, key, filename string) error {

	// Connection to s3 server
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

	// Return no error
	return nil
}

// DownloadObject
// Gets an object from s3 to a local file
func (app *App) DownloadObject(bucket, key, destination string) error {

	// Retrieve object
	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()

	// New s3 downloader
	downloader := s3manager.NewDownloader(app.Storage)
	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	// Return no error
	return nil
}

// DownloadBytes
// Get an object from a bucket and save the bytes
func (app *App) DownloadBytes(bucket, key string) ([]byte, error) {

	// Empty bytes for object
	var output []byte

	// Connection to s3 server
	storage_connection := s3.New(app.Storage)

	// Get object
	result, err := storage_connection.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return output, err
	}

	// Read object bytes
	output, err = ioutil.ReadAll(result.Body)
	if err != nil {
		return output, err
	}

	// Return the bytes with no error
	return output, nil
}

// ServeFile
// Takes bucket object and sends to user
func (app *App) ServeFile(w http.ResponseWriter, bucket, key, name string) {

	// Download object bytes
	data, err := app.DownloadBytes(bucket, key)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Set header for file download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", name))

	// Write headers and bytes to user
	w.Write(data)
}

// FindObject
// Take bytes and put into a bucket object
func (app *App) FindObject(bucket, key string) (string, error) {

	var mismatch string

	// Connection to s3 server
	storageConnection := s3.New(app.Storage)

	// Upload a new object
	resp, err := storageConnection.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		return mismatch, err
	}

	for _, item := range resp.Contents {
		id := strings.Split(*item.Key, ".")[0]
		if id == key {
			return *item.Key, nil
		}
	}

	// Return no error
	return mismatch, nil
}
