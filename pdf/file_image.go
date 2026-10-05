package pdf

import (
	"bytes"
	"compress/zlib"
	"errors"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"runtime/debug"

	"github.com/raceresult/gopdf/pdf/imagecache"
	"github.com/raceresult/gopdf/types"
	"github.com/raceresult/tiff"
	"golang.org/x/image/bmp"
)

// using the imageCache, images need to be decoded/encoded only once
var imageCache = imagecache.New(1024 * 1024 * 100) // 100 MB

// Image holds both the Image object and the reference to it
type Image struct {
	Reference types.Reference
	Image     *types.Image
}

// NewImage adds a new image as XObject to the file and returns the image object and its reference
func (q *File) NewImage(bts []byte) (*Image, error) {
	return q.newImage(bts, false)
}

// NewImageConcurrent adds a new image as XObject to the file and returns the image object and its reference. Can be called multiple times currently.
func (q *File) NewImageConcurrent(bts []byte) (*Image, error) {
	return q.newImage(bts, true)
}

// newImage adds a new image as XObject to the file and returns the image object and its reference
func (q *File) newImage(bts []byte, theadSafe bool) (*Image, error) {
	// read config
	im, name, err := image.DecodeConfig(bytes.NewReader(bts))
	if err != nil {
		return nil, err
	}

	// continue depending on type
	switch name {
	case "bmp":
		return q.newImageBmp(bts, im, theadSafe)

	case "jpg", "jpeg":
		return q.newImageJPG(bts, im, theadSafe)

	case "png":
		return q.newImagePNG(bts, im, theadSafe)

	case "gif":
		return q.newImageGIF(bts, im, theadSafe)

	case "tiff":
		return q.newImageTIFF(bts, im, theadSafe)

	default:
		return nil, errors.New("unsupported image type " + name)
	}
}

