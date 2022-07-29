package main

import (
	"encoding/csv"
	"log"
	"os"
)

func writedata(data [][]string, filename string) bool {
	csvFile, err := os.Create("report.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvFile.Close()
	csvwriter := csv.NewWriter(csvFile)
	for _, row := range data {
		_ = csvwriter.Write(row)
	}
	return true
}
