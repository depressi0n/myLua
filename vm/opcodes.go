package vm

import "fmt"

const (
	iABC = iota
	iABx
	iAsBx
	iAx
	isJ
)

// /*===========================================================================
//  We assume that instructions are unsigned 32-bit integers.
//  All instructions have an opcode in the first 7 bits.
//  Instructions can have the following formats:
//
//        3 3 2 2 2 2 2 2 2 2 2 2 1 1 1 1 1 1 1 1 1 1 0 0 0 0 0 0 0 0 0 0
//        1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0 9 8 7 6 5 4 3 2 1 0
//iABC          C(8)     |      B(8)     |k|     A(8)      |   Op(7)     |
//iABx                Bx(17)               |     A(8)      |   Op(7)     |
//iAsBx              sBx (signed)(17)      |     A(8)      |   Op(7)     |
//iAx                           Ax(25)                     |   Op(7)     |
//isJ                           sJ(25)                     |   Op(7)     |
//
//  A signed argument is represented in excess K: the represented value is
//  the written unsigned value minus K, where K is half the maximum for the
//  corresponding unsigned argument.
//===========================================================================*/

const (
	/* R[x] - 寄存器访问 */
	/* K[x] - 常量表访问 */
	/* RK(x) == if k(i) then K[x] else R[x] */
	/*----------------------------------------------------------------------
	  name		             args	 description
	------------------------------------------------------------------------*/

	OP_MOVE       = iota /*	A B	    R[A] := R[B]					*/
	OP_LOADI             /*	A sBx	R[A] := sBx					*/
	OP_LOADF             /*	A sBx	R[A] := (lua_Number)sBx				*/
	OP_LOADK             /*	A Bx	R[A] := K[Bx]					*/
	OP_LOADKX            /*	A	    R[A] := K[extra arg]				*/
	OP_LOADFALSE         /*	A	    R[A] := false					*/
	OP_LFALSESKIP        /* A	    R[A] := false; pc++	(*)			*/
	OP_LOADTRUE          /*	A	    R[A] := true					*/
	OP_LOADNIL           /*	A B	    R[A], R[A+1], ..., R[A+B] := nil		*/
	OP_GETUPVAL          /*	A B	    R[A] := UpValue[B]				*/
	OP_SETUPVAL          /*	A B	    UpValue[B] := R[A]				*/

	OP_GETTABUP /*	A B C	R[A] := UpValue[B][K[C]:string]			*/
	OP_GETTABLE /*	A B C	R[A] := R[B][R[C]]				*/
	OP_GETI     /*	A B C	R[A] := R[B][C]					*/
	OP_GETFIELD /*	A B C	R[A] := R[B][K[C]:string]			*/

	OP_SETTABUP /*	A B C	UpValue[A][K[B]:string] := RK(C)		*/
	OP_SETTABLE /*	A B C	R[A][R[B]] := RK(C)				*/
	OP_SETI     /*	A B C	R[A][B] := RK(C)				*/
	OP_SETFIELD /*	A B C	R[A][K[B]:string] := RK(C)			*/

	OP_NEWTABLE /*	A B C k	R[A] := {}					*/

	OP_SELF /*	A B C	R[A+1] := R[B]; R[A] := R[B][RK(C):string]	*/

	OP_ADDI /*	A B sC	R[A] := R[B] + sC				*/

	OP_ADDK  /*	A B C	R[A] := R[B] + K[C]:number			*/
	OP_SUBK  /*	A B C	R[A] := R[B] - K[C]:number			*/
	OP_MULK  /*	A B C	R[A] := R[B] * K[C]:number			*/
	OP_MODK  /*	A B C	R[A] := R[B] % K[C]:number			*/
	OP_POWK  /*	A B C	R[A] := R[B] ^ K[C]:number			*/
	OP_DIVK  /*	A B C	R[A] := R[B] / K[C]:number			*/
	OP_IDIVK /*	A B C	R[A] := R[B] // K[C]:number			*/

	OP_BANDK /*	A B C	R[A] := R[B] & K[C]:integer			*/
	OP_BORK  /*	A B C	R[A] := R[B] | K[C]:integer			*/
	OP_BXORK /*	A B C	R[A] := R[B] ~ K[C]:integer			*/

	OP_SHRI /*	A B sC	R[A] := R[B] >> sC				*/
	OP_SHLI /*	A B sC	R[A] := sC << R[B]				*/

	OP_ADD  /*	A B C	R[A] := R[B] + R[C]				*/
	OP_SUB  /*	A B C	R[A] := R[B] - R[C]				*/
	OP_MUL  /*	A B C	R[A] := R[B] * R[C]				*/
	OP_MOD  /*	A B C	R[A] := R[B] % R[C]				*/
	OP_POW  /*	A B C	R[A] := R[B] ^ R[C]				*/
	OP_DIV  /*	A B C	R[A] := R[B] / R[C]				*/
	OP_IDIV /*	A B C	R[A] := R[B] // R[C]				*/

	OP_BAND /*	A B C	R[A] := R[B] & R[C]				*/
	OP_BOR  /*	A B C	R[A] := R[B] | R[C]				*/
	OP_BXOR /*	A B C	R[A] := R[B] ~ R[C]				*/
	OP_SHL  /*	A B C	R[A] := R[B] << R[C]				*/
	OP_SHR  /*	A B C	R[A] := R[B] >> R[C]				*/

	OP_MMBIN  /*	A B C	call C metamethod over R[A] and R[B]	(*)	*/
	OP_MMBINI /*	A sB C k	call C metamethod over R[A] and sB	*/
	OP_MMBINK /*	A B C k		call C metamethod over R[A] and K[B]	*/

	OP_UNM  /*	A B	R[A] := -R[B]					*/
	OP_BNOT /*	A B	R[A] := ~R[B]					*/
	OP_NOT  /*	A B	R[A] := not R[B]				*/
	OP_LEN  /*	A B	R[A] := #R[B] (length operator)			*/

	OP_CONCAT /*	A B	R[A] := R[A].. ... ..R[A + B - 1]		*/

	OP_CLOSE /*	A	close all upvalues >= R[A]			*/
	OP_TBC   /*	A	mark variable A "to be closed"			*/
	OP_JMP   /*	sJ	pc += sJ					*/
	OP_EQ    /*	A B k	if ((R[A] == R[B]) ~= k) then pc++		*/
	OP_LT    /*	A B k	if ((R[A] <  R[B]) ~= k) then pc++		*/
	OP_LE    /*	A B k	if ((R[A] <= R[B]) ~= k) then pc++		*/

	OP_EQK /*	A B k	if ((R[A] == K[B]) ~= k) then pc++		*/
	OP_EQI /*	A sB k	if ((R[A] == sB) ~= k) then pc++		*/
	OP_LTI /*	A sB k	if ((R[A] < sB) ~= k) then pc++			*/
	OP_LEI /*	A sB k	if ((R[A] <= sB) ~= k) then pc++		*/
	OP_GTI /*	A sB k	if ((R[A] > sB) ~= k) then pc++			*/
	OP_GEI /*	A sB k	if ((R[A] >= sB) ~= k) then pc++		*/

	OP_TEST    /*	A k	if (not R[A] == k) then pc++			*/
	OP_TESTSET /*	A B k	if (not R[B] == k) then pc++ else R[A] := R[B] (*) */

	OP_CALL     /*	A B C	R[A], ... ,R[A+C-2] := R[A](R[A+1], ... ,R[A+B-1]) */
	OP_TAILCALL /*	A B C k	return R[A](R[A+1], ... ,R[A+B-1])		*/

	OP_RETURN  /*	A B C k	return R[A], ... ,R[A+B-2]	(see note)	*/
	OP_RETURN0 /*		return						*/
	OP_RETURN1 /*	A	return R[A]					*/

	OP_FORLOOP /*	A Bx	update counters; if loop continues then pc-=Bx; */
	OP_FORPREP /*	A Bx	<check values and prepare counters>;
		ifnot to run then pc+=Bx+1;			*/

	OP_TFORPREP /*	A Bx	create upvalue for R[A + 3]; pc+=Bx		*/
	OP_TFORCALL /*	A C	R[A+4], ... ,R[A+3+C] := R[A](R[A+1], R[A+2]);	*/
	OP_TFORLOOP /*	A Bx	if R[A+2] ~= nil then { R[A]=R[A+2]; pc -= Bx }	*/

	OP_SETLIST /*	A B C k	R[A][C+i] := R[A+i], 1 <= i <= B		*/

	OP_CLOSURE /*	A Bx	R[A] := closure(KPROTO[Bx])			*/

	OP_VARARG /*	A C	R[A], R[A+1], ..., R[A+C-2] = vararg		*/

	OP_VARARGPREP /*A	(adjust vararg parameters)			*/

	OP_EXTRAARG /*	Ax	extra (larger) argument for previous opcode	*/

	NUM_OPCODES
)

