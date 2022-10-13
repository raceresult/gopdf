package gopdf

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/raceresult/gopdf/builder"
	"github.com/raceresult/gopdf/types"
)

func TestExample1(t *testing.T) {
	// create new PDF Builder
	pb := builder.NewBuilder()

	// use a built-in standard font
	f, err := pb.NewStandardFont(types.StandardFont_Helvetica, types.EncodingWinAnsi)
	if err != nil {
		t.Error(err)
		return
	}

	// add first page
	p := pb.NewPage(builder.GetStandardPageSize(builder.PageSizeA4, false))

	// add "hello world" element
	p.AddElement(&builder.TextElement{
		Text:      "hello world",
		FontSize:  36,
		X:         builder.MM(105),
		Y:         builder.MM(100),
		TextAlign: builder.TextAlignCenter,
		Font:      f,
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
	pb := builder.NewBuilder()

	// add first page
	p := pb.NewPage(builder.GetStandardPageSize(builder.PageSizeA4, false))

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
	p.AddElement(&builder.ImageElement{
		Width:  builder.MM(170),
		Height: builder.MM(20),
		Left:   builder.MM(20),
		Top:    builder.MM(20),
		Img:    img,
	})

	// add rectangle
	p.AddElement(&builder.RectElement{
		X1:        builder.MM(20),
		Y1:        builder.MM(150),
		Width:     builder.MM(50),
		Height:    builder.MM(30),
		LineWidth: builder.MM(3),
		LineColor: builder.ColorRGB{0, 255, 0},
		FillColor: builder.ColorRGB{255, 0, 255},
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
	p.AddElement(&builder.TextElement{
		Text:         "hello world - 漢語",
		Font:         f,
		FontSize:     36,
		X:            builder.MM(20),
		Y:            builder.MM(100),
		Color:        builder.ColorRGB{200, 200, 200},
		OutlineColor: builder.ColorRGB{10, 20, 10},
		RenderMode:   types.RenderingModeFillAndStroke,
		OutlineWidth: builder.MM(0.5),
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
