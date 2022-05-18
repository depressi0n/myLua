package lexer

// 词素定义
const (
	TOKEN_EOF        = iota // end of file
	TOKEN_VARARG            // ...
	TOKEN_SEP_SEMI          // ;
	TOKEN_SEP_COMMA         // ,
	TOKEN_SEP_DOT           // .
	TOKEN_SEP_COLON         // :
	TOKEN_SEP_LABEL         // ::
	TOKEN_SEP_LPAREN        // (
	TOKEN_SEP_RPAREN        // )
	TOKEN_SEP_LBRACK        // [
	TOKEN_SEP_RBRACK        // ]
	TOKEN_SEP_LCURLY        // {
	TOKEN_SEP_RCURLY        // }

	// 操作符号

	TOKEN_OP_ASSIGN // =
	TOKEN_OP_MINUS  // -
	TOKEN_OP_WAVE   // ~
	TOKEN_OP_ADD    // +
	TOKEN_OP_MUL    // *
	TOKEN_OP_DIV    // /
	TOKEN_OP_IDIV   // //
	TOKEN_OP_POW    // ^
	TOKEN_OP_MOD    // %
	TOKEN_OP_BAND   // &
	TOKEN_OP_BOR    // |
	TOKEN_OP_SHR    // >>
	TOKEN_OP_SHL    // <<
	TOKEN_OP_CONCAT // ..
	TOKEN_OP_LT     // <
	TOKEN_OP_LE     // <=
	TOKEN_OP_GT     // >
	TOKEN_OP_GE     // >=
	TOKEN_OP_EQ     // ==
	TOKEN_OP_NE     // ~=
	TOKEN_OP_LEN    // #
	TOKEN_OP_AND    // and
	TOKEN_OP_OR     // or
	TOKEN_OP_NOT    // not

	TOKEN_KW_BREAK    // break
	TOKEN_KW_DO       // do
	TOKEN_KW_ELSE     // else
	TOKEN_KW_ELSEIF   // elseif
	TOKEN_KW_END      // end
	TOKEN_KW_FALSE    // false
	TOKEN_KW_FOR      // for
	TOKEN_KW_FUNCTION // function
	TOKEN_KW_GOTO     // goto
	TOKEN_KW_IF       // if
	TOKEN_KW_IN       // in
	TOKEN_KW_LOCAL    // local
	TOKEN_KW_NIL      // nil
	TOKEN_KW_REPEAT   // repeat
	TOKEN_KW_RETURN   // return
	TOKEN_KW_THEN     // then
	TOKEN_KW_TRUE     // true
	TOKEN_KW_UNTIL    // until
	TOKEN_KW_WHILE    // while

	TOKEN_IDENTIFIER // identifier
	TOKEN_NUMBER     // number literal
	TOKEN_STRING     // string literal

	TOKEN_OP_UNM  = TOKEN_OP_MINUS // unary minus
	TOKEN_OP_SUB  = TOKEN_OP_MINUS //
	TOKEN_OP_BNOT = TOKEN_OP_WAVE  //
	TOKEN_OP_BXOR = TOKEN_OP_WAVE  //
)

// 部分词素对应的保留关键字
// and    break  do        else    elseif  end
// false  for    function  goto    if      in
// local  nil    not       or      repeat  return
// then   true   until     while
var keyworkd = map[string]int{
	"and":      TOKEN_OP_AND,
	"break":    TOKEN_KW_BREAK,
	"do":       TOKEN_KW_DO,
	"else":     TOKEN_KW_ELSE,
	"elseif":   TOKEN_KW_ELSEIF,
	"end":      TOKEN_KW_END,
	"false":    TOKEN_KW_FALSE,
	"for":      TOKEN_KW_FOR,
	"function": TOKEN_KW_FUNCTION,
	"goto":     TOKEN_KW_GOTO,
	"if":       TOKEN_KW_IF,
	"in":       TOKEN_KW_IN,
	"local":    TOKEN_KW_LOCAL,
	"nil":      TOKEN_KW_NIL,
	"not":      TOKEN_OP_NOT,
	"or":       TOKEN_OP_OR,
	"repeat":   TOKEN_KW_REPEAT,
	"return":   TOKEN_KW_RETURN,
	"then":     TOKEN_KW_THEN,
	"true":     TOKEN_KW_TRUE,
	"until":    TOKEN_KW_UNTIL,
	"while":    TOKEN_KW_WHILE,
}

// 运算符优先级别
// 带()表示单个字符的词素
//  (+)     (-)     (*)     /     (%)     (^)     (#)
//  (&)     ~     (|)     <<    >>    //
//  ==    ~=    <=    >=    <     >     =
//  (;)     :     (,)     .     ..    ...
//  (()     ())     ({)     (})     [     (])     ::

// isNewLine判断当前字符是否是换行符
func isNewLine(c byte) bool {
	return c == '\r' || c == '\n'
}

// isWhiteSpace 判断当前字符是否是空白符号
func isWhiteSpace(c byte) bool {
	switch c {
	case '\t', '\n', '\v', '\f', '\r', ' ':
		return true
	default:
		return false
	}
}

// isDigit 判断当前字符是否是数字
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// 判断当前字符是否是字母
func isLetter(c uint8) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}
