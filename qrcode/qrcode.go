package qrcode

import "strings"

var projectMaps = map[string]string{"HIV": "HIV", "MLR": "Malaria", "A1": "Combi", "EgP": "EgP", "HBs": "HBs"}
var factoryMaps = map[string]string{"0001": "Factory1", "0002": "Factory2", "0004": "Factory4"}

func substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

//GetFactory 厂家名
func GetFactory(qrcode string) string {
	code := substr(qrcode, 2, 4)
	if v, ok := factoryMaps[code]; ok {
		return v
	}
	return code
}

//GetLotNo 批号
func GetLotNo(qrcode string) string {
	return substr(qrcode, 6, 8)
}

//GetIndex 序列号
func GetIndex(qrcode string) string {

	return substr(qrcode, 14, 6)
}

//GetProjectNames 项目名称描述
func GetProjectNames(qrcode string) string {
	pstr := strings.Split(qrcode, "|")
	pns := make([]string, 0, 0)
	for i := 1; i < len(pstr); i++ {
		code := substr(pstr[i], 0, 3)
		if v, ok := projectMaps[code]; ok {
			pns = append(pns, v)
		} else {
			pns = append(pns, code)
		}
	}
	return strings.Join(pns, ",")
}
