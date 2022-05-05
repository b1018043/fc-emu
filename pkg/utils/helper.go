package utils

func IsNegativeByte(v byte) bool {
	return v&(1<<7) != 0
}