// newImageBmp adds a new bmp file as XObject to the file
func (q *File) newImageBmp(bts []byte, conf image.Config, theadSafe bool) (*Image, error) {
	item, err := imageCache.Process(bts, func(bts []byte) (*imagecache.Item, error) {
		// decode image
		x, err := bmp.Decode(bytes.NewReader(bts))
		if err != nil {
			return nil, err
		}

		// build data
		var colorModel types.ColorSpaceFamily
		var data []byte
		var destData bytes.Buffer
		wData := zlib.NewWriter(&destData)
		if isGrayScales(conf, x) {
			colorModel = types.ColorSpace_DeviceGray
			for i := 0; i < conf.Height; i++ {
				data = make([]byte, 0, conf.Width)
				for j := 0; j < conf.Width; j++ {
					r, _, _, _ := x.At(j, i).RGBA()
					data = append(data, byte(r))
				}
				if _, err := wData.Write(data); err != nil {
					return nil, err
				}
			}
		} else {
			colorModel = types.ColorSpace_DeviceRGB
			for i := 0; i < conf.Height; i++ {
				data = make([]byte, 0, conf.Width*3)
				for j := 0; j < conf.Width; j++ {
					r, g, b, _ := x.At(j, i).RGBA()
					data = append(data, byte(r), byte(g), byte(b))
				}
				if _, err := wData.Write(data); err != nil {
					return nil, err
				}
			}
		}

		// free memory
		x = nil
		if conf.Width*conf.Height > 1024*1024 {
			debug.FreeOSMemory()
		}

		// finish zlib writers
		if err := wData.Close(); err != nil {
			return nil, err
		}

		// return item
		return &imagecache.Item{
			Data:       destData,
			ColorModel: colorModel,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	// create image stream
	imgStream := types.StreamObject{
		Dictionary: types.StreamDictionary{
			Filter: []types.Filter{types.Filter_FlateDecode},
			Length: item.Data.Len(),
		},
		Stream: item.Data.Bytes(),
	}
	img := types.Image{
		Stream:           imgStream.Stream,
		Dictionary:       imgStream.Dictionary.(types.StreamDictionary),
		Width:            types.Int(conf.Width),
		Height:           types.Int(conf.Height),
		BitsPerComponent: types.Int(8),
		ColorSpace:       item.ColorModel,
	}

	// finish
	if theadSafe {
		q.newImageMux.Lock()
		defer q.newImageMux.Unlock()
	}
	return &Image{
		Reference: q.creator.AddObject(img),
		Image:     &img,
	}, nil
}

// newImageJPG adds a new jpg image as XObject to the file
func (q *File) newImageJPG(bts []byte, conf image.Config, theadSafe bool) (*Image, error) {
	// for rgb, directly use the image
	if conf.ColorModel == color.YCbCrModel || conf.ColorModel == color.NRGBAModel {
		imgStream, err := types.NewStream(bts)
		if err != nil {
			return nil, err
		}
		img := types.Image{
			Width:            types.Int(conf.Width),
			Height:           types.Int(conf.Height),
			ColorSpace:       types.ColorSpace_DeviceRGB,
			BitsPerComponent: 8,
			Stream:           imgStream.Stream,
			Dictionary:       imgStream.Dictionary.(types.StreamDictionary),
		}
		img.Dictionary.Filter = []types.Filter{types.Filter_DCTDecode}
		if theadSafe {
			q.newImageMux.Lock()
			defer q.newImageMux.Unlock()
		}
		return &Image{
			Reference: q.creator.AddObject(img),
			Image:     &img,
		}, nil
	}

	item, err := imageCache.Process(bts, func(bts []byte) (*imagecache.Item, error) {
		// decode image
		x, err := jpeg.Decode(bytes.NewReader(bts))
		if err != nil {
			return nil, err
		}

		// build data
		var colorModel types.ColorSpaceFamily
		var data []byte
		var destData bytes.Buffer
		wData := zlib.NewWriter(&destData)
		switch conf.ColorModel {
		case color.CMYKModel:
			colorModel = types.ColorSpace_DeviceCMYK
			for i := 0; i < conf.Height; i++ {
				data = make([]byte, 0, conf.Width*4)
				for j := 0; j < conf.Width; j++ {
					c := x.At(j, i).(color.CMYK)
					data = append(data, c.C, c.M, c.Y, c.K)
				}
				if _, err := wData.Write(data); err != nil {
					return nil, err
				}
			}

		case color.GrayModel:
			colorModel = types.ColorSpace_DeviceGray
			for i := 0; i < conf.Height; i++ {
				data = make([]byte, 0, conf.Width)
				for j := 0; j < conf.Width; j++ {
					c := x.At(j, i).(color.Gray)
					data = append(data, c.Y)
				}
				if _, err := wData.Write(data); err != nil {
					return nil, err
				}
			}

		default:
			return nil, errors.New("unsupported color model")
		}

		// free memory
		x = nil
		if conf.Width*conf.Height > 1024*1024 {
			debug.FreeOSMemory()
		}

		// finish zlib writers
		if err := wData.Close(); err != nil {
			return nil, err
		}

		// return item
		return &imagecache.Item{
			Data:       destData,
			Mask:       bytes.Buffer{},
			ColorModel: colorModel,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	// create image stream
	imgStream := types.StreamObject{
		Dictionary: types.StreamDictionary{
			Filter: []types.Filter{types.Filter_FlateDecode},
			Length: item.Data.Len(),
		},
		Stream: item.Data.Bytes(),
	}
	img := types.Image{
		Width:            types.Int(conf.Width),
		Height:           types.Int(conf.Height),
		Stream:           imgStream.Stream,
		Dictionary:       imgStream.Dictionary.(types.StreamDictionary),
		ColorSpace:       item.ColorModel,
		BitsPerComponent: 8,
	}

	// finish
	if theadSafe {
		q.newImageMux.Lock()
		defer q.newImageMux.Unlock()
	}
	return &Image{
		Reference: q.creator.AddObject(img),
		Image:     &img,
	}, nil
}

// newImagePNG adds a new png image as XObject to the file
func (q *File) newImagePNG(bts []byte, conf image.Config, theadSafe bool) (*Image, error) {
	item, err := imageCache.Process(bts, func(bts []byte) (*imagecache.Item, error) {
		// decode image
		x, err := png.Decode(bytes.NewReader(bts))
		if err != nil {
			return nil, err
		}

		// build data
		var colorModel types.ColorSpaceFamily
		var destData, destMask bytes.Buffer
		var data, smask []byte
		wData := zlib.NewWriter(&destData)
		wMask := zlib.NewWriter(&destMask)
		if isGrayScales(conf, x) {
			colorModel = types.ColorSpace_DeviceGray
			for i := 0; i < conf.Height; i++ {
				data = make([]byte, 0, conf.Width)
				smask = make([]byte, 0, conf.Width)
				for j := 0; j < conf.Width; j++ {
					c := x.At(j, i)
					switch v := c.(type) {
					case color.NRGBA:
						data = append(data, v.R)
						smask = append(smask, v.A)
					case color.NRGBA64:
						data = append(data, byte(v.R/256))
						smask = append(smask, byte(v.A/256))
					default:
						r, _, _, a := c.RGBA()
						data = append(data, byte(r))
						smask = append(smask, byte(a))
					}
				}
				if _, err := wData.Write(data); err != nil {
					return nil, err
				}
				if _, err := wMask.Write(smask); err != nil {
					return nil, err
				}
			}

		} else {
			colorModel = types.ColorSpace_DeviceRGB
			for i := 0; i < conf.Height; i++ {
				data = make([]byte, 0, conf.Width*3)
				smask = make([]byte, 0, conf.Width)
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
				if _, err := wData.Write(data); err != nil {
					return nil, err
				}
				if _, err := wMask.Write(smask); err != nil {
					return nil, err
				}
			}
		}

		// free memory
		x = nil
		if conf.Width*conf.Height > 1024*1024 {
			debug.FreeOSMemory()
		}

		// finish zlib writers
		if err := wData.Close(); err != nil {
			return nil, err
		}
		if err := wMask.Close(); err != nil {
			return nil, err
		}
		return &imagecache.Item{
			Data:       destData,
			Mask:       destMask,
			ColorModel: colorModel,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	// create image stream
	imgStream := types.StreamObject{
		Dictionary: types.StreamDictionary{
			Filter: []types.Filter{types.Filter_FlateDecode},
			Length: item.Data.Len(),
		},
		Stream: item.Data.Bytes(),
	}
	img := types.Image{
		Stream:           imgStream.Stream,
		Dictionary:       imgStream.Dictionary.(types.StreamDictionary),
		Width:            types.Int(conf.Width),
		Height:           types.Int(conf.Height),
		BitsPerComponent: types.Int(8),
		ColorSpace:       item.ColorModel,
	}

	// create transparency mask
	smaskStream := types.StreamObject{
		Dictionary: types.StreamDictionary{
			Filter: []types.Filter{types.Filter_FlateDecode},
			Length: item.Mask.Len(),
		},
		Stream: item.Mask.Bytes(),
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
	if theadSafe {
		q.newImageMux.Lock()
		defer q.newImageMux.Unlock()
	}
	img.SMask = q.creator.AddObject(smaskImg)

	// finish
	// <<< newImageMux above
	return &Image{
		Reference: q.creator.AddObject(img),
		Image:     &img,
	}, nil
}

// newImageGIF adds a new gif image as XObject to the file
func (q *File) newImageGIF(bts []byte, conf image.Config, theadSafe bool) (*Image, error) {
	item, err := imageCache.Process(bts, func(bts []byte) (*imagecache.Item, error) {
		// decode image
		x, err := gif.Decode(bytes.NewReader(bts))
		if err != nil {
			return nil, err
		}

		// separate colors and transparency mask
		var colorModel types.ColorSpaceFamily
		var data, smask []byte
		var destData, destMask bytes.Buffer
		wData := zlib.NewWriter(&destData)
		wMask := zlib.NewWriter(&destMask)
		if isGrayScales(conf, x) {
			colorModel = types.ColorSpace_DeviceGray
			for i := 0; i < conf.Height; i++ {
				data = make([]byte, 0, conf.Width)
				smask = make([]byte, 0, conf.Width)
				for j := 0; j < conf.Width; j++ {
					c := x.At(j, i)
					switch v := c.(type) {
					case color.RGBA:
						data = append(data, v.R)
						smask = append(smask, v.A)
					default:
						r, _, _, a := c.RGBA()
						data = append(data, byte(r))
						smask = append(smask, byte(a))
					}
				}
				if _, err := wData.Write(data); err != nil {
					return nil, err
				}
				if _, err := wMask.Write(smask); err != nil {
					return nil, err
				}
			}
		} else {
			colorModel = types.ColorSpace_DeviceRGB
			for i := 0; i < conf.Height; i++ {
				data = make([]byte, 0, conf.Width*3)
				smask = make([]byte, 0, conf.Width)
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
				if _, err := wData.Write(data); err != nil {
					return nil, err
				}
				if _, err := wMask.Write(smask); err != nil {
					return nil, err
				}
			}
		}

		// free memory
		x = nil
		if conf.Width*conf.Height > 1024*1024 {
			debug.FreeOSMemory()
		}

		// finish zlib writers
		if err := wData.Close(); err != nil {
			return nil, err
		}
		if err := wMask.Close(); err != nil {
			return nil, err
		}

		// return item
		return &imagecache.Item{
			Data:       destData,
			Mask:       destMask,
			ColorModel: colorModel,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	// create image stream
	imgStream := types.StreamObject{
		Dictionary: types.StreamDictionary{
			Filter: []types.Filter{types.Filter_FlateDecode},
			Length: item.Data.Len(),
		},
		Stream: item.Data.Bytes(),
	}
	img := types.Image{
		Stream:           imgStream.Stream,
		Dictionary:       imgStream.Dictionary.(types.StreamDictionary),
		Width:            types.Int(conf.Width),
		Height:           types.Int(conf.Height),
		BitsPerComponent: types.Int(8),
		ColorSpace:       item.ColorModel,
	}

	// create transparency mask
	smaskStream := types.StreamObject{
		Dictionary: types.StreamDictionary{
			Filter: []types.Filter{types.Filter_FlateDecode},
			Length: item.Mask.Len(),
		},
		Stream: item.Mask.Bytes(),
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
	if theadSafe {
		q.newImageMux.Lock()
		defer q.newImageMux.Unlock()
	}
	img.SMask = q.creator.AddObject(sMaskImg)

	// finish
	// <<< newImageMux above
	return &Image{
		Reference: q.creator.AddObject(img),
		Image:     &img,
	}, nil
}

// newImageTIFF adds a new tif image as XObject to the file
func (q *File) newImageTIFF(bts []byte, conf image.Config, theadSafe bool) (*Image, error) {
	item, err := imageCache.Process(bts, func(bts []byte) (*imagecache.Item, error) {
		// decode image
		x, err := tiff.Decode(bytes.NewReader(bts))
		if err != nil {
			return nil, err
		}

		// separate colors and transparency mask
		var colorModel types.ColorSpaceFamily
		var data, smask []byte
		var destData, destMask bytes.Buffer
		wData := zlib.NewWriter(&destData)
		wMask := zlib.NewWriter(&destMask)
		if isGrayScales(conf, x) {
			colorModel = types.ColorSpace_DeviceGray
			for i := 0; i < conf.Height; i++ {
				data = make([]byte, 0, conf.Width)
				smask = make([]byte, 0, conf.Width)
				for j := 0; j < conf.Width; j++ {
					c := x.At(j, i)
					switch v := c.(type) {
					case color.RGBA:
						data = append(data, v.R)
						smask = append(smask, v.A)
					case color.NRGBA:
						data = append(data, v.R)
						smask = append(smask, v.A)
					default:
						r, _, _, a := c.RGBA()
						data = append(data, byte(r))
						smask = append(smask, byte(a))
					}
				}
				if _, err := wData.Write(data); err != nil {
					return nil, err
				}
				if _, err := wMask.Write(smask); err != nil {
					return nil, err
				}
			}
		} else {
			colorModel = types.ColorSpace_DeviceRGB
			for i := 0; i < conf.Height; i++ {
				data = make([]byte, 0, conf.Width*4)
				smask = make([]byte, 0, conf.Width)
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
						colorModel = types.ColorSpace_DeviceCMYK
					case tiff.CMYKA:
						c := x.At(j, i).(tiff.CMYKA)
						data = append(data, c.C, c.M, c.Y, c.K)
						smask = append(smask, c.A)
						colorModel = types.ColorSpace_DeviceCMYK
					default:
						r, g, b, a := c.RGBA()
						data = append(data, byte(r), byte(g), byte(b))
						smask = append(smask, byte(a))
					}
				}
				if _, err := wData.Write(data); err != nil {
					return nil, err
				}
				if _, err := wMask.Write(smask); err != nil {
					return nil, err
				}
			}
		}

		// free memory
		x = nil
		if conf.Width*conf.Height > 1024*1024 {
			debug.FreeOSMemory()
		}

		// finish zlib writers
		if err := wData.Close(); err != nil {
			return nil, err
		}
		if err := wMask.Close(); err != nil {
			return nil, err
		}

		// return item
		return &imagecache.Item{
			Data:       destData,
			Mask:       destMask,
			ColorModel: colorModel,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	// create image stream
	imgStream := types.StreamObject{
		Dictionary: types.StreamDictionary{
			Filter: []types.Filter{types.Filter_FlateDecode},
			Length: item.Data.Len(),
		},
		Stream: item.Data.Bytes(),
	}
	img := types.Image{
		Stream:           imgStream.Stream,
		Dictionary:       imgStream.Dictionary.(types.StreamDictionary),
		Width:            types.Int(conf.Width),
		Height:           types.Int(conf.Height),
		BitsPerComponent: types.Int(8),
		ColorSpace:       item.ColorModel,
	}

	// create transparency mask
	smaskStream := types.StreamObject{
		Dictionary: types.StreamDictionary{
			Filter: []types.Filter{types.Filter_FlateDecode},
			Length: item.Mask.Len(),
		},
		Stream: item.Mask.Bytes(),
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
	if theadSafe {
		q.newImageMux.Lock()
		defer q.newImageMux.Unlock()
	}
	img.SMask = q.creator.AddObject(sMaskImg)

	// finish
	// <<< newImageMux above
	return &Image{
		Reference: q.creator.AddObject(img),
		Image:     &img,
	}, nil
}

func isGrayScales(conf image.Config, img image.Image) bool {
	for i := 0; i < conf.Height; i++ {
		for j := 0; j < conf.Width; j++ {
			c := img.At(j, i)
			switch v := c.(type) {
			case color.NRGBA:
				if v.R != v.G || v.R != v.B || v.G != v.B {
					return false
				}
			case color.NRGBA64:
				if v.R != v.G || v.R != v.B || v.G != v.B {
					return false
				}
			case color.CMYK:
				return false
			case tiff.CMYKA:
				return false
			default:
				r, g, b, _ := c.RGBA()
				if r != g || r != b || g != b {
					return false
				}
			}
		}
	}
	return true
}
