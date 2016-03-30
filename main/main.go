package main

import (
	// "encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/qshuai162/MivdApi/account"
	"github.com/qshuai162/MivdApi/common/util"
	ana "github.com/qshuai162/MivdApi/imganalyze"
	"github.com/qshuai162/MivdApi/record"
	"gopkg.in/mgo.v2"
	"image"
	"image/draw"
	"io"
	// "io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var mongodbstr string = "121.41.46.25:27017"
var store = sessions.NewCookieStore([]byte("something"))

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.Static("/images", "./images")
	r.Use(CORSMiddleware())

	// r.POST("/uploadfile", func(c *gin.Context) {
	// 	savepath := getSavePath()
	// 	imgbase := c.PostForm("image")
	// 	ddd, _ := base64.StdEncoding.DecodeString(imgbase) //成图片文件并把文件写入到buffer
	// 	ioutil.WriteFile(savepath, ddd, 0666)

	// 	img := ana.DecodeImg(savepath)
	// 	if img == nil {
	// 		fmt.Println("解码失败")
	// 	}
	// 	subimg := ana.ImgCut(img, 0, 0, 19, 128)
	// 	//一条线，找出反应线中灰度值最小的线
	// 	ret := ana.TestLines(&subimg, 15, []int{3})
	// 	if ret[0] != 117.42133537989254 {
	// 		fmt.Println("失败：一条线，找出反应线中灰度值最小的线")
	// 	}

	// 	path, _ := filepath.Rel(basedir, savepath)                   //相对路径
	// 	c.JSON(200, gin.H{"data": ret[0], "path": "webapi/" + path}) // {data:ret[0]}}).String(200, "webapi/"+path)
	// })

	r.POST("/findwb", func(c *gin.Context) {
		file, head, err := c.Request.FormFile("fileToUpload")
		if err != nil {
			return
		}
		defer file.Close()
		basedir := util.GetCurrDir()
		savedir := basedir + "/images/temp/"
		if !util.IsExistFileOrDir(savedir) {
			os.MkdirAll(savedir, 0777) //创建文件夹
		}
		savepath := savedir + head.Filename
		fw, err := os.Create(savepath)
		io.Copy(fw, file)
		fmt.Println(err)
		if err != nil {
			return
		}
		defer fw.Close()
		img := ana.DecodeImg(savepath)
		if img == nil {
			fmt.Println("解码失败")
		}
		b, w := ana.GetWBGrayVal(img)
		path, _ := filepath.Rel(basedir, savepath)                              //相对路径
		c.JSON(200, gin.H{"black": b, "white": w, "path": "../webapi/" + path}) // {data:ret[0]}}).String(200, "webapi/"+path)
	})

	r.POST("/getline", func(c *gin.Context) {
		file, head, err := c.Request.FormFile("fileToUpload1")
		if err != nil {
			return
		}
		defer file.Close()
		basedir := util.GetCurrDir()
		savedir := basedir + "/images/temp/"
		if !util.IsExistFileOrDir(savedir) {
			os.MkdirAll(savedir, 0777) //创建文件夹
		}
		savepath := savedir + head.Filename
		fw, err := os.Create(savepath)
		io.Copy(fw, file)
		fmt.Println(err)
		if err != nil {
			return
		}
		defer fw.Close()
		img := ana.DecodeImg(savepath)
		if img == nil {
			fmt.Println("解码失败")
		}

		gray := image.NewGray((*img).Bounds())
		draw.Draw(gray, gray.Bounds(), *img, (*img).Bounds().Min, draw.Src) //原始图片转换为灰色图片

		arr := ana.ConvertoLine(gray)

		arr2 := make([]string, gray.Rect.Size().Y, gray.Rect.Size().Y)
		maxx := gray.Rect.Size().X
		maxy := gray.Rect.Size().Y
		for y := 0; y < maxy; y++ {
			temp := ""
			for x := 0; x < maxx; x++ {
				temp += strconv.Itoa(int(gray.GrayAt(x, y).Y)) + ","
			}
			arr2[y] = temp
		}

		path, _ := filepath.Rel(basedir, savepath)                                 //相对路径
		c.JSON(200, gin.H{"datas": arr, "all": arr2, "path": "../webapi/" + path}) // {data:ret[0]}}).String(200, "webapi/"+path)
	})

	r.GET("/apitest", func(c *gin.Context) {
		c.JSON(200, gin.H{"datas": 12345})
	})

	r.GET("/record/list/:page", func(c *gin.Context) {
		page := c.Param("page")
		idx, err := strconv.Atoi(page)
		if err != nil {
			c.JSON(200, gin.H{"code": 1, "data": err})
			return
		}
		session, err := mgo.Dial(mongodbstr)
		if err != nil {
			panic(err)
		}
		session.SetMode(mgo.Monotonic, true)
		defer session.Close()
		rs := record.GetList(session, idx, 20)
		c.JSON(200, gin.H{"code": 0, "data": rs})
	})

	r.GET("/record/detail/:id", func(c *gin.Context) {
		id := c.Param("id")
		session, err := mgo.Dial(mongodbstr)
		if err != nil {
			panic(err)
		}
		session.SetMode(mgo.Monotonic, true)
		defer session.Close()
		r := record.FindById(session, id)
		c.JSON(200, gin.H{"code": 0, "data": r})
	})

	r.POST("/account/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		session, err := mgo.Dial(mongodbstr)
		if err != nil {
			panic(err)
		}
		session.SetMode(mgo.Monotonic, true)
		defer session.Close()
		ac := account.Login(session, username, password)
		if ac.UserName != "" {
			c.JSON(200, gin.H{"code": 0, "data": ac})
		} else {
			c.JSON(200, gin.H{"code": 1, "data": "user name or password error"})
		}
	})

	r.Run(":8765")
}

func getSavePath() string {
	basedir := util.GetCurrDir()
	savedir := basedir + "images/temp/"
	if !util.IsExistFileOrDir(savedir) {
		os.MkdirAll(savedir, 0777) //创建文件夹
	}
	savepath := savedir + fmt.Sprintf("%d", time.Now().Unix()) + ".jpg"
	return savepath
}
