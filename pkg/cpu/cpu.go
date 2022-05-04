package cpu

import "log"

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
	"BRK", "ORA", "NONE", "NONE", "NONE", "ORA", "ASL", "NONE", "PHP", "ORA", "ASL", "NONE", "NONE", "ORA", "ASL", "NONE",
	"BPL", "ORA", "NONE", "NONE", "NONE", "ORA", "ASL", "NONE", "CLC", "ORA", "NONE", "NONE", "NONE", "ORA", "ASL", "NONE",
	"JSR", "AND", "NONE", "NONE", "BIT", "AND", "ROL", "NONE", "PLP", "AND", "POL", "NONE", "BIT", "AND", "ROL", "NONE",
	"BMI", "AND", "NONE", "NONE", "NONE", "AND", "ROL", "NONE", "SEC", "AND", "NONE", "NONE", "NONE", "AND", "ROL", "NONE",
	"RTI", "EOR", "NONE", "NONE", "NONE", "EOR", "LSR", "NONE", "PHA", "EOR", "LSR", "NONE", "JMP", "EOR", "LSR", "NONE",
	"BVC", "EOR", "NONE", "NONE", "NONE", "EOR", "LSR", "NONE", "CLI", "EOR", "NONE", "NONE", "NONE", "EOR", "LSR", "NONE",
	"RTS", "ADC", "NONE", "NONE", "NONE", "ADC", "ROR", "NONE", "PLA", "ADC", "ROR", "NONE", "JMP", "ADC", "ROR", "NONE",
	"BVS", "ADC", "NONE", "NONE", "NONE", "ADC", "ROR", "NONE", "SEI", "ADC", "NONE", "NONE", "NONE", "ADC", "ROR", "NONE",
	"NONE", "STA", "NONE", "NONE", "STY", "STA", "STX", "NONE", "DEY", "NONE", "TXA", "NONE", "STY", "STA", "STX", "NONE",
	"BCC", "STA", "NONE", "NONE", "STY", "STA", "STX", "NONE", "TYA", "STA", "TXS", "NONE", "NONE", "STA", "NONE", "NONE",
	"LDY", "LDA", "LDX", "NONE", "LDY", "LDA", "LDX", "NONE", "TAY", "LDA", "TAX", "NONE", "LDY", "LDA", "LDX", "NONE",
	"BCS", "LDA", "NONE", "NONE", "LDY", "LDA", "LDX", "NONE", "CLV", "LDA", "TSX", "NONE", "LDY", "LDA", "LDX", "NONE",
	"CPY", "CMP", "NONE", "NONE", "CPY", "CMP", "DEC", "NONE", "INY", "CMP", "DEX", "NONE", "CPY", "CMP", "DEC", "NONE",
	"BNE", "CMP", "NONE", "NONE", "NONE", "CMP", "DEC", "NONE", "CLD", "CMP", "NONE", "NONE", "NONE", "CMP", "DEC", "NONE",
	"CPX", "SBC", "NONE", "NONE", "CPX", "SBC", "INC", "NONE", "INX", "SBC", "NOP", "NONE", "CPX", "SBC", "INC", "NONE",
	"BEQ", "SBC", "NONE", "NONE", "NONE", "SBC", "INC", "NONE", "SED", "SBC", "NONE", "NONE", "NONE", "SBC", "INC", "NONE",
}

var operation_sizes = [256]int{
	1, 2, 0, 0, 0, 2, 2, 0, 1, 2, 1, 0, 0, 3, 3, 0,
	2, 2, 0, 0, 0, 2, 2, 0, 1, 3, 0, 0, 0, 3, 3, 0,
	3, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 0, 2, 2, 0, 1, 3, 0, 0, 0, 3, 3, 0,
	1, 2, 0, 0, 0, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 0, 2, 2, 0, 1, 3, 0, 0, 0, 3, 3, 0,
	1, 2, 0, 0, 0, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 0, 2, 2, 0, 1, 3, 0, 0, 0, 3, 3, 0,
	0, 2, 0, 0, 2, 2, 2, 0, 1, 0, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 0, 3, 0, 0,
	2, 2, 2, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 3, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 0, 2, 2, 0, 1, 3, 0, 0, 0, 3, 3, 0,
	2, 2, 0, 0, 2, 2, 2, 0, 1, 2, 1, 0, 3, 3, 3, 0,
	2, 2, 0, 0, 0, 2, 2, 0, 1, 3, 0, 0, 0, 3, 3, 0,
}