const (
	OpArgN = iota // argument is not used
	OpArgU        // argument is used
	OpArgR        // argument is a register or a jump offset, iABC下表示寄存器索引，iAsBx模式下表示跳转偏移
	OpArgK        // argument is a constant or register/constant，常量表索引或者寄存器索引
)

/*
** masks for instruction properties. The format is:
** bits 0-2: op mode
** bit 3: instruction set register A
** bit 4: operator is a test (next instruction must be a jump)
** bit 5: instruction uses 'L->top' set by previous instruction (when B == 0)
** bit 6: instruction sets 'L->top' for next instruction (when C == 0)
** bit 7: instruction is an MM instruction (call a metamethod)
 */

type opcode struct {
	setMMFlag byte // [7] MM, an MM instruction (call a metamethod)
	setOTFlag byte // [6] ot, sets 'L->top' for next instruction (when C == 0)
	setITFlag byte // [5] it, instruction uses 'L->top' set by previous instruction (when B == 0)
	testFlag  byte // [4] t, operator is a test (next instruction must be a jump)
	setAFlag  byte // [3] a, instruction set register A
	opMode    byte // [0-2] op mode
	name      string
}

// Excess-K编码
// 简单来说，这种整数编码方式将可表示的最小负整数编码为全0、
// 将可表示的最大正整数编码为全1、其余整数按大小顺序分布在中间。

