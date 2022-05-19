package ast

// chunk -> block 表明 chunk实际上等同于代码块

// Block 表示代码块，包含语句序列 Stats 和返回语句 RetExps 的表达式序列
// block -> {stat} [retstat]
// retstat -> return [explist] [';']
// explist -> exp {',' exp} => []Exp
type Block struct {
	LastLine int    // 代码块的末行行号，用于代码生成阶段
	Stats    []Stat // Lua语句
	RetExps  []Exp  // 表达式
}

// Stat 表示最基本的执行单位，只能执行不能用于求值
// stat -> ';'
//      -> varlist '=' explist
//      -> functioncall
//      -> label
//      -> break
//      -> goto Name
//      -> do block end
//      -> while exp do block end
//      -> repeat block until exp
//      -> if exp then block {elseif exp then block} [else block] end
//      -> for Name '=' exp ',' [',' exp] do block end
//      -> for namelist in explist do block end
//      -> function funcname funcbody
//      -> local function Name funcbody
//      -> local namelist ['=' explist]
type Stat interface{}

// EmptyStat stat -> ';'
type EmptyStat struct{}

// AssignStat stat -> varlist '=' explist
// varlist -> var {',' var} => []Exp
// explist -> exp {',' exp} => []Exp
type AssignStat struct {
	LastLine int // 代码生成阶段使用
	VarList  []Exp
	ExpList  []Exp
}

// FuncCallStat stat -> functioncall
// functioncall -> prefixexp args | prefixexp ':' Name args
// 函数调用既可以是表达式也可以是语句
type FuncCallStat = FunctionCallExp

// LabelStat label -> '::' Name '::'
type LabelStat struct {
	Name string /*记录标签名*/
}

// BreakStat stat -> break
type BreakStat struct {
	Line int /*用于代码生成阶段会产生一条跳转指令*/
}

// GotoStat stat -> goto Name，与 LabelStat 搭配使用
type GotoStat struct {
	Name string /*记录标签名*/
}

// DoStat stat -> do block end
type DoStat struct {
	Block *Block /*引入新的作用域*/
}

// WhileStat stat -> while exp do block end
// 用于实现条件循环
type WhileStat struct {
	Exp   Exp
	Block *Block
}

// RepeatStat stat -> repeat block until exp
// 用于实现条件循环
type RepeatStat struct {
	Block *Block
	Exp   Exp
}

// IfStat stat -> if exp then block {elseif exp then block} [else block] end
// => if exp then block {elseif exp then block} [elif true then block] end
// if exp then block {elseif exp then block} end
type IfStat struct {
	Exps   []Exp
	Blocks []*Block
}

// ForNumStat stat -> for Name '=' exp ',' [',' exp] do block end
type ForNumStat struct {
	LineOfFor int // 代码生成阶段使用
	LineOfDo  int // 代码生成阶段使用
	VarName   string
	InitExp   Exp
	LimitExp  Exp
	StepExp   Exp
	Block     *Block
}

// ForInStat stat -> for namelist in explist do block end
// namelist -> Name {',' Name} => []string
// explist -> exp {',' exp} => []Exp
type ForInStat struct {
	LineOfDo int // 代码生成阶段使用
	NameList []string
	ExpList  []Exp
	Block    *Block
}

// 非局部函数定义 stat -> function funcname funcbody
// funcname -> Name {'.' Name} [':' Name]
// funcbody -> '(' [paralist] ')' block end

// LocalFuncDefStat stat -> local function Name funcbody
// 定义局部函数定义语句
// funcbody -> '(' [paralist] ')' block end
type LocalFuncDefStat struct {
	Name string
	Exp  *FuncDefExp
}

// LocalVarDeclStat stat -> local namelist ['=' explist ]
// namelist -> Name {',' Name} => []string
// explist -> exp {',' exp} => []Exp
type LocalVarDeclStat struct {
	LastLine int // 代码生成阶段使用
	NameList []string
	ExpList  []Exp
}