var operation_modes = [256]int{
	modeImplied, modeIndirectX, 0, 0, 0, modeZeroPage, modeZeroPage, 0, modeImplied, modeImmediate, modeAccumulator, 0, 0, modeAbsolute, modeAbsolute, 0,
	modeRelative, modeIndirectY, 0, 0, 0, modeZeroPageX, modeZeroPageX, 0, modeImplied, modeAbsoluteX, 0, 0, 0, modeAbsoluteX, modeAbsoluteX, 0,
	modeAbsolute, modeIndirectX, 0, 0, modeZeroPage, modeZeroPage, modeZeroPage, 0, modeImplied, modeImmediate, modeAccumulator, 0, modeAbsolute, modeAbsolute, modeAbsolute, 0,
	modeRelative, modeIndirectY, 0, 0, 0, modeZeroPageX, modeZeroPageX, 0, modeImplied, modeAbsoluteY, 0, 0, 0, modeAbsoluteX, modeAbsoluteX, 0,
	modeImplied, modeIndirectX, 0, 0, 0, modeZeroPage, modeZeroPage, 0, modeImplied, modeImmediate, modeAccumulator, 0, modeAbsolute, modeAbsolute, modeAbsolute, 0,
	modeRelative, modeIndirectY, 0, 0, 0, modeZeroPageX, modeZeroPageX, 0, modeImplied, modeAbsoluteY, 0, 0, 0, modeAbsoluteX, modeAbsoluteX, 0,
	modeImplied, modeIndirectX, 0, 0, 0, modeZeroPage, modeZeroPage, 0, modeImplied, modeImmediate, modeAccumulator, 0, modeIndirect, modeAbsolute, modeAbsolute, 0,
	modeRelative, modeIndirectY, 0, 0, 0, modeZeroPageX, modeZeroPageX, 0, modeImplied, modeAbsoluteY, 0, 0, 0, modeAbsoluteX, modeAbsoluteX, 0,
	0, modeIndirectX, 0, 0, modeZeroPage, modeZeroPage, modeZeroPage, 0, modeImplied, 0, modeImplied, 0, modeAbsolute, modeAbsolute, modeAbsolute, 0,
	modeRelative, modeIndirectY, 0, 0, modeZeroPageX, modeZeroPageX, modeZeroPageY, 0, modeImplied, modeAbsoluteY, modeImplied, 0, 0, modeAbsoluteX, 0, 0,
	modeImmediate, modeIndirectX, modeImmediate, 0, modeZeroPage, modeZeroPage, modeZeroPage, 0, modeImplied, modeImmediate, modeImplied, 0, modeAbsolute, modeAbsolute, modeAbsolute, 0,
	modeRelative, modeIndirectY, 0, 0, modeZeroPageX, modeZeroPageX, modeZeroPageY, 0, modeImplied, modeAbsoluteY, modeImplied, 0, modeAbsoluteX, modeAbsoluteX, modeAbsoluteY, 0,
	modeImmediate, modeIndirectX, 0, 0, modeZeroPage, modeZeroPage, modeZeroPage, 0, modeImplied, modeImmediate, modeImplied, 0, modeAbsolute, modeAbsolute, modeAbsolute, 0,
	modeRelative, modeIndirectY, 0, 0, 0, modeZeroPageX, modeZeroPageX, 0, modeImplied, modeAbsoluteY, 0, 0, 0, modeAbsoluteX, modeAbsoluteX, 0,
	modeImmediate, modeIndirectX, 0, 0, modeZeroPage, modeZeroPage, modeZeroPage, 0, modeImplied, modeImmediate, modeImplied, 0, modeAbsolute, modeAbsolute, modeAbsolute, 0,
	modeRelative, modeIndirectY, 0, 0, 0, modeZeroPageX, modeZeroPageX, 0, modeImplied, modeAbsoluteY, 0, 0, 0, modeAbsoluteX, modeAbsoluteX, 0,
}

