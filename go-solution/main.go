package main

import (
	"bufio"
	"context"
	"encoding/json"
	"etl-pipeline/pkg"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type BatchReport struct {
	Duration     string `json:"duration"`
	RowsInserted int64  `json:"rowsInserted"`
	Errors       struct {
		Transform  int `json:"transform"`
		BulkInsert int `json:"bulkInsert"`
	} `json:"errors"`
}

func main() {
	done := make(chan bool, 10)
	log.Println("ETL pipeline starting 🚀")

	report := BatchReport{}
	startTime := time.Now()
	config, pgConfig := pkg.NewConfiguration()
	buffer := make([]pkg.Matchup, 0, config.BatchRecords)

	printConfig(config)

	f, err := os.Open(config.FilePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	ctx := context.Background()
	conn, err := pgxpool.Connect(ctx, pgConfig.ConnString)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		matchup, isValid := pkg.TransformDatum(scanner.Text())

		if isValid {
			buffer = append(buffer, matchup)
		}

		items := len(buffer)
		if items == config.BatchRecords {
			go func(buffer []pkg.Matchup) {
				log.Printf("Commiting Batch [size: %d]", items)
				_, err := pkg.BulkInsert(conn, pgConfig.Table, buffer)
				if err != nil {
					log.Fatalln(err)
					// FIXME:
				}

				// report.RowsInserted += rowsInserted // FIXME:

				log.Println("Batch was successfully committed")
				<-done
			}(buffer)

			buffer = make([]pkg.Matchup, 0, config.BatchRecords)
			done <- true
		}
	}

	close(done)

	if len(buffer) != 0 {
		_, err := pkg.BulkInsert(conn, pgConfig.Table, buffer)
		if err != nil {
			log.Fatalln(err)
		}
		// report.RowsInserted += rowsInserted
	}

	printReport(report, startTime)
	log.Println("ETL pipeline finished 🤖")
}

func printReport(report BatchReport, startTime time.Time) {
	report.Duration = fmt.Sprintf("%v ms", time.Since(startTime).Milliseconds())
	reportJSON, _ := json.MarshalIndent(report, "", " ")
	log.Println(string(reportJSON))
}

func printConfig(config *pkg.Configuration) {
	configJSON, _ := json.MarshalIndent(config, "", " ")
	log.Printf("[Starting with config]: %v", string(configJSON))
}
