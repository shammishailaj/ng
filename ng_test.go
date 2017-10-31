package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

var exeSuffix string // set to ".exe" on GOOS=windows
var testng string    // name of testng binary

func init() {
	if runtime.GOOS == "windows" {
		exeSuffix = ".exe"
	}
	testng = "./testng" + exeSuffix
}

// The TestMain function creates an ng  command for testing purposes.
func TestMain(m *testing.M) {
	out, err := exec.Command("go", "build", "-o", testng).CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "building testng failed: %v\n%s", err, out)
		os.Exit(2)
	}

	r := m.Run()

	os.Remove(testng)

	os.Exit(r)
}

func TestPrintf(t *testing.T) {
	// The printf builtin should not return any values, unlike fmt.Printf.
	out, err := exec.Command(testng, "-e", `printf("%x", 42)`).CombinedOutput()
	if err != nil {
		t.Errorf("testng failed: %v\n%s", err, out)
	}

	got := string(out)
	want := "2a"
	if got != want {
		t.Errorf("printf returned %q, want %q", got, want)
	}
}

func TestExitMsg(t *testing.T) {
	out, err := exec.Command(testng, "-e", `exit`).CombinedOutput()
	// TODO: this should return an non-zero exit code
	if err != nil {
		t.Errorf("testng failed: %v\n%s", err, out)
	}

	got := string(out)
	if !strings.Contains(got, "Ctrl-D") {
		t.Errorf("exit error does not mention Ctrl-D: %q", got)
	}
}
