package ast

// Exp 表示表达式，可以用于求值，但不能单独执行
// exp -> nil | false | true | Numeral | LiteralString | '...'
//     -> functiondef | prefixexp | tableconstructor
//     -> exp binop exp | unop exp
type Exp interface{}

// NilExp exp -> nil
type NilExp struct {
	Line int
}

// TrueExp exp -> true
type TrueExp struct {
	Line int
}

// FalseExp exp -> false
type FalseExp struct {
	Line int
}

// IntegerExp exp -> Numeral
type IntegerExp struct {
	Line int
	Val  int64
}

// FloatExp exp -> Numeral
type FloatExp struct {
	Line int
	Val  float64
}

// StringExp exp -> LiteralString
type StringExp struct {
	Line int
	Str  string
}

// VarargExp exp -> '...'
type VarargExp struct {
	Line int
}

// FuncDefExp exp -> functiondef
// functiondef -> function funcbody
// funcbody -> '(' [paralist] ')' block end
// paralist -> namelist [',' '...'] | '...' => []string, vararg
// namelist -> Name {',' Name} // []string
type FuncDefExp struct {
	Line     int // 花括号所在行号
	LastLine int // 'end' 所在行
	ParList  []string
	IsVararg bool
	Block    *Block
}

// 前缀表达式，可以作为表访问表达式、记录访问表达式、函数调用表达式的前缀
// 其中表访问表达式和记录访问表达式在语义上等价即t.k <=> t[k]
// exp -> prefixexp
// prefixexp -> var | functioncall | '(' exp ')'
// var -> Name | prefixexp '[' exp ']' | prefixexp '.' Name
// functioncall -> prefixexp args
//				-> prefixexp ':' Name args
// 等价于
// prefixexp -> Name => 名字表达式
//			 -> prefixexp '[' exp ']' => 表访问表达式
//           -> prefixexp '.' Name
//           -> prefixexp [':' Name] args
//			 -> '(' exp ')' => 圆括号表达式

// NameExp
// 名字表达式
type NameExp struct {
	Line int
	Name string
}

// TableAccessExp
// 表访问表达式
type TableAccessExp struct {
	LastLine  int // ']' 所在行
	PrefixExp Exp
	KeyExp    Exp
}

// ParenExp prefixexp -> '(' exp ')'
// 圆括号表达式
type ParenExp struct {
	Exp Exp
}

// FunctionCallExp functioncall -> prefixexp [':' Name] args
// args -> '(' [explist] ')' | tableconstructor | LiteralString
// 其中v:name(args) <=> v.name(v,args)
type FunctionCallExp struct {
	Line      int // '(' 所在行
	LastLine  int // ')' 所在行
	PrefixExp Exp
	NameExp   *StringExp
	Args      []Exp
}

// TableConstructorExp exp -> tableconstructor
// tableconstructor -> '{' fieldlist '}'
// fieldlist -> field { fieldsep field} [fieldsp]
// field -> '[' exp ']' '=' exp
//		 -> Name '=' exp
//		 -> exp
// fieldsep -> ',' | ';'
type TableConstructorExp struct {
	Line     int // '{' 所在行
	LastLine int // '}' 所在行
	KeyExps  []Exp
	ValExps  []Exp
}

// UnopExp exp -> unop exp
// unop -> '-' | not | '#' | '~'
type UnopExp struct {
	Line int
	Op   int
	Exp  Exp
}

// BinopExp exp -> exp binop exp
type BinopExp struct {
	Line     int
	Op       int
	LeftExp  Exp
	RightExp Exp
}

// ConcatExp
// 拼接运算符是一个特殊的运算符，便于代码生成阶段优化拼接操作
type ConcatExp struct {
	Line int // 最后一个'..'所在行
	Exps []Exp
}
