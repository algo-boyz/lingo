package parser

import (
	"reflect"
	"testing"
)

func TestLexerSimple(t *testing.T) {
	Flush()
	var tokLeftPar = HookToken(TokLabel(TokLeftPar))
	var tokRightPar = HookToken(TokLabel(TokRightPar))
	HookToken(TokLabel(TokQuote))
	HookToken(TokLabel(TokIdentifier))
	HookToken(TokLabel(TokEof))
	HookToken(TokLabel("match"))
	HookToken(TokLabel("name"))
	HookToken(TokLabel("string"))

	lex := NewLexer("(desc)")
	identifier, _ := NewToken(TokLabel(TokIdentifier), "desc")

	want := []Tok{
		*tokLeftPar,
		identifier,
		*tokRightPar,
	}

	got := []Tok{}
	var tok Tok
	var err error

	for tok, err = lex.NextToken(); !tok.IsEof() && err == nil; tok, err = lex.NextToken() {
		got = append(got, tok)
	}

	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Wrong result. Expected:\n%#v\nbut got:\n%#v", want, got)
	}

	Flush()
}

func TestLexerAdvanced(t *testing.T) {
	Flush()
	var tokLeftPar = HookToken(TokLabel(TokLeftPar))
	var tokRightPar = HookToken(TokLabel(TokRightPar))
	var tokMatch = HookToken(TokLabel("match"))
	var tokName = HookToken(TokLabel("name"))
	HookToken(TokLabel(TokQuote))
	HookToken(TokLabel(TokIdentifier))
	HookToken(TokLabel(TokEof))
	HookToken(TokLabel("string"))

	lex := NewLexer("(match (name X) \"(exec|match)\")")

	ident, _ := NewToken(TokLabel(TokIdentifier), "X")
	pat, _ := NewToken(TokLabel(TokString), "(exec|match)")

	want := []Tok{
		*tokLeftPar,
		*tokMatch,
		*tokLeftPar.Label.ToToken(),
		*tokName,
		ident,
		*tokRightPar,
		pat,
		*tokRightPar,
	}

	got := []Tok{}
	var tok Tok
	var err error

	for tok, err = lex.NextToken(); !tok.IsEof() && err == nil; tok, err = lex.NextToken() {
		got = append(got, tok)
	}

	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Wrong result. Expected:\n%#v\nbut got:\n%#v", want, got)
	}
	Flush()
}
