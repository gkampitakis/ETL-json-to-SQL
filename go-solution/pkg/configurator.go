package pkg

import (
	"fmt"
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
	FilePath     string `json:"filepath"`
	BatchRecords int    `json:"batchRecords"`
	ErrorLogPath string `json:"errorLogPath"`
}

func etlConfiguration() *Configuration {
	filepath := getEnv("FILE_PATH")
	batchRecords, err := strconv.Atoi(getEnv("BATCH_RECORDS", "5000"))
	if err != nil {
		batchRecords = 5000
	}
	errorLogPath := getEnv("ERROR_LOG_PATH")

	return &Configuration{
		BatchRecords: batchRecords,
		FilePath:     filepath,
		ErrorLogPath: errorLogPath,
	}
}

func getEnv(value string, defaultValue ...string) string {
	envValue := os.Getenv(value)
	if envValue != "" || len(defaultValue) == 0 {
		return envValue
	}

	return defaultValue[0]
}

type PGConfiguration struct {
	Table      string
	ConnString string
}

func newPGConfiguration() *PGConfiguration {
	table := getEnv("PG_TABLE", "matchups")
	pgPort := getEnv("PG_PORT", "5432")
	db := getEnv("PG_DATABASE", "ETL_db")
	user := getEnv("PG_USER", "ELT_user")
	password := getEnv("PG_PASS", "ETL_pass")

	return &PGConfiguration{
		Table:      table,
		ConnString: fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", user, password, pgPort, db),
	}
}

func NewConfiguration() (*Configuration, *PGConfiguration) {
	return etlConfiguration(), newPGConfiguration()
}
