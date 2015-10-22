package main

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/unirita/s3upadapter/testutil"
	"github.com/unirita/s3upadapter/upload"
)

func makeUploadSuccess() {
	doUpload = func(bucket string, key string, localPath string) error {
		return nil
	}
}

func makeUploadFail() {
	doUpload = func(bucket string, key string, localPath string) error {
		return errors.New("error")
	}
}

func restoreUploadFunc() {
	doUpload = upload.Do
}

func TestFetchArgs_コマンドラインオプションを取得できる(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.PanicOnError)
	os.Args = os.Args[:1]
	os.Args = append(os.Args, "-v", "-b", "bucket", "-k", "testkey", "-f", "localpath", "-c", "test.ini")
	args := fetchArgs()

	if args.versionFlag != flag_ON {
		t.Error("-vオプションの指定を検出できなかった。")
	}

	if args.bucketName != "bucket" {
		t.Error("-bオプションの指定を検出できなかった。")
	}

	if args.keyName != "testkey" {
		t.Error("-kオプションの指定を検出できなかった。")
	}

	if args.filePath != "localpath" {
		t.Error("-fオプションの指定を検出できなかった")
	}

	if args.configPath != "test.ini" {
		t.Error("-cオプションの指定を検出できなかった。")
	}

}

func TestFetchArgs_コマンドラインオプションに値が指定されなかった場合(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.PanicOnError)
	os.Args = os.Args[:1]
	args := fetchArgs()

	if args.versionFlag != flag_OFF {
		t.Error("-vオプションの値が想定と異なっている。")
	}

	if args.bucketName != "" {
		t.Error("-bオプションの値が想定と異なっている。")
	}

	if args.keyName != "" {
		t.Error("-kオプションの値が想定と異なっている。")
	}

	if args.filePath != "" {
		t.Error("-fオプションの値が想定と異なっている。")
	}

	if args.configPath != "" {
		t.Error("-cオプションの値が想定と異なっている。")
	}
}

func TestRealMain_バージョン出力オプションが指定された場合(t *testing.T) {
	c := testutil.NewStdoutCapturer()

	args := new(arguments)
	args.versionFlag = flag_ON

	c.Start()
	rc := realMain(args)
	out := c.Stop()

	if rc != rc_OK {
		t.Errorf("想定外のrc[%d]が返された。", rc)
	}
	if !strings.Contains(out, Version) {
		t.Error("出力内容が想定と違っている。")
		t.Logf("出力: %s", out)
	}
}

func TestRealMain_引数に何も指定されなかった場合(t *testing.T) {
	c := testutil.NewStdoutCapturer()

	args := new(arguments)

	c.Start()
	rc := realMain(args)
	out := c.Stop()

	if rc != rc_ERROR {
		t.Errorf("想定外のrc[%d]が返された。", rc)
	}
	if !strings.Contains(out, "Usage :") {
		t.Error("出力内容が想定と違っている。")
		t.Logf("出力: %s", out)
	}
}

func TestRealMain_引数にバケット名が指定されていない場合(t *testing.T) {
	c := testutil.NewStdoutCapturer()

	args := new(arguments)
	args.keyName = "uploadlocation"
	args.filePath = "localpath"
	args.configPath = "config.ini"

	c.Start()
	rc := realMain(args)
	out := c.Stop()

	if rc != rc_ERROR {
		t.Errorf("想定外のrc[%d]が返された。", rc)
	}
	if !strings.Contains(out, "Usage :") {
		t.Error("出力内容が想定と違っている。")
		t.Logf("出力: %s", out)
	}
}

func TestRealMain_引数に設定ファイルのパスが指定されていない場合(t *testing.T) {
	c := testutil.NewStdoutCapturer()

	args := new(arguments)
	args.bucketName = "bucket"
	args.keyName = "uploadlocation"
	args.filePath = "localpath"

	c.Start()
	rc := realMain(args)
	out := c.Stop()

	if rc != rc_ERROR {
		t.Errorf("想定外のrc[%d]が返された。", rc)
	}
	if !strings.Contains(out, "Usage :") {
		t.Error("出力内容が想定と違っている。")
		t.Logf("出力: %s", out)
	}
}

