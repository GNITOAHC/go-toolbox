package dotenv

import (
	"os"
	"testing"
)

func TestDotenv(t *testing.T) {
	err := Load("./env_test")
	if err != nil {
		t.Error(err)
	}

	// ENV=testenv
	// ENV_QUOTES="env_quotes"
	// ENV_SPACE = env_space
	// ENV_QUOTES_SPACE = "env_quotes_space"
	// A = a
	// B = b
	// ENV_VAR = env_${A}_${B}

	if os.Getenv("ENV") != "testenv" ||
		os.Getenv("ENV_QUOTES") != "env_quotes" ||
		os.Getenv("ENV_SPACE") != "env_space" ||
		os.Getenv("ENV_QUOTES_SPACE") != "env_quotes_space" ||
		os.Getenv("ENV_VAR") != "env_a_b" {
		t.Error("Error loading environment variables")
	}
	return
}

func TestLogError(t *testing.T) {
	err := Load("./env_err_test")
	if err.Error() != "Invalid line: DOTENV_ERROR error" {
		t.Error("Error detecting failed")
	}
	return
}
