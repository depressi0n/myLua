package binchunk

import (
	"io"
)

type binarychunk struct {
	header                  // 头部
	sizeUpvalues uint8      // 主函数upvalue数量
	mainFunc     *Prototype // 主函数原型
}
type header struct {
	signature [4]byte // 魔数，快速识别文件格式，0x1B4C7561
	version   byte    // 大版本号，小版本号，发布号
	format    byte    // 格式号
	luacData  [6]byte // 0x1993 0x0D 0x0A 0x1A 0x0A
	//cintSize       byte    // 4
	//sizetSize      byte    // 4
	luaIntegerSize byte    // 8
	luaNumberSize  byte    // 8
	luacInt        int64   // 8，0x5678，检查大小端模式
	luacNum        float64 // 8，存储浮点数370.5
}

type Upvalue struct {
	Instack byte
	Idx     byte
	Kind    byte
}
type AbsLineInfo struct {
	Pc   int
	Line int
}
type LocVar struct {
	VarName string
	StartPC int
	EndPC   int
}

// Prototype 定义函数原型，包括函数基本信息，指令集，常量表，
// upvalue表，子函数原型以及调试信息
type Prototype struct {
	// 基本信息
	Source          string // 源文件名
	LineDefined     int
	LastLineDefined int
	NumParams       byte // 固定参数个数
	IsVararg        byte // 是否vararg函数
	MaxStackSize    byte // 虚拟机中真正使用的是栈结构，除了支持push和pop外，还支持索引
	// 指令集
	Code []uint32
	// 常量表，存放字面量包括nil，布尔值，整数，浮点数，字符串
	// 每个常量都有一个tag，占1个字节
	// 0x00 -> nil 不存储
	// 0x01 -> boolean 字节0/1
	// 0x03 -> number
	// 0x04 -> 短字符串
	// 0x13 -> integer
	// 0x14 -> 长字符串
	Constants []interface{}
	// upvalue表
	Upvalues []Upvalue
	// 子函数原型
	Protos []*Prototype
	// 调试信息
	LineInfo     []byte        // 行号表
	AbsLineInfo  []AbsLineInfo // 行号表
	LocVars      []LocVar      // 局部变量表
	UpvalueNames []string      // _ENV
}

func Undump(reader io.Reader) *Prototype {
	r := NewLuaReader(reader)
	r.checkHeader()
	nupvals := r.loadByte()
	proto := r.loadProto("")
	if int(nupvals) != len(proto.Upvalues) {
		panic("unmatched nupvals and proto.Upvalues")
	}
	return proto
}
