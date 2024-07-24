package main

import (
	"fmt"
	"io"
	"strconv"

	"github.com/mbolis/mogo/config"
	"github.com/mbolis/mogo/template"
	"github.com/xuri/excelize/v2"
)

func GenerateXLSX(cfg config.Config, days []Day, out io.Writer) {
	tpl, err := excelize.OpenReader(template.XLSX())
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
			mustDuplicateRowTo(tpl, "Sheet1", sourceRow, appendRowIndex)

			mustSetCellValue(tpl, "Sheet1", appendRowIndex, 0, r.Date)
			mustSetCellValue(tpl, "Sheet1", appendRowIndex, 1, r.Date)
			if !r.Time.IsZero() {
				mustSetCellValue(tpl, "Sheet1", appendRowIndex, 2, r.Time)
			}

			phaseIcon, phaseName := r.PhaseText()
			mustSetCellStr(tpl, "Sheet1", appendRowIndex, 3, phaseIcon)
			mustSetCellStr(tpl, "Sheet1", appendRowIndex, 4, phaseName)

			signIcon, signName := r.SignText()
			mustSetCellStr(tpl, "Sheet1", appendRowIndex, 5, signIcon)
			mustSetCellStr(tpl, "Sheet1", appendRowIndex, 6, signName)

			mustSetCellStr(tpl, "Sheet1", appendRowIndex, 7, r.HaircutIcon())
			mustSetCellStr(tpl, "Sheet1", appendRowIndex, 8, r.NailsCutIcon())
			mustSetCellStr(tpl, "Sheet1", appendRowIndex, 9, r.EpilationIcon())
			mustSetCellStr(tpl, "Sheet1", appendRowIndex, 10, r.FacialCleansingIcon())
			mustSetCellStr(tpl, "Sheet1", appendRowIndex, 11, r.FaceMaskIcon())

			appendRowIndex++
		}
	}

	mustRemoveRow(tpl, "Sheet1", 2)
	mustRemoveRow(tpl, "Sheet1", 2)

	mustSetSheetName(tpl, "Sheet1", strconv.Itoa(cfg.Year))

	err = tpl.Write(out)
	if err != nil {
		panic(err)
	}
}

func mustDuplicateRowTo(tpl *excelize.File, sheet string, src, dst int) {
	err := tpl.DuplicateRowTo(sheet, src, dst)
	if err != nil {
		panic(err)
	}
}

func mustSetCellValue(tpl *excelize.File, sheet string, row, col int, value any) {
	err := tpl.SetCellValue(sheet, fmt.Sprintf("%c%d", 'A'+col, row), value)
	if err != nil {
		panic(err)
	}
}

func mustSetCellStr(tpl *excelize.File, sheet string, row, col int, value string) {
	err := tpl.SetCellStr(sheet, fmt.Sprintf("%c%d", 'A'+col, row), value)
	if err != nil {
		panic(err)
	}
}

func mustRemoveRow(tpl *excelize.File, sheet string, row int) {
	err := tpl.RemoveRow(sheet, row)
	if err != nil {
		panic(err)
	}
}

func mustSetSheetName(tpl *excelize.File, sheet, name string) {
	err := tpl.SetSheetName(sheet, name)
	if err != nil {
		panic(err)
	}
}
