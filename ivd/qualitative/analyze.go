package qualitative

import (
	"fmt"
	"image"
	"image/draw"
	"math"

	"github.com/qshuai162/MivdApi_Trail/imganalyze"
)

//Test ite
func Test(picpath string, hotWidth, cdis int, tdis []int, white, black int) (ret bool, judges []bool, grays []float64) {
	img := imganalyze.DecodeImg(picpath)
	gray := image.NewGray((*img).Bounds())
	draw.Draw(gray, gray.Bounds(), *img, (*img).Bounds().Min, draw.Src) //原始图片转换为灰色图片
	arr := imganalyze.ConverParttoLine(gray, 0.5)
	unit := float64((*img).Bounds().Max.X) / float64(hotWidth)
	offset := unit * 1.5
	fmt.Println("wb",white, black)
	cmin := int(math.Floor(float64(cdis)*unit - offset))
	cmax := int(math.Ceil(float64(cdis)*unit + offset))
	carr := (*arr)[cmin:cmax]
	cval, _ := imganalyze.FindMinValInLines(&carr) //寻找最低点，判定为C线
	cval = Adjust(cval, white, black)
	fmt.Println("cval", cval)
	//计算除C T线区域外的灰度平均值
	avg := imganalyze.CalcRangeAvg(arr, cmax, int(math.Floor(float64(tdis[0])*unit-offset)))
	avg = Adjust(avg, white, black)
	fmt.Println("avg", avg)
	//c线是否存在
	fmt.Println(Adjust(10, white, black))
	if avg-cval < Adjust(10, white, black) {
		return false, judges, grays
	}

	for i := 0; i < len(tdis); i++ {
		//在C线固定距离范围内寻找T线，最高点
		tmin := int(math.Floor(float64(tdis[i])*unit - offset))
		tmax := int(math.Ceil(float64(tdis[i])*unit + offset))
		tarr := (*arr)[tmin:tmax]
		tval, _ := imganalyze.FindMinValInLines(&tarr)
		tval = Adjust(tval, white, black)
		fmt.Println("tval", tval)
		fmt.Println("offset", Adjust(0.8, white, black))
		judges = append(judges, avg-tval > Adjust(0.7, white, black))
		grays = append(grays, tval)

	}
	return true, judges, grays
}


//BWValue 获取黑白色灰度值
func BWValue(picpath string) (int, int) {
	img := imganalyze.DecodeImg(picpath)
	gray := image.NewGray((*img).Bounds())
	draw.Draw(gray, gray.Bounds(), *img, (*img).Bounds().Min, draw.Src) //原始图片转换为灰色图片
	black, _, _ := imganalyze.FindMinGray(gray)
	white, _, _ := imganalyze.FindMaxGray(gray)
	return white, black
}

//Adjust 校准灰度
func Adjust(val float64, white, black int) float64 {
	return (float64((white-black))/255.0)*val + float64(black)
}
