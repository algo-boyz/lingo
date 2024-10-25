package parser

import (
	"fmt"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(input string) (*Sexpression, error) {
	root := NewSexpression(TokRoot, []*Sexpression{})

	stack := NewExpressionStack()
	stack.Push(root)
	lex := NewLexer(input)

	var err error
	var tok Tok

	for tok, err = lex.NextToken(); !tok.IsEof() && err == nil; tok, err = lex.NextToken() {
		parent, err := stack.Peek()
		if err != nil {
			return root, err
		}
		if tok.IsParOpen() {
			// descend, non-atomic
			nxttok, err := lex.NextToken()
			if err != nil {
				return root, err
			}
			if nxttok.IsUnknown() {
				return root, fmt.Errorf("unable to infer expression kind: %s", nxttok.value)
			}
			s := NewSexpression(TokLabel(nxttok.Label), []*Sexpression{})
			parent.Append(s)
			stack.Push(s)
			continue
		} else if tok.IsParClose() {
			// ascend, non-atomic
			stack.Pop()
			continue
		} else if tok.IsBracketOpen() {
			s := NewSexpression(TokLabel(TokVector), []*Sexpression{})
			parent.Append(s)
			stack.Push(s)
			continue
		} else if tok.IsBracketClose() {
			// ascend, non-atomic
			stack.Pop()
			continue
		} else if tok.IsCurlyOpen() {
			pair := NewSexpression(TokLabel(TokPair), []*Sexpression{})
			s := NewSexpression(TokLabel(TokDict), []*Sexpression{pair})
			parent.Append(s)
			stack.Push(s)
			stack.Push(pair)
			continue
		} else if tok.IsComma() {
			if parent.Kind != TokPair {
				return nil, fmt.Errorf("invalid dictionary format %s", parent.Kind)
			}
			stack.Pop()
			parent, err = stack.Peek()
			if err != nil {
				return nil, err
			}
			if parent.Kind != TokDict {
				return nil, fmt.Errorf("invalid dictionary format %s", parent.Kind)
			}
			pair := NewSexpression(TokLabel(TokPair), []*Sexpression{})
			parent.Append(pair)
			stack.Push(pair)
			continue
		} else if tok.IsCurlyClose() {
			parent, err := stack.Peek()
			if err != nil {
				return nil, err
			}
			if parent.Kind != TokPair {
				return nil, fmt.Errorf("invalid dictionary format")
			}
			stack.Pop()
			parent, err = stack.Peek()
			if err != nil {
				return nil, err
			}
			if parent.Kind != TokDict {
				return nil, fmt.Errorf("invalid dictionary format")
			}
			stack.Pop()
			continue
		}

		if tok.IsUnknown() {
			return root, fmt.Errorf("unable to infer expression kind: %s", tok.value)
		}

		expValue := tok.value
		expKind := tok.Label
		// get value
		atomic := NewAtomicSexpression(expKind, expValue)
		parent.Append(atomic)
	}

	// something went wrong in the parsing loop
	if err != nil {
		return nil, err
	}

	if stack.Len() != 1 {
		return nil, fmt.Errorf("expression invalid")
	}

	return root, err
}
