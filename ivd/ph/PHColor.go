//提供了三种算法判断PH值
//向量夹角余弦算法、颜色距离之和法、单色块颜色距离最小法
//单色块颜色距离最小法最准确

package ph

import (
	"fmt"
	"image"
	"image/draw"
	"math"	
	"github.com/qshuai162/MivdApi/imganalyze"

)

func Test(picpath string) float64 {
	fmt.Println(picpath)
	img := imganalyze.DecodeImg(picpath)
    gray := image.NewRGBA((*img).Bounds())
	draw.Draw(gray, gray.Bounds(), *img, (*img).Bounds().Min, draw.Src) //原始图片转换为灰色图片
	
	width := 30 //截图区域大小，毫米坐标系
	length := 125.4
	x := []int{5, 12, 19, 26} //phcolor7  毫米坐标系，每个色块最中心点的坐标
	y := []int{3, 11, 18, 25, 33, 41, 48, 56, 78, 85, 93, 100, 108, 115, 123, 63}
	
	xaxis, yaxis := coordinatesConvert(img,x, y, width, length)
	phMATRIX := colorfunc(img,xaxis, yaxis) //得到颜色矩阵
	slice := phMATRIX[15]

	return rbgDisMax(phMATRIX, slice)
}

//每个PH值四个方块的颜色值切片
func colorfunc(img *image.Image, xaxis []int, yaxis []int) [][]int {
	color := image.NewRGBA((*img).Bounds())
	draw.Draw(color, color.Bounds(), *img, (*img).Bounds().Min, draw.Src)
	ph := make([][]int, 16, 1024)
	for m := 0; m < len(yaxis); {
		arrsum := make([]int, 3)
		for n := 0; n < len(xaxis); {
			// 以下是计算单个目标区域的RGB值
			for i := yaxis[m] - 4; i <= yaxis[m]+4; i++ {
				for j := xaxis[n] - 4; j <= xaxis[n]+4; j++ {
					r := color.RGBAAt(j, i)
					arr := make([]int, 0)
					arr = append(arr, int(r.R), int(r.G), int(r.B))
					// fmt.Println(arr, i, j)
					// fmt.Println(r)

					for k := 0; k < len(arr); k++ {
						arrsum[k] += arr[k]
					}
				}
			}
			arrsum[0] = arrsum[0] / 81
			arrsum[1] = arrsum[1] / 81
			arrsum[2] = arrsum[2] / 81
			// fmt.Println(arrsum)
			ph[m] = append(ph[m], arrsum...)
			n++
		}
		m++
	}
	return ph
}

//颜色相似度算法一：向量夹角余弦
func vectorCOS(a []int, b []int) float64 {
	c := 0.0
	o := 0.0
	s := 0.0
	for i := 0; i < 12; i++ {
		c += float64(a[i] * b[i])
		o += float64((a[i] * a[i]))
		s += float64((b[i] * b[i]))
	}
	return c / (math.Sqrt(o) * math.Sqrt(s))
}

//颜色相似度算法二：PH颜色距离之和，和最小者为目标PH值
func rgbDistance(phControl []int, phTest []int) float64 {
	sumDis := 0.0

	for i := 0; i < len(phControl); {
		Dis := 0.0
		Dis += math.Pow(float64(phControl[i]-phTest[i]), 2)
		Dis += math.Pow(float64(phControl[i+1]-phTest[i+1]), 2)
		Dis += math.Pow(float64(phControl[i+2]-phTest[i+2]), 2)
		Dis = math.Sqrt(Dis)

		i += 3
		sumDis += Dis
	}
	return sumDis
}

//颜色相似度算法三：距离最小者为目标值
func rbgDisMax(phControl [][]int, phTest []int) float64 {
	arrDisMax := make([]float64, 0)
	for j := 0; j < 15; j++ {
		DisMax := 0.0
		for i := 0; i < len(phControl[j]); {
			Dis := 0.0
			Dis += math.Pow(float64(phControl[j][i]-phTest[i]), 2)
			Dis += math.Pow(float64(phControl[j][i+1]-phTest[i+1]), 2)
			Dis += math.Pow(float64(phControl[j][i+2]-phTest[i+2]), 2)
			Dis = math.Sqrt(Dis)
			if Dis > DisMax {
				DisMax = Dis
			}
			i += 3
		}
		arrDisMax = append(arrDisMax, DisMax)
	}
	var phVal, phvaltest int
	k := arrDisMax[0]
	for i := 0; i < len(arrDisMax)-1; i++ {
		if k > arrDisMax[i+1] {
			k = arrDisMax[i+1]
			phVal = i + 1
		}
	}
	disMintest := arrDisMax[0]
	for i := 0; i < len(arrDisMax)-1; i++ {
		if disMintest > arrDisMax[i] && arrDisMax[i] > k+5 {
			disMintest = arrDisMax[i]
			phvaltest = i
		}
	}
	var phResult float64
	if 14-phvaltest > 14-phVal {
		phResult = float64(14-phVal) + arrDisMax[phVal]/(arrDisMax[phVal]+arrDisMax[phvaltest])
	}
	if 14-phvaltest < 14-phVal {
		phResult = float64(14-phVal) - arrDisMax[phVal]/(arrDisMax[phVal]+arrDisMax[phvaltest])
	}

	phResult = math.Trunc(phResult*1e1+0.5) * 1e-1 //四舍五入，只保留一位小数

	return phResult

}

//PH值判定:距离判断法
func pHDisJudge(arrDisMax []float64) int {
	var phVal int
	k := arrDisMax[0]
	for i := 0; i < len(arrDisMax)-1; i++ {

		if k > arrDisMax[i+1] {
			k = arrDisMax[i+1]
			phVal = i + 1
		}
	}
	return 14 - phVal
}

//PH值判定：向量夹角余弦判断法
func pHvectCosJudge(arrVectCos []float64) int {
	var phVal int
	k := arrVectCos[0]
	for i := 0; i < len(arrVectCos)-1; i++ {

		if k < arrVectCos[i+1] {
			k = arrVectCos[i+1]
			phVal = i + 1
		}
	}
	return 14 - phVal
}

//自动计算采样坐标
func coordinatesConvert(img *image.Image, x []int, y []int, width int, length float64) ([]int, []int) {
	a := (*img).Bounds().Max.X
	b := (*img).Bounds().Max.Y

	for i := 0; i < len(x); i++ {
		x[i] = int(float64(x[i]) * float64(float64(a)/float64(width)))
	}
	for i := 0; i < len(y); i++ {

		y[i] = int(float64(y[i]) * float64(float64(b)/float64(length)))

	}

	return x, y
}
