// Copyright 2015 unirita Inc.
// Created 2015/10/20 kazami

package console

import (
	"fmt"
)

// USAGE表示用の定義メッセージ
const USAGE = `Usage :
    s3upadapter.exe [-v] [-b Bucket] [-k Key] [-l Localfile] [-c ConfigPath]
Option :
    -v                   :  　Print s3upadapter version.
    -b bucket name       :  　Designate a buket name.
    -k key name          :  　Designate a key name.
    -l local Path        :  　Designate a local Path.
    -c configFile Path   :  　Designate configFile Path.

    -b, -l, -c is a required input.
Copyright 2015 unirita Inc.
`

// コンソールメッセージ一覧
var msgs = map[string]string{
	"UPA001E": "DIRECTORY CAN NOT BE SPECIFIED TO -k OPTION.",
	"UPA002E": "FAILED TO READ CONFIG FILE. [%s]",
	"UPA003E": "CONFIG PARM IS NOT EXACT FORMAT. [%s]",
	"UPA004E": "UPLOAD FAILED. [%s]",
}

// 標準出力へメッセージコードcodeに対応したメッセージを表示する。
//
// param : code メッセージコードID。
//
// return : 出力文字数。
//
// return : エラー情報。
func Display(code string, a ...interface{}) (int, error) {
	msg := GetMessage(code, a...)

	return fmt.Println(msg)
}

// 出力メッセージを文字列型で取得する。
//
// param : code メッセージコードID。
//
//
// return : 取得したメッセージ
func GetMessage(code string, a ...interface{}) string {
	return fmt.Sprintf("%s %s", code, fmt.Sprintf(msgs[code], a...))
}
