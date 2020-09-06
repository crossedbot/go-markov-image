package markov

import (
	"crypto/rand"
	"fmt"
	"image"
	"image/color"
	// "image/gif"
	_ "image/gif"
	// "image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"math/big"
	"os"

	"github.com/crossedbot/collections/randomstack"
)

// Markov is an interface to a markov model of an image
type Markov interface {
	// GetNextColor returns a randomly selected tranistion color for the given
	// color.
	GetNextColor(c color.Color) color.Color

	// GetRandomColor returns a randomly selected color from the model.
	GetRandomColor() color.Color

	// AddColorTransition adds a color transition, c2, for the given color, c1,
	// to the model.
	AddColorTransition(c1 color.Color, c2 color.Color)

	// Generate returns a new image generated from the current model.
	Generate() *image.RGBA

	// ReadFile reads the given image file and sets the model accordingly.
	ReadFile(filename string) error

	// WriteFile generates a new image based on the current model and writes
	// it to the given file path.
	WriteFile(filename string) error
}

// stateSpace represents a markov state space of an image's colors
type stateSpace map[uint32][]uint32

// markov represents a markov model of an image
type markov struct {
	keys   []uint32        // encoded color keys of all distinct colors in the image
	model  stateSpace      // the state space of the image's colors
	format string          // the image's file format
	bounds image.Rectangle // the boundaries of the image
}

// adjacent are the relative difference between a given point and its directly
// adjacent points. IE. top, bottom, left, right.
var adjacent = []image.Point{
	image.Point{X: -1, Y: 0},
	image.Point{X: 0, Y: -1},
	image.Point{X: 1, Y: 0},
	image.Point{X: 0, Y: 1},
}

// New returns a new Markov instance.
func New() Markov {
	return &markov{model: make(stateSpace)}
}

// GetNextColor returns a randomly selected tranistion color for the given
// color.
func (m *markov) GetNextColor(c color.Color) color.Color {
	key := encodeColor(c)
	if values, ok := m.model[key]; ok {
		i, _ := rand.Int(rand.Reader, big.NewInt(int64(len(values))))
		return decodeColor(values[int(i.Int64())])
	}
	return nil
}

// GetRandomColor returns a randomly selected color from the model.
func (m *markov) GetRandomColor() color.Color {
	i, _ := rand.Int(rand.Reader, big.NewInt(int64(len(m.keys))))
	return decodeColor(m.keys[int(i.Int64())])
}

// AddColorTransition adds a color transition, c2, for the given color, c1, to
// the model.
func (m *markov) AddColorTransition(c1 color.Color, c2 color.Color) {
	key1 := encodeColor(c1)
	key2 := encodeColor(c2)
	if _, ok := m.model[key1]; !ok {
		m.keys = append(m.keys, key1)
	}
	m.model[key1] = append(m.model[key1], key2)
}

// Generate returns a new image generated from the current model.
func (m *markov) Generate() *image.RGBA {
	im := image.NewRGBA(image.Rect(m.MinX(), m.MinY(), m.MaxX(), m.MaxY()))
	stack := randomstack.New()
	x, _ := rand.Int(rand.Reader, big.NewInt(int64(m.MaxX())))
	y, _ := rand.Int(rand.Reader, big.NewInt(int64(m.MaxY())))
	p := image.Point{
		X: int(x.Int64()),
		Y: int(y.Int64()),
	}
	// add psuedo-random starting point to the new image
	c := m.GetRandomColor()
	im.SetRGBA(p.X, p.Y, colorToRGBA(c))
	stack.Push(p)
	// for each colored point in the stack, get its color, and set all adjacent
	// points to a new color
	for stack.Len() > 0 {
		// pop a randomly selected colored point
		p = stack.Pop().(image.Point)
		c = im.At(p.X, p.Y)
		for _, adj := range adjacent {
			p_ := p.Add(adj)
			if p_.X >= m.MinX() && p_.X < m.MaxX() &&
				p_.Y >= m.MinY() && p_.Y < m.MaxY() {
				if im.Pix[im.PixOffset(p_.X, p_.Y)] == 0 {
					// if the adjacent point fits within the image boundaries
					// and has not been set a color value: get the next color,
					// set the point's color, and add it to the stack for later
					// processing
					c = m.GetNextColor(c)
					im.SetRGBA(p_.X, p_.Y, colorToRGBA(c))
					stack.Push(p_)
				}
			}
		}
	}
	return im
}

// ReadFile reads the given image file and sets the model accordingly.
func (m *markov) ReadFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	im, format, err := image.Decode(f)
	if err != nil {
		return err
	}
	// TODO removes this check once more formats are supported
	if format != "png" {
		return fmt.Errorf(
			"file format \"%s\" not supported; supported formats are: \"png\"",
			format,
		)
	}
	m.format = format
	m.bounds = im.Bounds()
	// for each pixel sample its color and add its color tranistions to the
	// model
	for x := m.MinX(); x < m.MaxX(); x++ {
		for y := m.MinY(); y < m.MaxY(); y++ {
			p := image.Point{X: x, Y: y}
			c := im.At(p.X, p.Y)
			for _, adj := range adjacent {
				p := p.Add(adj)
				if p.X >= m.MinX() && p.X < m.MaxX() &&
					p.Y >= m.MinY() && p.Y < m.MaxY() {
					c_ := im.At(p.X, p.Y)
					m.AddColorTransition(c, c_)
				}
			}
		}
	}
	return nil
}

// WriteFile generates a new image based on the current model and writes it to
// the given file path.
func (m *markov) WriteFile(filename string) error {
	o, err := os.Create(filename)
	if err != nil {
		return err
	}
	d := m.Generate()
	switch m.format {
	case "png":
		png.Encode(o, d)
	// TODO readds these once they are supported... sorry :(
	// case "jpeg":
	//	jpeg.Encode(o, d, nil)
	// case "gif":
	//	gif.Encode(o, d, nil)
	default:
		return fmt.Errorf(
			"file format \"%s\" not supported; supported formats are: \"png\"",
			m.format,
		)
	}
	return nil
}

// MinX returns the lower X coordinate boundary of the image.
func (m *markov) MinX() int {
	return m.bounds.Min.X
}

// MinY returns the lower Y coordinate boundary of the image.
func (m *markov) MinY() int {
	return m.bounds.Min.Y
}

// MaxX returns the upper X coordinate boundary of the image.
func (m *markov) MaxX() int {
	return m.bounds.Max.X
}

// MaxY returns the upper Y coordinate boundary of the image.
func (m *markov) MaxY() int {
	return m.bounds.Max.Y
}
