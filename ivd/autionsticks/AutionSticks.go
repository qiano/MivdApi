package autionsticks
import (
	// "fmt"
	"image"
	"image/draw"
	"math"
)

func Autionsticks(img *image.Image) []string {
	width := 66                                                                         //截图区域大小，毫米坐标系
	length := 77                                                                        //AutionSticks  毫米坐标系，每个色块最中心点的坐标
	x := []float64{4, 18.5, 26.5, 35, 44, 52.5, 62.5, 10.5}                             //
	y := []float64{4.25, 11.25, 18.25, 25.25, 32.25, 39.25, 46.25, 53.25, 60.25, 67.25} //纵坐标一样
	TestResult := test(img, x, y, width, length)
	return TestResult
}

func test(img *image.Image, x, y []float64, width, length int) []string {
	xaxis, yaxis := coordinatesConvertAution(img, x, y, width, length)
	phMATRIX := colorfuncAution(img, xaxis, yaxis) //得到颜色矩阵
	testResult := make([]string, 0)

	ptt:=[]string{"normal","±（50）mg/dL","1+（100）mg/dL","2+（200）mg/dL","3+（500）mg/dL","4+（1000）mg/dL","Default"}
	testResult = append(testResult, "葡萄糖："+ ptt[rbgDisMaxAution(phMATRIX[0])])

	dbz:=[]string{"neg.","±（15）mg/dL","1+（30）mg/dL","2+（100）mg/dL","3+（300）mg/dL","4+（1000）mg/dL", "Default"}
	testResult = append(testResult, "蛋白质："+ dbz[rbgDisMaxAution(phMATRIX[1])])
	
	ndy:=[]string{"normal","Default","1+（2）mg/dL","2+（4）mg/dL","3+（8）mg/dL","4+（OVER）mg/dL","Default"}
	testResult = append(testResult, "尿胆原："+ ndy[rbgDisMaxAution(phMATRIX[2])])

	dhs:=[]string{"neg.","Default","1+（0.5）mg/dL","2+（2）mg/dL","3+（6）mg/dL","4+（OVER）mg/dL", "Default"}
	testResult = append(testResult, "胆红素："+ dhs[rbgDisMaxAution(phMATRIX[3])])

	jg:=[]string{"10 mg/dL","50 mg/dL", "100 mg/dL","200 mg/dL","300 mg/dL","Default","Default"}
	testResult = append(testResult, "肌酐："+ jg[rbgDisMaxAution(phMATRIX[4])])

	ph:=[]string{"5","6","7","8","9","Default","Default"}
	testResult = append(testResult, "PH："+ ph[rbgDisMaxAution(phMATRIX[5])])
	
	yx:=[]string{"neg.","Hemolysis ±（0.03）mg/dL","1+（0.06）mg/dL","2+（0.2）mg/dL","3+（1.0）mg/dL","Non Hemolysis 1+ mg/dL","2+ mg/dL","Non Hemolysis 1+ mg/dL"}
	testResult = append(testResult, "隐血："+ yx[BloodcolorfuncAution(img,xaxis,yaxis,width,length)])

	tt:=[]string{"neg.","±","1+(15)mg/dL","2+(40)mg/dL","3+(80)mg/dL","4+(150)mg/dL", "Default"}
	testResult = append(testResult, "酮体："+ tt[rbgDisMaxAution(phMATRIX[7])])
	
	yxsy:=[]string{"neg.","Default","1+","2+","Default","Default","Default"}
	testResult = append(testResult, "亚硝酸盐："+ yxsy[rbgDisMaxAution(phMATRIX[8])])

	bxb:=[]string{"neg.","25 Leu/uL","75 Leu/uL","250 Leu/uL","500 Leu/uL","Default","Default"}
	testResult = append(testResult, "白细胞："+ bxb[rbgDisMaxAution(phMATRIX[9])])
	
	return testResult
}

//单个测试项每个色块的颜色值切片
func colorfuncAution(img *image.Image, xaxis, yaxis []float64) [][]int {
	color := image.NewRGBA((*img).Bounds())
	draw.Draw(color, color.Bounds(), *img, (*img).Bounds().Min, draw.Src)

	ph := make([][]int, 10, 1024)
	for m := 0; m < len(yaxis); {
		arrsum := make([]int, 3)
		for n := 0; n < len(xaxis); {
			// 以下是计算单个目标区域的RGB值
			for i := int(yaxis[m]) - 4; i <= int(yaxis[m])+4; i++ {
				for j := int(xaxis[n]) - 4; j <= int(xaxis[n])+4; j++ {
					r := color.RGBAAt(j, i)
					arr := make([]int, 0)
					arr = append(arr, int(r.R), int(r.G), int(r.B))

					for k := 0; k < len(arr); k++ {
						arrsum[k] += arr[k]
					}
				}
			}
			arrsum[0] = arrsum[0] / 81
			arrsum[1] = arrsum[1] / 81
			arrsum[2] = arrsum[2] / 81

			ph[m] = append(ph[m], arrsum...)
			n++
		}
		m++
	}
	return ph
}

