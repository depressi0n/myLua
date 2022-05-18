package lexer

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Lexer 定义一个词法分析器，将输入流进行Token化即执行词法分析过程的数据结构
type Lexer struct {
	chunk     string // 源代码
	chunkName string // 源文件名称
	curLine   int    // 当前行号
	curColumn int    // 当前列号

	// 往后查看下一个token
	// 对当前状态进行备份，然后读取下一个token，记录类型
	// 恢复状态，并缓存这个token
	nextToken       string
	nextTokenKind   int
	nextTokenLine   int
	nextTokenColumn int
}

// NewLexer 创建一个词法分析器并初始化
func NewLexer(chunk string, chunkName string) *Lexer {
	return &Lexer{
		chunk:     chunk,
		chunkName: chunkName,
		curLine:   1,
		curColumn: 1,
	}
}

// hasPrefix 判断Lexer正在处理的当前位置是否以s作为前缀
func (l *Lexer) hasPrefix(prefix string) bool {
	return strings.HasPrefix(l.chunk, prefix)
}

// skipWhiteSpaces 跳过输入流中的空白符号和注释
// 以"--"开头表示注释
// 空白字符包括'\t', '\n', '\v', '\f', '\r', ' '
// 需要更新Lexer的当前处理行
func (l *Lexer) skipWhiteSpaces() {
	for len(l.chunk) > 0 {
		if l.hasPrefix("--") {
			l.skipComment()
		} else if l.hasPrefix("\r\n") || l.hasPrefix("\n\r") {
			l.next(2)
			l.curLine += 1
			l.curColumn = 0
		} else if isNewLine(l.chunk[0]) {
			l.next(1)
			l.curLine += 1
			l.curColumn = 0
		} else if isWhiteSpace(l.chunk[0]) {
			l.next(1)
		} else {
			break
		}
	}
}

