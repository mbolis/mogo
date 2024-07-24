package main

import (
	"io"
	"strconv"

	"github.com/mbolis/mogo/config"
	"github.com/mbolis/mogo/ods"
)

func GenerateODS(cfg config.Config, days []Day, out io.Writer) {
	doc, err := ods.LoadTemplate()
	if err != nil {
		panic(err)
	}

	header := doc.Row(0)
	header.SetCellString(0, T("Day"))
	header.SetCellString(1, T("Hour"))
	header.SetCellString(2, T("Phase"))
	header.SetCellString(3, T("Sign"))
	header.SetCellString(4, T("Haircut"))
	header.SetCellString(5, T("Nails cut"))
	header.SetCellString(6, T("Epilation"))
	header.SetCellString(7, T("Facial cleansing"))
	header.SetCellString(8, T("Face mask"))

	sourceRows := [2]*ods.Row{
		doc.Row(2).Remove(),
		doc.Row(1).Remove(),
	}

	prevRow := header

	for i, d := range days {
		sourceRow := sourceRows[i%2]

		for _, r := range d.Rows(cfg) {
			currRow := sourceRow.Duplicate()
			currRow.InsertAfter(prevRow)
			prevRow = currRow

			currRow.SetCellDate(0, r.Date)
			currRow.SetCellDate(1, r.Date)
			if !r.Time.IsZero() {
				currRow.SetCellTime(2, r.Time)
			}

			phaseIcon, phaseName := r.PhaseText()
			currRow.SetCellString(3, phaseIcon)
			currRow.SetCellString(4, phaseName)

			signIcon, signName := r.SignText()
			currRow.SetCellString(5, signIcon)
			currRow.SetCellString(6, signName)

			currRow.SetCellString(7, r.HaircutIcon())
			currRow.SetCellString(8, r.NailsCutIcon())
			currRow.SetCellString(9, r.EpilationIcon())
			currRow.SetCellString(10, r.FacialCleansingIcon())
			currRow.SetCellString(11, r.FaceMaskIcon())
		}
	}

	doc.SetWorksheetName(0, strconv.Itoa(cfg.Year))

	err = doc.Write(out)
	if err != nil {
		panic(err)
	}
}
