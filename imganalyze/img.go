package imganalyze

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

//图像解码
func DecodeImg(sourcepath string) *image.Image {
	ext := filepath.Ext(sourcepath)
	if ext == ".jpg" || ext == ".jpeg" {
		img, er := TryDecodeImg(sourcepath, "jpg")
		if er == nil {
			return img
		}
		img2, er2 := TryDecodeImg(sourcepath, "png")
		if er2 == nil {
			return img2
		}
	}
	if ext == ".png" {
		img2, er2 := TryDecodeImg(sourcepath, "png")
		if er2 == nil {
			return img2
		}
		img, er := TryDecodeImg(sourcepath, "jpg")
		if er == nil {
			return img
		}
	}
	return nil
}

//尝试图像解码
func TryDecodeImg(sourcepath string, code string) (*image.Image, error) {
	sourcefile, err := os.Open(sourcepath)
	if err != nil {
		fmt.Println(err)
	}
	defer sourcefile.Close()

	if code == "jpg" {
		img, er := jpeg.Decode(sourcefile)
		if er == nil {
			return &img, nil
		} else {
			return &img, er
		}
	}

	if code == "png" {
		img2, er2 := png.Decode(sourcefile)
		if er2 == nil {
			return &img2, nil
		} else {
			return &img2, er2
		}
	}
	return nil, errors.New("只支持格式：jpg和png")
}

//保存图片
func SaveImage(data *image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	png.Encode(file, *data)
	return nil
}

//图片裁剪
func ImgCut(img *image.Image, x0, y0, x1, y1 int) image.Image {

	gray := image.NewGray((*img).Bounds())
	draw.Draw(gray, gray.Bounds(), *img, (*img).Bounds().Min, draw.Src)
	rect := image.Rect(x0, y0, x1, y1)
	subimg := gray.SubImage(rect)
	return subimg
}
