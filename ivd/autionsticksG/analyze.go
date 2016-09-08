package autionsticksG

import (
	"github.com/qshuai162/MivdApi/imganalyze"
	"github.com/qshuai162/MivdApi/record"
	"strconv"
	"strings"
	"time"
)
//Test 德国尿十一项
func Test(picpath, plantpath, areapath,vendor, user,location string,lat,long float64) *record.Record {
	img := imganalyze.DecodeImg(picpath)
	result := AutionsticksG(img)
	re := new(record.Record)
	re.Type = "autionsticksg"
	re.Vendor=vendor
	re.Project = "autionsticksg"
	re.Result = strings.Join(result, ",")
	re.Location = location
	re.Longitude=long
	re.Latitude=lat
	re.Operator = user
	re.PhotoPath = picpath
	re.PlantPath=plantpath
	re.AreaPath = areapath
	re.LotNo = strconv.Itoa(time.Now().Year()) + time.Now().Month().String()
	re.DateTime = time.Now().Unix()
	return re
}
