package macro

import (
	"fmt"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
)

func NewQuoteMacro() *QuoteMacro {
	return &QuoteMacro{
		parser.NewSexpression(parser.TokRoot, []*parser.Sexpression{}),
	}
}

type QuoteMacro struct {
	sexp *parser.Sexpression
}

func (t *QuoteMacro) Identifier() string {
	return string(parser.TokSingleQuoteChar)
}

func (t *QuoteMacro) Result() (*parser.Sexpression, error) {
	if len(t.sexp.SubExpressions) == 0 {
		return nil, fmt.Errorf("quote macro translation failed")
	}
	return t.sexp, nil
}

func (t *QuoteMacro) Before() {
	t.sexp = parser.NewSexpression(parser.TokRoot, []*parser.Sexpression{})
}

func (t *QuoteMacro) After() {}

func translate_rec(expression *parser.Sexpression) (*parser.Sexpression, error) {
	if expression.IsAtomic() {
		return expression.DeepCopy(), nil
	}

	newExp := parser.NewSexpression(expression.Kind, []*parser.Sexpression{})
	var quoted *parser.Sexpression = nil

	for _, sub := range expression.SubExpressions {
		copyXp := sub.DeepCopy()
		if copyXp.Kind == parser.TokLabel(parser.TokSingleQuoteChar) {
			quoted = parser.NewSexpression(parser.TokSingleQuote, []*parser.Sexpression{})
			continue
		}
		translated, err := translate_rec(copyXp)
		if err != nil {
			return nil, err
		}
		if quoted != nil {
			quoted.Append(translated)
			translated = quoted
		}
		newExp.Append(translated)
	}
	return newExp, nil
}

func (t *QuoteMacro) Enter(expression *parser.Sexpression) (bool, error) {
	var err error
	if string(expression.Kind) == parser.TokRoot {
		t.sexp, err = translate_rec(expression)
	}
	return false, err
}

func (t *QuoteMacro) Leave(expression *parser.Sexpression) error {
	return nil
}
