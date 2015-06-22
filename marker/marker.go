package marker

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strconv"
)

// Marker represents the structure of a fiducial marker
type Marker struct {
	Code      int
	Size      int
	Division  int
	BlockSize int
	matrix    []int
	Marker    *image.NRGBA
	Name      string
}

// New returns a Marker with the given parameters.
// Using these parameters, it generates the visual representation of the marker
func New(code, division, blocksize int, name string, hasBorder bool) *Marker {

	if division < 3 || division > 8 {
		log.Fatal("The value of division must range from 3 to 8.")
	}

	if blocksize < 16 || blocksize > 32 {
		log.Fatal("The value of blocksize must range from 16 to 32")
	}

	size := blocksize * (division + 4)

	matrix := make([]int, division*division)
	matrix[0] = 1
	matrix[division-1] = -1
	matrix[division*(division-1)] = -1
	matrix[division*division-1] = -1

	binary := reverse(strconv.FormatInt(int64(code), 2))
	for i, j := 1, 0; i < len(matrix) && j < len(binary); i++ {

		if matrix[i] != -1 {
			matrix[i] = int(binary[j]) - 48
			j++
		}
	}

	m := image.NewNRGBA(image.Rect(0, 0, size, size))
	if hasBorder {
		draw.Draw(m, m.Bounds(), image.Black, image.ZP, draw.Src)
	} else {
		draw.Draw(m, m.Bounds(), image.Transparent, image.ZP, draw.Src)
	}

	fid := &Marker{
		code,
		size,
		division,
		blocksize,
		matrix,
		m,
		name}

	fid.draw(hasBorder)

	return fid
}

// Save saves the fiducial marker into a PNG image.
// If no name is specified, 'code-<code>.png' will be used as the filename.
func (m *Marker) Save() error {
	if m.Name == "" {
		m.Name = fmt.Sprintf("code-%d.png", m.Code)
	}

	f, err := os.Create(m.Name)
	if err != nil {
		return err
	}

	defer f.Close()

	if err = png.Encode(f, m.Marker); err != nil {
		return err
	}

	return nil
}

func (m *Marker) draw(hasBorder bool) {
	codeX := m.BlockSize * 2
	codeY := m.BlockSize * 2

	whiteBlock := image.Rect(m.BlockSize, m.BlockSize, m.Size-m.BlockSize, m.Size-m.BlockSize)
	blackBlock := image.Rect(m.BlockSize*2, m.BlockSize*2, m.Size-m.BlockSize*2, m.Size-m.BlockSize*2)

	draw.Draw(m.Marker, whiteBlock, image.White, image.ZP, draw.Src)
	draw.Draw(m.Marker, blackBlock, image.Black, image.ZP, draw.Src)

	for i, r := 0, 0; i < len(m.matrix); i++ {
		if m.matrix[i] == 1 {
			m.set(codeX, codeY, codeX+m.BlockSize, codeY+m.BlockSize, color.NRGBA{255, 255, 255, 255})
		}

		if r == m.Division-1 {
			r = 0
			codeY = codeY + m.BlockSize
			codeX = m.BlockSize * 2
			continue
		}

		codeX = codeX + m.BlockSize
		r = r + 1
	}
}

func (m *Marker) set(x1, y1, x2, y2 int, color color.Color) {
	for x := x1; x < x2; x++ {
		for y := y1; y < y2; y++ {
			m.Marker.Set(x, y, color)
		}
	}
}
