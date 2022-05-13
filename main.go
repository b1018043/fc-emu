package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/b1018043/fc-emu/pkg/cpu"
	"github.com/b1018043/fc-emu/pkg/logger"
	"github.com/b1018043/fc-emu/pkg/ppu"
	"github.com/b1018043/fc-emu/pkg/renderer"
	"github.com/b1018043/fc-emu/pkg/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	var (
		filename = flag.String("soft", "", "use soft name")
		d        = flag.Bool("d", false, "use debug(default false)")
	)
	flag.Parse()
	fmt.Println(*filename, *d)
	if *filename == "" {
		fmt.Println("invalid nes file")
		os.Exit(1)
	}
	if *d {
		logger.IsDebugMode = true
	}
	progROM, charROM, err := utils.LoadFCROM(*filename)
	if err != nil {
		log.Fatalln(err)
	}
	cpu := cpu.NewCPU()
	cpu.SetPRGROM(progROM)
	ppu := ppu.NewPPU(charROM)
	cpu.PPU = ppu
	game := renderer.NewEbitenRenderer(cpu, ppu)
	game.CPU.Reset()
	ebiten.SetWindowSize(renderer.WINDOW_WIDTH*2, renderer.WINDOW_HEIGHT*2)
	ebiten.SetWindowTitle("FC-emu")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatalln(err)
	}
}
