// Command genicon renders the Vigilant logo (shield + eye) and writes the
// application icons. Run it from the repository root:
//
//	go run ./tools/genicon
//
// The geometry mirrors assets/logo.svg so both stay visually identical.
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
)

type pt struct{ x, y float64 }

// shape is a closed polygon in a 256x256 coordinate space.
type shape []pt

func quad(p0, p1, p2 pt, n int) []pt {
	out := make([]pt, 0, n)
	for i := 1; i <= n; i++ {
		t := float64(i) / float64(n)
		u := 1 - t
		out = append(out, pt{u*u*p0.x + 2*u*t*p1.x + t*t*p2.x, u*u*p0.y + 2*u*t*p1.y + t*t*p2.y})
	}
	return out
}

func cubic(p0, p1, p2, p3 pt, n int) []pt {
	out := make([]pt, 0, n)
	for i := 1; i <= n; i++ {
		t := float64(i) / float64(n)
		u := 1 - t
		out = append(out, pt{
			u*u*u*p0.x + 3*u*u*t*p1.x + 3*u*t*t*p2.x + t*t*t*p3.x,
			u*u*u*p0.y + 3*u*u*t*p1.y + 3*u*t*t*p2.y + t*t*t*p3.y,
		})
	}
	return out
}

func (s shape) contains(x, y float64) bool {
	in := false
	for i, j := 0, len(s)-1; i < len(s); j, i = i, i+1 {
		a, b := s[i], s[j]
		if (a.y > y) != (b.y > y) && x < (b.x-a.x)*(y-a.y)/(b.y-a.y)+a.x {
			in = !in
		}
	}
	return in
}

func shield(inset float64) shape {
	// M128 12 L216 42 V120 C216 176 178 218 128 244 C78 218 40 176 40 120 V42 Z
	s := shape{{128, 12}, {216, 42}, {216, 120}}
	s = append(s, cubic(pt{216, 120}, pt{216, 176}, pt{178, 218}, pt{128, 244}, 24)...)
	s = append(s, cubic(pt{128, 244}, pt{78, 218}, pt{40, 176}, pt{40, 120}, 24)...)
	s = append(s, pt{40, 42})
	if inset == 0 {
		return s
	}
	// Shrink towards the centroid for the inner bevel.
	c := pt{128, 124}
	k := 1 - inset
	for i := range s {
		s[i] = pt{c.x + (s[i].x-c.x)*k, c.y + (s[i].y-c.y)*k}
	}
	return s
}

func lens(rx, ry float64) shape {
	// M(128-rx) 128 Q128 (128-ry) (128+rx) 128 Q128 (128+ry) (128-rx) 128 Z
	l := shape{{128 - rx, 128}}
	l = append(l, quad(pt{128 - rx, 128}, pt{128, 128 - ry}, pt{128 + rx, 128}, 32)...)
	l = append(l, quad(pt{128 + rx, 128}, pt{128, 128 + ry}, pt{128 - rx, 128}, 32)...)
	return l
}

type palette struct {
	shieldTop, shieldBottom, bevel color.RGBA
	eye, eyeDark                   color.RGBA
}

var (
	active = palette{
		shieldTop: color.RGBA{0x3b, 0x3d, 0x42, 255}, shieldBottom: color.RGBA{0x1c, 0x1d, 0x20, 255},
		bevel: color.RGBA{0x4a, 0x4c, 0x52, 255},
		eye:   color.RGBA{0x4c, 0xff, 0x1a, 255}, eyeDark: color.RGBA{0x16, 0x17, 0x19, 255},
	}
	paused = palette{
		shieldTop: color.RGBA{0x3b, 0x3d, 0x42, 255}, shieldBottom: color.RGBA{0x1c, 0x1d, 0x20, 255},
		bevel: color.RGBA{0x4a, 0x4c, 0x52, 255},
		eye:   color.RGBA{0x8b, 0x90, 0x99, 255}, eyeDark: color.RGBA{0x16, 0x17, 0x19, 255},
	}
)

