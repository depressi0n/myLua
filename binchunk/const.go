package binchunk

// Lua5.4中头部检查的相关标识
const (
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0x54
	LUAC_FORMAT      = 0 /* this is the official format */
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

const (
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_FALSE     = TAG_BOOLEAN
	TAG_TRUE      = TAG_BOOLEAN | 0x10
	TAG_NUMBER    = 0x03
	TAG_INTERER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

const (
	LUAI_MAXSHORTLEN = 40
)
