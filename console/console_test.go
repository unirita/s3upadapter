// Copyright 2015 unirita Inc.
// Created 2015/10/08 kazami

package console

import (
	"testing"

	"github.com/unirita/s3dladapter/testutil"
)

func TestDisplay_メッセージを出力できる(t *testing.T) {
	c := testutil.NewStdoutCapturer()
	c.Start()

	Display("ADP002E", "something error")

	output := c.Stop()

	if output != "ADP002E FAILED TO READ CONFIG FILE. [something error]\n" {
		t.Errorf("stdoutへの出力値[%s]が想定と違います。", output)
	}
}

func TestGetMessage_メッセージを文字列として取得できる(t *testing.T) {
	msg := GetMessage("ADP002E", "something error")
	if msg != "ADP002E FAILED TO READ CONFIG FILE. [something error]" {
		t.Errorf("取得値[%s]が想定と違います。", msg)
	}
}
