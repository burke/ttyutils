package ttyutils_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/burke/ttyutils"
	"github.com/kr/pty"
)

func TestIsTerminal(t *testing.T) {
	mpty, mtty, err := pty.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer mtty.Close()
	defer mpty.Close()

	if !ttyutils.IsTerminal(mtty.Fd()) {
		t.Error("tty should be reported as a terminal")
	}
	if !ttyutils.IsTerminal(mpty.Fd()) {
		t.Error("pty should be reported as a terminal")
	}

	null, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer null.Close()

	if ttyutils.IsTerminal(null.Fd()) {
		t.Error("/dev/null should not be reported as a terminal")
	}
}

func TestWinSize(t *testing.T) {
	mpty, mtty, err := pty.Open()
	if err != nil {
		t.Fatal(err)
	}
	defer mtty.Close()
	defer mpty.Close()

	size, err := ttyutils.Winsize(mtty)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(size, ttyutils.Ttysize{Lines: 0, Columns: 0}) {
		t.Errorf("expected 0x0 terminal, got %#v", size)
	}
}
