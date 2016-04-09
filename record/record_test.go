package record

import (
	"strconv"
	"testing"
	"time"
)

func Test_Add(t *testing.T) {
	re := new(Record)
	re.Type = "ph"
	re.Project = "PH"
	re.ResultMsg = strconv.Itoa(5)

	// re.Operator = "tester"
	re.PicData = ""
	re.AreaData = ""
	re.LotNo = "0"
	re.Date = time.Now().Unix()
	re.Save()
}
