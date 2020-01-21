package funx

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"os"
)

func AverageColor(img image.Image) [3]float64 {
	bounds := img.Bounds()
	r, g, b := 0.0, 0.0, 0.0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++{
		for x := bounds.Min.X; x < bounds.Max.X; x++{
			r1, g1, b1 ,_ := img.At(x, y).RGBA()
			r, g, b = r+ float64(r1), g +float64(g1), b+float64(b1)
		}
	}
	totalP := float64(bounds.Max.Y*bounds.Max.X)
	return [3]float64{r/totalP, g/totalP, b/totalP}
}

func Resize(in image.Image, newWidth int) image.NRGBA {
	bounds := in.Bounds()
	ratio := bounds.Dx() / newWidth
	out := image.NewNRGBA(image.Rect(bounds.Min.X/ratio, bounds.Min.Y/ratio, bounds.Max.X/ratio, bounds.Max.Y/ratio))
	for y, j := bounds.Min.Y, bounds.Min.Y; y < bounds.Max.Y; y, j = y+ratio, j+1{
		for x, k := bounds.Min.X, bounds.Min.X; x < bounds.Max.X; x, k = x+ratio, k+1{
			r, g, b, a := in.At(x, y).RGBA()
			out.SetNRGBA(k, j, color.NRGBA{uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)})
		}
	}
	return *out
}

func TilesDB() map[string][3]float64 {
	path := "D:\\go\\codes\\photomark\\tiles"
	fmt.Println("start populating tiles db....")
	db := make(map[string][3]float64)
	files, _ := ioutil.ReadDir(path)
	for _, f := range files{
		name := path + "\\" + f.Name()
		file, err := os.Open(name)
		if err == nil{
			img, _, err := image.Decode(file)
			if err == nil {
				db[name] = AverageColor(img)
			}else{
				fmt.Println("errr in populating : ", err, name)
			}
		}else{
			fmt.Println("can't open file: ", err, name)
		}
		file.Close()
	}
	fmt.Println("ok")
	return db
}

func  Nearest(tar [3]float64, db *map[string][3]float64) string {
	var filename string
	smallest := 1000000.0
	for k, v := range *db{
		dist := distance(tar, v)
		if dist < smallest{
			smallest = dist
			filename = k
		}
	}
	return filename
}

func distance(tar [3]float64, v [3]float64) float64 {
	return math.Sqrt(sq(v[0] - tar[0]) + sq(v[1] - tar[1]) + sq(v[2] - tar[2]))
}

func sq(n float64) float64 {
	return n*n
}



