package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jeandeaual/go-locale"
	"github.com/mbolis/mogo/icons"
	"golang.org/x/text/language"
)

type Config struct {
	Year   int
	Month  time.Month
	TZ     *time.Location
	Icons  icons.Style
	Output string
	Lang   language.Tag
}

type Format int

const (
	CSV Format = iota
	XLSX
	ODS
	PDF
)

func (c Config) Format() Format {
	ext := path.Ext(c.Output)
	switch ext {
	case ".csv", ".txt", "":
		return CSV
	case ".xlsx":
		return XLSX
	case ".ods":
		return ODS
	case ".pdf":
		return PDF
	default:
		panic("unrecognized file extension: " + ext)
	}
}

func (c Config) Range() (start, end time.Time) {
	start = time.Date(c.Year, c.Month, 1, 0, 0, 0, 0, c.TZ)
	if c.Month == 0 {
		start = start.AddDate(0, 1, 0)
		end = start.AddDate(1, 0, 0)
	} else {
		end = start.AddDate(0, 1, 0)
	}
	return
}

var monthRegex = regexp.MustCompile(
	`(?i)^(jan(uary)?|feb(ruary)?|mar(ch)?|apr(il)?|may|jun(e)?|` +
		`jul(y)?|aug(ust)?|sep(tember)?|oct(ober)?|nov(ember)?|dec(ember)?)$`,
)

var monthsByTrigram = map[string]time.Month{
	"jan": time.January,
	"feb": time.February,
	"mar": time.March,
	"apr": time.April,
	"may": time.May,
	"jun": time.June,
	"jul": time.July,
	"aug": time.August,
	"sep": time.September,
	"oct": time.October,
	"nov": time.November,
	"dec": time.December,
}

func (c *Config) SetMonth(s string) error {
	name := monthRegex.FindString(s)
	if name != "" {
		name = strings.ToLower(name[:3])
		c.Month = monthsByTrigram[name]
		return nil
	}

	m, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("unrecognized month '%s'", s)
	}

	c.Month = time.Month(m)
	return nil
}

func (c *Config) SetTZ(s string) (err error) {
	if c.TZ != nil {
		return errors.New("cannot mix -tz and -utc")
	}

	c.TZ, err = time.LoadLocation(s)
	return
}

func (c *Config) SetUTC(string) error {
	if c.TZ != nil {
		return errors.New("cannot mix --tz and --utc")
	}

	c.TZ = time.UTC
	return nil
}

var iconsPackRegex = regexp.MustCompile(`(?i)^(arrows|thumbs|semaphore)$`)

var iconsPacksByName = map[string]icons.Style{
	"arrows":    icons.Arrows,
	"thumbs":    icons.Thumbs,
	"semaphore": icons.Semaphore,
}

func (c *Config) SetIconPack(s string) error {
	name := iconsPackRegex.FindString(s)
	if name == "" {
		return fmt.Errorf("unrecognized icons style '%s'", s)
	}

	name = strings.ToLower(name)
	c.Icons = iconsPacksByName[name]
	return nil
}

func (c *Config) SetLang(s string) (err error) {
	c.Lang, err = language.Parse(s)
	return
}

const usage = `supported options:
    -y YEAR
    --year YEAR
        ephemeris will be calculated for the duration of YEAR (default: current year)
    -m MONTH
    --month MONTH
        if specified, the calculation will be restricted to MONTH
        can be either a number [1-12], or short or long name (jan/january, ...)
    -z TIMEZONE
    --tz TIMEZONE
        output time is local to TIMEZONE (default: system timezone)
        cannot be specified along with --utc
    -u
    --utc
        shortcut for '--tz UTC'
        cannot be specified along with --tz
    -i ICONS
    --icons ICONS
        the icons style to be used for indicators in the output
        can be one of: arrows, thumbs, semaphore (default: arrows)
    -o FILENAME
    --output FILENAME
        optional path to an output file, '-' for stdout (default: -)
        if the file name has an extension, it will be used to infer the format, otherwise CSV is assumed
    -l LANGUAGE
    --lang LANGUAGE
        translate the output into LANGUAGE if supported (default: system language)
        LANGUAGE must be a valid BCP 47 language string
    -h
    --help
        display this help message`

func Parse() (config Config) {
	currentYear := time.Now().Year()
	flag.IntVar(&config.Year, "y", currentYear, "")
	flag.IntVar(&config.Year, "year", currentYear, "")

	flag.Func("m", "", config.SetMonth)
	flag.Func("month", "", config.SetMonth)

	config.TZ = time.Local
	flag.Func("z", "", config.SetTZ)
	flag.Func("tz", "", config.SetTZ)
	flag.Func("timezone", "", config.SetTZ)

	flag.BoolFunc("u", "", config.SetUTC)
	flag.BoolFunc("utc", "", config.SetUTC)

	flag.Func("i", "", config.SetIconPack)
	flag.Func("icon", "", config.SetIconPack)
	flag.Func("icons", "", config.SetIconPack)

	flag.StringVar(&config.Output, "o", "-", "")
	flag.StringVar(&config.Output, "out", "-", "")
	flag.StringVar(&config.Output, "output", "-", "")

	lang, err := locale.GetLanguage()
	if err != nil {
		panic(err)
	}
	config.Lang, err = language.Parse(lang)
	if err != nil {
		panic(err)
	}
	flag.Func("l", "", config.SetLang)
	flag.Func("lang", "", config.SetLang)

	flag.Usage = func() { fmt.Fprintln(os.Stderr, usage) }
	flag.Parse()

	return config
}
