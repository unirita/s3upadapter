// Copyright 2015 unirita Inc.
// Created 2015/10/20 kazami

package upload

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type uploadFunc func(*s3manager.UploadInput) (*s3manager.UploadOutput, error)

var upload uploadFunc = s3manager.NewUploader(nil).Upload

func Do(bucket string, key string, localPath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	result, err := upload(&s3manager.UploadInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   file})
	if err != nil {
		return err
	}

	fmt.Printf("Uploaded [%s]", result.Location)
	return nil
}
