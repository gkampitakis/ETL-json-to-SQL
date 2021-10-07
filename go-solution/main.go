package main

import (
	"bufio"
	"etl-pipeline/pkg"
	"log"
	"os"
)

func main() {
	log.Println("ETL pipeline starting 🚀")

	config := pkg.NewConfiguration()
	buffer := make([]pkg.Matchup, 0, config.BatchRecords)

	f, err := os.Open(config.FilePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		matchup, isValid := pkg.TransformDatum(scanner.Text())

		if isValid {
			buffer = append(buffer, matchup)
		}
	}

	log.Printf("Number of rows %d", len(buffer))
}
