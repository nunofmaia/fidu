package marker

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strconv"
)

type Marker struct {
	Code       int
	Size       int
	Division   int
	Block_size int
	matrix     []int
	Marker     *image.RGBA
	Name       string
}

func New(code, size, division, blocksize int, name string) *Marker {

	if division < 3 || division > 8 {
		log.Fatal("The value of division must range from 3 to 8.")
	}

	s := blocksize * (division + 4)
	if size != s {
		if (size % (division + 4)) != 0 {
			log.Fatal("The proportions of the marker are incorrect.")
		}

		blocksize = size / (division + 4)
	} else {
		size = s
	}

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

	fid := &Marker{
		code,
		size,
		division,
		blocksize,
		matrix,
		image.NewRGBA(image.Rect(0, 0, size, size)),
		name}

	fid.draw()

	return fid
}

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

func (m *Marker) draw() {
	codeX := m.Block_size * 2
	codeY := m.Block_size * 2

	m.set(0, 0, m.Size, m.Size, color.NRGBA{0, 0, 0, 255})
	m.set(m.Block_size, m.Block_size, m.Size-m.Block_size, m.Size-m.Block_size, color.NRGBA{255, 255, 255, 255})
	m.set(m.Block_size*2, m.Block_size*2, m.Size-m.Block_size*2, m.Size-m.Block_size*2, color.NRGBA{0, 0, 0, 255})

	for i, r := 0, 0; i < len(m.matrix); i++ {
		if m.matrix[i] == 1 {
			m.set(codeX, codeY, codeX+m.Block_size, codeY+m.Block_size, color.NRGBA{255, 255, 255, 255})
		}

		if r == m.Division-1 {
			r = 0
			codeY = codeY + m.Block_size
			codeX = m.Block_size * 2
			continue
		}

		codeX = codeX + m.Block_size
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
