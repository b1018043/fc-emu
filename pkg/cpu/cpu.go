package cpu

// Interrupt type
const (
	_ = iota
	NONE
	NMI
	RESET
	IRQ
	BRK
)

const (
	_ = iota
	modeImplied
	modeAccumulator
	modeImmediate
	modeZeroPage
	modeZeroPageX
	modeZeroPageY
	modeAbsolute
	modeAbsoluteX
	modeAbsoluteY
	modeRelative
	modeIndirect
	modeIndexedIndirect
	modeIndirectIndexed
	modeIndirectX = modeIndexedIndirect
	modeIndirectY = modeIndirectIndexed
)

// TODO:命令がNONEの部分は変えていく
var operation_names = [256]string{
	"BRK", "ORA", "NONE", "NONE", "NONE", "ORA", "ASL", "NONE",
	"PHP", "ORA", "ASL", "NONE", "NONE", "ORA", "ASL", "NONE",
}

var operation_sizes = [256]int{
	1, 2, 0, 0, 0, 2, 2, 0,
	1, 2, 1, 0, 0, 3, 3, 0,
}

var operation_modes = [256]int{
	modeImplied, modeIndirectX, 0, 0, 0, modeZeroPage, modeZeroPage, 0,
	modeImplied, modeImmediate, modeAccumulator, 0, 0, modeAbsolute, modeAbsolute, 0,
}

var operation_cycles = [256]int{
	7, 6, 0, 0, 0, 3, 5, 0,
	3, 2, 2, 0, 0, 4, 6, 0,
}

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

func (c *CPU) Reset() {
	c.Registers = Registers{
		A:  0x00,
		X:  0x00,
		Y:  0x00,
		S:  0x01FD,
		P:  0 | (1 << 5) | (1 << 4) | (1 << 2),
		PC: 0xFFFC,
	}
}