var opcodes = []opcode{
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "MOVE"},       /* OP_MOVE */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iAsBx, name: "LOADI"},     /* OP_LOADI */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iAsBx, name: "LOADF"},     /* OP_LOADF */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABx, name: "LOADK"},      /* OP_LOADK */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABx, name: "LOADKX"},     /* OP_LOADKX */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "LOADFALSE"},  /* OP_LOADFALSE */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "LFALSESKIP"}, /* OP_LFALSESKIP */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "LOADTRUE"},   /* OP_LOADTRUE */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "LOADNIL"},    /* OP_LOADNIL */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "GETUPVAL"},   /* OP_GETUPVAL */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 0, opMode: iABC, name: "SETUPVAL"},   /* OP_SETUPVAL */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "GETTABUP"},   /* OP_GETTABUP */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "GETTABLE"},   /* OP_GETTABLE */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "GETI"},       /* OP_GETI */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "GETFIELD"},   /* OP_GETFIELD */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 0, opMode: iABC, name: "SETTABUP"},   /* OP_SETTABUP */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 0, opMode: iABC, name: "SETTABLE"},   /* OP_SETTABLE */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 0, opMode: iABC, name: "SETI"},       /* OP_SETI */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 0, opMode: iABC, name: "SETFIELD"},   /* OP_SETFIELD */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "NEWTABLE"},   /* OP_NEWTABLE */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "SELF"},       /* OP_SELF */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "ADDI"},       /* OP_ADDI */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "ADDK"},       /* OP_ADDK */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "SUBK"},       /* OP_SUBK */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "MULK"},       /* OP_MULK */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "MODK"},       /* OP_MODK */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "POWK"},       /* OP_POWK */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "DIVK"},       /* OP_DIVK */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "IDIVK"},      /* OP_IDIVK */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "BANDK"},      /* OP_BANDK */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "BORK"},       /* OP_BORK */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "BXORK"},      /* OP_BXORK */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "SHRI"},       /* OP_SHRI */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "SHLI"},       /* OP_SHLI */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "ADD"},        /* OP_ADD */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "SUB"},        /* OP_SUB */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "MUL"},        /* OP_MUL */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "MOD"},        /* OP_MOD */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "POW"},        /* OP_POW */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "DIV"},        /* OP_DIV */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "IDIV"},       /* OP_IDIV */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "BAND"},       /* OP_BAND */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "BOR"},        /* OP_BOR */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "BXOR"},       /* OP_BXOR */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "SHL"},        /* OP_SHL */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "SHR"},        /* OP_SHR */
	{setMMFlag: 1, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 0, opMode: iABC, name: "MMBIN"},      /* OP_MMBIN */
	{setMMFlag: 1, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 0, opMode: iABC, name: "MMBINI"},     /* OP_MMBINI*/
	{setMMFlag: 1, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 0, opMode: iABC, name: "MMBINK"},     /* OP_MMBINK*/
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "UNM"},        /* OP_UNM */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "BNOT"},       /* OP_BNOT */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "NOT"},        /* OP_NOT */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "LEN"},        /* OP_LEN */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "CONCAT"},     /* OP_CONCAT */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 0, opMode: iABC, name: "CLOSE"},      /* OP_CLOSE */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 0, opMode: iABC, name: "TBC"},        /* OP_TBC */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 0, opMode: isJ, name: "JMP"},         /* OP_JMP */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 1, setAFlag: 0, opMode: iABC, name: "EQ"},         /* OP_EQ */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 1, setAFlag: 0, opMode: iABC, name: "LT"},         /* OP_LT */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 1, setAFlag: 0, opMode: iABC, name: "LE"},         /* OP_LE */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 1, setAFlag: 0, opMode: iABC, name: "EQK"},        /* OP_EQK */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 1, setAFlag: 0, opMode: iABC, name: "EQI"},        /* OP_EQI */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 1, setAFlag: 0, opMode: iABC, name: "LTI"},        /* OP_LTI */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 1, setAFlag: 0, opMode: iABC, name: "LEI"},        /* OP_LEI */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 1, setAFlag: 0, opMode: iABC, name: "GTI"},        /* OP_GTI */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 1, setAFlag: 0, opMode: iABC, name: "GEI"},        /* OP_GEI */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 1, setAFlag: 0, opMode: iABC, name: "TEST"},       /* OP_TEST */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 1, setAFlag: 1, opMode: iABC, name: "TESTSET"},    /* OP_TESTSET */
	{setMMFlag: 0, setOTFlag: 1, setITFlag: 1, testFlag: 0, setAFlag: 1, opMode: iABC, name: "CALL"},       /* OP_CALL */
	{setMMFlag: 0, setOTFlag: 1, setITFlag: 1, testFlag: 0, setAFlag: 1, opMode: iABC, name: "TAILCALL"},   /* OP_TAILCALL */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 1, testFlag: 0, setAFlag: 0, opMode: iABC, name: "RETURN"},     /* OP_RETURN */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 0, opMode: iABC, name: "RETURN0"},    /* OP_RETURN0 */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 0, opMode: iABC, name: "RETURN1"},    /* OP_RETURN1 */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABx, name: "FORLOOP"},    /* OP_FORLOOP */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABx, name: "FORPREP"},    /* OP_FORPREP */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 0, opMode: iABx, name: "TFORPREP"},   /* OP_TFORPREP */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 0, opMode: iABC, name: "TFORCALL"},   /* OP_TFORCALL */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABx, name: "TFORLOOP"},   /* OP_TFORLOOP */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 1, testFlag: 0, setAFlag: 0, opMode: iABC, name: "SETLIST"},    /* OP_SETLIST */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABx, name: "CLOSURE"},    /* OP_CLOSURE */
	{setMMFlag: 0, setOTFlag: 1, setITFlag: 0, testFlag: 0, setAFlag: 1, opMode: iABC, name: "VARARG"},     /* OP_VARARG */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 1, testFlag: 0, setAFlag: 1, opMode: iABC, name: "VARARGPREP"}, /* OP_VARARGPREP */
	{setMMFlag: 0, setOTFlag: 0, setITFlag: 0, testFlag: 0, setAFlag: 0, opMode: iAx, name: "EXTRAARG"},    /* OP_EXTRAARG */
}

