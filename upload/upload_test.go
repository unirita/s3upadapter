package upload

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/unirita/s3upadapter/testutil"
)

func makeUploadSuccess() {
	upload = func(input *s3manager.UploadInput) (*s3manager.UploadOutput, error) {
		output := new(s3manager.UploadOutput)
		output.Location = "testlocation"

		return output, nil
	}
}

func makeUploadFail() {
	upload = func(input *s3manager.UploadInput) (*s3manager.UploadOutput, error) {
		return nil, errors.New("error")
	}
}

func restoreUploadFunc() {
	upload = s3manager.NewUploader(nil).Upload
}

func TestDo_ファイルのオープンに失敗した場合(t *testing.T) {
	err := Do("testbucket", "testkey", "noexists")
	if err == nil {
		t.Fatal("エラーが発生していない。")
	}
}

func TestDo_アップロードに失敗した場合(t *testing.T) {
	makeUploadFail()
	defer restoreUploadFunc()

	err := Do("testbucket", "testkey", filepath.Join("testdata", "correct.ini"))
	if err == nil {
		t.Fatal("エラーが発生していない。")
	}
}

func TestDo_正常系(t *testing.T) {
	makeUploadSuccess()
	defer restoreUploadFunc()

	c := testutil.NewStdoutCapturer()
	c.Start()
	err := Do("testbucket", "testkey", filepath.Join("testdata", "exists.txt"))
	out := c.Stop()

	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}
	if !strings.Contains(out, "testlocation") {
		t.Error("アップロード先のロケーションが出力されていない。")
	}
}
