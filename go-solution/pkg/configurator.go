package pkg

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	if os.Getenv("GO_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Printf("[Warning]: %v", err.Error())
		}
	}
}

type Configuration struct {
	_            struct{}
	FilePath     string
	BatchRecords int
	errorLogPath string
}

func NewConfiguration() *Configuration {
	filepath := getEnv("FILE_PATH")
	batchRecords, err := strconv.Atoi(getEnv("BATCH_RECORDS", "5000"))
	if err != nil {
		batchRecords = 5000
	}
	errorLogPath := getEnv("ERROR_LOG_PATH")

	return &Configuration{
		BatchRecords: batchRecords,
		FilePath:     filepath,
		errorLogPath: errorLogPath,
	}
}

func getEnv(value string, defaultValue ...string) string {
	envValue := os.Getenv(value)
	if envValue != "" || len(defaultValue) == 0 {
		return envValue
	}

	return defaultValue[0]
}

func NewPGConfiguration() {}
