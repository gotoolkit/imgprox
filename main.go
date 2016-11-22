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
			c.String(http.StatusNotFound, "Open file: ", err)
			return
		}
		defer file.Close()

		img, err := png.Decode(file)
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
	})
	router.Run(":8080")
}