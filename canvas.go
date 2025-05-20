package asciiphysics

import (
	"fmt"
	"image"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fogleman/gg"
	"github.com/qeesung/image2ascii/convert"
)

const (
	fps = 24
)

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
	image          image.Image
	fps            int64
	frameCount     int64
	start          int64
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
	return lipgloss.JoinVertical(lipgloss.Top, c.asciiConverter.Image2ASCIIString(ctx.Image(), &convert.DefaultOptions), fmt.Sprintf("%d", c.fps))
}

func (c *Canvas) AddDrawable(drawable Drawable) {
	c.drawable = append(c.drawable, drawable)
}

func (c Canvas) Update(msg tea.Msg) (Canvas, tea.Cmd) {
	if c.start == 0 {
		c.start = time.Now().Unix() - 1
	}
	if _, ok := msg.(canvasTick); ok {
		for i, circle := range c.drawable {
			c.drawable[i] = circle.Tick()
		}
		c.frameCount += 1
		c.fps = c.frameCount / (time.Now().Unix() - c.start)
		return c, newTick()
	}
	return c, nil
}
