package main

import (
	"bufio"
	"context"
	"encoding/json"
	"etl-pipeline/pkg"
	"fmt"
	"log"
	"os"
	"sync"
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

type Void = struct{}

func main() {
	log.Println("ETL pipeline starting 🚀")

	maxOperations := make(chan Void, 10)
	m := sync.Mutex{}
	wg := sync.WaitGroup{}
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
		if matchup, erred, ignore := pkg.TransformDatum(scanner.Text()); !ignore {
			buffer = append(buffer, matchup)
		} else if erred {
			report.Errors.Transform++
		}

		items := len(buffer)
		if items == config.BatchRecords {
			wg.Add(1)
			go func(buffer []pkg.Matchup) {
				log.Printf("Commiting Batch [size: %d]", items)
				rowsInserted, err := pkg.BulkInsert(conn, pgConfig.Table, buffer)
				if err != nil {
					log.Println(err)
					pkg.ErrorHandler(buffer, config.ErrorLogPath)
				}
				updateReport(&m, &report, rowsInserted, err)

				<-maxOperations
				wg.Done()
			}(buffer)

			buffer = make([]pkg.Matchup, 0, config.BatchRecords)
			maxOperations <- Void{}
		}
	}

	close(maxOperations)

	if len(buffer) != 0 {
		rowsInserted, err := pkg.BulkInsert(conn, pgConfig.Table, buffer)
		if err != nil {
			log.Println(err)
			pkg.ErrorHandler(buffer, config.ErrorLogPath)
		}

		updateReport(&m, &report, rowsInserted, err)
	}

	wg.Wait()

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

func updateReport(m *sync.Mutex, report *BatchReport, rows int64, err error) {
	if err == nil {
		m.Lock()
		report.RowsInserted += rows
		m.Unlock()
		log.Println("Batch was successfully committed")
		return
	}

	m.Lock()
	report.Errors.BulkInsert++
	m.Unlock()
}
