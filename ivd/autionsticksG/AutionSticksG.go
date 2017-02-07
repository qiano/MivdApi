package autionsticksG
import (
	 "fmt"
	"image"
	"image/draw"
	"math"
	 "github.com/qshuai162/MivdApi_Trail/imganalyze"
)

func AutionsticksG(picpath string) []string {
	fmt.Println(picpath)
	img := imganalyze.DecodeImg(picpath)
	fmt.Println(img)
	width := 55  //截图区域大小，毫米坐标系  X坐标
	length := 91 // Y坐标
	// 最后是测试项坐标                                                      //AutionSticks  毫米坐标系，每个色块最中心点的坐标
	x := []float64{4, 11, 18, 24.5, 32, 38.5, 45.5, 52.5}                    //
	y := []float64{3.5, 11, 18, 25.5, 32.5, 39.5, 48, 54.5, 61.25, 68.5, 76} //纵坐标一样
	TestResult := test(img, x, y, width, length)
	return TestResult
}

func test(img *image.Image, x, y []float64, width, length int) []string {
	xaxis, yaxis := coordinatesConvertAution(img, x, y, width, length)
	phMATRIX := colorfuncAution(img, xaxis, yaxis) //得到颜色矩
	fmt.Println(phMATRIX)
	var testResult = make([]string, 0)
		switch rbgDisMaxAution(phMATRIX[0]) { //胆红素
	case 0:
		// fmt.Println("胆红素：neg.")
		testResult = append(testResult, "胆红素：neg.")
	case 1:
		// fmt.Println("胆红素：1+（0.5）mg/dL")
		testResult = append(testResult, "胆红素：+")
	case 2:
		// fmt.Println("胆红素：2+（2）mg/dL")
		testResult = append(testResult, "胆红素：++")
	case 3:
		// fmt.Println("胆红素：3+（6）mg/dL")
		testResult = append(testResult, "胆红素：+++")
	default:
		// fmt.Println("胆红素：Default")
		testResult = append(testResult, "胆红素：Default")
	}

	ndy := []string{"normal", "2 (35)mg/dL(umol/l)", "4 (70)mg/dL(umol/l)", "8 (140)mg/dL(umol/l)", "12(200)mg/dL(umol/l)", "Default", "Default", "Default"}
	testResult = append(testResult, "尿胆原："+ndy[rbgDisMaxAution(phMATRIX[1])])

	switch rbgDisMaxAution(phMATRIX[2]) { // 酮体
	case 0:
		// fmt.Println("酮体：neg.")
		testResult = append(testResult, "酮体：neg.")
	case 1:
		// fmt.Println("酮体：正负")
		testResult = append(testResult, "酮体：trace")
	case 2:
		// fmt.Println("酮体：1+(15)mg/dL")
		testResult = append(testResult, "酮体：+")
	case 3:
		// fmt.Println("酮体：2+(40)mg/dL")
		testResult = append(testResult, "酮体：++")
	case 4:
		// fmt.Println("酮体：3+(80)mg/dL")
		testResult = append(testResult, "酮体：+++")
	default:
		// fmt.Println("酮体：Default")
		testResult = append(testResult, "酮体：Default")
	}

	switch rbgDisMaxAution(phMATRIX[3]) { // 维生素C
	case 0:
		// fmt.Println("维生素C：10 mg/dL")
		testResult = append(testResult, "维生素C：neg.")
	case 1:
		// fmt.Println("维生素C：50 mg/dL")
		testResult = append(testResult, "维生素C：+")
	case 2:
		// fmt.Println("维生素C：100 mg/dL")
		testResult = append(testResult, "维生素C：++")
	default:
		// fmt.Println("维生素C：Default")
		testResult = append(testResult, "维生素C：Default")
	}

	ptt := []string{"norm.", "50(2,8)mg/dL(nmol/l)", "100(5,6)mg/dL(nmol/l)", "250(14)mg/dL(nmol/l)", "500(28)mg/dL(nmol/l)", ">=1000(56)mg/dL(nmol/l)", "Default", "Default"}
	testResult = append(testResult, "葡萄糖："+ptt[rbgDisMaxAution(phMATRIX[4])])

	dbz := []string{"neg.", "trace", "30mg/dL", "100mg/dL", "500mg/dL", "Default", "Default", "Default"}
	fmt.Println(rbgDisMaxAution(phMATRIX[5]))
	testResult = append(testResult, "蛋白质："+dbz[rbgDisMaxAution(phMATRIX[5])])

	switch BloodcolorfuncAution(img,xaxis, yaxis, width, length) { //红细胞
	case 0:
		// fmt.Println("红细胞：neg.")
		testResult = append(testResult, "红细胞：neg.")
	case 1:
		// fmt.Println("红细胞：1+（0.06）mg/dL")
		testResult = append(testResult, "红细胞：+ca.5-10 Ery/ul")
	case 2:
		// fmt.Println("红细胞：2+（0.2）mg/dL")
		testResult = append(testResult, "红细胞：++ca.50 Ery/ul")
	case 3:
		// fmt.Println("红细胞：3+（1.0）mg/dL")
		testResult = append(testResult, "红细胞：+++ca.300 Ery/ul")
	case 4:
		// fmt.Println("红细胞：Non Hemolysis 1+ mg/dL")
		testResult = append(testResult, "红细胞：ca.5~10 Ery/ul")
	case 5:
		// fmt.Println("红细胞：2+ mg/dL")
		testResult = append(testResult, "红细胞：ca.50 Ery/ul")
	case 6:
		// fmt.Println("红细胞：2+ mg/dL")
		testResult = append(testResult, "红细胞：ca.300 Ery/ul")
	}

	switch rbgDisMaxAution(phMATRIX[7]) { // PH
	case 0:
		// fmt.Println("PH ：5")
		testResult = append(testResult, "PH ：5")
	case 1:
		// fmt.Println("PH ：6")
		testResult = append(testResult, "PH ：6")
	case 2:
		// fmt.Println("PH ：6")
		testResult = append(testResult, "PH ：6.5")
	case 3:
		// fmt.Println("PH ：7")
		testResult = append(testResult, "PH ：7")
	case 4:
		// fmt.Println("PH ：8")
		testResult = append(testResult, "PH ：8")
	case 5:
		// fmt.Println("PH ：9")
		testResult = append(testResult, "PH ：9")
	default:
		// fmt.Println("PH ：Default")
		testResult = append(testResult, "PH ：Default")
	}

	switch rbgDisMaxAution(phMATRIX[8]) { //亚硝酸盐
	case 0:
		// fmt.Println("亚硝酸盐：neg.")
		testResult = append(testResult, "亚硝酸盐：neg.")
	case 1:
		// fmt.Println("亚硝酸盐：1+")
		testResult = append(testResult, "亚硝酸盐：pink")
	case 2:
		// fmt.Println("亚硝酸盐：2+")
		testResult = append(testResult, "亚硝酸盐：rose")
	default:
		// fmt.Println("亚硝酸盐：Default")
		testResult = append(testResult, "亚硝酸盐：Default")
	}

	switch rbgDisMaxAution(phMATRIX[9]) { //白细胞
	case 0:
		// fmt.Println("白细胞：neg.")
		testResult = append(testResult, "白细胞：neg.")
	case 1:
		// fmt.Println("白细胞：25 Leu/uL")
		testResult = append(testResult, "白细胞：ca.25 Leuko/uL")
	case 2:
		// fmt.Println("白细胞：75 Leu/uL")
		testResult = append(testResult, "白细胞：ca.75 Leuko/uL")
	case 3:
		// fmt.Println("白细胞：500 Leu/uL")
		testResult = append(testResult, "白细胞：ca.500 Leuko/uL")
	default:
		// fmt.Println("白细胞：Default")
		testResult = append(testResult, "白细胞：Default")
	}

	gravity := []string{"1.000", "1.005", "1.010", "1.015", "1.020", "1.025", "1.030", "Default"}
	testResult = append(testResult, "比重："+gravity[rbgDisMaxAution(phMATRIX[10])])
	return testResult
}

//单个测试项每个色块的颜色值切片
func colorfuncAution(img *image.Image, xaxis, yaxis []float64) [][]int {
	color := image.NewRGBA((*img).Bounds())
	draw.Draw(color, color.Bounds(), *img, (*img).Bounds().Min, draw.Src)

	ph := make([][]int, 11, 1024)
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
