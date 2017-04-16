package qualitative

import (
	"github.com/qshuai162/MivdApi_Trail/imganalyze"
    "math"
    "image"
    "image/draw"
    //  "fmt"
)

func Test(picpath string,hotWidth,cdis,tdis int) (bool,float64){
    img := imganalyze.DecodeImg(picpath)
    gray := image.NewGray((*img).Bounds())
	draw.Draw(gray, gray.Bounds(), *img, (*img).Bounds().Min, draw.Src) //原始图片转换为灰色图片
    

	arr := imganalyze.ConverParttoLine(gray,0.5)
    unit:=float64((*img).Bounds().Max.X)/float64(hotWidth)
    offset:=unit*1.2
    
    cmin:=int(math.Floor(float64(cdis)*unit-offset))
    cmax:=int(math.Ceil(float64(cdis)*unit+offset))
    carr:=(*arr)[cmin:cmax]
	cval, _ := imganalyze.FindMinValInLines(&carr) //寻找最低点，判定为C线
    
    //在C线固定距离范围内寻找T线，最高点
    tmin:=int(math.Floor(float64(tdis)*unit-offset))
    tmax:=int(math.Ceil(float64(tdis)*unit+offset))
    tarr:=(*arr)[tmin:tmax]
    tval,_:=imganalyze.FindMinValInLines(&tarr)
    
    //计算除C T线区域外的灰度平均值
    avg:=imganalyze.CalcRangeAvg(arr,cmax,tmin)
    //c线是否存在
    if avg-cval<10{
        return false,-1
    }
    black,_,_:=imganalyze.FindMinGray(gray)
    white,_,_:=imganalyze.FindMaxGray(gray)
    tl:=(float64((white-black))/255.0)*tval+float64(black)
    return true,tl;
}

