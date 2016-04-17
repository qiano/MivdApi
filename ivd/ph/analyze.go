package ph


import(
    
	"github.com/qshuai162/MivdApi/imganalyze"
    "github.com/qshuai162/MivdApi/record"
    "strconv"
    "time"
)



func Do(picpath,areapath,user string ) *record.Record{
    img:=imganalyze.DecodeImg(picpath)
    ph := TestPH(img)
    re := new(record.Record)
    re.Type = "ph"
    re.Project ="ph"
    re.ResultMsg = strconv.Itoa(ph)
    re.Location = "shanghai"
    re.UserName = user
    re.PicData = picpath
    re.AreaData = areapath
    re.LotNo = strconv.Itoa(time.Now().Year()) + time.Now().Month().String()
    re.Date = time.Now().Unix()
    return re
}