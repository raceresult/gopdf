package pdf

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/raceresult/gopdf/types"
	"golang.org/x/image/bmp"
)

// Image holds both the Image object and the reference to it
type Image struct {
	Reference types.Reference
	Image     *types.Image
}

// NewImage adds a new image as XObject to the file and returns the image object and its reference
func (q *File) NewImage(bts []byte) (*Image, error) {
	// read config
	im, name, err := image.DecodeConfig(bytes.NewReader(bts))
	if err != nil {
		return nil, err
	}

	// continue depending on type
	switch name {
	case "bmp":
		return q.newImageBmp(bts, im)

	case "jpg", "jpeg":
		return q.newImageJPG(bts, im)

	case "png":
		return q.newImagePNG(bts, im)

	case "gif":
		return q.newImageGIF(bts, im)

	default:
		return nil, errors.New("unsupported image type " + name)
	}
}

// newImageBmp adds a new bmp file as XObject to the file
func (q *File) newImageBmp(bts []byte, conf image.Config) (*Image, error) {
	// decode image
	x, err := bmp.Decode(bytes.NewReader(bts))
	if err != nil {
		return nil, err
	}

	// build data
	var data []byte
	for i := 0; i < conf.Height; i++ {
		for j := 0; j < conf.Width; j++ {
			r, g, b, _ := x.At(j, i).RGBA()
			data = append(data, byte(r), byte(g), byte(b))
		}
	}

	// create image stream
	imgStream, err := types.NewStream(data, types.Filter_FlateDecode)
	if err != nil {
		return nil, err
	}
	img := types.Image{
		Stream:           imgStream.Stream,
		Dictionary:       imgStream.Dictionary.(types.StreamDictionary),
		Width:            types.Int(conf.Width),
		Height:           types.Int(conf.Height),
		BitsPerComponent: types.Int(8),
		ColorSpace:       types.ColorSpace_DeviceRGB,
	}

	// finish
	return &Image{
		Reference: q.creator.AddObject(img),
		Image:     &img,
	}, nil
}

// newImageJPG adds a new jpg image as XObject to the file
func (q *File) newImageJPG(bts []byte, conf image.Config) (*Image, error) {
	// decode image
	x, err := jpeg.Decode(bytes.NewReader(bts))
	if err != nil {
		return nil, err
	}

	// prepare Image object
	img := types.Image{
		Width:  types.Int(conf.Width),
		Height: types.Int(conf.Height),
	}

	// build data
	var data []byte
	switch conf.ColorModel {
	case color.YCbCrModel:
		img.ColorSpace = types.ColorSpace_DeviceRGB
		img.BitsPerComponent = 8
		for i := 0; i < conf.Height; i++ {
			for j := 0; j < conf.Width; j++ {
				c := x.At(j, i).(color.YCbCr)
				r, g, b := color.YCbCrToRGB(c.Y, c.Cb, c.Cr)
				data = append(data, r, g, b)
			}
		}

	case color.CMYKModel:
		img.ColorSpace = types.ColorSpace_DeviceCMYK
		img.BitsPerComponent = 8
		for i := 0; i < conf.Height; i++ {
			for j := 0; j < conf.Width; j++ {
				c := x.At(j, i).(color.CMYK)
				data = append(data, c.C, c.M, c.Y, c.K)
			}
		}

	case color.GrayModel:
		img.ColorSpace = types.ColorSpace_DeviceGray
		img.BitsPerComponent = 8
		for i := 0; i < conf.Height; i++ {
			for j := 0; j < conf.Width; j++ {
				c := x.At(j, i).(color.Gray)
				data = append(data, c.Y)
			}
		}

	case color.NRGBAModel:
		img.ColorSpace = types.ColorSpace_DeviceRGB
		img.BitsPerComponent = 8
		for i := 0; i < conf.Height; i++ {
			for j := 0; j < conf.Width; j++ {
				c := x.At(j, i).(color.NRGBA)
				data = append(data, c.R, c.G, c.B)
			}
		}

	default:
		return nil, errors.New("unsupported color model")
	}

	// create image stream
	imgStream, err := types.NewStream(data, types.Filter_FlateDecode)
	if err != nil {
		return nil, err
	}
	img.Stream = imgStream.Stream
	img.Dictionary = imgStream.Dictionary.(types.StreamDictionary)

	// finish
	return &Image{
		Reference: q.creator.AddObject(img),
		Image:     &img,
	}, nil
}

