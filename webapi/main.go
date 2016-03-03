package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"image"
	"image/draw"
	"io"
	"os"
	"path/filepath"
	"qian/util"
	ana "qian/vitrodiag/imganalyze"
	"strconv"
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

	r.POST("/uploadfile", func(c *gin.Context) {
		file, head, err := c.Request.FormFile("fileToUpload")
		if err != nil {
			return
		}
		defer file.Close()
		basedir := util.GetCurrDir()
		savedir := basedir + "/images/temp/"
		fmt.Println(savedir)
		if !util.IsExistFileOrDir(savedir) {
			os.MkdirAll(savedir, 0777) //创建文件夹
		}
		savepath := savedir + head.Filename
		fmt.Println(savepath)
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
		subimg := ana.ImgCut(img, 0, 0, 19, 128)
		//一条线，找出反应线中灰度值最小的线
		ret := ana.TestLines(&subimg, 15, []int{3})
		if ret[0] != 117.42133537989254 {
			fmt.Println("失败：一条线，找出反应线中灰度值最小的线")
		}

		path, _ := filepath.Rel(basedir, savepath)                   //相对路径
		c.JSON(200, gin.H{"data": ret[0], "path": "webapi/" + path}) // {data:ret[0]}}).String(200, "webapi/"+path)
	})

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

	r.Run(":8888")

}
