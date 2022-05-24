package vm

// opcode中相关位置和大小
const (
	SIZE_C  = 8
	SIZE_B  = 8
	SIZE_Bx = SIZE_C + SIZE_B + 1
	SIZE_A  = 8
	SIZE_Ax = SIZE_Bx + SIZE_A
	SIZE_sJ = SIZE_Bx + SIZE_A

	SIZE_OP = 7

	POS_OP = 0
	POS_A  = POS_OP + SIZE_OP
	POS_k  = POS_A + SIZE_A
	POS_B  = POS_k + 1
	POS_C  = POS_B + SIZE_B
	POS_Bx = POS_k
	POS_Ax = POS_A
	POS_sJ = POS_A
)

const (
	MAXARG_A   = (1 << SIZE_A) - 1
	MAXARG_B   = (1 << SIZE_B) - 1
	MAXARG_C   = (1 << SIZE_C) - 1
	MAXARG_Ax  = (1 << SIZE_Ax) - 1
	MAXARG_Bx  = (1 << SIZE_Bx) - 1
	MAXARG_sBx = MAXARG_Bx >> 1
	MAXARG_sJ  = (1 << SIZE_sJ) - 1

	OFFSET_sC  = MAXARG_C >> 1
	OFFSET_sBx = MAXARG_Bx >> 1
	OFFSET_sJ  = MAXARG_sJ >> 1
)

const (
	MAXINDEXRK = MAXARG_B
	NO_REG     = MAXARG_A
)
