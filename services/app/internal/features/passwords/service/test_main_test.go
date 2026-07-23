package passwords_service_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setup()
	exitCode := m.Run()
	teardown()
	os.Exit(exitCode)
}

func setup() {}

func teardown() {}
