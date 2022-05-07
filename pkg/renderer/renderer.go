package renderer

import (
	"github.com/b1018043/fc-emu/pkg/cpu"
	"github.com/b1018043/fc-emu/pkg/ppu"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WINDOW_WIDTH  = 256
	WINDOW_HEIGHT = 224
)

type EbitenRenderer struct {
	pixels []byte
	CPU    cpu.CPU
	PPU    ppu.PPU
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
	for {
		cycle := 0
		cycle += e.CPU.Run()
		if ok := e.PPU.Run(cycle * 3); ok {
			e.updatePixels()
			break
		}
	}
	return nil
}

func (e *EbitenRenderer) updatePixels() {
	// for i, v := range e.PPU.Background {
	// }
}

func (e *EbitenRenderer) Draw(screen *ebiten.Image) {
	if e.pixels == nil {
		e.pixels = make([]byte, WINDOW_WIDTH*WINDOW_HEIGHT*4)
	}
	screen.ReplacePixels(e.pixels)
}

var _ ebiten.Game = &EbitenRenderer{}
