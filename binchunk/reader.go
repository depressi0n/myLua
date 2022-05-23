package binchunk

import (
	"bufio"
	"encoding/binary"
	"io"
	"math"
	"unsafe"
)

type LuaReader struct {
	*bufio.Reader
}

var (
	litterEndian = binary.LittleEndian
	bigEndian    = binary.BigEndian
)

func NewLuaReader(reader io.Reader) *LuaReader {
	return &LuaReader{bufio.NewReader(reader)}
}
func (r *LuaReader) loadBytes(n uint) []byte {
	buf := make([]byte, n)
	if cnt, err := r.Read(buf); err != nil || uint(cnt) != n {
		panic(err)
	}
	return buf
}

func (r *LuaReader) loadByte() byte {
	b, err := r.ReadByte()
	if err != nil {
		panic(err)
	}
	return b
}
func (r *LuaReader) loadUint64() uint64 {
	return litterEndian.Uint64(r.loadBytes(LUA_INTEGER_SIZE))
}
func (r *LuaReader) loadUint32() uint32 {
	return litterEndian.Uint32(r.loadBytes(uint(unsafe.Sizeof(uint32(0)))))
}
func (r *LuaReader) loadLuaInteger() int64 {
	return int64(r.loadUint64())
}

func (r *LuaReader) loadLuaNumber() float64 {
	return math.Float64frombits(r.loadUint64())
}

// Lua5.4 引入LEB128编码的变种，使用大端序
// 利用MSB是否为0来标识是否有后续字节
// data :          xxxxxxxx yyyyyyyy zzzzzzzz
// step1: 00000xxx 0xxxxxyy 0yyyyyyz 0zzzzzzz
// step2: 00000xxx 0xxxxxyy 0yyyyyyz 1zzzzzzz
// Binary  : 00000001 00101100
// Lua 5.4 : 00000010 10101100ƒ
func (r *LuaReader) loadUnsigned(limit uint) uint {
	x := uint(0)
	var b byte
	limit >>= 7
	for (b & 0x80) == 0 {
		b = r.loadByte()
		if x >= limit {
			panic("integer overflow")
		}
		x = (x << 7) | uint(b&0x7f)
	}
	return x
}
func (r *LuaReader) loadInt() int {
	return int(r.loadUnsigned(math.MaxInt))
}
func (r *LuaReader) loadString() string {
	var res string
	size := r.loadUnsigned(math.MaxInt)
	if size == 0 {
		return res
	}
	size--
	if size <= LUAI_MAXSHORTLEN { // short string
		res = string(r.loadBytes(size))
	} else { // long string
		res = string(r.loadBytes(size))
	}
	return res
}

func (r *LuaReader) loadCode() []uint32 {
	code := make([]uint32, r.loadInt())
	for i := range code {
		code[i] = r.loadUint32()
	}
	return code
}

func (r *LuaReader) loadConstants() []interface{} {
	constants := make([]interface{}, r.loadInt())
	for i := 0; i < len(constants); i++ {
		tag := r.loadByte()
		switch tag {
		case TAG_NIL:
			constants[i] = nil
		case TAG_BOOLEAN:
			constants[i] = tag != 0
		case TAG_NUMBER:
			constants[i] = r.loadLuaNumber()
		case TAG_INTERER:
			constants[i] = r.loadLuaInteger()
		case TAG_SHORT_STR, TAG_LONG_STR:
			constants[i] = r.loadString()
		default:
			panic("unknown tag")
		}
	}
	return constants
}

func (r *LuaReader) loadUpvalues() []Upvalue {
	upvalues := make([]Upvalue, r.loadInt())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			Instack: r.loadByte(),
			Idx:     r.loadByte(),
			Kind:    r.loadByte(),
		}
	}
	return upvalues
}

func (r *LuaReader) loadProtos(parentSource string) []*Prototype {
	protos := make([]*Prototype, r.loadInt())
	for i := range protos {
		protos[i] = r.loadProto(parentSource)
	}
	return protos
}

func (r *LuaReader) loadLineInfo() []byte {
	lineInfo := make([]byte, r.loadInt())
	for i := 0; i < len(lineInfo); i++ {
		lineInfo[i] = r.loadByte()
	}
	return lineInfo
}
func (r *LuaReader) loadAbsLineInfo() []AbsLineInfo {
	lineInfo := make([]AbsLineInfo, r.loadInt())
	for i := range lineInfo {
		lineInfo[i].Line = r.loadInt()
		lineInfo[i].Pc = r.loadInt()
	}
	return lineInfo
}

func (r *LuaReader) loadLocVars() []LocVar {
	locVars := make([]LocVar, r.loadInt())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: r.loadString(),
			StartPC: r.loadInt(),
			EndPC:   r.loadInt(),
		}
	}
	return locVars
}

func (r *LuaReader) loadUpvalueNames() []string {
	names := make([]string, r.loadInt())
	for i := range names {
		names[i] = r.loadString()
	}
	return names
}
