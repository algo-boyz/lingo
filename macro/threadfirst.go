package macro

import (
	"fmt"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
)

func NewThreadFirstMacro() *ThreadFirstMacro {
	return &ThreadFirstMacro{
		parser.NewSexpression(parser.TokRoot, []*parser.Sexpression{}),
		parser.NewExpressionStack(),
	}
}

type ThreadFirstMacro struct {
	sexp  *parser.Sexpression
	stack *parser.ExpressionStack
}

func (t *ThreadFirstMacro) Identifier() string {
	return "->"
}

func (t *ThreadFirstMacro) Result() (*parser.Sexpression, error) {
	if len(t.sexp.SubExpressions) == 0 {
		return nil, fmt.Errorf("threading macro translation failed")
	}
	return t.sexp, nil
}

func (t *ThreadFirstMacro) Before() {
	t.sexp = parser.NewSexpression(parser.TokRoot, []*parser.Sexpression{})
	t.stack = parser.NewExpressionStack()
}

func (t *ThreadFirstMacro) After() {}

func (t *ThreadFirstMacro) Enter(expression *parser.Sexpression) (bool, error) {
	if expression.Kind == "->" {
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
			copyXp.SubExpressions = append([]*parser.Sexpression{last.DeepCopy()}, copyXp.SubExpressions...)
			t.stack.Push(copyXp)
		}
		return false, nil
	}
	return true, nil
}

func (t *ThreadFirstMacro) Leave(expression *parser.Sexpression) error {
	if expression.Kind == "->" {
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
