package autionsticks

import(
	"github.com/qshuai162/MivdApi/imganalyze"
    "github.com/qshuai162/MivdApi/record"
    "strconv"
    "time"
    "strings"
)

func Do(picpath,areapath,user string ) *record.Record{
    img:=imganalyze.DecodeImg(picpath)
    result := Autionsticks(img)
    re := new(record.Record)
    re.Type = "autionsticks"
    re.Project ="autionsticks"
    re.ResultMsg =strings.Join(result,",")
    re.Location = "shanghai"
    re.UserName = user
    re.PicData = picpath
    re.AreaData = areapath
    re.LotNo = strconv.Itoa(time.Now().Year()) + time.Now().Month().String()
    re.Date = time.Now().Unix()
    return re
}