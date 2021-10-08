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

	"github.com/jackc/pgx/v4"
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
	log.Println("ETL pipeline starting 🚀")

	report := BatchReport{}
	startTime := time.Now()
	config, pgConfig := pkg.NewConfiguration()
	configJSON, _ := json.MarshalIndent(config, "", " ")
	buffer := make([]pkg.Matchup, 0, config.BatchRecords)

	log.Printf("[Starting with config]: %v", string(configJSON))

	f, err := os.Open(config.FilePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, pgConfig.ConnString)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close(ctx)

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		matchup, isValid := pkg.TransformDatum(scanner.Text())

		if isValid {
			buffer = append(buffer, matchup)
		}

		items := len(buffer)
		if items == config.BatchRecords {
			log.Printf("Commiting Batch [size: %d]", items)
			rowsInserted, err := pkg.BulkInsert(conn, pgConfig.Table, buffer)
			if err != nil {
				log.Fatalln(err)
			}

			report.RowsInserted += rowsInserted
			buffer = make([]pkg.Matchup, 0, config.BatchRecords)

			log.Println("Batch was successfully committed")
		}
	}

	if len(buffer) != 0 {
		rowsInserted, err := pkg.BulkInsert(conn, pgConfig.Table, buffer)
		if err != nil {
			log.Fatalln(err)
		}
		report.RowsInserted += rowsInserted
	}

	report.Duration = fmt.Sprintf("%v ms", time.Since(startTime).Milliseconds())
	reportJSON, _ := json.MarshalIndent(report, "", " ")
	log.Println(string(reportJSON))

	log.Println("ETL pipeline finished 🤖")
}
