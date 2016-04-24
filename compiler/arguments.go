package compiler

type argumentPosition struct {
	lower uint32
	upper uint32
}

var argument = []argumentPosition{
	argumentPosition{lower: 0x10, upper: 0x14},
	argumentPosition{lower: 0x18, upper: 0x1c},
	argumentPosition{lower: 0x20, upper: 0x24},
	argumentPosition{lower: 0x28, upper: 0x2c},
	argumentPosition{lower: 0x30, upper: 0x34},
	argumentPosition{lower: 0x38, upper: 0x3c},
}