// NextToken 返回下一个token
func (l *Lexer) NextToken() (int, int, int, string) {
	// 查看当前是否已经解析过下一个token即查看缓存
	if l.nextTokenLine > 0 {
		line := l.nextTokenLine
		column := l.nextTokenColumn
		kind := l.nextTokenKind
		token := l.nextToken

		l.curLine = l.nextTokenLine
		l.curColumn = l.nextTokenColumn
		l.nextTokenLine = 0
		l.nextTokenColumn = 0
		return line, column, kind, token
	}

	// 跳过空白符号和注释
	l.skipWhiteSpaces()
	if len(l.chunk) == 0 {
		return l.curLine, l.curColumn, TOKEN_EOF, "EOF"
	}

	// 根据当前首个字符进行处理
	// 主要包括两种情况：
	// 1. 根据第一个字符即可知道词素 如各种分隔符号
	// 2. 根据第一个字符不可知道词素 如< 和 <= 都是以<开始的 需要区分
	switch l.chunk[0] {
	case '+':
		l.next(1)
		return l.curLine, l.curColumn, TOKEN_OP_ADD, "+"
	case '-':
		l.next(1)
		return l.curLine, l.curColumn, TOKEN_OP_MINUS, "-"
	case '*':
		l.next(1)
		return l.curLine, l.curColumn, TOKEN_OP_MUL, "*"
	case '%':
		l.next(1)
		return l.curLine, l.curColumn, TOKEN_OP_MOD, "%"
	case '^':
		l.next(1)
		return l.curLine, l.curColumn, TOKEN_OP_POW, "^"
	case '#':
		l.next(1)
		return l.curLine, l.curColumn, TOKEN_OP_LEN, "#"
	case '&':
		l.next(1)
		return l.curLine, l.curColumn, TOKEN_OP_BAND, "&"
	case '|':
		l.next(1)
		return l.curLine, l.curColumn, TOKEN_OP_BOR, "|"
	case ';':
		l.next(1)
		return l.curLine, l.curColumn, TOKEN_SEP_SEMI, ";"
	case ',':
		l.next(1)
		return l.curLine, l.curColumn, TOKEN_SEP_COMMA, ","
	case '(':
		l.next(1)
		return l.curLine, l.curColumn, TOKEN_SEP_LPAREN, "("
	case ')':
		l.next(1)
		return l.curLine, l.curColumn, TOKEN_SEP_RPAREN, ")"
	case '{':
		l.next(1)
		return l.curLine, l.curColumn, TOKEN_SEP_LCURLY, "{"
	case '}':
		l.next(1)
		return l.curLine, l.curColumn, TOKEN_SEP_RCURLY, "}"
	case ']':
		l.next(1)
		return l.curLine, l.curColumn, TOKEN_SEP_RBRACK, "]"

	case '/':
		if l.hasPrefix("//") {
			l.next(2)
			return l.curLine, l.curColumn, TOKEN_OP_IDIV, "//"
		} else {
			l.next(1)
			return l.curLine, l.curColumn, TOKEN_OP_DIV, "/"
		}
	case '~':
		if l.hasPrefix("~=") {
			l.next(2)
			return l.curLine, l.curColumn, TOKEN_OP_NE, "~="
		} else {
			l.next(1)
			return l.curLine, l.curColumn, TOKEN_OP_WAVE, "~"
		}
	case '<':
		if l.hasPrefix("<<") {
			l.next(2)
			return l.curLine, l.curColumn, TOKEN_OP_SHL, "<<"
		} else if l.hasPrefix("<=") {
			l.next(2)
			return l.curLine, l.curColumn, TOKEN_OP_LE, "<="
		} else {
			l.next(1)
			return l.curLine, l.curColumn, TOKEN_OP_LT, "<"
		}
	case '>':
		if l.hasPrefix(">>") {
			l.next(2)
			return l.curLine, l.curColumn, TOKEN_OP_SHR, ">>"
		} else if l.hasPrefix(">=") {
			l.next(2)
			return l.curLine, l.curColumn, TOKEN_OP_GE, ">="
		} else {
			l.next(1)
			return l.curLine, l.curColumn, TOKEN_OP_GT, ">"
		}
	case '=':
		if l.hasPrefix("==") {
			l.next(2)
			return l.curLine, l.curColumn, TOKEN_OP_EQ, "=="
		} else {
			l.next(1)
			return l.curLine, l.curColumn, TOKEN_OP_ASSIGN, "="
		}
	case ':':
		if l.hasPrefix("::") {
			l.next(2)
			return l.curLine, l.curColumn, TOKEN_SEP_LABEL, "::"
		} else {
			l.next(1)
			return l.curLine, l.curColumn, TOKEN_SEP_COLON, ":"
		}

	case '.':
		if l.hasPrefix("...") {
			l.next(3)
			return l.curLine, l.curColumn, TOKEN_VARARG, "..."
		} else if l.hasPrefix("..") {
			l.next(2)
			return l.curLine, l.curColumn, TOKEN_OP_CONCAT, ".."
		} else if len(l.chunk) == 1 || // 仅有一个.号
			!isDigit(l.chunk[1]) { // 后续跟随有非数字字符，表示成员？
			l.next(1)
			return l.curLine, l.curColumn, TOKEN_SEP_DOT, "."
		}
	case '[':
		if l.hasPrefix("[[") || l.hasPrefix("[=") {
			return l.curLine, l.curColumn, TOKEN_STRING, l.scanLongString()
		} else {
			l.next(1)
			return l.curLine, l.curColumn, TOKEN_SEP_LBRACK, "["
		}
	case '\'', '"':
		return l.curLine, l.curColumn, TOKEN_STRING, l.scanShortString()
	}
	// 如果不是以上字符，则考虑是数字或者标识符
	c := l.chunk[0]
	if c == '.' || isDigit(c) {
		token := l.scanNumber()
		return l.curLine, l.curColumn, TOKEN_NUMBER, token
	}
	// 标识符可以以_或大小写字母开始
	if c == '_' || isLetter(c) {
		token := l.scanIdentifier()
		if kind, ok := keyworkd[token]; ok {
			return l.curLine, l.curColumn, kind, token
		} else {
			return l.curLine, l.curColumn, TOKEN_IDENTIFIER, token
		}
	}
	l.error("unexpected symbol near %q", c)
	return -1, -1, TOKEN_EOF, ""
}

