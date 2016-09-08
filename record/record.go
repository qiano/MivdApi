package record

import (
	. "github.com/qshuai162/common/config"
	// "fmt"	

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

type Record struct {
	Id        int    //编号
	Vendor    string //供应商
	Type      string //类型
	DateTime  int64  //时间
	PhotoPath string //照片路径
	PlantPath string //板条图片路径
	AreaPath  string //反应区图片路径
	Project   string //测试项目
	Result    string //判定
	LotNo     string //批号
	Index     int    //序号
	Operator  string //操作员
	Location  string //地点
	Longitude float64 //经度
	Latitude  float64 //纬度
	Remark    string //备注
	PatientName string //病人姓名
	PatientNo string //病人编号
	QrCode    string   //二维码
}

var mongodbstr = Config["mongodbHost"]
var dbName = Config["mongodbDbName"]

// var mongodbstr = "121.41.46.25:27017"
// var dbName = "mivd_dev"
var collectionName = "record"

func NewRecord(qrcode,picpath, plantpath, areapath,vendor,ty,project, operator,location string,lat,long float64,lotno string,index int,patientName,patientNo string) *Record {
	re := new(Record)
	re.Type = ty
	re.Vendor=vendor
	re.Project =project
	re.Location = location
	re.Longitude=long
	re.Latitude=lat
	re.Operator = operator
	re.PhotoPath = picpath
	re.PlantPath=plantpath
	re.AreaPath = areapath
	re.LotNo=lotno
	re.Index=index
	re.DateTime = time.Now().Unix()
	re.PatientName=patientName
	re.PatientNo=patientNo
	re.QrCode=qrcode
	return re
}

func (r *Record) Save() {
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	col := session.DB(dbName).C(collectionName)
	var temp Record
	col.Find(bson.M{"type": r.Type}).Sort("-$natural").One(&temp)
	r.Id = temp.Id + 1
	col.Insert(r)
}

func GetList(pageIdx, pageSize int, user, role, ty,pname string) (records []Record) {
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	col := session.DB(dbName).C(collectionName)
	query:=bson.M{"type":ty}
	if role=="user"{
		query["operator"]=user
	}
	if pname!=""{
		reg:=bson.M{"$regex":pname,"$options":"i"}
		query["patientname"]=reg
	}
	col.Find(query).Sort("-$natural").Skip((pageIdx - 1) * pageSize).Limit(pageSize).All(&records)
	
	// if role != "user" {
	// 	col.Find(bson.M{"type": ty}).Sort("-$natural").Skip((pageIdx - 1) * pageSize).Limit(pageSize).All(&records)
	// } else {
	// 	col.Find(bson.M{"operator": user, "type": ty}).Sort("-$natural").Skip((pageIdx - 1) * pageSize).Limit(pageSize).All(&records)
	// }
	return
}

func FindByID(tp,id string) (record Record) {
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	ID, _ := strconv.Atoi(id)
	session.DB(dbName).C(collectionName).Find(bson.M{"id": ID,"type":tp}).One(&record)
	return
}

//Exist 二维码是否已经存在
func Exist(qrcode string) bool{
	session, err := mgo.Dial(mongodbstr)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	 count,err:=session.DB(dbName).C(collectionName).Find(bson.M{"qrcode": qrcode}).Count()
	 if err!=nil{
		 return false
	 }
	return count>1
}

// func FindByOid(oid string) (record Record) {
// 	session, err := mgo.Dial(mongodbstr)
// 	if err != nil {
// 		panic(err)
// 	}
// 	session.SetMode(mgo.Monotonic, true)
// 	defer session.Close()
// 	session.DB(dbName).C(collectionName).FindId(oid)
// 	return
// }
