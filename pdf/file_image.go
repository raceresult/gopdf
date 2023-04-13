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
	"github.com/raceresult/tiff"
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

	case "tiff":
		return q.newImageTIFF(bts, im)

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

	// is actually grayscale?
	colorspace := types.ColorSpace_DeviceRGB
	if data2, isGray := reduceRGBToGray(data); isGray {
		data = data2
		colorspace = types.ColorSpace_DeviceGray
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
		ColorSpace:       colorspace,
	}

	// finish
	return &Image{
		Reference: q.creator.AddObject(img),
		Image:     &img,
	}, nil
}

// newImageJPG adds a new jpg image as XObject to the file
func (q *File) newImageJPG(bts []byte, conf image.Config) (*Image, error) {
	// prepare Image object
	img := types.Image{
		Width:  types.Int(conf.Width),
		Height: types.Int(conf.Height),
	}

	// for rgb, directly use the image
	if conf.ColorModel == color.YCbCrModel || conf.ColorModel == color.NRGBAModel {
		imgStream, err := types.NewStream(bts)
		if err != nil {
			return nil, err
		}
		img.ColorSpace = types.ColorSpace_DeviceRGB
		img.BitsPerComponent = 8
		img.Stream = imgStream.Stream
		img.Dictionary = imgStream.Dictionary.(types.StreamDictionary)
		img.Dictionary.Filter = []types.Filter{types.Filter_DCTDecode}
		return &Image{
			Reference: q.creator.AddObject(img),
			Image:     &img,
		}, nil
	}

	// decode image
	x, err := jpeg.Decode(bytes.NewReader(bts))
	if err != nil {
		return nil, err
	}

	// build data
	var data []byte
	switch conf.ColorModel {
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
				data = append(data, v.R, v.G, v.B)
				smask = append(smask, v.A)
			case color.NRGBA64:
				data = append(data, byte(v.R/256), byte(v.G/256), byte(v.B/256))
				smask = append(smask, byte(v.A/256))
			default:
				r, g, b, a := c.RGBA()
				data = append(data, byte(r), byte(g), byte(b))
				smask = append(smask, byte(a))
			}
		}
	}

	// is actually grayscale?
	colorspace := types.ColorSpace_DeviceRGB
	if data2, isGray := reduceRGBToGray(data); isGray {
		data = data2
		colorspace = types.ColorSpace_DeviceGray
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
		ColorSpace:       colorspace,
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

	// is actually grayscale?
	colorspace := types.ColorSpace_DeviceRGB
	if data2, isGray := reduceRGBToGray(data); isGray {
		data = data2
		colorspace = types.ColorSpace_DeviceGray
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
		ColorSpace:       colorspace,
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

// newImageTIFF adds a new tif image as XObject to the file
func (q *File) newImageTIFF(bts []byte, conf image.Config) (*Image, error) {
	// decode image
	x, err := tiff.Decode(bytes.NewReader(bts))
	if err != nil {
		return nil, err
	}

	// separate colors and transparency mask
	var data, smask []byte
	colorSpace := types.ColorSpace_DeviceRGB
	for i := 0; i < conf.Height; i++ {
		for j := 0; j < conf.Width; j++ {
			c := x.At(j, i)
			switch v := c.(type) {
			case color.RGBA:
				data = append(data, v.R, v.G, v.B)
				smask = append(smask, v.A)
			case color.NRGBA:
				data = append(data, v.R, v.G, v.B)
				smask = append(smask, v.A)
			case color.CMYK:
				c := x.At(j, i).(color.CMYK)
				data = append(data, c.C, c.M, c.Y, c.K)
				smask = append(smask, 255)
				colorSpace = types.ColorSpace_DeviceCMYK
			case tiff.CMYKA:
				c := x.At(j, i).(tiff.CMYKA)
				data = append(data, c.C, c.M, c.Y, c.K)
				smask = append(smask, c.A)
				colorSpace = types.ColorSpace_DeviceCMYK
			default:
				r, g, b, a := c.RGBA()
				data = append(data, byte(r), byte(g), byte(b))
				smask = append(smask, byte(a))
			}
		}
	}

	// is actually grayscale?
	if colorSpace == types.ColorSpace_DeviceRGB {
		if data2, isGray := reduceRGBToGray(data); isGray {
			data = data2
			colorSpace = types.ColorSpace_DeviceGray
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
		ColorSpace:       colorSpace,
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

func reduceRGBToGray(data []byte) ([]byte, bool) {
	if len(data)%3 != 0 {
		return nil, false
	}

	for i := 0; i < len(data); i += 3 {
		if data[i] != data[i+1] || data[i] != data[i+2] || data[i+1] != data[i+2] {
			return nil, false
		}
	}

	data2 := make([]byte, 0, len(data)/3)
	for i := 0; i < len(data); i += 3 {
		data2 = append(data2, data[i])
	}
	return data2, true
}