var (
	outer     = shield(0)
	bevelS    = shield(0.07)
	innerS    = shield(0.075)
	eyeOuter  = lens(88, 72)
	eyeInner  = lens(70, 52)
	irisCx    = 128.0
	irisCy    = 128.0
	irisOuter = 31.0
	irisInner = 21.0
	pupil     = 11.0
	hlx, hly  = 142.0, 114.0
	hlr       = 5.0
)

func lerp(a, b color.RGBA, t float64) color.RGBA {
	f := func(x, y uint8) uint8 { return uint8(float64(x) + (float64(y)-float64(x))*t + 0.5) }
	return color.RGBA{f(a.R, b.R), f(a.G, b.G), f(a.B, b.B), 255}
}

func dist(x, y, cx, cy float64) float64 { return math.Hypot(x-cx, y-cy) }

// sample returns the colour at a point of the 256x256 logo space (alpha 0 outside).
func sample(p palette, x, y float64) color.RGBA {
	if !outer.contains(x, y) {
		return color.RGBA{}
	}
	c := lerp(p.shieldTop, p.shieldBottom, math.Min(1, math.Max(0, (y-12)/232)))
	if bevelS.contains(x, y) && !innerS.contains(x, y) {
		c = p.bevel
	}
	if eyeOuter.contains(x, y) {
		c = p.eye
		if eyeInner.contains(x, y) {
			c = p.eyeDark
			d := dist(x, y, irisCx, irisCy)
			if d <= irisOuter && d >= irisInner {
				c = p.eye
			}
			if dist(x, y, hlx, hly) <= hlr {
				c = p.eye
			}
		}
	}
	return c
}

func render(p palette, size int) *image.NRGBA {
	const ss = 4
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	scale := 256.0 / float64(size)
	for py := 0; py < size; py++ {
		for px := 0; px < size; px++ {
			var r, g, b, a float64
			for sy := 0; sy < ss; sy++ {
				for sx := 0; sx < ss; sx++ {
					x := (float64(px) + (float64(sx)+0.5)/ss) * scale
					y := (float64(py) + (float64(sy)+0.5)/ss) * scale
					c := sample(p, x, y)
					if c.A > 0 {
						r += float64(c.R)
						g += float64(c.G)
						b += float64(c.B)
						a++
					}
				}
			}
			if a > 0 {
				img.SetNRGBA(px, py, color.NRGBA{uint8(r/a + 0.5), uint8(g/a + 0.5), uint8(b/a + 0.5), uint8(a / (ss * ss) * 255)})
			}
		}
	}
	return img
}

func encodePNG(img image.Image) []byte {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

// writeICO writes a multi-resolution ICO containing PNG-compressed images.
func writeICO(path string, p palette, sizes []int) {
	var imgs [][]byte
	for _, s := range sizes {
		imgs = append(imgs, encodePNG(render(p, s)))
	}
	var buf bytes.Buffer
	w := func(v any) { _ = binary.Write(&buf, binary.LittleEndian, v) }
	w(uint16(0))
	w(uint16(1))
	w(uint16(len(sizes)))
	offset := 6 + 16*len(sizes)
	for i, s := range sizes {
		dim := uint8(s)
		if s >= 256 {
			dim = 0
		}
		buf.Write([]byte{dim, dim, 0, 0})
		w(uint16(1))
		w(uint16(32))
		w(uint32(len(imgs[i])))
		w(uint32(offset))
		offset += len(imgs[i])
	}
	for _, d := range imgs {
		buf.Write(d)
	}
	mustWrite(path, buf.Bytes())
}

func mustWrite(path string, data []byte) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote", path)
}

func main() {
	sizes := []int{16, 24, 32, 48, 64, 128, 256}
	writeICO("assets/icon.ico", active, sizes)
	writeICO("assets/icon-paused.ico", paused, sizes)
	mustWrite("assets/icon.png", encodePNG(render(active, 512)))
}
