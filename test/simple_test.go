package test

import (
	"testing"

	env "github.com/terawatthour/rigid-env"
)

func TestSimple(t *testing.T) {
	type Environment struct {
		ENV         string
		TARGET_PROD string
		TARGET_DEV  string
		VERSION     int
	}

	var environment Environment

	err := env.Load(&environment, ".env")
	if err != nil {
		panic(err)
	}
	if environment.VERSION != 10 {
		panic("wrong value, should be (int)10")
	}
}