func TestRealMain_引数にキー名が指定されていない場合(t *testing.T) {
	c := testutil.NewStdoutCapturer()

	args := new(arguments)
	args.bucketName = "bucket"
	args.filePath = "localpath"
	args.configPath = "config.ini"

	c.Start()
	rc := realMain(args)
	out := c.Stop()

	if rc != rc_ERROR {
		t.Errorf("想定外のrc[%d]が返された。", rc)
	}
	if !strings.Contains(out, "Usage :") {
		t.Error("出力内容が想定と違っている。")
		t.Logf("出力: %s", out)
	}
}

func TestRealMain_キー名としてディレクトリが指定された場合(t *testing.T) {
	c := testutil.NewStdoutCapturer()

	args := new(arguments)
	args.bucketName = "bucket"
	args.keyName = "test/"
	args.filePath = "localpath"
	args.configPath = "config.ini"

	c.Start()
	rc := realMain(args)
	out := c.Stop()

	if rc != rc_ERROR {
		t.Errorf("想定外のrc[%d]が返された。", rc)
	}
	if !strings.Contains(out, "DIRECTORY CAN NOT BE SPECIFIED") {
		t.Error("出力内容が想定と違っている。")
		t.Logf("出力: %s", out)
	}
}

func TestRealMain_存在しない設定ファイルが指定された場合(t *testing.T) {
	c := testutil.NewStdoutCapturer()

	args := new(arguments)
	args.bucketName = "testbucket"
	args.keyName = "testuploadlocation/test.txt"
	args.filePath = "testlocalpath"
	args.configPath = "noexistsconf.ini"

	c.Start()
	rc := realMain(args)
	out := c.Stop()

	if rc != rc_ERROR {
		t.Errorf("想定外のrc[%d]が返された。", rc)
	}
	if !strings.Contains(out, "FAILED TO READ CONFIG FILE.") {
		t.Error("出力内容が想定と違っている。")
		t.Logf("出力: %s", out)
	}
}

func TestRealMain_不正な内容の設定ファイルが指定された場合(t *testing.T) {
	c := testutil.NewStdoutCapturer()

	args := new(arguments)
	args.bucketName = "testbucket"
	args.keyName = "testuploadlocation/test.txt"
	args.filePath = "testlocalpath"
	args.configPath = filepath.Join("testdata", "error.ini")

	c.Start()
	rc := realMain(args)
	out := c.Stop()

	if rc != rc_ERROR {
		t.Errorf("想定外のrc[%d]が返された。", rc)
	}
	if !strings.Contains(out, "CONFIG PARM IS NOT EXACT FORMAT.") {
		t.Error("出力内容が想定と違っている。")
		t.Logf("出力: %s", out)
	}
}

func TestRealMain_アップロードに失敗した場合(t *testing.T) {
	makeUploadFail()
	defer restoreUploadFunc()

	c := testutil.NewStdoutCapturer()

	args := new(arguments)
	args.bucketName = "testbucket"
	args.keyName = "testuploadlocation/test.txt"
	args.filePath = "testlocalpath"
	args.configPath = filepath.Join("testdata", "correct.ini")

	c.Start()
	rc := realMain(args)
	out := c.Stop()

	if rc != rc_ERROR {
		t.Errorf("想定外のrc[%d]が返された。", rc)
	}
	if !strings.Contains(out, "UPA004E") {
		t.Error("出力内容が想定と違っている。")
		t.Logf("出力: %s", out)
	}
}

func TestRealMain_正常系(t *testing.T) {
	makeUploadSuccess()
	defer restoreUploadFunc()

	c := testutil.NewStdoutCapturer()

	args := new(arguments)
	args.bucketName = "testbucket"
	args.keyName = "testuploadlocation/test.txt"
	args.filePath = "testlocalpath"
	args.configPath = filepath.Join("testdata", "correct.ini")

	c.Start()
	rc := realMain(args)
	c.Stop()

	if rc == rc_ERROR {
		t.Errorf("想定外のrc[%d]が返された。", rc)
	}
}
