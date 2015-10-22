package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"

	"github.com/unirita/s3upadapter/config"
	"github.com/unirita/s3upadapter/console"
	"github.com/unirita/s3upadapter/upload"
)

// 実行時引数のオプション
type arguments struct {
	versionFlag bool   // バージョン情報表示フラグ
	bucketName  string //バケット名
	keyName     string //S3にアップロードする場所のキー名
	filePath    string //アップロード元のローカルファイル
	configPath  string //設定ファイルのパス
}

//バッチプログラムの戻り値
const (
	rc_OK    = 0
	rc_ERROR = 1
)

// フラグ系実行時引数のON/OFF
const (
	flag_ON  = true
	flag_OFF = false
)

type uploadFunc func(string, string, string) error

var doUpload uploadFunc = upload.Upload

func main() {
	args := fetchArgs()
	rc := realMain(args)
	os.Exit(rc)
}

func realMain(args *arguments) int {
	if args.versionFlag == flag_ON {
		showVersion()
		return rc_OK
	}

	if args.bucketName == "" || args.keyName == "" || args.filePath == "" || args.configPath == "" {
		showUsage()
		return rc_ERROR
	}

	if strings.HasSuffix(args.keyName, "/") {
		console.Display("UPA001E")
		return rc_ERROR
	}

	if err := config.Load(args.configPath); err != nil {
		console.Display("UPA002E", err)
		return rc_ERROR
	}

	if err := config.DetectError(); err != nil {
		console.Display("UPA003E", err)
		return rc_ERROR
	}

	defaults.DefaultConfig.Credentials = credentials.NewStaticCredentials(config.Aws.AccessKeyId, config.Aws.SecletAccessKey, "")
	defaults.DefaultConfig.Region = &config.Aws.Region
	if err := doUpload(args.bucketName, args.keyName, args.filePath); err != nil {
		console.Display("UPA004E", err)
		return rc_ERROR
	}

	return rc_OK
}

// コマンドライン引数を解析し、arguments構造体を返す。
func fetchArgs() *arguments {
	flag.Usage = showUsage
	args := new(arguments)
	flag.BoolVar(&args.versionFlag, "v", false, "version option")
	flag.StringVar(&args.bucketName, "b", "", "Designate bucket option")
	flag.StringVar(&args.keyName, "k", "", "Designate key name option")
	flag.StringVar(&args.filePath, "f", "", "Designate file path option")
	flag.StringVar(&args.configPath, "c", "", "Designate config file option")
	flag.Parse()

	return args
}

// バージョンを表示する。
func showVersion() {
	fmt.Printf("%s\n", Version)
}

// オンラインヘルプを表示する。
func showUsage() {
	fmt.Print(console.USAGE)
}
