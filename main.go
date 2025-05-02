package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

var palette = []color.Color{
	color.White,                  // 0
	color.Black,                  // 1
	color.RGBA{255, 0, 0, 255},   // 2 - Rouge
	color.RGBA{0, 255, 0, 255},   // 3 - Vert
	color.RGBA{0, 0, 255, 255},   // 4 - Bleu
	color.RGBA{255, 255, 0, 255}, // 5 - Jaune
	color.RGBA{0, 255, 255, 255}, // 6 - Cyan
	color.RGBA{255, 0, 255, 255}, // 7 - Magenta
}

func lissajous(out io.Writer, cycles int, nframes int, delay int, size int) {
	const (
		res = 0.001 // angular resolution
	)
	fmt.Printf("Cycles : %d | nframes : %d | delay : %d | size : %d\n", cycles, nframes, delay, size)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	color := 1.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < float64(cycles)*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*float64(size)+0.5), size+int(y*float64(size)+0.5), uint8(color))
		}
		color += 0.1
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	gif.EncodeAll(out, &anim)

}

var mu sync.Mutex
var count int

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/count", counter)
	http.HandleFunc("/img", img)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
func formLissajous(r *http.Request) (int, int, int, int) {
	if err:=r.ParseForm() ; err!=nil{
		log.Print(err)
	}
	cycles := 5
	nframes := 64
	demiSize := 100
	delay := 8
	formchecker("cycles", r, &cycles)
	formchecker("nframes", r, &nframes)
	formchecker("demiSize", r, &demiSize)
	formchecker("delay", r, &delay)
	return cycles, nframes, demiSize, delay
}
func formchecker(str string, r *http.Request, toChange *int) {
	if r.Form[str] != nil {
		fmt.Printf("%q = %q\n", str, r.Form[str])
		*toChange, _ = strconv.Atoi(r.Form[str][0])
	}
}
func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintf(w, "%s %s %s \n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q]= %q\n", k, v)
	}

	fmt.Fprintf(w, "Host = %q \n", r.Host)
	fmt.Fprintf(w, "RemoteAddr = %q \n", r.RemoteAddr)

	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}

	for k, v := range r.Form {
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
	}
	cycles, nframes, demiSize, delay := formLissajous(r)
	fmt.Fprintf(w, `
        <html>
        <body>
            <h1>Welcome</h1>
            <p>This is a debug page with a GIF:</p>
            <img src="/img?cycles=%d&nframes=%d&demiSize=%d&delay=%d" alt="Lissajous animation">
        </body>
        </html>
    `, cycles, nframes, demiSize, delay)
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w, "Url.Path=%q\n", r.URL.Path)
}
func img(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	w.Header().Set("Content-Type", "image/gif")
	cycles, nframes, demiSize, delay := formLissajous(r)

	lissajous(w, cycles, nframes, delay, demiSize)
}
func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	fmt.Fprintf(w, "Count of Requests : %d", count)
	mu.Unlock()
}
