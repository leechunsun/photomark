package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	"image/draw"
	"image/jpeg"
	"net/http"
	"os"
	"strconv"
	"time"

)
import "./funx"

var TILESDB map[string][3]float64


func CloneTilesDB() map[string][3]float64 {
	db := make(map[string][3]float64)
	for k, v := range TILESDB {
		db[k] = v
	}
	return db
}



func main() {
	//mux := http.NewServeMux()
	//files := http.FileServer(http.Dir("D:\\go\\codes\\photomark\\static"))
	//mux.Handle("/static/", http.StripPrefix("/static/", files))
	//mux.HandleFunc("/", upload)
	//mux.HandleFunc("/mosaic", mosaic)
	//server := &http.Server{
	//	Addr:              "127.0.0.1:8779",
	//	Handler:           mux,
	//}
	TILESDB = funx.TilesDB()
	//fmt.Println("server start....")
	//server.ListenAndServe()
	gen()
}


func upload(w http.ResponseWriter, r *http.Request)  {
	t, _ := template.ParseFiles("D:\\go\\codes\\photomark\\upload.html")
	t.Execute(w, nil)
}

func mosaic(w http.ResponseWriter, r *http.Request)  {
	t0 := time.Now()
	r.ParseMultipartForm(10485760)
	file, _, _ := r.FormFile("image")
	defer file.Close()
	tileSize, _ := strconv.Atoi(r.FormValue("title_size"))
	original, _, _ := image.Decode(file)
	bounds := original.Bounds()
	new_image := image.NewNRGBA(image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y))
	db := CloneTilesDB()
	sp := image.Point{0,0}
	for y := bounds.Min.Y; y < bounds.Max.Y; y=y+ tileSize{
		for x := bounds.Min.X; x < bounds.Max.X; x=x+ tileSize{
			r, g, b, _ := original.At(x, y).RGBA()
			color := [3]float64{float64(r), float64(g), float64(b)}
			nearst := funx.Nearest(color, &db)
			file, err := os.Open(nearst)
			if err == nil {
				img, _, err := image.Decode(file)
				if err == nil{
					t := funx.Resize(img, tileSize)
					tile := t.SubImage(t.Bounds())
					tileBounds := image.Rect(x, y, x+tileSize, y+tileSize)
					draw.Draw(new_image, tileBounds, tile, sp, draw.Src)
				}else {
					fmt.Println("nearst....:" , nearst)
				}
			}else{
				fmt.Println("error:", nearst)
			}
			file.Close()
		}
	}
	buf1 := new(bytes.Buffer)
	jpeg.Encode(buf1, original, nil)
	originalStr := base64.StdEncoding.EncodeToString(buf1.Bytes())
	buf2 := new(bytes.Buffer)
	jpeg.Encode(buf2, new_image, nil)
	masaicStr := base64.StdEncoding.EncodeToString(buf2.Bytes())
	t1 := time.Now()
	_ = map[string]string{
		"original": originalStr,
		"masaic": masaicStr,
		"duration": fmt.Sprintf("%v ", t1.Sub(t0)),
	}
	file, _ = os.Create("D:\\go\\codes\\photomark\\aa.jpg")

	//t, _ := template.ParseFiles("D:\\go\\codes\\photomark\\result.html")
	w.Write([]byte("aaaaa"))
}

func gen()  {
	t0 := time.Now()
	file, _ := os.OpenFile("D:\\go\\codes\\photomark\\oo.jpg", os.O_RDWR, 0666)
	defer file.Close()
	tileSize := 40
	original, _, _ := image.Decode(file)
	bounds := original.Bounds()
	new_image := image.NewNRGBA(image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y))
	db := CloneTilesDB()
	sp := image.Point{0,0}
	for y := bounds.Min.Y; y < bounds.Max.Y; y=y+ tileSize{
		for x := bounds.Min.X; x < bounds.Max.X; x=x+ tileSize{
			r, g, b, _ := original.At(x, y).RGBA()
			color := [3]float64{float64(r), float64(g), float64(b)}
			nearst := funx.Nearest(color, &db)
			file, err := os.Open(nearst)
			if err == nil {
				img, _, err := image.Decode(file)
				if err == nil{
					t := funx.Resize(img, tileSize)
					tile := t.SubImage(t.Bounds())
					tileBounds := image.Rect(x, y, x+tileSize, y+tileSize)
					draw.Draw(new_image, tileBounds, tile, sp, draw.Src)
					fmt.Println(x, y)
				}else {
					fmt.Println("nearst....:" , nearst)
				}
			}else{
				fmt.Println("error:", nearst)
			}
			file.Close()
		}
	}
	buf1 := new(bytes.Buffer)
	jpeg.Encode(buf1, original, nil)
	originalStr := base64.StdEncoding.EncodeToString(buf1.Bytes())
	buf2 := new(bytes.Buffer)
	jpeg.Encode(buf2, new_image, nil)
	masaicStr := base64.StdEncoding.EncodeToString(buf2.Bytes())
	t1 := time.Now()
	_ = map[string]string{
		"original": originalStr,
		"masaic": masaicStr,
		"duration": fmt.Sprintf("%v ", t1.Sub(t0)),
	}
	fileob, _ := os.OpenFile("D:\\go\\codes\\photomark\\aa.jpeg", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	fileob.Write(buf2.Bytes())
	//t, _ := template.ParseFiles("D:\\go\\codes\\photomark\\result.html")
	//w.Write([]byte("aaaaa"))
}
