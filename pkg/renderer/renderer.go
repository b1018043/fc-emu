package renderer

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WINDOW_WIDTH  = 256
	WINDOW_HEIGHT = 224
)

type EbitenRenderer struct {
}

func (e *EbitenRenderer) renderSprite(screen *ebiten.Image, sprite [][]byte, spriteNum int) {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
		}
	}
}
