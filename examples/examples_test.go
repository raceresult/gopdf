package examples

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/raceresult/gopdf"
	"github.com/raceresult/gopdf/types"
)

func TestExample1(t *testing.T) {
	// create new PDF Builder
	pb := gopdf.New()

	// use a built-in standard font
	f, err := pb.NewStandardFont(types.StandardFont_Helvetica, types.EncodingWinAnsi)
	if err != nil {
		t.Error(err)
		return
	}

	// add first page
	p := pb.NewPage(gopdf.GetStandardPageSize(gopdf.PageSizeA4, false))

	// add "hello world" element
	p.AddElement(&gopdf.TextElement{
		TextChunk: gopdf.TextChunk{
			Text:     "hello world",
			Font:     f,
			FontSize: 36,
		},
		Left:      gopdf.MM(105),
		Top:       gopdf.MM(100),
		TextAlign: gopdf.HorizontalAlignCenter,
	})

	// output
	bts, err := pb.Build()
	if err != nil {
		t.Error(err)
		return
	}
	err = ioutil.WriteFile("example1.pdf", bts, os.ModePerm)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestExample2(t *testing.T) {
	// create new PDF Builder
	pb := gopdf.New()

	// add first page
	p := pb.NewPage(gopdf.GetStandardPageSize(gopdf.PageSizeA4, false))

	// add image
	bts, err := ioutil.ReadFile("myImage.jpg")
	if err != nil {
		t.Error(err)
		return
	}
	img, err := pb.NewImage(bts)
	if err != nil {
		t.Error(err)
		return
	}
	p.AddElement(&gopdf.ImageElement{
		Width:  gopdf.MM(170),
		Height: gopdf.MM(20),
		Left:   gopdf.MM(20),
		Top:    gopdf.MM(20),
		Img:    img,
	})

	// add rectangle
	p.AddElement(&gopdf.RectElement{
		Left:      gopdf.MM(20),
		Top:       gopdf.MM(150),
		Width:     gopdf.MM(50),
		Height:    gopdf.MM(30),
		LineWidth: gopdf.MM(3),
		LineColor: gopdf.NewColorRGB(0, 255, 0),
		FillColor: gopdf.NewColorRGB(255, 0, 255),
	})

	// add composite font
	bts, err = ioutil.ReadFile("arialuni.ttf")
	if err != nil {
		t.Error(err)
		return
	}
	f, err := pb.NewCompositeFont(bts)
	if err != nil {
		t.Error(err)
		return
	}

	// add text using composite font
	p.AddElement(&gopdf.TextElement{
		TextChunk: gopdf.TextChunk{
			Text:         "hello world - 漢語",
			Font:         f,
			FontSize:     36,
			Color:        gopdf.NewColorRGB(200, 200, 200),
			OutlineColor: gopdf.NewColorRGB(10, 20, 10),
			OutlineWidth: gopdf.MM(0.5),
		},
		Left: gopdf.MM(20),
		Top:  gopdf.MM(100),
	})

	// output
	bts, err = pb.Build()
	if err != nil {
		t.Error(err)
		return
	}
	err = ioutil.WriteFile("example2.pdf", bts, os.ModePerm)
	if err != nil {
		t.Error(err)
		return
	}
}
