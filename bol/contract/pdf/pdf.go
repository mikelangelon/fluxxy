package pdf

import (
	"bol/contract/contract_template"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"bol/contract/section"
	"text/template"
	"bytes"
	"bol/contract/contract"
)

func GeneratePdf() *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world")
	//err := pdf.OutputFileAndClose("hello.pdf")
	//if err != nil {
	//	fmt.Errorf("Error! %s", err)
	//}
	return pdf
}

type ExtraParams struct {
	ParamsMap map[string]string
}

func GeneratePdfFromContract(contractId int64) *gofpdf.Fpdf {
	con, err := contract.Lookup(contractId)
	if err != nil {
		fmt.Errorf("Error retrieving contract! %s", err)
	}
	var temp *contract_template.Template

	if con.TemplateVersion.Version == 0 {
		temp, err = contract_template.Lookup(con.TemplateVersion.ID)
	} else {
		temp, err = contract_template.LookupByVersion(con.TemplateVersion.ID, con.TemplateVersion.Version)
	}
	if err != nil {
		fmt.Errorf("Error retrieving contract_template! %s", err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	titleStr := con.Title
	pdf.SetTitle(titleStr, false)
	pdf.SetAuthor(con.Actor, false)

	pdf.SetHeaderFunc(func() {
		// Arial bold 15
		pdf.SetFont("Arial", "B", 16)
		// Calculate width of title and position
		wd := pdf.GetStringWidth(titleStr) + 6
		pdf.SetX((210 - wd) / 2)
		// Thickness of frame (1 mm)
		pdf.SetLineWidth(1)
		// Background color
		pdf.SetFillColor(100, 220, 255)
		// Title
		pdf.CellFormat(wd, 9, titleStr, "1", 1, "C", true, 0, "")
		// Line break
		pdf.Ln(10)
	})

	for i, secInfo := range temp.Sections {
		var sec *section.Section
		if secInfo.Version == 0 {
			sec, err = section.Lookup(secInfo.ID)
			if err != nil {
				fmt.Errorf("Error retrieving section! %s", err)
			}
		} else {
			sec, err = section.LookupByVersion(secInfo.ID, secInfo.Version)
			if err != nil {
				fmt.Errorf("Error retrieving section! %s", err)
			}
		}
		printSection(pdf, i+1, sec, con)
	}

	return pdf
}

func printSection(pdf *gofpdf.Fpdf, chapNum int, section *section.Section, con *contract.Contract) {
	pdf.AddPage()
	defineSectionTitle(pdf, chapNum, section.Title)
	pdf.SetFont("Arial", "", 12)
	t := template.New("template example")
	t, _ = t.Parse(section.Text)

	buf := make([]byte, 0)
	w := bytes.NewBuffer(buf)
	t.Execute(w, con)
	body := w.String()

	_, lineHt := pdf.GetFontSize()
	html := pdf.HTMLBasicNew()
	html.Write(lineHt, body)

}

func defineSectionTitle(pdf *gofpdf.Fpdf, chapNum int, titleStr string) {
	pdf.SetFont("Helvetica", "", 14)
	pdf.SetFillColor(0, 0, 0)
	pdf.CellFormat(0, 6, fmt.Sprintf("Section %d : %s", chapNum, titleStr),
		"", 1, "L", false, 0, "")
	// Line break
	pdf.Ln(4)
}