type Instruction uint32

// 指令解码 Lua 5.4

func (i Instruction) Opcode() int {
	return int(i & 0x7F)
}
func (i Instruction) OpName() string {
	return opcodes[i.Opcode()].name
}
func (i Instruction) OpMode() byte {
	return opcodes[i.Opcode()].opMode
}
func (i Instruction) testAMode() byte {
	return opcodes[i.Opcode()].setAFlag
}
func (i Instruction) testTMode() byte {
	return opcodes[i.Opcode()].testFlag
}
func (i Instruction) testITMode() byte {
	return opcodes[i.Opcode()].setITFlag
}
func (i Instruction) testOTMode() byte {
	return opcodes[i.Opcode()].setOTFlag
}

// IABC
/* ┌------┐┌------┐-┌------┐┌-----┐
   |  C:8 ||  B:8 |k|  A:8 || Op:7|
   └------┘└------┘-└------┘└-----┘ */

func (i Instruction) IABC() (int, int, int, int) {
	a := int(i >> POS_A & 0xFF)
	k := int(i >> POS_k & 0x1)
	b := int(i >> POS_B & 0xFF)
	c := int(i >> POS_C & 0xFF)
	return a, k, b, c
}

// IABx (7)
/* ┌---------------┐┌------┐┌-----┐
   |     Bx:17     ||  A:8 || Op:7|
   └---------------┘└------┘└-----┘ */

