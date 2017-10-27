package main

import (
	"fmt"
	"testing"
)

// ファイルがないとき
func Test_newGrcron_NothingStateFile(t *testing.T) {
	if _, err := newGrcron("active", "test_fixture/nothing"); err == nil {
		t.Error("nothing statefile")
	}
}

// ファイルの中身がおかしいときはdefault stateになる
func Test_newGrcron_InvalidState01(t *testing.T) {
	gr, err := newGrcron("active", "test_fixture/state_hoge")
	if err != nil {
		t.Error("expect no error")
	}
	if gr.CurrentState != "active" {
		t.Errorf("expect status active, gr.CurrentState:%s", gr.CurrentState)
	}
}

// ファイルの中身がおかしいときはdefault stateになる
func Test_newGrcron_InvalidState02(t *testing.T) {
	gr, err := newGrcron("passive", "test_fixture/state_hoge")
	if err != nil {
		t.Error("expect no error")
	}
	if gr.CurrentState != "passive" {
		t.Errorf("expect status passive, gr.CurrentState:%s", gr.CurrentState)
	}
}

// ファイルの中身が空のときもdefault statusになる
func Test_newGrcron_EmptyFile(t *testing.T) {
	gr, err := newGrcron("passive", "test_fixture/state_empty")
	if err != nil {
		t.Error("expect no error")
	}
	if gr.CurrentState != "passive" {
		t.Errorf("expect status passive, gr.CurrentState:%s", gr.CurrentState)
	}
}

func Test_canRun_DownKeepalivedAndActive(t *testing.T) {
	testKeepalivedActive = func() (bool, error) {
		return false, fmt.Errorf("keepalived is probably down")
	}
	gr, err := newGrcron("passive", "test_fixture/state_active")
	if err != nil {
		t.Error("expect no error")
	}
	if gr.CurrentState != "active" {
		t.Error("expect state active")
	}
	canrun, err := gr.canRun()
	if !(canrun == false && err != nil) {
		t.Error("oops ???")
	}
}

func Test_canRun_UpKeepalivedAndActive(t *testing.T) {
	testKeepalivedActive = func() (bool, error) {
		return true, nil
	}
	gr, err := newGrcron("passive", "test_fixture/state_active")
	if err != nil {
		t.Error("expect no error")
	}
	if gr.CurrentState != "active" {
		t.Error("expect state active")
	}
	canrun, err := gr.canRun()
	if !(canrun == true && err == nil) {
		t.Error("oops ???")
	}
}

func Test_canRun_UpKeepalivedAndPassive(t *testing.T) {
	testKeepalivedActive = func() (bool, error) {
		return true, nil
	}
	gr, err := newGrcron("passive", "test_fixture/state_passive")
	if err != nil {
		t.Error("expect no error")
	}
	if gr.CurrentState != "passive" {
		t.Error("expect state passive")
	}
	canrun, err := gr.canRun()
	if !(canrun == false && err == nil) {
		t.Error("oops ???")
	}
}
