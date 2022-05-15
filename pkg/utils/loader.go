package utils

import (
	"bufio"
	"fmt"
	"os"
	"reflect"

	"github.com/b1018043/fc-emu/pkg/logger"
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

	if !reflect.DeepEqual(bytes[0:3], []byte("NES")) {
		return nil, nil, fmt.Errorf("[error]: %s is not nes file", filename)
	}

	progROM, charROM := parseROMBuffer(bytes)

	logger.DebugLog(logger.PRINT, "progrom size: 0x%04x\ncharrom size: 0x%04x\n", bytes[4], bytes[5])

	return progROM, charROM, nil
}

func parseROMBuffer(bytes []byte) ([]byte, []byte) {
	charROMPages := int(bytes[5])
	progROMPages := int(bytes[4])
	charROMStart := NES_HEADER_SIZE + progROMPages*PROG_ROM_PAGE_SIZE
	charROMEnd := charROMStart + charROMPages*CHAR_ROM_PAGE_SIZE
	progROM := bytes[NES_HEADER_SIZE:charROMStart]
	charROM := bytes[charROMStart:charROMEnd]

	return progROM, charROM
}