func (i Instruction) IABx() (int, int) {
	a := int(i >> POS_A & 0xFF)
	bx := int(i >> POS_Bx)
	return a, bx
}

// IAsBx (2)
/* ┌---------------┐┌------┐┌-----┐
   |    sBx:17     ||  A:8 || Op:7|
   └---------------┘└------┘└-----┘ */

func (i Instruction) IAsBx() (int, int) {
	a, bx := i.IABx()
	return a, bx - MAXARG_sBx
}

// IAx (EXTRAARG)
/* ┌-----------------------┐┌-----┐
   |          Ax:25        || Op:7|
   └-----------------------┘└-----┘*/

func (i Instruction) IAx() int {
	return int(i >> POS_Ax)
}

// IsJ (JMP)
/* ┌-----------------------┐┌-----┐
   |          sJ:25        || Op:7|
   └-----------------------┘└-----┘*/

func (i Instruction) IsJx() int {
	return int(i >> POS_sJ)
}

// PrintOprands provide debug information for testing
func PrintOprands(i Instruction) {
	switch i.OpMode() {
	case iABC:
		a, k, b, c := i.IABC()
		fmt.Printf("%d", a)
		if k == 0 { // 常数
			fmt.Printf(" %d", b)
			fmt.Printf(" %d", c)
		} else { // 寄存器
			fmt.Printf(" R[%d]", b)
			fmt.Printf(" R[%d]", c)
		}

	case iABx:
		a, bx := i.IABx()
		fmt.Printf("%d", a)
		fmt.Printf(" %d", bx)
	case iAsBx:
		a, sbx := i.IABx()
		fmt.Printf("%d", a)
		fmt.Printf(" %d", sbx)
	case iAx:
		ax := i.IAx()
		fmt.Printf("%d", ax)
	case isJ:
		jx := i.IsJx()
		fmt.Printf("%d", jx)
	default:
		panic("unreachable code")
	}
}
