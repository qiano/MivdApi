package main

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"io"
	"io/ioutil"
	"image/jpeg"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/qshuai162/account"
	. "github.com/qshuai162/common/config"
	"github.com/qshuai162/common/util"
	ana "github.com/qshuai162/MivdApi/imganalyze"
	"github.com/qshuai162/MivdApi/ivd/ph"
	"github.com/qshuai162/MivdApi/ivd/autionsticks"
	"github.com/qshuai162/MivdApi/record"

	// "io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"math"
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
		x,_:=strconv.ParseFloat(c.PostForm("x"),64)
		y,_:=strconv.ParseFloat(c.PostForm("y"),64)
		w,_:=strconv.ParseFloat(c.PostForm("w"),64)
		h,_:=strconv.ParseFloat(c.PostForm("h"),64)
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
		if err != nil {
			return
		}
		defer fw.Close()
		img := ana.DecodeImg(savepath)
		if img == nil {
			fmt.Println("解码失败")
		}
		mx:=float64((*img).Bounds().Max.X)
		my:=float64((*img).Bounds().Max.Y)
		subimg := ana.ImgCut(img, int(math.Ceil(mx*x)), int(math.Ceil(my*y)), int(math.Floor(mx*(x+w))),int(math.Floor(my*(y+h))))
			
		gray := image.NewGray(subimg.Bounds())
		draw.Draw(gray, gray.Bounds(), subimg, subimg.Bounds().Min, draw.Src) //原始图片转换为灰色图片
		fmt.Println(gray)
		arr := ana.ConvertoLine(gray)
		arr2 := make([]string, gray.Rect.Size().Y, gray.Rect.Size().Y)
		maxx := gray.Rect.Size().X
		maxy := gray.Rect.Size().Y
		fmt.Println(maxx,maxy)
		for y := 0; y < maxy; y++ {
			temp := ""
			for x := 0; x < maxx; x++ {
				temp += strconv.Itoa(int(gray.Pix[maxx*y+x])) + ","
			}
			arr2[y] = temp
		}

		path, _ := filepath.Rel(basedir, savepath)                                 //相对路径
		c.JSON(200, gin.H{"datas": arr, "all": arr2,"sub":ana.ImgToBase64(&subimg), "path": "../webapi/" + path}) // {data:ret[0]}}).String(200, "webapi/"+path)
	})

	r.GET("/record/list/:type/:page", func(c *gin.Context) {
		page := c.Param("page")
		user,_ := c.GetQuery("user")
		role,_ := c.GetQuery("role")
		ty:=c.Param("type")
		fmt.Println(page,user,role,ty)
		idx, err := strconv.Atoi(page)
		if err != nil {
			c.JSON(200, gin.H{"code": 1, "data": err})
			return
		}

		rs := record.GetList(idx, 20, user, role,ty)
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
		Type := c.Param("type")
		switch Type {
		case "ph":
			picpath:=getSavePath("ph")
			ddd, _ := base64.StdEncoding.DecodeString(imgbase) //成图片文件并把文件写入到buffer
			ioutil.WriteFile(picpath, ddd, 0666) 
		
			img := ana.Base64ToImg(imgbase)
    		area := ana.ImgCut(img, 0, (*img).Bounds().Max.Y*8/17.0, (*img).Bounds().Max.X, (*img).Bounds().Max.Y*9/17.0)
			areapath:=getSavePath("ph")
			f, _ := os.Create(areapath)     //创建文件
			defer f.Close()                   //关闭文件
			jpeg.Encode(f, area, nil)
			re:=ph.Do(picpath,areapath,user)
			re.Save()
			c.JSON(200, gin.H{"code": 0, "data": re.Id, "msg": ""})
			break
		case "autionsticks":
			picpath:=getSavePath("autionsticks")
			ddd, _ := base64.StdEncoding.DecodeString(imgbase) //成图片文件并把文件写入到buffer
			ioutil.WriteFile(picpath, ddd, 0666) 
		
			img := ana.Base64ToImg(imgbase)
    		area := ana.ImgCut(img, (*img).Bounds().Max.X/8 , 0, (*img).Bounds().Max.X*2/8, (*img).Bounds().Max.Y)
			areapath:=getSavePath("autionsticks")
			f, _ := os.Create(areapath)     //创建文件
			defer f.Close()                   //关闭文件
			jpeg.Encode(f, area, nil)
		
			re:=autionsticks.Do(picpath,areapath,user)
			re.Save()
			c.JSON(200, gin.H{"code": 0, "data": re.Id, "msg": ""})
			break	
		default:
			c.JSON(200, gin.H{"code": 1, "data":nil, "msg": "not surpport"})
		}
	})

	r.Run(":" + Config["port"])
}

func getSavePath(flag string) string {
	basedir := util.GetCurrDir()
	savedir := basedir + "images/"+flag+"/"
	if !util.IsExistFileOrDir(savedir) {
		os.MkdirAll(savedir, 0777) //创建文件夹
	}
	savepath := savedir + fmt.Sprintf("%d", time.Now().UnixNano()) + ".jpg"
	return savepath
}