// next 将Lexer的当前处理行往后移动
func (l *Lexer) next(n int) {
	l.chunk = l.chunk[n:]
	l.curColumn += n
}

// skipComment 跳过注释
// 注释以"--"开始，包含以下两种：
// 		长注释：[=*[ xxx ]=*]，其中"xxx"表示一个长字符串
//		短注释：注释到当前行末截至
func (l *Lexer) skipComment() {
	l.next(2) // 跳过"--"
	// 长注释 [[some comment]]
	if l.hasPrefix("[") {
		if !(reOpeningLongBracket.FindString(l.chunk) == "") {
			l.scanLongString()
			return
		}
	}
	// 短注释遇到换行符即表示注释结束（可以理解为单行注释）
	for len(l.chunk) > 0 && !isNewLine(l.chunk[0]) {
		l.next(1)
	}
}

var reNewLine = regexp.MustCompile("\r\n|\n\r|\n|\r")

// 字面字符串包括一种由长括号包含的方式定义，[[之间可以包含若干个=
// 两个正的方括号间插入 n 个等号定义为 第 n 级开长括号
// 不受分行限制，不处理任何转义符，并且忽略掉任何不同级别的长括号。
// 其中碰到的任何形式的换行串（回车、换行、回车加换行、换行加回车），
// 都会被转换为单个换行符。
var reOpeningLongBracket = regexp.MustCompile(`\[=*\[`)

// 寻找左右长方括号，如果任何一个都找不到，则语法错误
// 提取字符串字面量，把左方括号和右方括号去掉，换行符序列统一换成换行符
// 将第一个换行符去掉后得到最终字符串
func (l *Lexer) scanLongString() string {
	openingLongBracket := reOpeningLongBracket.FindString(l.chunk)
	if openingLongBracket == "" {
		l.error("invalid long string delimiter near '%s", l.chunk[0:2])
	}
	// 结尾的字符串必须以相同级别的闭长括号作为结尾
	closingLongBracket := strings.Replace(openingLongBracket, "[", "]", -1)
	closingLongBracketIdx := strings.Index(l.chunk, closingLongBracket)
	// 没有找到则表示长括号没有正常结束
	if closingLongBracketIdx < 0 {
		l.error("unfinished long string or comment")
	}
	str := l.chunk[len(openingLongBracket):closingLongBracketIdx]
	l.next(closingLongBracketIdx + len(closingLongBracket))

	// 将长字符串中所有的换行符统一表示为'\n'
	str = reNewLine.ReplaceAllString(str, "\n")
	// 更新行号和列号
	l.curLine += strings.Count(str, "\n")
	l.curColumn = 0
	// 如果有首个空行，则去掉
	if len(str) > 0 && str[0] == '\n' {
		str = str[1:]
	}
	return str
}

// 字面串 可以用单引号或双引号括起。
// 字面串内部可以包含下列 C 风格的转义串： '\a' （响铃）， '\b' （退格）， '\f' （换页）， '\n' （换行），
//									 '\r' （回车）， '\t' （横项制表）， '\v' （纵向制表），
//									 '\\' （反斜杠）， '\"' （双引号）， 以及 '\'' (单引号)。
// 在反斜杠后跟一个真正的换行等价于在字符串中写一个换行符。
// 转义串 '\z' 会忽略其后的一系列空白符，包括换行； 在需要对一个很长的字符串常量断行为多行并希望在每个新行保持缩进时非常有用。
// (?s) 表示单行模式，将更改.的含义，使它与每一个字符匹配（包括换行符\n）。
var reShortStr = regexp.MustCompile(`(?s)(^'(\\\\|\\'|\\\n|\\z\s*|[^'\n])*')|(^"(\\\\|\\"|\\\n|\\z\s*|[^"\n])*")`)

