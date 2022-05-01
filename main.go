package main

import (
	"image/color"
	"log"

	"github.com/b1018043/fc-emu/pkg/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	PIXEL_RATIO          = 1
	DEFAULT_WINDOW_WIDTH = 800
	CHAR_ROM_SIZE        = 0x2000
)

type Game struct {
	charROM      []byte
	height       int
	width        int
	spritePerRow int
	spriteNum    int
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x00, 0, 0, 0xff})
	for i := 0; i < g.spriteNum; i++ {
		sprite := g.buildSprite(i)
		g.renderSprite(screen, sprite, i)
	}
}

func (g *Game) buildSprite(spriteNum int) [][]byte {
	sprite := make([][]byte, 8)
	for i := 0; i < 8; i++ {
		sprite[i] = make([]byte, 8)
	}
	// NOTE: 本来は16までいける？
	for i := 0; i < 15; i++ {
		for j := 0; j < 8; j++ {
			if g.charROM[spriteNum*16+i]&(0x80>>j) != 0 {
				sprite[i%8][j] += 0x01 << (i / 8)
			}
		}
	}
	return sprite
}

func (g *Game) renderSprite(screen *ebiten.Image, sprite [][]byte, spriteNum int) {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			x := (j + (spriteNum%g.spritePerRow)*8) * PIXEL_RATIO
			y := (i + (spriteNum/g.spritePerRow)*8) * PIXEL_RATIO
			colorVal := 85 * sprite[i][j]
			ebitenutil.DrawRect(screen, float64(x), float64(y), PIXEL_RATIO, PIXEL_RATIO, color.RGBA{
				colorVal, colorVal, colorVal, 0xff,
			})
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.width, g.height
}

func (g *Game) Update() error {
	return nil
}

var _ ebiten.Game = &Game{}

func main() {
	filename := "./static/roms/sample1.nes"
	pages, _, charROM, err := utils.LoadFCROM(filename)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("length: %d\n", len(charROM))

	spritePerRow := DEFAULT_WINDOW_WIDTH / (8 * PIXEL_RATIO)
	spriteNum := utils.CHAR_ROM_PAGE_SIZE * pages / 16
	rowNum := (spriteNum / spritePerRow) + 1
	height := rowNum * 8 * PIXEL_RATIO

	game := &Game{charROM: charROM, height: height, width: DEFAULT_WINDOW_WIDTH, spritePerRow: spritePerRow, spriteNum: spriteNum}
	ebiten.SetWindowSize(game.width, game.height)
	ebiten.SetWindowTitle("Sprite")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatalln(err)
	}
}
