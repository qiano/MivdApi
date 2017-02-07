package track

import (
	"strconv"
	"time"

	"github.com/qshuai162/MivdApi_Trail/qrcode"
	. "github.com/qshuai162/common/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TrackEntity struct {
	QrCode  string //二维码
	Factory string //厂商
	LotNo   string //批号
	Index   int    //序号
	Project string //测试项目
}

//TrackRecord 追踪记录
type TrackRecord struct {
	QrCode    string
	DateTime  int64
	Operator  string  //操作员
	TrueName  string  //操作员真实姓名
	Location  string  //地点
	Longitude float64 //经度
	Latitude  float64 //纬度
}

var mongodbstr = Config["mongodbHost"]
var dbName = Config["mongodbDbName"]

// var mongodbstr = "127.0.0.1:27017"
// var dbName = "mivd_dev"
var collectionName = "trackentity"

//AddTrackRecord 新增追踪记录
func AddTrackRecord(qr, operator, truename, location string, lat, long float64) bool {
	if !Exist(qr) {
		addTrackEntity(qr)
	}
	return addTrack(qr, operator, truename, location, lat, long)
}

func addTrack(qrcode, operator, truename, location string, lat, long float64) bool {
	re := new(TrackRecord)
	re.QrCode = qrcode
	re.Location = location
	re.Longitude = long
	re.Latitude = lat
	re.Operator = operator
	re.TrueName = truename
	re.DateTime = time.Now().Unix()
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	col := session.DB(dbName).C("trackrecord")
	col.Insert(re)
	return true
}

func addTrackEntity(qr string) *TrackEntity {
	re := new(TrackEntity)
	re.QrCode = qr
	re.Factory = qrcode.GetFactory(qr)
	re.LotNo = qrcode.GetLotNo(qr)
	re.Index, _ = strconv.Atoi(qrcode.GetIndex(qr))
	re.Project = qrcode.GetProjectNames(qr)
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	col := session.DB(dbName).C(collectionName)
	col.Insert(re)
	return re
}

//GetList list
func GetList(pageIdx, pageSize int, user, role string) (records []TrackEntity) {
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	col := session.DB(dbName).C(collectionName)
	query := bson.M{}
	
	col.Find(query).Sort("-$natural").Skip((pageIdx - 1) * pageSize).Limit(pageSize).All(&records)
	return
}

//FindTrackEntity detail
func FindTrackEntity(qr string) (record TrackEntity) {
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	session.DB(dbName).C(collectionName).Find(bson.M{"qrcode": qr}).One(&record)
	return
}

//GetTrackRecordByQrCode 追踪记录
func GetTrackRecordByQrCode(qr string) (records []TrackRecord) {
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	col := session.DB(dbName).C("trackrecord")
	col.Find(bson.M{"qrcode": qr}).Sort("-$natural").All(&records)
	return
}

//Exist 二维码是否已经存在
func Exist(qr string) bool {
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	count, err := session.DB(dbName).C(collectionName).Find(bson.M{"qrcode": qr}).Count()
	if err != nil {
		return false
	}
	return count >= 1
}

//Query website query
// func Query(pageIdx, pageSize, id int, pname, pno, test, factory, lotno, result, location, operator, startdate, enddate, sort string) (total int, records []Record) {
// 	session, err := mgo.Dial(mongodbstr)
// 	if err != nil {
// 		panic(err)
// 	}
// 	session.SetMode(mgo.Monotonic, true)
// 	defer session.Close()

// 	col := session.DB(dbName).C(collectionName)
// 	query := bson.M{}
// 	// query := bson.M{"type": ty}
// 	// if role == "user" {
// 	// 	query["operator"] = user
// 	// }
// 	if id != 0 {
// 		query["id"] = id

