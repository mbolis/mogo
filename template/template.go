package template

import (
	"bytes"
	_ "embed"
)

//go:embed template.ods
var ods []byte

//go:embed template.xlsx
var xlsx []byte

func ODS() (*bytes.Reader, int64) {
	return bytes.NewReader(ods), int64(len(ods))
}

func XLSX() *bytes.Reader {
	return bytes.NewReader(xlsx)
}
