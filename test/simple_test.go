package test

import (
	"fmt"
	"os"
	"testing"

	env "github.com/terawatthour/rigid-env"
)

func TestSimple(t *testing.T) {
	type Environment struct {
		ENV         string
		TARGET_PROD *string
		TARGET_DEV  string
		VERSION     int
	}

	var environment Environment

	files := []string{".env"}
	err := env.Load(&environment, files...)
	if err != nil {
		panic(err)
	}
	fmt.Println(environment)
	if environment.VERSION != 10 {
		panic("wrong value, should be (int)10")
	}

	if os.Getenv("VERSION") != "10" {
		panic("wrong value, should be (string)10")
	}
}
