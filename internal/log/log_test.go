package log

import (
	"bytes"

	"os"
	"os/exec"
	"testing"
)

// captureOutput captures the output of a logging API and returns the string
func captureOutput(f func()) string {
	var buf bytes.Buffer

	// Set the output to save to buffer
	SetOutput(&buf)

	// Run the test
	f()

	// Reset the output to stderr
	SetOutput(os.Stderr)

	return buf.String()
}

// TestLevels tests setting and retrieving the levels, and the Stringer interface
func TestLevels(t *testing.T) {
	levels := []struct {
		l Level
		s string
	}{
		{DEBUG, "DEBUG"},
		{INFO, "INFO"},
		{WARNING, "WARNING"},
		{ERROR, "ERROR"},
		{FATAL, "FATAL"},
	}

	for _, level := range levels {

		if level.s != level.l.String() {
			t.Errorf("Stringify failed, expected %v, got %v\n", level.s, level.l.String())
		}

		SetLevel(level.l)

		retlevel := GetLevel()
		if level.l != retlevel {
			t.Errorf("Set level %v, got %v instead\n", level.l, retlevel)
		}
	}
}

// TestOutputs verifies that the outputs match the expected string
func TestOutputs(t *testing.T) {
	tests := []struct {
		f    func(...interface{})
		exp  string
		args []interface{}
	}{
		{Info, "INFO info message\n", []interface{}{"info message"}},
		{Error, "ERROR error message\n", []interface{}{"error message"}},
		{Warning, "WARNING Warning message\n", []interface{}{"Warning message"}},
		{Debug, "DEBUG log_test.go:74 - Debug message\n", []interface{}{"Debug message"}},
		{Info, "INFO info message 100 3.14 (2+3i)\n", []interface{}{"info message", 100, 3.14, 2 + 3i}},
	}

	// Make sure everything is captured
	SetLevel(DEBUG)

	for _, test := range tests {
		output := captureOutput(func() {
			test.f(test.args...)
		})

		if output != test.exp {
			t.Errorf("Expected '%v', got '%v'\n", test.exp, output)
		}
	}
}

// TestFilter ensures that only messages at or above the set log level are captured
func TestFilter(t *testing.T) {
	tests := []struct {
		level Level
		f     func(...interface{})
		exp   string
		args  []interface{}
	}{
		{DEBUG, Debug, "DEBUG log_test.go:120 - test\n", []interface{}{"test"}},
		{DEBUG, Info, "INFO test\n", []interface{}{"test"}},
		{DEBUG, Warning, "WARNING test\n", []interface{}{"test"}},
		{DEBUG, Error, "ERROR test\n", []interface{}{"test"}},

		{INFO, Debug, "", []interface{}{"test"}},
		{INFO, Info, "INFO test\n", []interface{}{"test"}},
		{INFO, Warning, "WARNING test\n", []interface{}{"test"}},
		{INFO, Error, "ERROR test\n", []interface{}{"test"}},

		{WARNING, Debug, "", []interface{}{"test"}},
		{WARNING, Info, "", []interface{}{"test"}},
		{WARNING, Warning, "WARNING test\n", []interface{}{"test"}},
		{WARNING, Error, "ERROR test\n", []interface{}{"test"}},

		{ERROR, Debug, "", []interface{}{"test"}},
		{ERROR, Info, "", []interface{}{"test"}},
		{ERROR, Warning, "", []interface{}{"test"}},
		{ERROR, Error, "ERROR test\n", []interface{}{"test"}},

		{FATAL, Debug, "", []interface{}{"test"}},
		{FATAL, Info, "", []interface{}{"test"}},
		{FATAL, Warning, "", []interface{}{"test"}},
		{FATAL, Error, "", []interface{}{"test"}},
	}

	for _, test := range tests {
		SetLevel(test.level)
		output := captureOutput(func() {
			test.f(test.args...)
		})

		if output != test.exp {
			t.Errorf("Expected '%v', got '%v'\n", test.exp, output)
		}
	}
}

// TestFatal verifies that it kills the program
func TestFatal(t *testing.T) {
	if os.Getenv("OVERLORD_TEST_LOG_TESTFATAL") == "xyzzy" {
		Fatal("die die die")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFatal")
	cmd.Env = append(os.Environ(), "OVERLORD_TEST_LOG_TESTFATAL=xyzzy")
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Errorf("Fatal did not exit with status 1")
}
