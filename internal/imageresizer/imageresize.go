package imageresizer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/nfnt/resize"
)

// DecodeImageBase64 ...
func DecodeImageBase64(base64Image string) (image.Image, error) {

	bytesImage, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		return nil, fmt.Errorf("could not decode image form string to bytes")
	}
	img, _, err := image.Decode(bytes.NewReader(bytesImage))

	return img, nil
}

// EncodeImageBase64 ...
func EncodeImageBase64(img []byte) (string, error) {

	Base64String := base64.StdEncoding.EncodeToString(img)
	return Base64String, nil
}

// ResizeImage ...
func ResizeImage(img image.Image, maxWidth uint, maxHeight uint, format string) ([]byte, error) {

	thumb := resize.Thumbnail(maxWidth, maxHeight, img, resize.Lanczos3)

	var err error
	buf := new(bytes.Buffer)
	switch format {
	case "jpg":
		err = jpeg.Encode(buf, thumb, nil)
	case "jpeg":
		err = jpeg.Encode(buf, thumb, nil)
	case "png":
		err = png.Encode(buf, thumb)
	default:
		return nil, fmt.Errorf("unknown format")
	}

	if err != nil {
		return nil, fmt.Errorf("could not create thumnail from image")
	}

	return []byte(buf.Bytes()), nil
}
