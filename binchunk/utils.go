package binchunk

func (r *LuaReader) checkHeader() {
	// signature [4]byte // 魔数，快速识别文件格式，0x1B4C7561
	if string(r.loadBytes(uint(len(LUA_SIGNATURE)))) != LUA_SIGNATURE {
		panic("unmatched signature")
	}
	//	version   byte    // 大版本号，小版本号，发布号
	if r.loadByte() != LUAC_VERSION {
		panic("unmatched version")
	}
	//	format    byte    // 格式号
	if r.loadByte() != LUAC_FORMAT {
		panic("unmatched format")
	}
	//	luacData  [6]byte // 0x19 0x93 0x0D 0x0A 0x1A 0x0A
	if string(r.loadBytes(uint(len(LUAC_DATA)))) != LUAC_DATA {
		panic("unmatched luac data")
	}
	//	InstructionSize  byte    // 8
	if r.loadByte() != INSTRUCTION_SIZE {
		panic("unmatched instruction size")
	}
	//	luaIntegerSize byte    // 8
	if r.loadByte() != LUA_INTEGER_SIZE {
		panic("unmatched Lua_Integer size")
	}
	//	luaNumberSize  byte    // 8
	if r.loadByte() != LUA_NUMBER_SIZE {
		panic("unmatched Lua_Number size")
	}
	//	luacInt        int64   // 8，0x5678，检查大小端模式
	if r.loadLuaInteger() != LUAC_INT {
		panic("unmatched endianness")
	}
	//	luacNum        float64 // 8，存储浮点数370.5
	if r.loadLuaNumber() != LUAC_NUM {
		panic("unmatched float format")
	}
}
func (r *LuaReader) loadProto(parentSource string) *Prototype {
	source := r.loadString()
	if source == "" {
		source = parentSource
	}
	return &Prototype{
		Source:          source,
		LineDefined:     r.loadInt(),
		LastLineDefined: r.loadInt(),
		NumParams:       r.loadByte(),
		IsVararg:        r.loadByte(),
		MaxStackSize:    r.loadByte(),
		Code:            r.loadCode(),
		Constants:       r.loadConstants(),
		Upvalues:        r.loadUpvalues(),
		Protos:          r.loadProtos(source),
		// Debug information
		LineInfo:     r.loadLineInfo(),
		AbsLineInfo:  r.loadAbsLineInfo(),
		LocVars:      r.loadLocVars(),
		UpvalueNames: r.loadUpvalueNames(),
	}
}