var operation_cycles = [256]int{
	7, 6, 0, 0, 0, 3, 5, 0, 3, 2, 2, 0, 0, 4, 6, 0,
	2, 5, 0, 0, 0, 4, 6, 0, 2, 4, 0, 0, 0, 4, 7, 0,
	6, 6, 0, 0, 3, 3, 5, 0, 4, 2, 2, 0, 4, 4, 6, 0,
	2, 5, 0, 0, 0, 4, 6, 0, 2, 4, 0, 0, 0, 4, 7, 0,
	6, 6, 0, 0, 0, 3, 5, 0, 3, 2, 2, 0, 3, 4, 6, 0,
	2, 5, 0, 0, 0, 4, 6, 0, 2, 4, 0, 0, 0, 4, 7, 0,
	6, 6, 0, 0, 0, 3, 5, 0, 4, 2, 2, 0, 5, 4, 6, 0,
	2, 5, 0, 0, 0, 4, 6, 0, 2, 4, 0, 0, 0, 4, 7, 0,
	0, 6, 0, 0, 3, 3, 3, 0, 2, 0, 2, 0, 4, 4, 4, 0,
	2, 6, 0, 0, 4, 4, 4, 0, 2, 5, 2, 0, 0, 5, 0, 0,
	2, 6, 2, 0, 3, 3, 3, 0, 2, 2, 2, 0, 4, 4, 4, 0,
	2, 5, 0, 0, 4, 4, 4, 0, 2, 4, 2, 0, 4, 4, 4, 0,
	2, 6, 0, 0, 3, 3, 5, 0, 2, 2, 2, 0, 4, 4, 6, 0,
	2, 5, 0, 0, 0, 4, 6, 0, 2, 4, 0, 0, 0, 4, 7, 0,
	2, 6, 0, 0, 3, 3, 5, 0, 2, 2, 2, 0, 4, 4, 6, 0,
	2, 5, 0, 0, 0, 4, 6, 0, 2, 4, 0, 0, 0, 4, 7, 0,
}

type CPU struct {
	Registers
	Interrupt int
	MemoryMap [0xffff]byte
}

// レジスタ内容に関しては http://hp.vector.co.jp/authors/VA042397/nes/6502.html を参照
type Registers struct {
	A  byte           // アキュムレーター
	X  byte           // インデックスレジスタ
	Y  byte           // インデックスレジスタ
	S  uint16         // スタックポインタ
	P  statusRegister // ステータスレジスタ 上位8bitは0x01で固定 7:N 6:V 5:R=1 4:B 3:D 2:I 1:Z 0:C
	PC uint16         // プログラムカウンタ
}

type statusRegister struct {
	N bool // 演算結果がマイナス(bit7=1)の時にセット
	V bool // オーバーフロー時にセット
	R bool // 予約済み 常にtrue
	B bool // ブレークモード BRK発生時にtrue,IRQ発生時にfalseにセット
	D bool // D とりあえず falseにする
	I bool // false IRQ許可 true IRQ禁止
	Z bool // 演算結果が0の時にtrue
	C bool // キャリー発生時にtrue
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
		A: 0x00,
		X: 0x00,
		Y: 0x00,
		S: 0x01FD,
		P: statusRegister{
			N: false,
			V: false,
			R: true,
			B: true,
			D: false,
			I: true,
			Z: false,
			C: false,
		},
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
		A: 0x00,
		X: 0x00,
		Y: 0x00,
		S: 0x01FD,
		P: statusRegister{
			N: false,
			V: false,
			R: true,
			B: true,
			D: false,
			I: true,
			Z: false,
			C: false,
		},
		PC: 0xFFFC,
	}
}

func (c *CPU) exec(opecode int) {
	switch operation_names[opecode] {
	case "LDA":
	case "LDX":
	case "LDY":
	case "STA":
	case "STX":
	case "STY":
	case "TAX":
	case "TAY":
	case "TXS":
	case "TYA":
	case "ADC":
	case "AND":
	case "ASL":
	case "BIT":
	case "CMP":
	case "CPX":
	case "CPY":
	case "DEC":
	case "DEX":
	case "DEY":
	case "EOR":
	case "INC":
	case "INX":
	case "INY":
	case "LSR":
	case "ORA":
	case "ROL":
	case "ROR":
	case "SBC":
	case "PHA":
	case "PHP":
	case "PLA":
	case "PLP":
	case "JMP":
	case "JSR":
	case "RTS":
	case "RTI":
	case "BCC":
	case "BCS":
	case "BEQ":
	case "BMI":
	case "BNE":
	case "BPL":
	case "BVC":
	case "BVS":
	case "CLC":
		c.P.C = false
	case "CLD":
		// nothing to do
	case "CLI":
		c.P.I = false
	case "CLV":
		c.P.V = true
	case "SEC":
		c.P.C = true
	case "SED":
		// nothing to do
	case "SEI":
		c.P.I = true
	case "BRK":
		c.Interrupt = BRK
		c.P.B = true
	case "NOP":
		// nothing to do
	default:
		if operation_names[opecode] != "NONE" {
			log.Printf("%s operation called, but not implemented.\n", operation_names[opecode])
		}
	}
}
