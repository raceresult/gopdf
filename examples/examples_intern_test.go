package examples

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/raceresult/gopdf"
	"github.com/raceresult/gopdf/parser"
)

func TestRead(t *testing.T) {
	bts, err := ioutil.ReadFile("C:\\users\\soenke\\downloads\\startnummer.pdf")
	//bts, err := ioutil.ReadFile("C:\\users\\soenke\\downloads\\test.pdf")
	if err != nil {
		t.Error(err)
		return
	}

	f, err := parser.New(bts)
	if err != nil {
		t.Error(err)
		return
	}

	allPages, err := f.GetAllPages()
	if err != nil {
		t.Error(err)
		return
	}

	newPDF := gopdf.New()

	p := allPages[11]
	newPage := newPDF.NewPage(gopdf.PageSize{gopdf.Pt(float64(p.MediaBox.URX)), gopdf.Pt(float64(p.MediaBox.URY))})
	cp, err := newPDF.NewCapturedPage(p, f.File())
	cp.Left = gopdf.MM(10)
	if err != nil {
		t.Error(err)
		return
	}
	newPage.AddElement(cp)

	p = allPages[10]
	newPage = newPDF.NewPage(gopdf.PageSize{gopdf.Pt(float64(p.MediaBox.URX)), gopdf.Pt(float64(p.MediaBox.URY))})
	cp, err = newPDF.NewCapturedPage(p, f.File())
	cp.Top = gopdf.MM(10)
	if err != nil {
		t.Error(err)
		return
	}
	newPage.AddElement(cp)

	// output
	bts, err = newPDF.Build()
	if err != nil {
		t.Error(err)
		return
	}
	err = ioutil.WriteFile("exampleRead.pdf", bts, os.ModePerm)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestFedex(t *testing.T) {
	files := []string{
		"C:\\users\\soenke\\downloads\\FedEx test\\794677158844_794677158844_original.pdf",
		"C:\\users\\soenke\\downloads\\FedEx test\\794677158844_794677158855_original.pdf",
		"C:\\users\\soenke\\downloads\\FedEx test\\794677158844_794677158866_original.pdf",
	}

	newPDF := gopdf.New()

	for _, f := range files {
		bts, err := ioutil.ReadFile(f)
		if err != nil {
			t.Error(err)
			return
		}

		f, err := parser.New(bts)
		if err != nil {
			t.Error(err)
			return
		}

		p, err := f.GetPage(1)
		if err != nil {
			t.Error(err)
			return
		}

		newPage := newPDF.NewPage(gopdf.GetStandardPageSize(gopdf.PageSizeA4, true))
		newPage.Rotate = int(p.Rotate)
		cp, err := newPDF.NewCapturedPage(p, f.File())
		if err != nil {
			t.Error(err)
			return
		}
		//cp.Top = gopdf.Pt(30)
		cp.Left = gopdf.Pt(25)
		newPage.AddElement(cp)
	}

	// output
	bts, err := newPDF.Build()
	if err != nil {
		t.Error(err)
		return
	}
	err = ioutil.WriteFile("C:\\users\\soenke\\downloads\\FedEx test\\output.pdf", bts, os.ModePerm)
	if err != nil {
		t.Error(err)
		return
	}
}
