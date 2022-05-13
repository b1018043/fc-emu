package cpu

import (
	"github.com/b1018043/fc-emu/pkg/logger"
	"github.com/b1018043/fc-emu/pkg/ppu"
	"github.com/b1018043/fc-emu/pkg/utils"
)

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

const (
	stackStart = 0x0100
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
	MemoryMap [0xffff + 1]byte
	PPU       *ppu.PPU
}

// レジスタ内容に関しては http://hp.vector.co.jp/authors/VA042397/nes/6502.html を参照
type Registers struct {
	A  byte           // アキュムレーター
	X  byte           // インデックスレジスタ
	Y  byte           // インデックスレジスタ
	S  byte           // スタックポインタレジスタ スタックポインタは0x01|S
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
		MemoryMap: [0xFFFF + 1]byte{},
	}
}

func NewRegisters() *Registers {
	return &Registers{
		A: 0x00,
		X: 0x00,
		Y: 0x00,
		S: 0xFD,
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

func (c *CPU) SetPRGROM(progROM []byte) {
	// log.Printf("progROM size: 0x%x\n", len(progROM))
	for i := 0; i < len(progROM); i++ {
		c.setMemoryValue(uint16(i+0x8000), progROM[i])
		// log.Printf("ROM[0x%x]: 0x%02x, opename: %s\n", i+0x8000, progROM[i], operation_names[progROM[i]])
	}
}

func (c *CPU) SetPPU(ppu *ppu.PPU) {
	c.PPU = ppu
}

func (c *CPU) getAddress(addr uint16) uint16 {
	return uint16(c.getMemoryValue(addr)) | uint16(c.getMemoryValue(addr+1))<<8
}

func (c *CPU) setAddress(addr, val uint16) {
	c.setMemoryValue(addr, byte(val))
	c.setMemoryValue(addr+1, byte(val>>8))
}

func (c *CPU) Push(v uint8) {
	// 本来 0x0100以下には入らない
	// if c.S < 0 {
	// 	// TODO: なんか処理をかく
	// 	return
	// }
	c.setMemoryValue(stackStart|uint16(c.S), v)
	c.S--
}

func (c *CPU) PushAddress(val uint16) {
	// if c.S < 0 {
	// 	return
	// }
	c.setAddress(stackStart|uint16(c.S-1), val)
	c.S -= 2
}

func (c *CPU) PushStatusRegister() {
	var t uint8 = 0
	if c.P.N {
		t |= 1 << 7
	}
	if c.P.V {
		t |= 1 << 6
	}
	if c.P.R {
		t |= 1 << 5
	}
	if c.P.B {
		t |= 1 << 4
	}
	if c.P.D {
		t |= 1 << 3
	}
	if c.P.I {
		t |= 1 << 2
	}
	if c.P.Z {
		t |= 1 << 1
	}
	if c.P.C {
		t |= 1 << 0
	}
	c.Push(t)
}

func (c *CPU) Pop() uint8 {
	c.Registers.S++
	return c.getMemoryValue(stackStart | uint16(c.S))
}

func (c *CPU) PopAddress() uint16 {
	addr := c.getAddress(stackStart | uint16(c.S+1))
	c.S += 2
	return addr
}

func (c *CPU) PopStatusRegister() {
	t := c.Pop()
	if t>>7&1 != 0 {
		c.P.N = true
	}
	if t>>6&1 != 0 {
		c.P.V = true
	}
	if t>>5&1 != 0 {
		c.P.R = true
	}
	if t>>4&1 != 0 {
		c.P.B = true
	}
	if t>>3&1 != 0 {
		c.P.D = true
	}
	if t>>2&1 != 0 {
		c.P.I = true
	}
	if t>>1&1 != 0 {
		c.P.Z = true
	}
	if t&1 != 0 {
		c.P.C = true
	}
}

func (c *CPU) SetZN(val uint8) {
	c.P.Z = val == 0
	c.P.N = utils.IsNegativeByte(val)
}

func (c *CPU) Reset() {
	newPC := c.getAddress(0xfffc)
	if newPC == 0x0000 {
		newPC = 0x8000
	}

	c.Registers = Registers{
		A: 0x00,
		X: 0x00,
		Y: 0x00,
		S: 0xFD,
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
		PC: newPC,
	}
}

// TODO: アドレスの取得処理を追加して置き換える
func (c *CPU) detectAddress(mode int) uint16 {
	switch mode {
	case modeAbsolute:
		return uint16(c.getMemoryValue(c.PC)) | uint16(c.getMemoryValue(c.PC+1))<<8
	case modeAbsoluteX:
		return (uint16(c.getMemoryValue(c.PC)) | uint16(c.getMemoryValue(c.PC+1))<<8) + uint16(c.X)
	case modeAbsoluteY:
		return (uint16(c.getMemoryValue(c.PC)) | uint16(c.getMemoryValue(c.PC+1))<<8) + uint16(c.Y)
	case modeAccumulator:
		return 0
	case modeImmediate:
		return c.PC
	case modeImplied:
		return 0
	case modeIndirect:
		abs := uint16(c.getMemoryValue(c.PC)) | (uint16(c.getMemoryValue(c.PC+1)) << 8)
		return uint16(c.getMemoryValue(abs)) | (uint16(c.getMemoryValue(abs+1)) << 8)
	case modeIndirectX:
		t := uint16(c.getMemoryValue(c.PC)) + uint16(c.X)
		return uint16(c.getMemoryValue(t)) | (uint16(c.getMemoryValue(t+1)) << 8)
	case modeIndirectY:
		t := uint16(c.getMemoryValue(c.PC))
		return ((uint16(c.getMemoryValue(t))) | (uint16(c.getMemoryValue(t+1)) << 8)) + uint16(c.Y)
	case modeRelative:
		offset := uint16(c.getMemoryValue(c.PC))
		var t uint16 = 0
		if offset >= 0x80 {
			t = 0x100
		}
		return c.PC + 1 + offset - t
	case modeZeroPage:
		return uint16(c.getMemoryValue(c.PC))
	case modeZeroPageX:
		return uint16(c.getMemoryValue(c.PC)) + uint16(c.X)
	case modeZeroPageY:
		return uint16(c.getMemoryValue(c.PC)) + uint16(c.Y)
	default:
		logger.DebugLog(logger.FATAL, "unknown operation mode %d, c.PC=0x%x\n", mode, c.PC-1)
	}
	return 1
}

func (c *CPU) Run() int {
	switch c.Interrupt {
	case RESET:
	case NMI:
	case IRQ:
		c.P.B = false
	case BRK:
		if !c.P.I {
			c.P.B = true
			c.PC++
			c.PushAddress(c.PC)
			c.PushStatusRegister()
			c.P.I = true
			c.PC = c.getAddress(0xfffe)
		}
	case NONE:
	default:
		logger.DebugLog(logger.FATAL, "unknown interrupt: %d\n", c.Interrupt)
	}
	opecode := c.getMemoryValue(c.PC)
	logger.DebugLog(logger.PRINT, "PC: %x, opecode: %d, size: %d, name: %s", c.PC, opecode, operation_sizes[opecode], operation_names[opecode])
	c.PC++
	address := c.detectAddress(operation_modes[opecode])
	c.PC += uint16(operation_sizes[opecode] - 1)
	logger.DebugLog(logger.PRINT, "address: 0x%04x\n", address)
	c.exec(opecode, address)
	// log.Printf("after PC: %x", c.PC)
	return operation_cycles[opecode]
}

func (c *CPU) setMemoryValue(address uint16, val byte) {
	// NOTE: 0x2000 ~ 0x2008 はPPUのレジスタへのアクセスを行う
	if address >= 0x2000 && address < 0x2008 {
		// log.Printf("address: 0x%x\n", address)
		switch {
		// TODO: remove magic number
		case address == 0x2006:
			c.PPU.SetAddress(val)
		case address == 0x2007:
			c.PPU.SetData(val)
		default:
			// NOTE: 現状では0x2006と0x2007にのみ対応している
			// もし、他のレジスタへの書き込みも行いたい場合には、初めにPPUのレジスタ用のスライスを
			// 初期化するところから実装を始めるのが吉
			// c.PPU.Registers[address-0x2000] = val
			logger.DebugLog(logger.PRINT, "address: 0x%04x, val: 0x%02x\n", address, val)
		}
	} else {
		c.MemoryMap[address] = val
	}
}

func (c *CPU) getMemoryValue(address uint16) byte {
	if address >= 0x2000 && address < 0x2008 {
		switch {
		case address == 0x2007:
			return c.PPU.GetData()
		default:
			// return c.PPU.Registers[address-0x2000]
			return c.MemoryMap[address]
		}
	} else {
		return c.MemoryMap[address]
	}
}

func (c *CPU) exec(opecode byte, address uint16) {
	switch operation_names[opecode] {
	case "LDA":
		c.A = c.getMemoryValue(address)
		c.SetZN(c.A)
	case "LDX":
		c.X = c.getMemoryValue(address)
		c.SetZN(c.X)
	case "LDY":
		c.Y = c.getMemoryValue(address)
		c.SetZN(c.Y)
	case "STA":
		c.setMemoryValue(address, c.A)
	case "STX":
		c.setMemoryValue(address, c.X)
	case "STY":
		c.setMemoryValue(address, c.Y)
	case "TAX":
		c.X = c.A
		c.SetZN(c.X)
	case "TAY":
		c.Y = c.A
		c.SetZN(c.Y)
	case "TSX":
		c.X = c.S
		c.SetZN(c.X)
	case "TXA":
		c.A = c.X
		c.SetZN(c.A)
	case "TXS":
		c.S = c.X
		c.SetZN(c.S)
	case "TYA":
		c.A = c.Y
		c.SetZN(c.A)
	case "ADC":
		v1 := c.A
		v2 := c.getMemoryValue(address)
		var v3 byte
		if c.P.C {
			v3 = 0
		} else {
			v3 = 1
		}
		c.A = v1 + v2 + v3

		c.SetZN(c.A)

		// v1 v2の符号が同じ && v1と演算結果の符号が異なる
		c.P.V = (!((v1^v2)&0x80 != 0)) && ((c.A^v1)&0x80 != 1)
		c.P.C = int(v1)+int(v2)+int(v3) > 0xFF
	case "AND":
		c.A &= c.getMemoryValue(address)
		c.SetZN(c.A)
	case "ASL":
		var v uint8
		if operation_modes[opecode] == modeAccumulator {
			v = c.A
		} else {
			v = c.getMemoryValue(address)
		}
		c.P.C = (v>>7)&1 != 0
		v <<= 1

		c.SetZN(v)

		if operation_modes[opecode] == modeAccumulator {
			c.A = v
		} else {
			c.setMemoryValue(address, v)
		}
	case "BIT":
		v := c.getMemoryValue(address)
		c.P.N = utils.IsNegativeByte(v)
		c.P.V = (v>>6&1 != 0)
		c.P.Z = v&c.A == 0
	case "CMP":
		v := c.A - c.getMemoryValue(address)
		c.P.C = !utils.IsNegativeByte(v)
		c.SetZN(v)
	case "CPX":
		v := c.X - c.getMemoryValue(address)
		c.P.C = !utils.IsNegativeByte(v)
		c.SetZN(v)
	case "CPY":
		v := c.Y - c.getMemoryValue(address)
		// P.Zがtrueであれば等しい、P.C がtrueであればc.Yの方が大きい、P.Nがtrueであればアドレスで示された値が大きい
		c.P.C = !utils.IsNegativeByte(v)
		c.SetZN(v)
	case "DEC":
		v := c.getMemoryValue(address) - 1
		c.SetZN(v)
		c.setMemoryValue(address, v)
	case "DEX":
		c.X--
		c.SetZN(c.X)
	case "DEY":
		c.Y--
		c.SetZN(c.Y)
	case "EOR":
		c.A ^= c.getMemoryValue(address)
		c.SetZN(c.A)
	case "INC":
		v := c.getMemoryValue(address) + 1
		c.SetZN(v)
		c.setMemoryValue(address, v)
	case "INX":
		c.X++
		c.SetZN(c.X)
	case "INY":
		c.Y++
		c.SetZN(c.Y)
	case "LSR":
		var v uint8
		if operation_modes[opecode] == modeAccumulator {
			v = c.A
		} else {
			v = c.getMemoryValue(address)
		}
		c.P.C = v&1 != 0
		v >>= 1

		c.SetZN(v)

		if operation_modes[opecode] == modeAccumulator {
			c.A = v
		} else {
			c.setMemoryValue(address, v)
		}
	case "ORA":
		c.A |= c.getMemoryValue(address)
		c.SetZN(c.A)
	case "ROL":
		var v uint8
		if operation_modes[opecode] == modeAccumulator {
			v = c.A
		} else {
			v = c.getMemoryValue(address)
		}

		carry := (v & 0x80) >> 7
		v <<= 1
		v |= carry
		c.P.C = (carry != 0)
		c.SetZN(v)

		if operation_modes[opecode] == modeAccumulator {
			c.A = v
		} else {
			c.setMemoryValue(address, v)
		}

	case "ROR":
		var v uint8
		if operation_modes[opecode] == modeAccumulator {
			v = c.A
		} else {
			v = c.getMemoryValue(address)
		}

		carry := v & 1
		v >>= 1
		v |= carry << 7
		c.P.C = (carry != 0)

		c.SetZN(v)

		if operation_modes[opecode] == modeAccumulator {
			c.A = v
		} else {
			c.setMemoryValue(address, v)
		}
	case "SBC":
		v1 := c.A
		v2 := c.getMemoryValue(address)
		var v3 byte
		if c.P.C {
			v3 = 0
		} else {
			v3 = 1
		}
		c.A = v1 - v2 - v3
		// c.a + v2 = v1 (v3は無視) が成り立ち、c.aとv2のどちらもv1と符号が異なればオーバーフロー
		// NOTE: これであっているか不安
		c.P.V = utils.IsNegativeByte(v1^v2) && utils.IsNegativeByte(v1^c.A)

		c.SetZN(c.A)

		c.P.C = int(v1)-int(v2)-int(v3) < 0x00
	case "PHA":
		c.Push(c.A)
	case "PHP":
		c.PushStatusRegister()
	case "PLA":
		v := c.Pop()
		c.A = v
		c.SetZN(v)
	case "PLP":
		c.PopStatusRegister()
	case "JMP":
		c.PC = address
	case "JSR":
		// NOTE: JSR命令の最後のアドレスを格納する
		c.PushAddress(c.PC - 1)
		c.PC = address
	case "RTS":
		c.PC = c.PopAddress() + 1
	case "RTI":
		c.PopStatusRegister()
		c.PC = c.PopAddress()
	case "BCC":
		if !c.P.C {
			c.PC = address
		}
	case "BCS":
		if c.P.C {
			c.PC = address
		}
	case "BEQ":
		if c.P.Z {
			c.PC = address
		}
	case "BMI":
		if c.P.N {
			c.PC = address
		}
	case "BNE":
		if !c.P.Z {
			c.PC = address
		}
	case "BPL":
		if !c.P.N {
			c.PC = address
		}
	case "BVC":
		if !c.P.V {
			c.PC = address
		}
	case "BVS":
		if c.P.V {
			c.PC = address
		}
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
			logger.DebugLog(logger.PRINT, "%s operation called, but not implemented.\n", operation_names[opecode])
		}
	}
}
