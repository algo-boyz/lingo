package parser

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type TokMatcher interface {
	Match(s string) TokLabel
	Id() string
}

var tokMatcherMap = map[string]TokMatcher{}
var tokMatcherSlice = []TokMatcher{}

type Lexer struct {
	input           string
	currentPosition int
	nextPosition    int
	char            byte
}

type TokLabel string

type Tok struct {
	Label   TokLabel
	value   string
	Keyword bool
}

var label2Token = map[TokLabel]Tok{}

const (
	TokSingleQuoteChar = '\''
	TokSingleQuote     = "quote"
	TokComment         = ';'
	TokRoot            = "root"
	TokEof             = "eof"
	TokLeftPar         = '('
	TokRightPar        = ')'
	TokLeftBracket     = '['
	TokRightBracket    = ']'
	TokLeftCurly       = '{'
	TokRightCurly      = '}'
	TokUnknown         = ""
	TokIllegal         = "illegal"
	TokString          = "string"
	TokInt             = "int"
	TokFloat           = "float"
	TokIdentifier      = "identifier"
	TokKeyword         = "keyword"
	TokVector          = "vec"
	TokDict            = "dict"
	TokPair            = "pair"
	TokQuote           = '"'
	TokComma           = ','
)

func HookToken(label TokLabel) *Tok {
	if _, ok := label2Token[label]; ok {
		log.Fatalf("token '%s' already linked", label)
	}
	tok := Tok{
		Label: label,
		value: string(label),
	}
	label2Token[label] = tok
	return &tok
}

func HookMatcher(matcher TokMatcher) {
	if _, ok := tokMatcherMap[matcher.Id()]; ok {
		log.Fatalf("matcher for %s already registered", matcher.Id())
	}
	tokMatcherSlice = append([]TokMatcher{matcher}, tokMatcherSlice...)
	tokMatcherMap[matcher.Id()] = matcher
}

type KeyWordMatcher struct{}

func (i KeyWordMatcher) Match(s string) TokLabel {
	if strings.HasPrefix(s, ":") {
		return TokKeyword
	}
	return TokUnknown
}
func (i KeyWordMatcher) Id() string {
	return "keyword"
}

type IdentifierMatcher struct{}

func (i IdentifierMatcher) Match(s string) TokLabel {
	return TokIdentifier
}
func (i IdentifierMatcher) Id() string {
	return "identifier"
}

type HookedTokenMatcher struct{}

func (i HookedTokenMatcher) Match(s string) TokLabel {
	tok, ok := label2Token[TokLabel(s)]
	if !ok {
		return TokUnknown
	}
	return tok.Label
}
func (i HookedTokenMatcher) Id() string {
	return "hooked"
}

type NumberMatcher struct{}

func (i NumberMatcher) Match(s string) TokLabel {
	_, err := strconv.Atoi(s)
	if err == nil {
		return TokInt
	}
	_, err = strconv.ParseFloat(s, 64)
	if err == nil {
		return TokFloat
	}

	return TokUnknown
}

func (i NumberMatcher) Id() string {
	return "number"
}

func init() {
	HookMatcher(IdentifierMatcher{})
	HookMatcher(NumberMatcher{})
	HookMatcher(KeyWordMatcher{})
	HookMatcher(HookedTokenMatcher{})
}

func Flush() {
	label2Token = map[TokLabel]Tok{}
}

func NewToken(label TokLabel, value string) (Tok, error) {
	tok, ok := label2Token[label]
	if !ok {
		return tok, fmt.Errorf("could not find token for '%s'", label)
	}
	return tok.CloneWithValue(value), nil
}

func (label TokLabel) ToToken() *Tok {
	if val, ok := label2Token[label]; ok {
		return &val
	}
	unknown := label2Token[TokUnknown]
	return &unknown
}

func (t Tok) Clone() Tok {
	return Tok{t.Label, t.value, t.Keyword}
}

func (t Tok) CloneWithValue(value string) Tok {
	return Tok{t.Label, value, t.Keyword}
}

func (t Tok) IsSingleQuoteChar() bool {
	return string(t.Label) == string(TokSingleQuoteChar)
}

func (t Tok) IsSingleQuote() bool {
	return string(t.Label) == string(TokSingleQuote)
}

