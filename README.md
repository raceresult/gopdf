
race result AG is always looking for smart, passionate and hard-working Golang and Javascript 
developers with both strong attention to the detail and an entrepreneurial approach to tasks.

We are a Germany-based company building technology for timing sports events such London Marathon, 
Challenge Roth, Tour Down Under and thousands of other races in more than 80 countries.
Check out [www.raceresult.com](https://www.raceresult.com) and 
[karriere.raceresult.com](https://karriere.raceresult.com) for more information.

gopdf - Free PDF creator in pure Golang
================================

Go code (golang) to create PDF documents with several layers of abstraction that allow both, 
easy placement of elements on pages and accessing lower layers to implement any type of PDF 
object / operator.

Does support Composite Fonts for full Unicode support.

Getting Started
-------------------------------------------------------------------------------------------

The highest level of abstraction is provided by the gopdf package. Simply create a new Builder object,
add pages and add elements to the pages:

```go
package yours

import (
    "github.com/raceresult/gopdf"
    "github.com/raceresult/gopdf/types"
)

func TestExample1(t *testing.T) {
    // create new PDF Builder
    pb := gopdf.New()
    
    // use a built-in standard fontm
    f, err := pb.NewStandardFont(types.StandardFont_Helvetica, types.EncodingWinAnsi)
    if err != nil {
        t.Error(err)
        return
    }
    
    // add first page
    p := pb.NewPage(gopdf.GetStandardPageSize(gopdf.PageSizeA4, false))
    
    // add "hello world" element
    p.AddElement(&gopdf.TextElement{
        Text:      "hello world",
        FontSize:  36,
        X:         gopdf.MM(105),
        Y:         gopdf.MM(100),
        TextAlign: gopdf.TextAlignCenter,
        Font:      f,
    })
    
    // output
    bts, err := pb.Build()
    ...
}
```

More advanced: let's add an image, a rectangle and a text using a composite font. 
Using a composite font, any unicode character can be mapped to any glyph in one or more fonts.
Non-composite fonts, on the contrary, only support 256 different characters, so in a world of 
UTF-8/Unicode, composite fonts are the only thing you want to use.
Character encoding and subsetting (embed fonts reduced to only those characters you are using) is 
quite sophisticated but completely handled by the pdf package. When using Composite Fonts you only 
need to handover normal (UTF-8 encoded) Go strings.

```go
package yours

import (
    "github.com/raceresult/gopdf/builder"
    "github.com/raceresult/gopdf/types"
)

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
        Width:  gopdf.MM(70),
        Height: gopdf.MM(70 * float64(img.Image.Height) / float64(img.Image.Width)),
        Left:   gopdf.MM(20),
        Top:    gopdf.MM(20),
        Img:    img,
    })
    
    // add rectangle
    p.AddElement(&gopdf.RectElement{
        X1:        gopdf.MM(20),
        Y1:        gopdf.MM(150),
        Width:     gopdf.MM(50),
        Height:    gopdf.MM(30),
        LineWidth: gopdf.MM(3),
        LineColor: gopdf.ColorRGB{0, 255, 0},
        FillColor: gopdf.ColorRGB{255, 0, 255},
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
        Text:         "hello world - 漢語",
        Font:         f,
        FontSize:     36,
        X:            gopdf.MM(20),
        Y:            gopdf.MM(100),
        Color:        gopdf.ColorRGB{200, 200, 200},
        OutlineColor: gopdf.ColorRGB{10, 20, 10},
        RenderMode:   types.RenderingModeFillAndStroke,
        OutlineWidth: gopdf.MM(0.5),
    })
    
    
    // output
    bts, err := pb.Build()
    ...
}
```

This way, the builder supports the following element types with various attributes:
* images (supports jpg, png, gif, bmp)
* lines 
* rectangles (filled or unfilled)
* texts 
* textboxes (text with max width, opt. max height, word wrap, ..)

Advanced: Add your own functionality
-------------------------------------------------------------------------------------------
The following section describes how to go one level deeper and add custom functionality that 
is not yet provided by the builder package.

First of all, a quick overview of the PDF file format: PDF files are a set of "objects" that are 
referenced by other objects. For example, the root is the "Document Catalog" which references to 
the "Page Tree" which references to the individual page objects. Each page references to resources
like images or fonts, and also to a "Content Stream" that contains the commands to draw text, images,
or graphical elements.

The content stream is a list of "Operators" with "Operands", you could also call them function calls.
Some functions change a state, for example set the drawing position or change the drawing color, others
actually draw the text/image/line/...

Please take a look at the 
[PDF Reference manual](https://stuff.mit.edu/afs/sipb/contrib/doc/specs/software/adobe/pdf/PDFReference.pdf).
You may have to use it look up the functions you need to use for your needs. 

Let's assume you need a custom function to draw a cubic Bézier curve. The operators for this are: 
* "m": move to position
* "c": add bezier curve to current path line
* "S": stroke the path line.

Most operators have been implemented in the package "pdf" (even if there are no elements in the
builder package using them). In this example, the functions needed are:

```go
func (q *Page) Path_m(x, y float64)

func (q *Page) Path_c(x1, y1, x2, y2, x3, y3 float64)

func (q *Page) Path_S()
```

You can easily create your own element type for the Builder, it only needs to fulfill the Element interface:

```go
type Element interface {
    Build(page *pdf.Page)
}
```

For example:
```go
type BezierCurve struct {
    X, Y Length
    X1, Y1, X2, Y2, X3, Y3 Length
}

func (q *BezierCurve) Build(page *pdf.Page) {
    page.Path_m(q.X.Pt(), q.Y.Pt())
    page.Path_c(q.X1.Pt(), q.Y1.Pt(), q.X2.Pt(), q.Y2.Pt(), q.X3.Pt(), q.Y3.Pt())
    page.Path_S()
}
```

An instance of BezierCurve can now be added to a page using the AddElement method.

If you need to use operators that are not implemented in the pdf package, you can use the general 
function AddCommand:

```go
func (q *Page) AddCommand(operator string, args ...types.Object)
```

If, however, the operator is implemented, please use the associated function to avoid unexpected behavior.
In order to minimize the size of the PDF file, the pdf package keeps tracking of the current text and 
graphics state and ignores function call, that would not change the state, for example:

```go
func (q *Page) TextState_Tc(charSpace float64) {
    if q.textState.Tc == charSpace {
        return
    }
    q.textState.Tc = charSpace
    
    q.AddCommand("Tc", types.Int(charSpace))
}
```

If you would call AddCommand("Tc", ...) instead of TextState_Tc, the internal textState would not be updated
and you may see unexpected behavior later in your program.

------

Installation
============

To install GoPDF, use `go get`:

    go get github.com/raceresult/gopdf

------

Staying up to date
==================

To update GoPDF to the latest version, use `go get -u github.com/raceresult/gopdf`.

------

Supported go versions
==================

We currently support the most recent major Go versions from 1.16 onward.

------

Contributing
============

Please feel free to submit issues, fork the repository and send pull requests!

When submitting an issue, we ask that you please include a complete test function that demonstrates the issue.

------

License
=======

This project is licensed under the terms of the MIT license.