//距离最小者为目标值
func rbgDisMaxAution(phControl []int) int {
	var j int
	DisMin := 1000.0
	slice := phControl[21:24]
	for i := 0; i < len(phControl)-3; {
		Dis := 0.0
		Dis += math.Pow(float64(phControl[i]-slice[0]), 2)
		Dis += math.Pow(float64(phControl[i+1]-slice[1]), 2)
		Dis += math.Pow(float64(phControl[i+2]-slice[2]), 2)
		Dis = math.Sqrt(Dis)

		if Dis < DisMin {
			DisMin = Dis
			j = i
		}

		i += 3
	}
	return j / 3

}

//自动计算采样坐标
func coordinatesConvertAution(img *image.Image, x, y []float64, width, length int) ([]float64, []float64) {
	color := image.NewRGBA((*img).Bounds())
	draw.Draw(color, color.Bounds(), *img, (*img).Bounds().Min, draw.Src)
	a := (*img).Bounds().Max.X
	b := (*img).Bounds().Max.Y

	for i := 0; i < len(x); i++ {
		x[i] = float64(x[i]) * float64(float64(a)/float64(width))
	}
	for i := 0; i < len(y); i++ {

		y[i] = float64(y[i]) * float64(float64(b)/float64(length))

	}

	return x, y
}

//Bloodcolorfunc 隐血：测试项目色块颜色不均匀性判断
func BloodcolorfuncAution(img *image.Image, xaxis, yaxis []float64, width, length int) int {
	var j, colorTag int
	color := image.NewRGBA((*img).Bounds())
	draw.Draw(color, color.Bounds(), *img, (*img).Bounds().Min, draw.Src)

	rMid := color.RGBAAt(int(xaxis[7]), int(yaxis[6]))
	arrMid := make([]int, 0)
	arrMid = append(arrMid, int(rMid.R), int(rMid.G), int(rMid.B))
	//判断测试色块颜色均匀性
	arrsum := make([]int, 3)
	for i := int(yaxis[6]) - 4; i <= int(yaxis[6])+4; i++ {
		for j := int(xaxis[7]) - 4; j <= int(xaxis[7])+4; j++ {
			r := color.RGBAAt(j, i)
			arr := make([]int, 0)
			arr = append(arr, int(r.R), int(r.G), int(r.B))
			if arr[0]-arrMid[0] > 20 {
				colorTag = 1
			}
			for k := 0; k < len(arr); k++ {
				arrsum[k] += arr[k]
			}
		}
	}
	phMATRIX := colorfuncAution(img, xaxis, yaxis) //得到颜色矩阵
	if colorTag == 1 {                             //隐血测试色块不均匀，则与5，6色块对比
		DisMin := 1000.0
		slice := phMATRIX[6][21:24]
		phControl := phMATRIX[6][15:21]
		for i := 0; i < len(phControl)-3; {
			Dis := 0.0
			Dis += math.Pow(float64(phControl[i]-slice[0]), 2)
			Dis += math.Pow(float64(phControl[i+1]-slice[1]), 2)
			Dis += math.Pow(float64(phControl[i+2]-slice[2]), 2)
			Dis = math.Sqrt(Dis)

			if Dis < DisMin {
				DisMin = Dis
				j = i
			}
			i += 3
		}
		//		fmt.Println(DisMin, j/3, colorTag)
		return j / 3
	} else { //隐血测试色块均匀，则与前面4个色块对比
		DisMin := 1000.0
		slice := phMATRIX[6][21:24]
		phControl := phMATRIX[6][0:15]
		for i := 0; i < len(phControl)-3; {
			Dis := 0.0
			Dis += math.Pow(float64(phControl[i]-slice[0]), 2)
			Dis += math.Pow(float64(phControl[i+1]-slice[1]), 2)
			Dis += math.Pow(float64(phControl[i+2]-slice[2]), 2)
			Dis = math.Sqrt(Dis)

			if Dis < DisMin {
				DisMin = Dis
				j = i
			}

			i += 3
		}
		//		fmt.Println(DisMin, j/3, colorTag)
		return j / 3
	}
}