func (l *Lexer) scanShortString() string {
	if str := reShortStr.FindString(l.chunk); str != "" {
		l.next(len(str))
		// 去掉'或"
		str = str[1 : len(str)-1]
		// 检查str中是否需要完成转义
		if strings.Index(str, `\`) >= 0 {
			// 寻找其中的换行符，更新行号和列号
			l.curLine += len(reNewLine.FindAllString(str, -1))
			l.curColumn = 0
			// 对获取对字符串进行转义表示
			str = l.escape(str)
		}
		return str
	}
	l.error("unfinished string")
	return ""
}

// error用于抛出错误信息
func (l *Lexer) error(f string, a ...interface{}) {
	err := fmt.Sprintf(f, a...)
	err = fmt.Sprintf("%s:%d:%d: %s", l.chunkName, l.curLine, l.curColumn, err)
	panic(err)
}

// \ddd ， 这里的 ddd 是一到三个十进制数字。
var reDecEscapeSeq = regexp.MustCompile(`^\\\d{1,3}`)

// \xXX， 此处的 XX 必须是恰好两个字符的 16 进制数。
var reHexEscapeSeq = regexp.MustCompile(`^\\x[\da-fA-F]{2}`)

// 转义符 \u{XXX} 来表示 （这里必须有一对花括号）用 UTF-8 编码的 Unicode 字符，
// 此处的 XXX 是用 16 进制表示的字符编号。
var reUnicodeEscapeSeq = regexp.MustCompile(`^\\u{[\da-fA-F]+}`)

// escape 对字符串完成转义
// '\a' （响铃）， '\b' （退格）， '\f' （换页）， '\n' （换行），
// '\r' （回车）， '\t' （横项制表）， '\v' （纵向制表），
// '\\' （反斜杠）， '\"' （双引号）， 以及 '\'' (单引号)。
// \xXX， 此处的 XX 必须是恰好两个字符的 16 进制数。
// \ddd ， 这里的 ddd 是一到三个十进制数字。
// 注意，如果在转义符后接着恰巧是一个数字符号的话， 必须在这个转义形式中写满三个数字。
// \u{XXX} 来表示 （这里必须有一对花括号）用 UTF-8 编码的 Unicode 字符，
func (l *Lexer) escape(str string) string {
	var buf bytes.Buffer
	for len(str) > 0 {
		if str[0] != '\\' {
			buf.WriteByte(str[0])
			str = str[1:]
			continue
		}
		// 此时以\开头但没有后续字符
		if len(str) == 1 {
			l.error("unfinished string")
		}
		switch str[1] {
		case 'a':
			buf.WriteByte('\a')
			str = str[2:]
			continue
		case 'b':
			buf.WriteByte('\b')
			str = str[2:]
			continue
		case 'f':
			buf.WriteByte('\f')
			str = str[2:]
			continue
		case 'n':
			buf.WriteByte('\n')
			str = str[2:]
			continue
		case '\n':
			buf.WriteByte('\n')
			str = str[2:]
			continue
		case 'r':
			buf.WriteByte('\r')
			str = str[2:]
			continue
		case 't':
			buf.WriteByte('\t')
			str = str[2:]
			continue
		case 'v':
			buf.WriteByte('\v')
			str = str[2:]
			continue
		case '"':
			buf.WriteByte('"')
			str = str[2:]
			continue
		case '\'':
			buf.WriteByte('\'')
			str = str[2:]
			continue
		case '\\':
			buf.WriteByte('\\')
			str = str[2:]
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // \ddd
			if found := reDecEscapeSeq.FindString(str); found != "" {
				d, _ := strconv.ParseInt(found[1:], 10, 32)
				if d < 0xFF {
					buf.WriteByte(byte(d))
					str = str[len(found):]
					continue
				}
				l.error("decimal escape too large near '%s'", found)
			}
		case 'x': // \xXX
			if found := reHexEscapeSeq.FindString(str); found != "" {
				d, _ := strconv.ParseInt(found[2:], 16, 32)
				buf.WriteByte(byte(d))
				str = str[len(found):]
				continue
			}
		case 'u': // \u{XXX}
			if found := reUnicodeEscapeSeq.FindString(str); found != "" {
				d, err := strconv.ParseInt(found[3:len(found)-1], 16, 32)
				if err == nil && d < 0x10FFFF {
					buf.WriteRune(rune(d))
					str = str[len(found):]
					continue
				}
				l.error("UTF-8 value too large near '%s'", found)
			}
		case 'z':
			str = str[2:]
			for len(str) > 0 && isWhiteSpace(str[0]) {
				str = str[1:]
			}
			continue
		default:
			l.error("invalid escape sequence near '\\%c'", str[1])
		}
	}
	return buf.String()
}

// 数字字面量由可选的小数部分和可选的十为底的指数部分构成：
// 指数部分用字符 'e' 或 'E' 来标记。
// Lua 也接受以 0x 或 0X 开头的 16 进制常量。
// 16 进制常量也接受小数加指数部分的形式，指数部分是以二为底， 用字符 'p' 或 'P' 来标记。
// 数字常量中包含小数点或指数部分时，被认为是一个浮点数； 否则被认为是一个整数。
var reNumber = regexp.MustCompile(`^0[xX][\da-fA-F]*(\.[\da-fA-F]*)?([pP][+\-]?\d+)?|^\d*(\.\d*)?([eE][+\-]?\d+)?`)

func (l *Lexer) scanNumber() string {
	return l.scan(reNumber)
}

// 标识符可以是由非数字打头的任意字母下划线和数字构成的字符串。
// 标识符可用于对变量、表的域、以及标签命名。^[a-z_][a-z\d_]*
var reIdentifier = regexp.MustCompile(`(?i)^[a-z_][a-z\d_]*`)

func (l *Lexer) scanIdentifier() string {
	return l.scan(reIdentifier)
}

func (l *Lexer) scan(re *regexp.Regexp) string {
	if token := re.FindString(l.chunk); token != "" {
		l.next(len(token))
		return token
	}
	panic("unreachable")
}

// LookAhead 预读并缓存下一个token，同时返回下一个token的类型
func (l *Lexer) LookAhead() int {
	// 查看当前是否已经有缓存
	if l.nextTokenLine > 0 {
		return l.nextTokenKind
	}
	// 保存词法分析器当前状态
	currentLine := l.curLine
	currentColumn := l.curColumn
	line, column, kind, token := l.NextToken()
	// 恢复词法分析器状态，并缓存下一个token
	l.curLine = currentLine
	l.curColumn = currentColumn
	l.nextTokenLine = line
	l.nextTokenColumn = column
	l.nextTokenKind = kind
	l.nextToken = token
	return kind
}

// NextTokenOfKind 读取下一个token，期待类型是_kind
func (l *Lexer) NextTokenOfKind(_kind int) (int, int, string) {
	line, column, kind, token := l.NextToken()
	if kind != _kind {
		l.error("syntax error near '%s'", token)
	}
	return line, column, token
}

// NexIdentifier 读取下一个标识符
func (l *Lexer) NexIdentifier() (line int, column int, token string) {
	return l.NextTokenOfKind(TOKEN_IDENTIFIER)
}

// Line 返回Lexer正在处理的当前行行号
func (l *Lexer) Line() int {
	return l.curLine
}
