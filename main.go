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
)

func main() {
	router := gin.Default()

	router.GET("/images/:name", func(c *gin.Context) {
		inputname := c.Param("name")
		outx := c.Query("x")
		outy := c.Query("y")
		outwidth := c.Query("w")
		outHeight := c.Query("h")
		file, err := os.Open(fmt.Sprint("images/", inputname))
		if err != nil {
			log.Fatal("Open file: ", err)
		}
		defer file.Close()

		img, err := png.Decode(file)
		if err != nil {
			log.Fatal("Decode image: ", err)
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
			log.Fatal("create file: ", err)
		}
		defer outfile.Close()

		png.Encode(outfile, tmpImg)

		c.File(outfilename)
	})
	router.Run(":8080")
}