package main

import (
	"encoding/csv"
	"io"

	"github.com/mbolis/mogo/config"
)

func GenerateCSV(cfg config.Config, days []Day, out io.Writer) {
	rows := [][]string{
		{
			T("Month"), T("Day"), T("Hour"), T("Phase"), "", T("Sign"), "",
			T("Haircut"), T("Nails cut"), T("Epilation"), T("Facial cleansing"), T("Face mask"),
		},
	}
	for _, d := range days {
		for _, r := range d.Rows(cfg) {
			rows = append(rows, r.Strings())
		}
	}

	csv := csv.NewWriter(out)
	err := csv.WriteAll(rows)
	if err != nil {
		panic(err)
	}
}
