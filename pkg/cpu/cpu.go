package cpu

const (
	_ = iota
	NONE
	NMI
	RESET
	IRQ
	BRK
)

type CPU struct {
	Registers
	Interrupt int
	MemoryMap [0xffff]byte
}

// レジスタ内容に関しては http://hp.vector.co.jp/authors/VA042397/nes/6502.html を参照
type Registers struct {
	A  uint8  // アキュムレーター
	X  uint8  // インデックスレジスタ
	Y  uint8  // インデックスレジスタ
	S  uint16 // スタックポインタ 上位8bitは0x01で固定 7:N 6:V 5:R=1 4:B 3:D 2:I 1:Z 0:C
	P  uint8  // ステータスレジスタ
	PC uint16 // プログラムカウンタ
}

func NewCPU() *CPU {
	return &CPU{
		Registers: *NewRegisters(),
		Interrupt: NONE,
		MemoryMap: [0xffff]byte{},
	}
}

func NewRegisters() *Registers {
	return &Registers{
		A:  0x00,
		X:  0x00,
		Y:  0x00,
		S:  0x01FD,
		P:  0 | (1 << 5) | (1 << 4) | (1 << 2),
		PC: 0x8000,
	}
}

func (c *CPU) Push(v uint8) {
	// 本来 0x0100以下には入らない
	if c.Registers.S < 0x0100 {
		// TODO: なんか処理をかく
		return
	}
	c.MemoryMap[c.Registers.S] = v
	c.Registers.S--
}

func (c *CPU) Pop() uint8 {
	c.Registers.S++
	return c.MemoryMap[c.Registers.S]
}
