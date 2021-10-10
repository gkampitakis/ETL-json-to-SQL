package pkg

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func ErrorHandler(matchups []Matchup, path string) {
	log.Println(path)
	csvFile, err := createFile(path)
	if err != nil {
		log.Println(err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	for _, item := range matchups {
		row := []string{
			item.Champion,
			fmt.Sprint(item.DamageDealtToChampions),
			item.GameVersion,
			fmt.Sprint(item.GoldEarned),
			fmt.Sprint(item.Win),
			fmt.Sprint(item.MinionsKilled),
			fmt.Sprint(item.KDA),
			item.Lane,
			item.Region,
			item.SummonerName,
			fmt.Sprint(item.VisionScore),
		}

		err := csvWriter.Write(row)
		if err != nil {
			log.Printf("[Save Logger]: %s", err.Error())
		}
	}
}

func createFile(f string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(f), 0770); err != nil {
		return nil, err
	}
	return os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770)
}
