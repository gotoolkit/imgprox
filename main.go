package main

import (
	"os"
	"log"
	"image/png"
	"image"
	"image/draw"
	"gopkg.in/gin-gonic/gin.v1"
	"fmt"
	"strconv"
	"net/http"
	"strings"
	"image/jpeg"
	"path/filepath"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", indexPage)
	router.StaticFS("/images", http.Dir("images"))
	router.GET("/image/:name/cut", imageProx)
	router.Run(":8080")
}

func indexPage(c *gin.Context) {

	files, err := fileList("images")
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Main website",
		"files": files,
	})
}

type ImageFile struct {
	FileInfo os.FileInfo
	RealPath string
}

func fileList(path string) ([]ImageFile, error) {
	var imgFiles []ImageFile

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if isImageFile(info) {
			img := ImageFile{
				FileInfo: info,
				RealPath: path,
			}
			imgFiles = append(imgFiles, img)
		}
		return nil
	})
	return imgFiles, err
}

func isImageFile(info os.FileInfo) bool {
	if  info.IsDir() {
		return false
	} else {
		name := info.Name()
		return strings.Contains(name, ".png") || strings.Contains(name, ".jpeg") || strings.Contains(name, ".jpg")
	}

}


func imageProx(c *gin.Context) {
	inputname := c.Param("name")
	outx := c.Query("x")
	outy := c.Query("y")
	outwidth := c.Query("w")
	outHeight := c.Query("h")
	file, err := os.Open(fmt.Sprint("images/", inputname))
	if err != nil {
		c.String(http.StatusNotFound, "Open file: ", err)
		return
	}
	defer file.Close()

	var img image.Image

	if strings.Contains(inputname, ".png") {
		img, err = png.Decode(file)
	} else {
		img, err = jpeg.Decode(file)
	}

	if err != nil {
		c.String(http.StatusConflict, "Decode image: ", err)
		return
	}

	bounds := img.Bounds()
	x, y := bounds.Max.X, bounds.Max.Y
	log.Println("width: ", x, " height: ", y)
	var w, h int
	if outwidth == "" {
		w = x
	} else {
		w, _ = strconv.Atoi(outwidth)
	}

	if outHeight == "" {
		h = y
	} else {
		h, _ = strconv.Atoi(outHeight)
	}

	tmpImg := image.NewRGBA(image.Rect(0, 0, w, h))

	pointX, _ := strconv.Atoi(outx)
	pointY, _ := strconv.Atoi(outy)
	draw.Draw(tmpImg, tmpImg.Bounds(), img, image.Pt(pointX, pointY), draw.Src)
	outfilename := fmt.Sprint("images/result_", inputname)
	outfile, err := os.Create(outfilename)
	if err != nil {
		c.String(http.StatusInternalServerError, "Create file: ", err)
		return
	}
	defer outfile.Close()

	png.Encode(outfile, tmpImg)

	c.File(outfilename)
}
