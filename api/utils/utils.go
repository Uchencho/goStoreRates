package utils

import (
	"os"
)

func localServerAddress() string {
	serverAddress, present := os.LookupEnv("SERVER_ADDRESS")
	if present {
		return serverAddress
	}
	const defaultServerAddress = "127.0.0.1:8000"
	return defaultServerAddress
}

func GetServerddress() string {
	serverAddress, present := os.LookupEnv("PORT")
	if present {
		return ":" + serverAddress
	}
	return localServerAddress()
}
