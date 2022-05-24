package binchunk

import (
	"fmt"
	"github.com/depressi0n/myLua/vm"
	"os"
	"testing"
)

func TestUndump(t *testing.T) {
	filename := "binchunk_test"
	file, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}
	proto := Undump(file)
	list(proto)
}

func list(fc *Prototype) {
	printHeader(fc)
	printCode(fc)
	printDetail(fc)
	for _, p := range fc.Protos {
		list(p)
	}
}

func printDetail(fc *Prototype) {
	fmt.Printf("constants (%d):\n", len(fc.Constants))
	for i, constant := range fc.Constants {
		fmt.Printf("\t%d\t%s\n", i+1, constantToString(constant))
	}
	fmt.Printf("locals (%d):\n", len(fc.LocVars))
	for i, locVar := range fc.LocVars {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, locVar.VarName, locVar.StartPC+1, locVar.EndPC+1)
	}
	fmt.Printf("upvalues (%d):\n", len(fc.Upvalues))
	for i, upvalue := range fc.Upvalues {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, upvalname(fc, i), upvalue.Instack, upvalue.Idx)
	}
}

func upvalname(fc *Prototype, idx int) string {
	if len(fc.UpvalueNames) > 0 {
		return fc.UpvalueNames[idx]
	}
	return "-"
}

func constantToString(constant interface{}) string {
	switch constant.(type) {
	case nil:
		return "nil"
	case bool:
		return fmt.Sprintf("%t", constant)
	case float64:
		return fmt.Sprintf("%g", constant)
	case int64:
		return fmt.Sprintf("%d", constant)
	case string:
		return fmt.Sprintf("%q", constant)
	default:
		return "?"
	}
}

func printCode(fc *Prototype) {
	for pc, c := range fc.Code {
		line := "-"
		if len(fc.LineInfo) > 0 {
			line = fmt.Sprintf("%d", fc.LineInfo[pc])
		}
		i := vm.Instruction(c)
		fmt.Printf("\t%d\t[%s]\t%s \t", pc+1, line, i.OpName())
		vm.PrintOprands(i)
		fmt.Println()
	}
}

func printHeader(fc *Prototype) {
	funcType := "main"
	if fc.LineDefined > 0 {
		funcType = "function"
	}

	varargFlag := ""
	if fc.IsVararg > 0 {
		varargFlag = "+"
	}

	fmt.Printf("\n%s <%s:%d,%d> (%d instructions)\n", funcType, fc.Source, fc.LineDefined, fc.LastLineDefined, len(fc.Code))
	fmt.Printf("%d%s params, %d slots, %d upvalues, ", fc.NumParams, varargFlag, fc.MaxStackSize, len(fc.Upvalues))
	fmt.Printf("%d locals, %d constants, %d functions\n", len(fc.LocVars), len(fc.Constants), len(fc.Protos))
}
