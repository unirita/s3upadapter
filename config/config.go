// Copyright 2015 unirita Inc.
// Created 2015/10/20 kazami

package config

import (
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

type config struct {
	Aws awsTable
	Log logTable
}

const (
	Log_Flag_ON  = 1
	Log_Flag_OFF = 0
)

// 設定ファイルのawsテーブル
type awsTable struct {
	AccessKeyId     string `toml:"access_key_id"`
	SecletAccessKey string `toml:"secret_access_key"`
	Region          string `toml:"region"`
}

// 設定ファイルのlogテーブル
type logTable struct {
	LogDebug          int `toml:"log_on"`
	LogSigning        int `toml:"signing_on"`
	LogHTTPBody       int `toml:"httpbody_on"`
	LogRequestRetries int `toml:"request_retries_on"`
	LogRequestErrors  int `toml:"request_errors_on"`
}

var Aws = new(awsTable)
var Log = new(logTable)

// 設定ファイルをロードする。
//
// 引数: filePath ロードする設定ファイルのパス
//
// 戻り値： エラー情報
func Load(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	return loadReader(f)
}

func loadReader(reader io.Reader) error {
	c := new(config)

	if _, err := toml.DecodeReader(reader, c); err != nil {
		return err
	}

	Aws = &c.Aws
	Log = &c.Log

	return nil
}

// 設定値のエラー検出を行う。
//
// return : エラー情報
func DetectError() error {
	if Aws.AccessKeyId == "" {
		return fmt.Errorf("Aws.access_key_id is blank.")
	}

	if Aws.SecletAccessKey == "" {
		return fmt.Errorf("Aws.seclet_access_key value is not set.")
	}

	if Aws.Region == "" {
		return fmt.Errorf("Aws.region value is not set.")
	}

	if Log.LogDebug != Log_Flag_ON && Log.LogDebug != Log_Flag_OFF {
		return fmt.Errorf("Log.log_debug (%d) must be a 1 or 0.", Log.LogDebug)
	}

	if Log.LogSigning != Log_Flag_ON && Log.LogSigning != Log_Flag_OFF {
		return fmt.Errorf("Log.log_signing (%d) must be a 1 or 0.", Log.LogSigning)
	}

	if Log.LogHTTPBody != Log_Flag_ON && Log.LogHTTPBody != Log_Flag_OFF {
		return fmt.Errorf("Log.log_httpbody (%d) must be a 1 or 0.", Log.LogHTTPBody)
	}

	if Log.LogRequestRetries != Log_Flag_ON && Log.LogRequestRetries != Log_Flag_OFF {
		return fmt.Errorf("Log.log_request_retries (%d) must be a 1 or 0.", Log.LogRequestRetries)
	}

	if Log.LogRequestErrors != Log_Flag_ON && Log.LogRequestErrors != Log_Flag_OFF {
		return fmt.Errorf("Log.log_request_errors (%d) must be a 1 or 0.", Log.LogRequestErrors)
	}

	return nil
}
