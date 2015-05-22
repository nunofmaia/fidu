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

// Marker represents the structure of a fiducial marker
type Marker struct {
	Code      int
	Size      int
	Division  int
	BlockSize int
	matrix    []int
	Marker    *image.RGBA
	Name      string
}

// New returns a Marker with the given parameters.
// Using these parameters, it generates the visual representation of the marker
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

func (m *Marker) draw() {
	codeX := m.BlockSize * 2
	codeY := m.BlockSize * 2

	m.set(0, 0, m.Size, m.Size, color.NRGBA{0, 0, 0, 255})
	m.set(m.BlockSize, m.BlockSize, m.Size-m.BlockSize, m.Size-m.BlockSize, color.NRGBA{255, 255, 255, 255})
	m.set(m.BlockSize*2, m.BlockSize*2, m.Size-m.BlockSize*2, m.Size-m.BlockSize*2, color.NRGBA{0, 0, 0, 255})

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
