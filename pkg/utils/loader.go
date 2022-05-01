package utils

import (
	"bufio"
	"os"
)

// NOTE: サイズの定数に関しては https://www.nesdev.org/wiki/INES を参照
const (
	NES_HEADER_SIZE    = 0x0010
	PROG_ROM_PAGE_SIZE = 0x4000
	CHAR_ROM_PAGE_SIZE = 0x2000
)

// param: filename ret: progROM charROM error
func LoadFCROM(filename string) ([]byte, []byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	stats, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}
	size := stats.Size()
	bytes := make([]byte, size)
	b := bufio.NewReader(file)
	_, err = b.Read(bytes)
	if err != nil {
		return nil, nil, err
	}

	progROM, charROM := parseROMBuffer(bytes)

	return progROM, charROM, nil
}

func parseROMBuffer(bytes []byte) ([]byte, []byte) {
	charROMPages := int(bytes[5])
	progROMPages := int(bytes[4])
	charROMStart := NES_HEADER_SIZE + progROMPages*PROG_ROM_PAGE_SIZE
	charROMEnd := charROMStart + charROMPages*CHAR_ROM_PAGE_SIZE
	progROM := bytes[NES_HEADER_SIZE : charROMStart-1]
	charROM := bytes[charROMStart : charROMEnd-1]

	return progROM, charROM
}