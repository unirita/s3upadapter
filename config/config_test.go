package config

import (
	//"runtime"
	"strings"
	"testing"
)

func generateTestConfig() {
	Aws.AccessKeyId = `testkeyid`
	Aws.SecletAccessKey = `seclettestkey`
	Aws.Region = `ap-northeast-1`
	Log.LogDebug = 0
	Log.LogSigning = 0
	Log.LogHTTPBody = 0
	Log.LogRequestRetries = 0
	Log.LogRequestErrors = 0
}

func TestLoad_存在しないファイルをロードしようとした場合はエラー(t *testing.T) {
	if err := Load("noexistfilepath"); err == nil {
		t.Error("エラーが発生していない。")
	}
}

func TestLoadByReader_Readerから設定値を取得できる(t *testing.T) {
	conf := `
[aws]
access_key_id='testkeyid'
secret_access_key='seclettestkey'
region='ap-northeast-1'
[download]
download_dir='c:\TEST'
[log]
log_on=0
signing_on=0
httpbody_on=0
request_retries_on=0
request_errors_on=0
`

	r := strings.NewReader(conf)
	err := loadReader(r)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した[%s]", err)
	}

	if Aws.AccessKeyId != `testkeyid` {
		t.Errorf("access_key_idの値[%s]は想定と違っている。", Aws.AccessKeyId)
	}
	if Aws.SecletAccessKey != `seclettestkey` {
		t.Errorf("seclet_access_keyの値[%s]は想定と違っている。", Aws.SecletAccessKey)
	}
	if Aws.Region != `ap-northeast-1` {
		t.Errorf("regionの値[%s]は想定と違っている。", Aws.Region)
	}

}

func TestLoadByReader_tomlの書式に沿っていない場合はエラーが発生する(t *testing.T) {
	conf := `
[aws]
access_key_id=testkeyid
seclet_access_key=seclettestkey
region='ap-northeast-1'
[download]
download_dir='c:\TEST'
[log]
log_on=0
signing_on=0
httpbody_on=0
request_retries_on=0
request_errors_on=0
`

	r := strings.NewReader(conf)
	err := loadReader(r)
	if err == nil {
		t.Error("エラーが発生しなかった")
	}

}

func TestDetectError_パラメータの値が設定されていない場合はエラー(t *testing.T) {
	conf := `
[aws]
access_key_id='testkeyid'
seclet_access_key='seclettestkey'
region='ap-northeast-1'
[download]
download_dir='c:\TEST'
[log]
log_on=
signing_on=
httpbody_on=
request_retries_on=
request_errors_on=
`

	r := strings.NewReader(conf)
	err := loadReader(r)
	if err == nil {
		t.Error("期待しないエラーが発生していない")
	}

}

func TestDetectError_設定内容にエラーが無い場合はnilを返す(t *testing.T) {
	generateTestConfig()
	if err := DetectError(); err != nil {
		t.Errorf("想定外のエラーが発生した： %s", err)
	}
}

func TestDetectError_設定ファイルのアクセスキーIDが空の場合はエラー(t *testing.T) {
	generateTestConfig()
	Aws.AccessKeyId = ``
	if err := DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestDetectError_設定ファイルのシークレットアクセスキーが空の場合はエラー(t *testing.T) {
	generateTestConfig()
	Aws.SecletAccessKey = ``
	if err := DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestDetectError_設定ファイルのリージョンが空の場合はエラー(t *testing.T) {
	generateTestConfig()
	Aws.Region = ``
	if err := DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestDetectError_ログレベルのlog_onの値が不正(t *testing.T) {
	generateTestConfig()

	Log.LogDebug = -1
	if err := DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}

	Log.LogDebug = 2
	if err := DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestDetectError_ログレベルのsigning_onの値が不正(t *testing.T) {
	generateTestConfig()

	Log.LogSigning = -1

	if err := DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}

	Log.LogSigning = 2
	if err := DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestDetectError_ログレベルのhttpbody_onの値が不正(t *testing.T) {
	generateTestConfig()

	Log.LogHTTPBody = -1

	if err := DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}

	Log.LogHTTPBody = 2
	if err := DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestDetectError_ログレベルのrequest_retries_onの値が不正(t *testing.T) {
	generateTestConfig()

	Log.LogRequestRetries = -1

	if err := DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}

	Log.LogRequestRetries = 2
	if err := DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestDetectError_ログレベルのrequest_errors_onの値が不正(t *testing.T) {
	generateTestConfig()

	Log.LogRequestErrors = -1

	if err := DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}

	Log.LogRequestErrors = 2
	if err := DetectError(); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestDetectError_ログレベルが全て正常値はエラーが発生しない(t *testing.T) {
	//全てのログレベルのキーが0
	generateTestConfig()
	if err := DetectError(); err != nil {
		t.Error("期待していないエラーが発生している。全て0の場合")
	}

	//全てのログレベルのキーが1
	Log.LogDebug = 1
	Log.LogSigning = 1
	Log.LogHTTPBody = 1
	Log.LogRequestRetries = 1
	Log.LogRequestErrors = 1

	if err := DetectError(); err != nil {
		t.Error("予期していないエラーが発生している。全て1の場合")
	}
}
