package main

import (
	"log"

	"github.com/b1018043/fc-emu/pkg/cpu"
	"github.com/b1018043/fc-emu/pkg/ppu"
	"github.com/b1018043/fc-emu/pkg/renderer"
	"github.com/b1018043/fc-emu/pkg/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	progROM, charROM, err := utils.LoadFCROM("./static/roms/sample1.nes")
	if err != nil {
		log.Fatalln(err)
	}
	cpu := cpu.NewCPU()
	cpu.SetPRGROM(progROM)
	ppu := ppu.NewPPU(charROM)
	game := renderer.NewEbitenRenderer(cpu, ppu)
	ebiten.SetWindowSize(renderer.WINDOW_WIDTH, renderer.WINDOW_HEIGHT)
	ebiten.SetWindowTitle("FC-emu")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatalln(err)
	}
}
