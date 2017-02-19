package main

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	ana "github.com/qshuai162/MivdApi_Trail/imganalyze"
	"github.com/qshuai162/MivdApi_Trail/ivd/autionsticksG"
	"github.com/qshuai162/MivdApi_Trail/ivd/ph"
	"github.com/qshuai162/MivdApi_Trail/ivd/qualitative"
	"github.com/qshuai162/MivdApi_Trail/manage"
	"github.com/qshuai162/MivdApi_Trail/record"
	"github.com/qshuai162/MivdApi_Trail/track"
	"github.com/qshuai162/account"
	"github.com/qshuai162/common/config"
	"github.com/qshuai162/common/util"

	// "io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var store = sessions.NewCookieStore([]byte("something"))

//CORSMiddleware 跨域
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
		x, _ := strconv.ParseFloat(c.PostForm("x"), 64)
		y, _ := strconv.ParseFloat(c.PostForm("y"), 64)
		w, _ := strconv.ParseFloat(c.PostForm("w"), 64)
		h, _ := strconv.ParseFloat(c.PostForm("h"), 64)
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
		mx := float64((*img).Bounds().Max.X)
		my := float64((*img).Bounds().Max.Y)
		subimg := ana.ImgCut(img, int(math.Ceil(mx*x)), int(math.Ceil(my*y)), int(math.Floor(mx*(x+w))), int(math.Floor(my*(y+h))))

		gray := image.NewGray(subimg.Bounds())
		draw.Draw(gray, gray.Bounds(), subimg, subimg.Bounds().Min, draw.Src) //原始图片转换为灰色图片
		fmt.Println(gray)
		arr := ana.ConvertoLine(gray)
		arr2 := make([]string, gray.Rect.Size().Y, gray.Rect.Size().Y)
		maxx := gray.Rect.Size().X
		maxy := gray.Rect.Size().Y
		fmt.Println(maxx, maxy)
		for y := 0; y < maxy; y++ {
			temp := ""
			for x := 0; x < maxx; x++ {
				temp += strconv.Itoa(int(gray.Pix[maxx*y+x])) + ","
			}
			arr2[y] = temp
		}

		path, _ := filepath.Rel(basedir, savepath)                                                                  //相对路径
		c.JSON(200, gin.H{"datas": arr, "all": arr2, "sub": ana.ImgToBase64(&subimg), "path": "../webapi/" + path}) // {data:ret[0]}}).String(200, "webapi/"+path)
	})

	r.GET("/record/list/:type/:page", func(c *gin.Context) {
		page := c.Param("page")
		ty := c.Param("type")
		user, _ := c.GetQuery("user")
		role, _ := c.GetQuery("role")
		pname, _ := c.GetQuery("pname")
		test, _ := c.GetQuery("test")
		date, _ := c.GetQuery("date")
		// fmt.Println(page, user, role, ty)
		idx, err := strconv.Atoi(page)
		if err != nil {
			c.JSON(200, gin.H{"code": 1, "data": err})
			return
		}

		rs := record.GetList(idx, 20, user, role, ty, pname, test, date)
		c.JSON(200, gin.H{"code": 0, "data": rs})
	})

	r.GET("/record/detail/:type/:id", func(c *gin.Context) {
		id := c.Param("id")
		tp := c.Param("type")
		re := record.FindByID(tp, id)
		c.JSON(200, gin.H{"code": 0, "data": re})
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

	r.POST("/account/register", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		truename := c.PostForm("truename")
		organization := c.PostForm("organization")
		code, msg := account.Add(username, password, truename, organization, "user")
		c.JSON(200, gin.H{"code": code, "msg": msg})
	})

	r.POST("/upload/:type", func(c *gin.Context) {
		remark := c.PostForm("remark")
		imgbase := c.PostForm("image")
		operator := c.PostForm("operator")
		truename := c.PostForm("truename")
		vendor := c.PostForm("vendor")
		project := c.PostForm("project")
		location := c.PostForm("location")
		lat, _ := strconv.ParseFloat(c.PostForm("lat"), 64)
		long, _ := strconv.ParseFloat(c.PostForm("long"), 64)
		Type := c.Param("type")
		lotno := c.PostForm("lotno")
		index, _ := strconv.Atoi(c.PostForm("index"))
		pname := c.PostForm("patientName")
		pno := c.PostForm("patientNo")
		qrcode := c.PostForm("qrcode")
		px, _ := strconv.ParseFloat(c.PostForm("px"), 64)
		py, _ := strconv.ParseFloat(c.PostForm("py"), 64)
		pwidth, _ := strconv.ParseFloat(c.PostForm("pwidth"), 64)
		pheight, _ := strconv.ParseFloat(c.PostForm("pheight"), 64)
		projects := c.PostForm("projects")
		ps := strings.Split(projects, "|")
		photopath, plantpath, hotpath := savePictures(vendor, operator, Type, imgbase, px, py, pwidth, pheight, ps)
		record := record.NewRecord(qrcode, photopath, plantpath, hotpath, vendor, Type, project, operator, truename, location, lat, long, lotno, index, pname, pno, remark)
		switch Type {
		case "PH":
			result := ph.Test(util.GetCurrDir() + plantpath)
			record.Result = fmt.Sprintf("%.1f", result)
			record.Save()
			c.JSON(200, gin.H{"code": 0, "data": record, "msg": ""})
			break
		case "Combi":
			re := autionsticksG.AutionsticksG(util.GetCurrDir() + plantpath)
			record.Result = strings.Join(re, ",")
			record.Save()
			c.JSON(200, gin.H{"code": 0, "data": record, "msg": ""})
			break
		case "Qualitative":
			hpaths := strings.Split(hotpath, "|")
			fmt.Println(ps)
			for i := 0; i < len(ps); i++ {
				params := strings.Split(ps[i], ",")
				plen := len(params)
				hwidthmm, _ := strconv.Atoi(params[plen-1])
				cdis, _ := strconv.Atoi(params[4])
				tdis := make([]int, 0, 0)
				for j := 5; j <= plen-2; j++ {
					temp, _ := strconv.Atoi(params[j])
					tdis = append(tdis, temp)
				}

				fmt.Println(params)
				white, black := qualitative.BWValue(util.GetCurrDir() + photopath)
				isvalid, judges, grays := qualitative.Test(util.GetCurrDir()+hpaths[i], hwidthmm, cdis, tdis, white, black)
				if isvalid {
					for k := 0; k < len(judges); k++ {
						if judges[k] {
							record.Result += "+"
						} else {
							record.Result += "-"
						}
						if k < len(judges)-1 {
							record.Result += ","
						}
					}
					for k := 0; k < len(grays); k++ {
						record.GrayVal += fmt.Sprintf("%.2f", grays[k])
						if k < len(grays)-1 {
							record.GrayVal += ","
						}
					}
				} else {
					record.Result += "Invalid"
				}
				if i < len(ps)-1 {
					record.Result += " "
				}
			}
			record.Save()
			c.JSON(200, gin.H{"code": 0, "data": record, "msg": ""})
			break
		default:
			c.JSON(200, gin.H{"code": 1, "data": nil, "msg": "not surpport"})
		}
	})

	r.GET("/recode/exist", func(c *gin.Context) {
		code, _ := c.GetQuery("qrcode")
		result := record.Exist(code)

		c.JSON(200, gin.H{"code": 0, "data": result, "msg": ""})
	})

	r.GET("/record/query", func(c *gin.Context) {
		pageIndex, _ := c.GetQuery("pageIndex")
		pageSize, _ := c.GetQuery("pageSize")
		pname, _ := c.GetQuery("patientname")
		pno, _ := c.GetQuery("patientno")
		test, _ := c.GetQuery("test")
		factory, _ := c.GetQuery("factory")
		lotno, _ := c.GetQuery("lotno")
		result, _ := c.GetQuery("result")
		location, _ := c.GetQuery("location")
		operator, _ := c.GetQuery("operator")
		startdate, _ := c.GetQuery("startdate")
		enddate, _ := c.GetQuery("enddate")
		sort, _ := c.GetQuery("sort")
		id, _ := c.GetQuery("id")
		// fmt.Println(page, user, role, ty)
		idx, err := strconv.Atoi(pageIndex)
		size, _ := strconv.Atoi(pageSize)
		idi, _ := strconv.Atoi(id)
		if err != nil {
			c.JSON(200, gin.H{"code": 1, "data": err})
			return
		}

		total, rs := record.Query(idx, size, idi, pname, pno, test, factory, lotno, result, location, operator, startdate, enddate, sort)
		c.JSON(200, gin.H{"code": 0, "total": total, "data": rs})

	})

	//recall
	r.GET("/recall/add", func(c *gin.Context) {
		qrcode, _ := c.GetQuery("qrcode")
		fmt.Println(qrcode)
		ret, msg := manage.AddRecall(qrcode)
		c.JSON(200, gin.H{"ret": ret, "msg": msg})
	})

	r.GET("/recall/delete", func(c *gin.Context) {
		fmt.Println("start")
		vendor, _ := c.GetQuery("vendor")
		lotno, _ := c.GetQuery("lotno")
		fmt.Println(vendor, lotno)
		manage.Delete(vendor, lotno)
		c.JSON(200, gin.H{"code": 1, "msg": ""})
	})

	r.GET("/recall/query", func(c *gin.Context) {
		pageIndex, _ := c.GetQuery("pageIndex")
		pageSize, _ := c.GetQuery("pageSize")
		vendor, _ := c.GetQuery("vendor")
		lotno, _ := c.GetQuery("lotno")
		startdate, _ := c.GetQuery("startdate")
		enddate, _ := c.GetQuery("enddate")
		sort, _ := c.GetQuery("sort")
		idx, err := strconv.Atoi(pageIndex)
		size, _ := strconv.Atoi(pageSize)
		if err != nil {
			c.JSON(200, gin.H{"code": 1, "data": err})
			return
		}
		total, rs := manage.Query(idx, size, vendor, lotno, startdate, enddate, sort)
		c.JSON(200, gin.H{"code": 0, "total": total, "data": rs})

	})

	r.GET("/recall/exist/:vendor/:lotno", func(c *gin.Context) {
		vendor := c.Param("vendor")
		lotno := c.Param("lotno")
		result := manage.Exist(vendor, lotno)
		c.JSON(200, gin.H{"code": 0, "data": result, "msg": ""})
	})

	//tarck
	r.POST("/track/add", func(c *gin.Context) {
		operator := c.PostForm("operator")
		truename := c.PostForm("truename")
		location := c.PostForm("location")
		lat, _ := strconv.ParseFloat(c.PostForm("lat"), 64)
		long, _ := strconv.ParseFloat(c.PostForm("long"), 64)
		qr := c.PostForm("qrcode")

		result := track.AddTrackRecord(qr, operator, truename, location, lat, long)
		c.JSON(200, gin.H{"code": 0, "data": result, "msg": ""})
	})

	r.GET("/track/detail/:qrcode", func(c *gin.Context) {
		qrcode := c.Param("qrcode")
		entity := track.FindTrackEntity(qrcode)
		records := track.GetTrackRecordByQrCode(qrcode)
		c.JSON(200, gin.H{"code": 0, "entity": entity, "records": records})
	})

	r.GET("/track/list/:page", func(c *gin.Context) {
		page := c.Param("page")
		user, _ := c.GetQuery("user")
		role, _ := c.GetQuery("role")
		idx, err := strconv.Atoi(page)
		if err != nil {
			c.JSON(200, gin.H{"code": 1, "data": err})
			return
		}
		rs := track.GetList(idx, 20, user, role)
		c.JSON(200, gin.H{"code": 0, "data": rs})
	})

	r.GET("/track/query", func(c *gin.Context) {
		pageIndex, _ := c.GetQuery("pageIndex")
		pageSize, _ := c.GetQuery("pageSize")
		project, _ := c.GetQuery("project")
		factory, _ := c.GetQuery("factory")
		lotno, _ := c.GetQuery("lotno")
		index, _ := c.GetQuery("index")
		operator, _ := c.GetQuery("operator")
		sort, _ := c.GetQuery("sort")
		id, _ := c.GetQuery("id")
		// fmt.Println(page, user, role, ty)
		idx, err := strconv.Atoi(pageIndex)
		size, _ := strconv.Atoi(pageSize)
		idi, _ := strconv.Atoi(id)
		if err != nil {
			c.JSON(200, gin.H{"code": 1, "data": err})
			return
		}
		total, rs := track.Query(idx, size, idi, project, factory, lotno, index, operator, sort)
		c.JSON(200, gin.H{"code": 0, "total": total, "data": rs})

	})

	r.Run(":" + config.Config["port"])
}

