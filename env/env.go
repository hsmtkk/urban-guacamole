package env

import (
	"fmt"
	"os"
	"strconv"
)

func LoadEnvVar(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", fmt.Errorf("you must define %s env var", name)
	}
	return val, nil
}

func GetPort() (int, error) {
	portStr, err := LoadEnvVar("PORT")
	if err != nil {
		return 0, err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s as int; %w", portStr, err)
	}
	return port, nil
}
