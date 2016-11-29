package main

import (
	"fmt"
	"gopkg.in/gin-gonic/gin.v1"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", indexPage)
	router.StaticFS("/images", http.Dir("images"))
	router.StaticFS("/outputs", http.Dir("outputs"))
	router.GET("/image/:name/cut", imageProx)
	router.GET("/image/:name/size", imageSize)
	router.Run(":8080")
}

func indexPage(c *gin.Context) {

	files, err := fileList("images")
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Image processing",
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
	if info.IsDir() {
		return false
	} else {
		name := info.Name()
		return strings.Contains(name, ".png") || strings.Contains(name, ".jpeg") || strings.Contains(name, ".jpg")
	}

}

func imageSize(c *gin.Context) {
	inputname := c.Param("name")
	img, err := imageFile(inputname)
	if err != nil {
		c.AbortWithError(http.StatusConflict, err)
		return
	}
	bounds := img.Bounds()
	x, y := bounds.Max.X, bounds.Max.Y
	c.JSON(http.StatusOK, gin.H{
		"width": x,
		"height": y,
	})
}

func imageFile(filename string) (image.Image, error) {
	var img image.Image

	file, err := os.Open(fmt.Sprint("images/", filename))
	if err != nil {
		return img, fmt.Errorf("Open file: %v", err)
	}

	defer file.Close()

	if strings.Contains(filename, ".png") {
		img, err = png.Decode(file)
	} else {
		img, err = jpeg.Decode(file)
	}

	if err != nil {
		return img, fmt.Errorf("Decode image: %v", err)
	}
	return img, nil
}

func imageProx(c *gin.Context) {
	inputname := c.Param("name")
	outx := c.Query("x")
	outy := c.Query("y")
	outwidth := c.Query("w")
	outHeight := c.Query("h")

	img, err := imageFile(inputname)
	if err != nil {
		c.AbortWithError(http.StatusConflict, err)
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
	outfilename := fmt.Sprint("outputs/", pointX, "_", pointY, "_", w, "_", h, "_", inputname)
	outfile, err := os.Create(outfilename)
	if err != nil {
		c.String(http.StatusInternalServerError, "Create file: ", err)
		return
	}
	defer outfile.Close()

	png.Encode(outfile, tmpImg)
	log.Println(outfilename)
	c.JSON(http.StatusOK, gin.H{
		"filename": outfilename,
	})
}