// newImagePNG adds a new png image as XObject to the file
func (q *File) newImagePNG(bts []byte, conf image.Config) (*Image, error) {
	// decode image
	x, err := png.Decode(bytes.NewReader(bts))
	if err != nil {
		return nil, err
	}

	// separate colors and transparency mask
	var data, smask []byte
	for i := 0; i < conf.Height; i++ {
		for j := 0; j < conf.Width; j++ {
			c := x.At(j, i)
			switch v := c.(type) {
			case color.NRGBA:
				data = append(data, byte(v.R), byte(v.G), byte(v.B))
				smask = append(smask, byte(v.A))
			default:
				r, g, b, a := c.RGBA()
				data = append(data, byte(r), byte(g), byte(b))
				smask = append(smask, byte(a))
			}
		}
	}

	// create image stream
	imgStream, err := types.NewStream(data, types.Filter_FlateDecode)
	if err != nil {
		return nil, err
	}
	img := types.Image{
		Stream:           imgStream.Stream,
		Dictionary:       imgStream.Dictionary.(types.StreamDictionary),
		Width:            types.Int(conf.Width),
		Height:           types.Int(conf.Height),
		BitsPerComponent: types.Int(8),
		ColorSpace:       types.ColorSpace_DeviceRGB,
	}

	// create transparency mask
	smaskStream, err := types.NewStream(smask, types.Filter_FlateDecode)
	if err != nil {
		return nil, err
	}
	dict := smaskStream.Dictionary.(types.StreamDictionary)
	dict.DecodeParms = types.Dictionary{
		"Colors":           types.Int(1),
		"BitsPerComponent": types.Int(8),
		"Columns":          types.Int(conf.Width),
	}
	smaskStream.Dictionary = dict
	smaskImg := types.Image{
		Stream:           smaskStream.Stream,
		Dictionary:       smaskStream.Dictionary.(types.StreamDictionary),
		Width:            types.Int(conf.Width),
		Height:           types.Int(conf.Height),
		ColorSpace:       types.ColorSpace_DeviceGray,
		BitsPerComponent: types.Int(8),
	}
	img.SMask = q.creator.AddObject(smaskImg)

	// finish
	return &Image{
		Reference: q.creator.AddObject(img),
		Image:     &img,
	}, nil
}

// newImageGIF adds a new gif image as XObject to the file
func (q *File) newImageGIF(bts []byte, conf image.Config) (*Image, error) {
	// decode image
	x, err := gif.Decode(bytes.NewReader(bts))
	if err != nil {
		return nil, err
	}

	// separate colors and transparency mask
	var data, smask []byte
	for i := 0; i < conf.Height; i++ {
		for j := 0; j < conf.Width; j++ {
			c := x.At(j, i)
			switch v := c.(type) {
			case color.RGBA:
				data = append(data, v.R, v.G, v.B)
				smask = append(smask, v.A)
			default:
				r, g, b, a := c.RGBA()
				data = append(data, byte(r), byte(g), byte(b))
				smask = append(smask, byte(a))
			}
		}
	}

	// create image stream
	imgStream, err := types.NewStream(data, types.Filter_FlateDecode)
	if err != nil {
		return nil, err
	}
	img := types.Image{
		Stream:           imgStream.Stream,
		Dictionary:       imgStream.Dictionary.(types.StreamDictionary),
		Width:            types.Int(conf.Width),
		Height:           types.Int(conf.Height),
		BitsPerComponent: types.Int(8),
		ColorSpace:       types.ColorSpace_DeviceRGB,
	}

	// create transparency mask
	smaskStream, err := types.NewStream(smask, types.Filter_FlateDecode)
	if err != nil {
		return nil, err
	}
	dict := smaskStream.Dictionary.(types.StreamDictionary)
	dict.DecodeParms = types.Dictionary{
		"Colors":           types.Int(1),
		"BitsPerComponent": types.Int(8),
		"Columns":          types.Int(conf.Width),
	}
	smaskStream.Dictionary = dict
	sMaskImg := types.Image{
		Stream:           smaskStream.Stream,
		Dictionary:       smaskStream.Dictionary.(types.StreamDictionary),
		Width:            types.Int(conf.Width),
		Height:           types.Int(conf.Height),
		ColorSpace:       types.ColorSpace_DeviceGray,
		BitsPerComponent: types.Int(8),
	}
	img.SMask = q.creator.AddObject(sMaskImg)

	// finish
	return &Image{
		Reference: q.creator.AddObject(img),
		Image:     &img,
	}, nil
}
