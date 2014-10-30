package aspect

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

// callerInfo returns a string containing the file and line number of the
// assert call that failed.
// https://github.com/stretchr/testify/blob/master/assert/assertions.go
// Copyright (c) 2012 - 2013 Mat Ryer and Tyler Bunnell
func callerInfo() string {
	file := ""
	line := 0
	ok := false

	for i := 0; ; i++ {
		_, file, line, ok = runtime.Caller(i)
		if !ok {
			return ""
		}
		parts := strings.Split(file, "/")
		file = parts[len(parts)-1]

		// dir := parts[len(parts)-2]
		if file == "testing.go" {
			continue
		}
		break
	}
	return fmt.Sprintf("%s:%d", file, line)
}

type sqlTest struct {
	t       *testing.T
	dialect Dialect
}

// Error tests that the given Compiles instances generates an error for the
// current dialect.
// TODO Allow a specific error
func (t *sqlTest) Error(stmt Compiles) {
	if _, err := stmt.Compile(t.dialect, Params()); err == nil {
		t.t.Errorf("%s: expected error, received nil", callerInfo())
	}
}

// SQL tests that the given Compiles instance matches the expected string for
// the current dialect.
func (t *sqlTest) SQL(expect string, stmt Compiles, ps ...interface{}) {
	// Get caller information in case of failure
	caller := callerInfo()

	// Start a new parameters instance
	params := Params()

	// Compile the given stmt with the tester's dialect
	actual, err := stmt.Compile(t.dialect, params)
	if err != nil {
		t.t.Errorf("%s: unexpected error from compile: %s", caller, err)
		return
	}

	if expect != actual {
		t.t.Errorf(
			"%s: unexpected SQL: expect %s, got %s",
			caller,
			expect,
			actual,
		)
	}
	// Test that the parameters are equal
	if params.Len() != len(ps) {
		t.t.Errorf(
			"%s: unexpected number of parameters for %s: expect %d, got %d",
			caller,
			actual,
			params.Len(),
			len(ps),
		)
		return
	}

	// TODO Examine individual parameters for equality
}

func NewTester(t *testing.T, d Dialect) *sqlTest {
	return &sqlTest{t: t, dialect: d}
}
