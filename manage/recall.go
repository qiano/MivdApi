package manage

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/qshuai162/common/config"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Recall recall model
type Recall struct {
	Vendor   string
	Lotno    string
	DateTime int64 //时间
}

var mongodbstr = config.Config["mongodbHost"]
var dbName = config.Config["mongodbDbName"]

// var mongodbstr = "121.41.46.25:27017"
// var dbName = "mivd_dev"
var collectionName = "recall"

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

//AddRecall add recall setting
func AddRecall(qrcode string) (bool, string) {
	if len(qrcode)<14{
		return false,"that's not a qrcoce"
	}
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	vendor := substr(qrcode, 2, 4)
	lotno := substr(qrcode, 6, 8)
	if !Exist(vendor, lotno) {
		col := session.DB(dbName).C(collectionName)
		temp := &Recall{Vendor: vendor, Lotno: lotno, DateTime: time.Now().Unix()}
		e := col.Insert(temp)
		if e == nil {
			return true, ""
		}
		return false, "err"
	}
	return false, "The recall information have been added!"
}

//Exist exist?
func Exist(vendor, lotno string) bool {
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	count, err := session.DB(dbName).C(collectionName).Find(bson.M{"vendor": vendor, "lotno": lotno}).Count()
	if err != nil {
		return false
	}
	return count >= 1
}

//Query query
func Query(pageIdx, pageSize int, vendor, lotno, startdate, enddate, sort string) (total int, recalls []Recall) {
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	col := session.DB(dbName).C(collectionName)
	query := bson.M{}
	if vendor != "" {
		query["vendor"] = vendor
	}
	if lotno != "" {
		query["lotno"] = lotno
	}

	if startdate != "" || enddate != "" {
		start := (int64)(0)
		end := (int64)(0)
		if startdate != "" {
			sStrs := strings.Split(startdate, "-")
			y, _ := strconv.Atoi(sStrs[0])
			m, _ := strconv.Atoi(sStrs[1])
			d, _ := strconv.Atoi(sStrs[2])
			start = time.Date(y, (time.Month)(m), d-1, 0, 0, 0, 0, time.Local).Unix()
		}
		if enddate != "" {
			eStrs := strings.Split(enddate, "-")
			y, _ := strconv.Atoi(eStrs[0])
			m, _ := strconv.Atoi(eStrs[1])
			d, _ := strconv.Atoi(eStrs[2])
			end = time.Date(y, (time.Month)(m), d-1, 23, 59, 59, 999, time.Local).Unix()
		}
		fmt.Println(start, end)
		if start > 0 && end > 0 {
			var c = []bson.M{bson.M{"datetime": bson.M{"$gte": start}}, bson.M{"datetime": bson.M{"$lte": end}}}
			query["$and"] = c
		} else {
			if start > 0 {
				fmt.Println(start)
				query["datetime"] = bson.M{"$gte": start}
			}
			if end > 0 {
				fmt.Println(end)
				query["datetime"] = bson.M{"$lte": end}
			}
		}
	}

	q := col.Find(query)
	if sort != "" {
		q = q.Sort(sort)
	}
	total, _ = q.Count()
	q.Skip((pageIdx - 1) * pageSize).Limit(pageSize).All(&recalls)
	return
}

//Delete delete
func Delete(vendor, lotno string) {
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	col := session.DB(dbName).C(collectionName)
	col.Remove(bson.M{"vendor": vendor, "lotno": lotno})
}
