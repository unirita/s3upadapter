// Copyright 2015 unirita Inc.
// Created 2015/10/20 kazami

package upload

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/unirita/s3upadapter/config"
)

//モックのインタフェース
type UploadMock interface {
	Upload(bucket string, uploadKey string, localPath string) error
	getS3Instance() *s3.S3
	uploadFile(bucket string, uploadKey string, localPath string) error
}

func Upload(bucket string, uploadKey string, localPath string) error {
	//設定ファイルの情報を与えてS3のインスタンスを作成する
	client := getS3Instance()

	var upLocation string
	_, fileName := filepath.Split(localPath)

	if uploadKey != "" {
		params := &s3.GetObjectInput{Bucket: &bucket, Key: &uploadKey}
		_, connectErr := client.GetObject(params)
		if connectErr != nil {
			return connectErr
		}
	}

	upLocation = uploadKey + fileName

	if err := uploadFile(bucket, upLocation, localPath); err != nil {
		return err
	}

	return nil
}

func uploadFile(bucket string, uploadKey string, localPath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	u := s3manager.NewUploader(nil)
	result, err := u.Upload(&s3manager.UploadInput{
		Bucket: &bucket,
		Key:    &uploadKey,
		Body:   file})

	if err != nil {
		return err
	}

	fmt.Printf("Uploaded [%s]", result.Location)
	return nil
}

//S3のインスタンスを取得する
func getS3Instance() *s3.S3 {
	defaults.DefaultConfig.Credentials = credentials.NewStaticCredentials(config.Aws.AccessKeyId, config.Aws.SecletAccessKey, "")
	defaults.DefaultConfig.Region = &config.Aws.Region

	return s3.New(createConf())
}

func createConf() *aws.Config {
	conf := aws.NewConfig()

	if config.Log.LogDebug == config.Log_Flag_OFF {
		conf.WithLogLevel(aws.LogOff)
	} else {

		loglevel := aws.LogDebug

		if config.Log.LogSigning == config.Log_Flag_ON {
			loglevel |= aws.LogDebugWithSigning
		}

		if config.Log.LogHTTPBody == config.Log_Flag_ON {
			loglevel |= aws.LogDebugWithHTTPBody
		}

		if config.Log.LogRequestRetries == config.Log_Flag_ON {
			loglevel |= aws.LogDebugWithRequestRetries
		}

		if config.Log.LogRequestErrors == config.Log_Flag_ON {
			loglevel |= aws.LogDebugWithRequestErrors
		}

		conf.WithLogLevel(loglevel)
	}
	return conf
}
