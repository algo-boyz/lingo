package macro

import (
	"fmt"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
)

const (
	First = iota
	Last
)

type Mode int

func NewThreadMacro(mode Mode) *ThreadMacro {
	return &ThreadMacro{
		mode,
		parser.NewSexpression(parser.TokRoot, []*parser.Sexpression{}),
		parser.NewExpressionStack(),
	}
}

type ThreadMacro struct {
	mode  Mode
	sexp  *parser.Sexpression
	stack *parser.ExpressionStack
}

func (t *ThreadMacro) Identifier() string {
	switch t.mode {
	case First:
		return "->"
	case Last:
		return "->>"
	}
	return ""
}

func (t *ThreadMacro) Result() (*parser.Sexpression, error) {
	if len(t.sexp.SubExpressions) == 0 {
		return nil, fmt.Errorf("threading macro translation failed")
	}
	return t.sexp, nil
}

func (t *ThreadMacro) Before() {
	t.sexp = parser.NewSexpression(parser.TokRoot, []*parser.Sexpression{})
	t.stack = parser.NewExpressionStack()
}

func (t *ThreadMacro) After() {}

func (t *ThreadMacro) Enter(expression *parser.Sexpression) (bool, error) {
	if string(expression.Kind) == t.Identifier() {
		for _, sub := range expression.SubExpressions {
			copyXp := sub.DeepCopy()
			if t.stack.Empty() {
				t.stack.Push(copyXp)
				continue
			}

			last, err := t.stack.Pop()
			if err != nil {
				return false, err
			}
			if copyXp.IsAtomic() {
				return false, fmt.Errorf("threading macro not applicable to atomic expressions")
			}
			switch t.mode {
			case First:
				copyXp.SubExpressions = append([]*parser.Sexpression{last.DeepCopy()}, copyXp.SubExpressions...)
			case Last:
				copyXp.SubExpressions = append(copyXp.SubExpressions, last.DeepCopy())
			}
			t.stack.Push(copyXp)
		}
		return false, nil
	}
	return true, nil
}

func (t *ThreadMacro) Leave(expression *parser.Sexpression) error {
	if string(expression.Kind) == t.Identifier() {
		if t.stack.Len() != 1 {
			return fmt.Errorf("threading macro translation failed")
		}

		last, err := t.stack.Pop()
		if err != nil {
			return err
		}
		t.sexp.SubExpressions = append(t.sexp.SubExpressions, last)
	}
	return nil
}
