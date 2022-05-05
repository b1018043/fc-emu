package ppu

const (
	PPU_PATTERN_TABLE0   = 0x0000
	PPU_PATTERN_TABLE1   = 0x1000
	PPU_NAME_TABLE0      = 0x2000
	PPU_ELEMENT_TABLE0   = 0x23C0
	PPU_NAME_TABLE1      = 0x2400
	PPU_ELEMENT_TABLE1   = 0x27C0
	PPU_NAME_TABLE2      = 0x2800
	PPU_ELEMENT_TABLE2   = 0x2BC0
	PPU_NAME_TABLE3      = 0x2C00
	PPU_ELEMENT_TABLE3   = 0x2FC0
	PPU_MIRROR_NE_TABLES = 0x3000
	PPU_BG_PALLET        = 0x3F00
	PPU_SPRITE_PALLET    = 0x3F10
	PPU_MIRROR_PALLETS   = 0x3F20
)

type PPU struct {
	// Registers []byte // 0x2000~0x2007
	MemoryMap [0x3FFF]byte
	Cycle     int
	Line      int
}

func (p *PPU) Run(cycle int) {
	p.Cycle += cycle
}