func (t Tok) IsParOpen() bool {
	return string(t.Label) == string(TokLeftPar)
}

func (t Tok) IsParClose() bool {
	return string(t.Label) == string(TokRightPar)
}

func (t Tok) IsBracketOpen() bool {
	return string(t.Label) == string(TokLeftBracket)
}

func (t Tok) IsBracketClose() bool {
	return string(t.Label) == string(TokRightBracket)
}

func (t Tok) IsCurlyOpen() bool {
	return string(t.Label) == string(TokLeftCurly)
}

func (t Tok) IsCurlyClose() bool {
	return string(t.Label) == string(TokRightCurly)
}

func (t Tok) IsComma() bool {
	return string(t.Label) == string(TokComma)
}

func (t Tok) IsComment() bool {
	return string(t.Label) == string(TokComment)
}

func (t Tok) IsUnknown() bool {
	return string(t.Label) == string(TokUnknown)
}

func (t Tok) IsIdentifier() bool {
	return string(t.Label) == string(TokIdentifier)
}

func (t Tok) IsEof() bool {
	return t.Label == TokEof
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() (Tok, error) {
	var tok Tok

	l.skipWS()

	var err error
	switch l.char {
	case TokSingleQuoteChar:
		tok, err = NewToken(TokLabel(TokSingleQuoteChar), string(TokSingleQuoteChar))
		if err != nil {
			return tok, err
		}
	case TokComment:
		tok, err = NewToken(TokLabel(TokComment), l.readComment())
	case TokLeftPar:
		tok, err = NewToken(TokLabel(TokLeftPar), string(TokLeftPar))
	case TokRightPar:
		tok, err = NewToken(TokLabel(TokRightPar), string(TokRightPar))
	case TokLeftBracket:
		tok, err = NewToken(TokLabel(TokLeftBracket), string(TokLeftBracket))
	case TokRightBracket:
		tok, err = NewToken(TokLabel(TokRightBracket), string(TokRightBracket))
	case TokLeftCurly:
		tok, err = NewToken(TokLabel(TokLeftCurly), string(TokLeftCurly))
	case TokRightCurly:
		tok, err = NewToken(TokLabel(TokRightCurly), string(TokRightCurly))
	case TokComma:
		tok, err = NewToken(TokLabel(TokComma), string(TokComma))
	case TokQuote:
		tok, err = NewToken(TokLabel(TokString), l.readString())
	case 0:
		tok, err = NewToken(TokLabel(TokEof), "")
		if err != nil {
			return tok, err
		}
	default:
		if isLetterSequence(l.char) {
			value := l.readLetterSequence()
			for _, tm := range tokMatcherSlice {
				tok := tm.Match(value)
				if tok != TokUnknown {
					return NewToken(tok, value)
				}
			}
			return NewToken(TokLabel(TokUnknown), value)
		} else {
			return NewToken(TokLabel(TokUnknown), string(l.char))
		}
	}

	l.readChar()
	return tok, err
}

func (l *Lexer) readChar() {
	l.char = l.peekChar()
	l.currentPosition = l.nextPosition
	l.nextPosition += 1
}

func (l *Lexer) readString() string {
	pos := l.currentPosition + 1
	for {
		l.readChar()
		if l.char == '"' || l.char == 0 {
			break
		}
	}
	return l.input[pos:l.currentPosition]
}

func (l *Lexer) readComment() string {
	pos := l.currentPosition + 1
	for {
		l.readChar()
		if l.char == '\n' {
			break
		}
	}
	return l.input[pos:l.currentPosition]
}

func (l *Lexer) readLetterSequence() string {
	pos := l.currentPosition
	for isLetterSequence(l.char) {
		l.readChar()
	}
	return l.input[pos:l.currentPosition]
}

func isLetterSequence(ch byte) bool {
	return isInteger(ch) ||
		'a' <= ch && ch <= 'z' ||
		'A' <= ch && ch <= 'Z' ||
		ch == '_' ||
		ch == ':' ||
		ch == '+' ||
		ch == '-' ||
		ch == '=' ||
		ch == '*' ||
		ch == '/' ||
		ch == '>' ||
		ch == '<' ||
		ch == '?' ||
		ch == '!'
}

func isInteger(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipWS() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.nextPosition]
	}
}
