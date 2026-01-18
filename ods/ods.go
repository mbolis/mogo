package ods

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/mbolis/mogo/template"
	"github.com/psanford/memfs"
)

type Document struct {
	fs  *memfs.FS
	xml *etree.Document
}

func LoadTemplate() (*Document, error) {
	fs := memfs.New()

	files, err := zip.NewReader(template.ODS())
	if err != nil {
		return nil, err
	}

	err = unzipInto(fs, files)
	if err != nil {
		return nil, err
	}

	doc, err := readDoc(fs)
	if err != nil {
		return nil, err
	}

	return &Document{fs, doc}, nil
}

func unzipInto(fs *memfs.FS, files *zip.Reader) error {
	for _, f := range files.File {
		if f.FileInfo().IsDir() {
			err := fs.MkdirAll(strings.TrimRight(f.Name, "/"), os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		bytes, err := readZipEntry(f)
		if err != nil {
			return err
		}

		err = writeEntryInto(fs, f.Name, bytes, f.Mode())
		if err != nil {
			return err
		}
	}
	return nil
}

func readZipEntry(f *zip.File) ([]byte, error) {
	z, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer z.Close()

	return io.ReadAll(z)
}

func writeEntryInto(fs *memfs.FS, filename string, bytes []byte, mode os.FileMode) error {
	err := fs.MkdirAll(filepath.Dir(filename), os.ModePerm)
	if err != nil {
		return err
	}

	return fs.WriteFile(filename, bytes, mode)
}

func readDoc(fs *memfs.FS) (*etree.Document, error) {
	contentFile, err := fs.Open("content.xml")
	if err != nil {
		return nil, err
	}
	defer contentFile.Close()

	doc := etree.NewDocument()
	_, err = doc.ReadFrom(contentFile)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (doc *Document) Row(r int) *Row {
	xmlRows := doc.xml.FindElements("//table:table-row")
	if len(xmlRows) <= r {
		return nil
	}
	xml := xmlRows[r]
	return &Row{xml}
}

func (doc *Document) SetWorksheetName(i int, name string) {
	sheets := doc.xml.FindElements("//table:table")
	if len(sheets) <= i {
		return
	}

	sheet := sheets[i]
	sheet.CreateAttr("table:name", name)
}

func (doc *Document) Write(out io.Writer) error {
	bytes, err := doc.xml.WriteToBytes()
	if err != nil {
		return err
	}
	err = doc.fs.WriteFile("content.xml", bytes, 0)
	if err != nil {
		return err
	}

	zip := zip.NewWriter(out)
	defer zip.Close()
	return zip.AddFS(doc.fs)
}

// SetHeaderRows moves the first n table rows into a
// table:table-header-rows element so they are treated as
// repeating header rows by Calc/LibreOffice when printing
// or exporting to PDF.
func (doc *Document) SetHeaderRows(n int) {
	sheets := doc.xml.FindElements("//table:table")
	if len(sheets) == 0 || n <= 0 {
		return
	}

	sheet := sheets[0]
	rows := sheet.SelectElements("table:table-row")
	if len(rows) < n {
		return
	}

	hdr := etree.NewElement("table:table-header-rows")
	for i := 0; i < n; i++ {
		r := rows[i]
		r.Parent().RemoveChild(r)
		hdr.AddChild(r)
	}

	sheet.InsertChildAt(0, hdr)
}

type Row struct {
	xml *etree.Element
}

func (row *Row) Remove() *Row {
	row.xml.Parent().RemoveChild(row.xml)
	return row
}

func (row *Row) Duplicate() *Row {
	return &Row{row.xml.Copy()}
}

func (row *Row) InsertAfter(prev *Row) {
	prev.xml.Parent().InsertChildAt(prev.xml.Index()+1, row.xml)
}

func (row *Row) SetCellString(c int, value string) {
	cell := row.getCell(c)

	cell.CreateAttr("office:value-type", "string")
	cell.CreateAttr("calcext:value-type", "string")
	p := cell.SelectElement("text:p")
	if p == nil {
		p = cell.CreateElement("text:p")
	}
	p.SetText(value)
}

func (row *Row) SetCellDate(c int, value time.Time) {
	cell := row.getCell(c)

	cell.CreateAttr("office:value-type", "date")
	cell.CreateAttr("calcext:value-type", "date")

	date := value.Format("2006-01-02")
	cell.CreateAttr("office:date-value", date)
}

func (row *Row) SetCellTime(c int, value time.Time) {
	cell := row.getCell(c)

	cell.CreateAttr("office:value-type", "time")
	cell.CreateAttr("calcext:value-type", "time")

	time := value.Format("PT15H04M05S")
	cell.CreateAttr("office:time-value", time)
}

func (row *Row) getCell(c int) *etree.Element {
	cells := row.xml.SelectElements("table:table-cell")
	for i := 0; i < len(cells); i++ {
		cell := cells[i]

		repeatAttr := cell.SelectAttr("table:number-columns-repeated")
		if repeatAttr != nil {
			repeat, err := strconv.Atoi(repeatAttr.Value)
			if err != nil {
				panic(err)
			}

			if repeat > 1 {
				cell.RemoveAttr("table:number-columns-repeated")

				var additionalCells []*etree.Element
				for i := 1; i < repeat; i++ {
					c := cell.Copy()
					additionalCells = append(additionalCells, c)
					cell.Parent().InsertChildAt(cell.Index()+1, c)
				}

				leadingCells := cells[:i+1]
				trailingCells := cells[i+1:]
				cells = append(leadingCells, additionalCells...)
				cells = append(cells, trailingCells...)
			}
		}
		if i == c {
			return cell
		}
	}
	return nil
}
