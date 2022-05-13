package ppu

const (
	PPU_PATTERN_TABLE0   = 0x0000
	PPU_PATTERN_TABLE1   = 0x1000
	PPU_NAME_TABLE0      = 0x2000
	PPU_ATTR_TABLE0      = 0x23C0
	PPU_NAME_TABLE1      = 0x2400
	PPU_ATTR_TABLE1      = 0x27C0
	PPU_NAME_TABLE2      = 0x2800
	PPU_ATTR_TABLE2      = 0x2BC0
	PPU_NAME_TABLE3      = 0x2C00
	PPU_ATTR_TABLE3      = 0x2FC0
	PPU_MIRROR_NE_TABLES = 0x3000
	PPU_BG_PALLET        = 0x3F00
	PPU_SPRITE_PALLET    = 0x3F10
	PPU_MIRROR_PALLETS   = 0x3F20
)

type BackgroundContent struct {
	Tile      [][]byte
	PaletteID int
}

type PPU struct {
	Registers     []byte // 0x2000~0x2007
	MemoryMap     [0x3FFF + 1]byte
	Cycle         int
	Line          int
	charROM       []byte
	Background    []BackgroundContent
	addressBuffer []byte
	PaletteRAM    []byte
}

func NewPPU(charROM []byte) *PPU {
	return &PPU{
		MemoryMap:     [0x3fff + 1]byte{},
		Cycle:         0,
		Line:          0,
		charROM:       charROM,
		Background:    make([]BackgroundContent, 30*32),
		addressBuffer: make([]byte, 0, 2),
	}
}

func (p *PPU) Run(cycle int) bool {
	p.Cycle += cycle
	if p.Line == 0 {
		p.Background = p.Background[:0]
	}

	if p.Cycle >= 341 {
		p.Cycle -= 341

		if p.Line <= 240 && p.Line%8 == 0 {
			p.buildBackground()
		}
		p.Line++
		if p.Line == 262 {
			p.Line = 0
			p.PaletteRAM = p.MemoryMap[PPU_BG_PALLET:]
			return true
		}
	}
	return false
}

// タイルの座標からどのスプライトを表示すれば良いか判断する
func (p *PPU) getSpriteID(tileX, tileY uint16) byte {
	// TODO: 現状ではname table 0 だけ対応
	return p.MemoryMap[PPU_NAME_TABLE0+tileY*32+tileX]
}

/*  _ _
 * |0|1|
 * |_|_|
 * |2|3|
 * |_|_|
 */

// どのブロックに属しているかを判定している
func (p *PPU) getBlockID(tileX, tileY uint16) byte {
	return (byte(tileX%4) / 2) + (byte(tileY%4)/2)*2
}

func (p *PPU) buildSprite(spriteNum int) [][]byte {
	sprite := make([][]byte, 8)
	for i := 0; i < 8; i++ {
		sprite[i] = make([]byte, 8)
	}

	for i := 0; i < 16; i++ {
		for j := 0; j < 8; j++ {
			if p.charROM[spriteNum*16+i]&(0x80>>j) != 0 {
				sprite[i%8][j] += 0x01 << (i / 8)
			}
		}
	}
	return sprite
}

func (p *PPU) getAttr(tileX, tileY uint16) byte {
	address := PPU_ATTR_TABLE0 + tileX/4 + (tileY/4)*8
	return p.MemoryMap[address]
}

func (p *PPU) buildTile(tileX, tileY uint16) ([][]byte, int) {
	blockID := p.getBlockID(tileX, tileY)
	spriteID := p.getSpriteID(tileX, tileY)
	attr := p.getAttr(tileX, tileY)
	paletteID := (attr >> (blockID * 2)) & 0x03
	sprite := p.buildSprite(int(spriteID))
	return sprite, int(paletteID)
}

func (p *PPU) buildBackground() {
	tileY := (p.Line / 8) % 30
	for x := 0; x < 32; x++ {
		tileX := x % 32
		tile, palette := p.buildTile(uint16(tileX), uint16(tileY))
		p.Background = append(p.Background, BackgroundContent{
			Tile:      tile,
			PaletteID: palette,
		})
	}
}

func (p *PPU) setAddress(addr uint16) {
	p.addressBuffer[0] = byte(addr >> 8)
	p.addressBuffer[1] = byte(addr)
}

func (p *PPU) getAddress() uint16 {
	addr := uint16(p.addressBuffer[0])<<8 | uint16(p.addressBuffer[1])
	return addr
}

func (p *PPU) GetData() byte {
	v := p.MemoryMap[p.getAddress()]
	p.setAddress(p.getAddress() + 1)
	return v
}

func (p *PPU) SetData(val byte) {
	p.MemoryMap[p.getAddress()] = val
	p.setAddress(p.getAddress() + 1)
}

func (p *PPU) SetAddress(addr byte) {
	if len(p.addressBuffer) >= 2 {
		p.addressBuffer = p.addressBuffer[:0]
	}
	p.addressBuffer = append(p.addressBuffer, addr)
}