func getSavePath(flag string) string {
	basedir := util.GetCurrDir()
	savedir := basedir + "images/" + flag + "/"
	if !util.IsExistFileOrDir(savedir) {
		os.MkdirAll(savedir, 0777) //创建文件夹
	}
	savepath := savedir + fmt.Sprintf("%d", time.Now().UnixNano()) + ".jpg"
	return savepath
}

func savePictures(vendor, operator, flag, base64Str string, px, py, pwidth, pheight float64, ps []string) (photoPath, plantPath, hotPath string) {
	basedir := util.GetCurrDir()
	savedir := "/images/" + operator + "/" + vendor + "/" + flag + "/"
	if !util.IsExistFileOrDir(basedir + savedir) {
		os.MkdirAll(basedir+savedir, 0777) //创建文件夹
	}
	tnow := fmt.Sprintf("%d", time.Now().UnixNano())
	photoPath = savedir + tnow + "_photo" + ".jpg"
	plantPath = savedir + tnow + "_plant" + ".jpg"
	for i := 0; i < len(ps); i++ {
		hotPath += savedir + tnow + "_hot" + strconv.Itoa(i) + ".jpg"
		if i < len(ps)-1 {
			hotPath += "|"
		}
	}
	ddd, _ := base64.StdEncoding.DecodeString(base64Str) //成图片文件并把文件写入到buffer
	ioutil.WriteFile(basedir+photoPath, ddd, 0666)

	img := ana.Base64ToImg(base64Str)

	plant := ana.ImgCutByRate(img, px, py, pwidth, pheight)
	f, _ := os.Create(basedir + plantPath) //创建文件
	defer f.Close()                        //关闭文件
	jpeg.Encode(f, plant, nil)

	hpaths := strings.Split(hotPath, "|")
	for i := 0; i < len(ps); i++ {
		params := strings.Split(ps[i], ",")
		x, _ := strconv.ParseFloat(params[0], 64)
		y, _ := strconv.ParseFloat(params[1], 64)
		w, _ := strconv.ParseFloat(params[2], 64)
		h, _ := strconv.ParseFloat(params[3], 64)
		hot := ana.ImgCutByRate(img, x, y, w, h)
		f1, _ := os.Create(basedir + hpaths[i])
		defer f1.Close()
		jpeg.Encode(f1, hot, nil)
	}
	return
}
