package main

import (
	// "encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/qshuai162/MivdApi/account"
	. "github.com/qshuai162/MivdApi/common/config"
	"github.com/qshuai162/MivdApi/common/util"
	ana "github.com/qshuai162/MivdApi/imganalyze"
	"github.com/qshuai162/MivdApi/record"
	"image"
	"image/draw"
	"io"
	// "io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

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

	r.GET("/record/list/:page", func(c *gin.Context) {
		page := c.Param("page")
		user := c.Param("user")
		role := c.Param("role")

		idx, err := strconv.Atoi(page)
		if err != nil {
			c.JSON(200, gin.H{"code": 1, "data": err})
			return
		}

		rs := record.GetList(idx, 20, user, role)
		for i := 0; i < len(rs); i++ {
			rs[i].PicData = ""
		}
		c.JSON(200, gin.H{"code": 0, "data": rs})
	})

	r.GET("/record/detail/:id", func(c *gin.Context) {
		id := c.Param("id")
		r := record.FindById(id)
		c.JSON(200, gin.H{"code": 0, "data": r})
	})

	r.POST("/account/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		ac := account.Login(username, password)
		if ac.UserName != "" {
			c.JSON(200, gin.H{"code": 0, "data": ac})
		} else {
			c.JSON(200, gin.H{"code": 1, "data": "user name or password error"})
		}
	})

	r.POST("/upload/:type", func(c *gin.Context) {
		imgbase := c.PostForm("image")
		user := c.PostForm("user")

		img := ana.Base64ToImg(imgbase)
		area := ana.ImgCut(img, 0, (*img).Bounds().Max.Y*8/17.0, (*img).Bounds().Max.X, (*img).Bounds().Max.Y*9/17.0)
		ph := ana.TestPH(img)
		re := new(record.Record)
		re.Type = "ph"
		re.Project = "PH"
		re.ResultMsg = strconv.Itoa(ph)
		re.Location = "shanghai"
		re.UserName = user
		re.PicData = imgbase
		re.AreaData = ana.ImgToBase64(&area)
		re.LotNo = strconv.Itoa(time.Now().Year()) + time.Now().Month().String()
		re.Date = time.Now().Unix()

		re.Save()

		c.JSON(200, gin.H{"code": 0, "data": re.Id, "msg": ""})
	})

	r.Run(":" + Config["port"])
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
