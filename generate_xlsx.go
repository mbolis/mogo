package main

import (
	"io"
	"strconv"

	"github.com/mbolis/mogo/config"
	"github.com/xuri/excelize/v2"
)

func GenerateXLSX(cfg config.Config, days []Day, out io.Writer) {
	tpl, err := excelize.OpenFile("template.xlsx")
	if err != nil {
		panic(err)
	}
	defer tpl.Close()

	tpl.SetCellStr("Sheet1", "A1", T("Month"))
	tpl.SetCellStr("Sheet1", "B1", T("Day"))
	tpl.SetCellStr("Sheet1", "C1", T("Hour"))
	tpl.SetCellStr("Sheet1", "D1", T("Phase"))
	tpl.SetCellStr("Sheet1", "F1", T("Sign"))
	tpl.SetCellStr("Sheet1", "H1", T("Haircut"))
	tpl.SetCellStr("Sheet1", "I1", T("Nails cut"))
	tpl.SetCellStr("Sheet1", "J1", T("Epilation"))
	tpl.SetCellStr("Sheet1", "K1", T("Facial cleansing"))
	tpl.SetCellStr("Sheet1", "L1", T("Face mask"))

	appendRowIndex := 4
	for i, d := range days {
		sourceRow := 3 - i%2 // odd=2 even=3

		for _, r := range d.Rows(cfg) {
			targetCellRow := strconv.Itoa(appendRowIndex)

			err = tpl.DuplicateRowTo("Sheet1", sourceRow, appendRowIndex)
			if err != nil {
				panic(err)
			}

			// TODO handle / panic on all errors
			tpl.SetCellValue("Sheet1", "A"+targetCellRow, r.Date)
			tpl.SetCellValue("Sheet1", "B"+targetCellRow, r.Date)
			if !r.Time.IsZero() {
				tpl.SetCellValue("Sheet1", "C"+targetCellRow, r.Time)
			}

			phaseIcon, phaseName := r.PhaseText()
			tpl.SetCellStr("Sheet1", "D"+targetCellRow, phaseIcon)
			tpl.SetCellStr("Sheet1", "E"+targetCellRow, phaseName)

			signIcon, signName := r.SignText()
			tpl.SetCellStr("Sheet1", "F"+targetCellRow, signIcon)
			tpl.SetCellStr("Sheet1", "G"+targetCellRow, signName)

			tpl.SetCellStr("Sheet1", "H"+targetCellRow, r.HaircutIcon())
			tpl.SetCellStr("Sheet1", "I"+targetCellRow, r.NailsCutIcon())
			tpl.SetCellStr("Sheet1", "J"+targetCellRow, r.EpilationIcon())
			tpl.SetCellStr("Sheet1", "K"+targetCellRow, r.FacialCleansingIcon())
			tpl.SetCellStr("Sheet1", "L"+targetCellRow, r.FaceMaskIcon())

			appendRowIndex++
		}
	}

	err = tpl.RemoveRow("Sheet1", 2)
	if err != nil {
		panic(err)
	}
	err = tpl.RemoveRow("Sheet1", 2)
	if err != nil {
		panic(err)
	}

	err = tpl.SetSheetName("Sheet1", strconv.Itoa(cfg.Year))
	if err != nil {
		panic(err)
	}

	err = tpl.Write(out)
	if err != nil {
		panic(err)
	}
}
