package main

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/imgo"
	"image"
	"log"
	"math"
	"time"
)

var (
	StartX = 46
	StartY = 390
	Width  = 275
	Height = 340
)

var (
	knife  = fmt.Sprintf("%x%x%x", 189, 200, 225)
	target = fmt.Sprintf("%x%x%x", 234, 61, 72)
)

func main() {
	for i := 1; i <= 1000; i++ {
		fmt.Printf("第%d次执行 ...\n", i)

		time.Sleep(1500 * time.Millisecond)
		run()
	}
}

func run() {
	if _, err := capture("游戏区域.png", StartX, StartY, Width, Height); err != nil {
		log.Fatalf("capture(0) fail :: %s \n", err.Error())
	}

	img1, err := capture("游戏区域_左侧.png", StartX, StartY, 22, Height)
	if err != nil {
		log.Fatalf("capture(1) fail :: %s \n", err.Error())
	}

	img2, err := capture("游戏区域_右侧.png", StartX+Width-21, StartY, 22, Height)
	if err != nil {
		log.Fatalf("capture(2) fail :: %s \n", err.Error())
	}

	var o1, o2 image.Point

	xy1, ok1 := recognizeKnifeEx(img1, StartX, StartY, knife)
	if ok1 {
		o1 = xy1
		fmt.Printf("飞刀在左侧，绝对位置为%s \n", xy1)

		o2 = recognizeTargetEx(img2, StartX+Width-21, StartY, target)
		fmt.Printf("目标在右侧，绝对位置为%s \n", o2)
	} else {
		xy2, ok2 := recognizeKnifeEx(img2, StartX+Width-21, StartY, knife)
		if ok2 {
			o1 = xy2
			fmt.Printf("飞刀在右侧，绝对位置为%s \n", xy2)

			o2 = recognizeTargetEx(img1, StartX, StartY, target)
			fmt.Printf("目标在左侧，绝对位置为%s \n", o2)
		} else {
			log.Fatalf("在游戏区域的左侧和右侧均没有找到对应的颜色值 #%s \n", knife)
		}
	}

	k := float64(o2.Y-o1.Y) / float64(o2.X-o1.X)
	ln := math.Sqrt(float64((o2.Y-o1.Y)*(o2.Y-o1.Y) + (o2.X-o1.X)*(o2.X-o1.X)))
	fmt.Printf("标准斜率为 %.2f，直线距离为 %.2f \n", k, ln)

	//p0 := image.Point{
	//	X: o1.X + int(float64(o2.X-o1.X)*(0.175)),
	//	Y: o1.Y + int(float64(o2.Y-o1.Y)*(0.175)),
	//}
	//robotgo.Move(p0.X, p0.Y)
	//
	//colors := make(map[string]int)
	//for i := 0; i < 12; i++ {
	//	c0 := robotgo.GetPixelColor(p0.X, p0.Y)
	//
	//	colors[c0]++
	//	time.Sleep(100 * time.Millisecond)
	//}
	//
	//var max int
	//var def string
	//for c, n := range colors {
	//	if n > max {
	//		max = n
	//		def = c
	//	}
	//}
	//fmt.Printf("背景颜色为 #%s \n", def)
	//
	//for {
	//	c0 := robotgo.GetPixelColor(p0.X, p0.Y)
	//	if !strings.EqualFold(c0, def) {
	//		fmt.Printf("方向箭头颜色为 #%s \n", c0)
	//
	//		robotgo.Toggle("left")
	//		robotgo.Toggle("left", "up")
	//
	//		return
	//	}
	//}

	s1 := fmt.Sprintf("%x%x%x", 224, 224, 224)
	s2 := fmt.Sprintf("%x%x%x", 254, 254, 254)

	start := time.Now()
	for {

		diff := time.Now().Sub(start).Seconds()
		p0 := image.Point{
			X: o1.X + int(float64(o2.X-o1.X)*(0.78-diff*0.02)),
			Y: o1.Y + int(float64(o2.Y-o1.Y)*(0.78-diff*0.02)),
		}
		robotgo.Move(p0.X, p0.Y)

		c0 := robotgo.GetPixelColor(p0.X, p0.Y)
		if c0 >= s1 && c0 <= s2 {
			fmt.Printf("方向箭头颜色为 %s \n", c0)

			robotgo.Toggle("left")
			robotgo.Toggle("left", "up")

			return
		}
	}
}

func capture(name string, x, y, dx, dy int) (image.Image, error) {
	bit := robotgo.CaptureScreen(x, y, dx, dy)
	defer robotgo.FreeBitmap(bit)

	img := robotgo.ToImage(bit)
	if err := imgo.Save(name, img); err != nil {
		return nil, err
	}

	return img, nil
}

func recognizeKnifeEx(img image.Image, x, y int, color string) (image.Point, bool) {
	rect := img.Bounds()

	for h0 := rect.Min.Y; h0 < rect.Max.Y; h0++ {
		for w0 := rect.Min.X; w0 < rect.Max.X; w0++ {
			r0, g0, b0, _ := img.At(w0, h0).RGBA()
			if r0>>8 == 189 && g0>>8 == 200 && b0>>8 == 225 {
				return image.Point{X: x + w0/2, Y: y + h0/2}, true
			}
		}
	}

	return image.Point{}, false
}

func recognizeTargetEx(img image.Image, x, y int, color string) image.Point {
	rect := img.Bounds()

	wMin, wMax, hMin, hMax := 9999, 0, 9999, 0
	for w0 := rect.Min.X; w0 < rect.Max.X; w0++ {
		for h0 := rect.Min.Y; h0 < rect.Max.Y; h0++ {
			r0, g0, b0, _ := img.At(w0, h0).RGBA()

			if r0>>8 == 234 && g0>>8 == 61 && b0>>8 == 72 {
				if wMin > w0 {
					wMin = w0
				}

				if hMin > h0 {
					hMin = h0
				}

				if wMax < w0 {
					wMax = w0
				}

				if hMax < h0 {
					hMax = h0
				}
			}
		}
	}

	return image.Point{X: x + (wMin+wMax)/4, Y: y + (hMin+hMax)/4}
}
