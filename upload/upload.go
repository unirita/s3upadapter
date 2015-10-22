// Copyright 2015 unirita Inc.
// Created 2015/10/20 kazami

package upload

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func Do(bucket string, key string, localPath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	u := s3manager.NewUploader(nil)
	result, err := u.Upload(&s3manager.UploadInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   file})
	if err != nil {
		return err
	}

	fmt.Printf("Uploaded [%s]", result.Location)
	return nil
}
