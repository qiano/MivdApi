package imganalyze

import (
	// "image"
	// "image/draw"
	"testing"
)

func TestTestLines(t *testing.T) {

	img := DecodeImg("test.jpg")
	if img == nil {
		t.Fatal("解码失败")
	}
	subimg := ImgCut(img, 0, 0, 19, 128)
	//一条线，找出反应线中灰度值最小的线
	ret := TestLines(&subimg, 15, []int{3})
	if ret[0] != 117.42133537989254 {
		t.Fatal("失败：一条线，找出反应线中灰度值最小的线")
	}
}
