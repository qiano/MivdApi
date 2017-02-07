package imganalyze

import (
	"fmt"
	"image"
	"image/draw"
	"math"
)

func GetWBGrayVal(img *image.Image) (int, int) {
	gray := image.NewGray((*img).Bounds())
	draw.Draw(gray, gray.Bounds(), *img, (*img).Bounds().Min, draw.Src) //原始图片转换为灰色图片
	maxx := gray.Rect.Size().X
	maxy := gray.Rect.Size().Y

	min := 255
	max := 0

	for y := 0; y < maxy; y++ {
		for x := 0; x < maxx; x++ {
			v := int(gray.GrayAt(x, y).Y)
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
		}
	}
	return min, max
}

//取得反应区线的灰度值：原图地址，反应区高(mm)，线与基线的距离(mm)
func TestLines(img *image.Image, height int, lines []int) []float64 {
	gray := image.NewGray((*img).Bounds())
	draw.Draw(gray, gray.Bounds(), *img, (*img).Bounds().Min, draw.Src) //原始图片转换为灰色图片
	fmt.Println("size=", gray.Bounds().Size().X, gray.Bounds().Size().Y)

	arr := ConvertoLine(gray)
	fmt.Println(arr)
	bval, bpos := FindMinValInLines(arr)
	wval, _ := FindMaxValInLines(arr)
	LightCompensation(arr, bval, wval) //光线补偿
	fmt.Println("补偿后：", arr)

	//反应线范围
	lineranges := CalcLineRange(len(*arr), height, bpos, lines)
	results := make([]float64, 0)
	for _, v := range lineranges {
		m2 := (*arr)[v[0]:v[1]]
		val, _ := FindMinValInLines(&m2)
		results = append(results, val)
	}

	fmt.Println("results=", results)
	return results
}

//找到灰度值最大的点
func FindMaxGray(gray *image.Gray) (val, x, y int) {
	maxx := gray.Rect.Size().X
	maxy := gray.Rect.Size().Y
	var max uint8 = gray.GrayAt(0, 0).Y
	var px, py int = 0, 0
	for y := 0; y < maxy; y++ {
		for x := 0; x < maxx; x++ {
			if gray.GrayAt(x, y).Y > max {
				max = gray.GrayAt(x, y).Y
				px = x
				py = y
			}
		}
	}
	return int(max), px, py
}

//找到灰度值最小的点
func FindMinGray(gray *image.Gray) (val, x, y int) {
	maxx := gray.Rect.Size().X
	maxy := gray.Rect.Size().Y
	var min uint8 = gray.GrayAt(0, 0).Y
	var px, py int = 0, 0
	for y := 0; y < maxy; y++ {
		for x := 0; x < maxx; x++ {
			if gray.GrayAt(x, y).Y < min {
				min = gray.GrayAt(x, y).Y
				px = x
				py = y
			}
		}
	}
	return int(min), px, py
}

//转换为线数组
func ConvertoLine(gray *image.Gray) *[]float64 {
	arr := make([]float64, 0, gray.Rect.Size().Y)
	maxx := gray.Rect.Size().X
	maxy := gray.Rect.Size().Y
	for y := 0; y < maxy; y++ {
		xsum := 0
		for x := 0; x < maxx; x++ {
			xsum += int(gray.Pix[maxx*y+x])
		}
		arr = append(arr, float64(xsum)/float64(maxx))
	}
	return &arr
}

//转换中间部分数据为线数组
func ConverParttoLine(gray *image.Gray,rate float64) *[]float64 {
	arr := make([]float64, 0, gray.Rect.Size().Y)
	minx:=int(math.Ceil((1-rate)/2*float64(gray.Rect.Size().X)))
	maxx := int(math.Floor((1+rate)/2*float64(gray.Rect.Size().X)))
	maxy := gray.Rect.Size().Y
	for y := 0; y < maxy; y++ {
		xsum := 0
		for x := minx; x < maxx; x++ {
			xsum += int(gray.Pix[gray.Rect.Size().X*y+x])
		}
		arr = append(arr, float64(xsum)/float64(maxx-minx))
	}
	return &arr
}

//查找灰度值最小的线
func FindMinValInLines(arr *[]float64) (val float64, idx int) {
	for i, v := range *arr {
		if i == 0 {
			val = v
			idx = i
		} else {
			if v < val {
				val = v
				idx = i
			}
		}
	}
	return val, idx
}

//查找灰度值最大的线
func FindMaxValInLines(arr *[]float64) (val float64, idx int) {
	for i, v := range *arr {
		if i == 0 {
			val = v
			idx = i
		} else {
			if v > val {
				val = v
				idx = i
			}
		}
	}
	return val, idx
}

//计算反应线范围
func CalcLineRange(pixheight int, height int, basepix int, lines []int) map[int][]int {
	unit := float64(pixheight) / float64(height) //1mm=?pix
	fmt.Println(unit)
	m1 := make(map[int][]int)
	for i := 0; i < len(lines); i++ {
		min := basepix + int(math.Floor(float64(lines[i]-1)*unit))
		max := basepix + int(math.Ceil(float64(lines[i]+1)*unit))
		m1[i] = []int{min, max}
	}
	return m1
}

//计算范围内的平均值
func CalcRangeAvg(arr *[]float64, min, max int) float64 {
	var sum float64
	for i := min; i < max; i++ {
		sum += (*arr)[i]
	}
	return sum / float64(max-min)
}

//光线补偿
func LightCompensation(arr *[]float64, black float64, white float64) {
	for i, v := range *arr {
		(*arr)[i] = float64(255) * (v - black) / (white - black)
	}
}
