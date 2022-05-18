package lexer

import (
	"fmt"
	"io/ioutil"
	"testing"
)

var chunk string
var chunkName string

func kindToCategory(kind int) string {
	switch {
	case kind < TOKEN_SEP_SEMI:
		return "other"
	case kind <= TOKEN_SEP_RCURLY:
		return "separator"
	case kind <= TOKEN_OP_NOT:
		return "operator"
	case kind <= TOKEN_KW_WHILE:
		return "keyword"
	case kind == TOKEN_IDENTIFIER:
		return "identifier"
	case kind == TOKEN_NUMBER:
		return "number"
	case kind == TOKEN_STRING:
		return "string"
	default:
		return "other"
	}
}
func TestLexer(t *testing.T) {
	chunkName = "hello_world.lua"
	data, err := ioutil.ReadFile(chunkName)
	if err != nil {
		t.Fatalf("fail to open file named %s: %v", chunkName, err)
	}
	chunk = string(data)
	lexer := NewLexer(chunk, chunkName)
	for {
		line, column, kind, token := lexer.NextToken()
		fmt.Printf("[%2d:%4d] [%-10s] %s\n",
			line, column, kindToCategory(kind), token)
		if kind == TOKEN_EOF {
			break
		}
	}
}
