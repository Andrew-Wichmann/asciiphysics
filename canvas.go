package asciiphysics

import (
	"fmt"
	"image"
	"strings"
	"time"

	imgManip "github.com/TheZoraiz/ascii-image-converter/image_manipulation"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fogleman/gg"
	"github.com/qeesung/image2ascii/convert"
)

const (
	fps = 30
)

func flattenAscii(asciiSet [][]imgManip.AsciiChar) string {
	var ascii []string

	for _, line := range asciiSet {
		var tempAscii string

		for _, char := range line {
			tempAscii += char.OriginalColor
		}

		ascii = append(ascii, tempAscii)
	}
	result := strings.Join(ascii, "\n")

	return result
}

type canvasTick struct{}
type Drawable interface {
	Tick() Drawable
	Draw(*gg.Context)
}

func newTick() tea.Cmd {
	return tea.Tick(time.Second/fps, func(time.Time) tea.Msg {
		return canvasTick{}
	})
}

type Canvas struct {
	drawable       []Drawable
	width, height  int
	asciiConverter *convert.ImageConverter
}

func NewCanvas(width, height int) Canvas {
	converter := convert.NewImageConverter()
	return Canvas{
		width:          width,
		height:         height,
		asciiConverter: converter,
	}
}

func (c Canvas) Init() tea.Cmd {
	return newTick()
}

func (c Canvas) View() string {
	i := image.NewRGBA(image.Rect(0, 0, c.width, c.height))
	ctx := gg.NewContextForImage(i)
	for _, drawable := range c.drawable {
		drawable.Draw(ctx)
	}
	dimensions := []int{c.width, c.height}
	imgSet, err := imgManip.ConvertToAsciiPixels(ctx.Image(), dimensions, c.width, c.height, false, false, false, false, true)
	if err != nil {
		return fmt.Sprintf("Error rendering canvas: %s", err.Error())
	}
	asciiSet, err := imgManip.ConvertToAsciiChars(imgSet, false, true, false, true, false, "", [3]int{0, 0, 0})
	if err != nil {
		return fmt.Sprintf("Error rendering canvas: %s", err.Error())
	}
	return flattenAscii(asciiSet)
}

func (c *Canvas) AddDrawable(drawable Drawable) {
	c.drawable = append(c.drawable, drawable)
}

func (c Canvas) Update(msg tea.Msg) (Canvas, tea.Cmd) {
	if _, ok := msg.(canvasTick); ok {
		for i, circle := range c.drawable {
			c.drawable[i] = circle.Tick()
		}
		return c, newTick()
	}
	return c, nil
}
