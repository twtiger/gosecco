package compiler

type argumentPosition struct {
	lower uint32
	upper uint32
}

var argument = []argumentPosition{
	argumentPosition{upper: 0x10, lower: 0x14},
	argumentPosition{upper: 0x18, lower: 0x1c},
	argumentPosition{upper: 0x20, lower: 0x24},
	argumentPosition{upper: 0x28, lower: 0x2c},
	argumentPosition{upper: 0x30, lower: 0x34},
	argumentPosition{upper: 0x38, lower: 0x3c},
}