// 	}
// 	if pname != "" {
// 		reg := bson.M{"$regex": pname, "$options": "i"}
// 		query["patientname"] = reg
// 	}
// 	if pno != "" {
// 		reg := bson.M{"$regex": pno, "$options": "i"}
// 		query["patientno"] = reg
// 	}
// 	if test != "" {
// 		reg := bson.M{"$regex": test, "$options": "i"}
// 		query["project"] = reg
// 	}
// 	if factory != "" {
// 		reg := bson.M{"$regex": factory, "$options": "i"}
// 		query["vendor"] = reg
// 	}
// 	if lotno != "" {
// 		reg := bson.M{"$regex": lotno, "$options": "i"}
// 		query["lotno"] = reg
// 	}
// 	if result != "" {
// 		reg := bson.M{"$regex": result, "$options": "i"}
// 		query["result"] = reg
// 	}
// 	if location != "" {
// 		reg := bson.M{"$regex": location, "$options": "i"}
// 		query["location"] = reg
// 	}
// 	if operator != "" {
// 		reg := bson.M{"$regex": operator, "$options": "i"}
// 		query["operator"] = reg
// 	}

// 	if startdate != "" || enddate != "" {
// 		start := (int64)(0)
// 		end := (int64)(0)
// 		if startdate != "" {
// 			sStrs := strings.Split(startdate, "-")
// 			y, _ := strconv.Atoi(sStrs[0])
// 			m, _ := strconv.Atoi(sStrs[1])
// 			d, _ := strconv.Atoi(sStrs[2])
// 			start = time.Date(y, (time.Month)(m), d-1, 0, 0, 0, 0, time.Local).Unix()
// 		}
// 		if enddate != "" {
// 			eStrs := strings.Split(enddate, "-")
// 			y, _ := strconv.Atoi(eStrs[0])
// 			m, _ := strconv.Atoi(eStrs[1])
// 			d, _ := strconv.Atoi(eStrs[2])
// 			end = time.Date(y, (time.Month)(m), d-1, 23, 59, 59, 999, time.Local).Unix()
// 		}
// 		fmt.Println(start, end)
// 		if start > 0 && end > 0 {
// 			var c = []bson.M{bson.M{"datetime": bson.M{"$gte": start}}, bson.M{"datetime": bson.M{"$lte": end}}}
// 			query["$and"] = c
// 		} else {
// 			if start > 0 {
// 				fmt.Println(start)
// 				query["datetime"] = bson.M{"$gte": start}
// 			}
// 			if end > 0 {
// 				fmt.Println(end)
// 				query["datetime"] = bson.M{"$lte": end}
// 			}
// 		}

// 		// var strs = strings.Split(date, "-")
// 		// fmt.Println(strs)
// 		// y, _ := strconv.Atoi(strs[0])
// 		// m, _ := strconv.Atoi(strs[1])
// 		// d, _ := strconv.Atoi(strs[2])
// 		// start := time.Date(y, (time.Month)(m), d, 0, 0, 0, 0, time.Local).Unix()
// 		// end := time.Date(y, (time.Month)(m), d, 23, 59, 59, 999, time.Local).Unix()
// 		// var c = []bson.M{bson.M{"datetime": bson.M{"$gte": start}}, bson.M{"datetime": bson.M{"$lte": end}}}
// 		// query["$and"] = c
// 	}

// 	q := col.Find(query)
// 	if sort != "" {
// 		q = q.Sort(sort)
// 	}
// 	total, _ = q.Count()
// 	q.Skip((pageIdx - 1) * pageSize).Limit(pageSize).All(&records)

// 	// if role != "user" {
// 	// 	col.Find(bson.M{"type": ty}).Sort("-$natural").Skip((pageIdx - 1) * pageSize).Limit(pageSize).All(&records)
// 	// } else {
// 	// 	col.Find(bson.M{"operator": user, "type": ty}).Sort("-$natural").Skip((pageIdx - 1) * pageSize).Limit(pageSize).All(&records)
// 	// }
// 	return
// }
