package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/mbolis/mogo/config"
	"github.com/mbolis/mogo/i18n"
	"github.com/mbolis/mogo/model"
	"github.com/mbolis/mogo/phase"
	"github.com/mbolis/mogo/sign"
	"github.com/mbolis/mogo/status"
	"github.com/mshafiee/swephgo"
)

var T = i18n.T

func main() {
	cfg := config.Parse()
	defer swephgo.Close()

	i18n.SetLang(cfg.Lang)

	start, end := cfg.Range()

	var days []Day
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		ph := phase.ForDay(d)
		sign := sign.ForDay(d)
		days = append(days, Day{d, ph, sign})
	}

	var out io.WriteCloser
	if cfg.Output == "-" {
		out = os.Stdout
	} else {
		var err error
		out, err = os.Create(cfg.Output)
		if err != nil {
			panic(err)
		}
		defer out.Close()
	}

	switch cfg.Format() {
	case config.CSV:
		GenerateCSV(cfg, days, out)
	case config.XLSX:
		GenerateXLSX(cfg, days, out)
	case config.ODS:
		GenerateODS(cfg, days, out)
	case config.PDF:
		tmp, err := os.MkdirTemp("", "mogo-")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(tmp)

		odsOut, err := os.Create(filepath.Join(tmp, "out.ods"))
		if err != nil {
			panic(err)
		}
		defer odsOut.Close()

		GenerateODS(cfg, days, odsOut)
		soffice, err := exec.LookPath("soffice")
		if err == nil {
			soffice, err = filepath.Abs(soffice)
		}
		if err != nil {
			panic(err)
		}

		cmd := exec.Command(soffice, "--headless",
			"-env:UserInstallation=file:///tmp/LibreOffice_Conversion_mogo",
			"--convert-to", "pdf:calc_pdf_Export", "--outdir", tmp, odsOut.Name())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			panic(err)
		}

		pdf, err := os.Open(filepath.Join(tmp, "out.pdf"))
		if err != nil {
			panic(err)
		}
		defer pdf.Close()

		_, err = io.Copy(out, pdf)
		if err != nil {
			panic(err)
		}
	}
}

type Day struct {
	Time  time.Time
	Phase model.DailyValue[phase.Phase]
	Sign  model.DailyValue[sign.Sign]
}

type Row struct {
	status.Entry
	cfg config.Config
}

func (r Row) PhaseText() (icon, name string) {
	if r.Phase < 0 {
		return "", ""
	}
	return string(r.cfg.Icons.Phase(r.Phase)), T("phase." + r.Phase.String())
}

func (r Row) SignText() (icon, name string) {
	if r.Sign < 0 {
		return "", ""
	}
	return string(r.cfg.Icons.Sign(r.Sign)), T("zodiac." + r.Sign.String())
}

func (r Row) HaircutIcon() string {
	return r.cfg.Icons.Status(r.Haircut())
}
func (r Row) NailsCutIcon() string {
	return r.cfg.Icons.Status(r.NailsCut())
}
func (r Row) EpilationIcon() string {
	return r.cfg.Icons.Status(r.Epilation())
}
func (r Row) FacialCleansingIcon() string {
	return r.cfg.Icons.Status(r.FacialCleansing())
}
func (r Row) FaceMaskIcon() string {
	return r.cfg.Icons.Status(r.FaceMask())
}

func (r Row) Strings() []string {
	month := T("month." + r.Date.Format("Jan"))
	day := fmt.Sprintf("%s %d", T("weekday."+r.Date.Format("Mon")), r.Date.Day())

	time := ""
	if !r.Time.IsZero() {
		time = r.Time.Format("15:04")
	}

	phaseIcon, phaseName := r.PhaseText()
	signIcon, signName := r.SignText()

	return []string{
		month, day, time, phaseIcon, phaseName, signIcon, signName,
		r.HaircutIcon(), r.NailsCutIcon(), r.EpilationIcon(), r.FacialCleansingIcon(), r.FaceMaskIcon(),
	}
}

func (d Day) Rows(cfg config.Config) []Row {
	var entries []status.Entry
	switch {
	case d.Phase.Event != nil && d.Sign.Event != nil:
		switch {
		case d.Phase.Event.Time.Before(d.Sign.Event.Time):
			entries = []status.Entry{
				{Date: d.Time, Time: d.Phase.Event.Time, Phase: d.Phase.Event.Value, Sign: d.Sign.Curr},
				{Date: d.Time, Time: d.Sign.Event.Time, Phase: -1, Sign: d.Sign.Event.Value},
			}
		case d.Sign.Event.Time.Before(d.Phase.Event.Time):
			entries = []status.Entry{
				{Date: d.Time, Time: d.Sign.Event.Time, Phase: -1, Sign: d.Sign.Event.Value},
				{Date: d.Time, Time: d.Phase.Event.Time, Phase: d.Phase.Event.Value, Sign: d.Sign.Event.Value},
			}
		case d.Phase.Event.Time == d.Sign.Event.Time:
			entries = []status.Entry{
				{Date: d.Time, Time: d.Phase.Event.Time, Phase: d.Phase.Event.Value, Sign: d.Sign.Event.Value},
			}
		}

	case d.Phase.Event != nil:
		entries = []status.Entry{
			{Date: d.Time, Time: d.Phase.Event.Time, Phase: d.Phase.Event.Value, Sign: d.Sign.Curr},
		}

	case d.Sign.Event != nil:
		entries = []status.Entry{
			{Date: d.Time, Time: d.Sign.Event.Time, Phase: d.Phase.Curr, Sign: d.Sign.Event.Value},
		}

	default:
		entries = []status.Entry{
			{Date: d.Time, Phase: d.Phase.Curr, Sign: d.Sign.Curr},
		}
	}

	var rows []Row
	for _, e := range entries {
		rows = append(rows, Row{e, cfg})
	}
	return rows
}
