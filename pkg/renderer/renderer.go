package renderer

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WINDOW_WIDTH  = 256
	WINDOW_HEIGHT = 224
)

type EbitenRenderer struct {
	pixels []byte
}

func (e *EbitenRenderer) renderSprite(screen *ebiten.Image, sprite [][]byte, spriteNum int) {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
		}
	}
}

func (e *EbitenRenderer) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WINDOW_WIDTH, WINDOW_HEIGHT
}

func (e *EbitenRenderer) Update() error {
	return nil
}

func (e *EbitenRenderer) Draw(screen *ebiten.Image) {
	if e.pixels == nil {
		e.pixels = make([]byte, WINDOW_WIDTH*WINDOW_HEIGHT*4)
	}
	screen.ReplacePixels(e.pixels)
}

var _ ebiten.Game = &EbitenRenderer{}
