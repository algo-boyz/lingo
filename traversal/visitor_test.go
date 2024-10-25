package traversal

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
)

func TestVisitor(t *testing.T) {
	//(root
	// (match
	//  (name
	//   (identifier X)
	//  )
	//  (string "(exec|match)")
	// )
	//)
	expression := parser.NewSexpression(
		"root",
		[]*parser.Sexpression{
			parser.NewSexpression(
				"def",
				[]*parser.Sexpression{
					parser.NewAtomicSexpression(parser.TokIdentifier, "a"),
					parser.NewAtomicSexpression(parser.TokString, "b"),
				},
			),
			parser.NewSexpression(
				"match",
				[]*parser.Sexpression{
					parser.NewSexpression(
						"name",
						[]*parser.Sexpression{
							parser.NewAtomicSexpression(parser.TokIdentifier, "a"),
						},
					),
					parser.NewAtomicSexpression(parser.TokString, "(exec|match)"),
				},
			),
		},
	)

	visitor := DemoVisitor{[]string{}, []string{}}
	walker := NewExpressionWalker(&visitor)
	walker.Walk(expression)
	wantBefore := []string{parser.TokRoot, "def", parser.TokIdentifier, parser.TokString, "match", "name", parser.TokIdentifier, parser.TokString}
	assert.Equal(t, wantBefore, visitor.trackBefore)
	wantAfter := []string{parser.TokIdentifier, parser.TokString, "def", parser.TokIdentifier, "name", parser.TokString, "match", parser.TokRoot}
	assert.Equal(t, wantAfter, visitor.trackAfter)
}

type DemoVisitor struct {
	trackBefore []string
	trackAfter  []string
}

func (t *DemoVisitor) Before() {}

func (t *DemoVisitor) After() {}

func (t *DemoVisitor) Enter(expression *parser.Sexpression) (bool, error) {
	t.trackBefore = append(t.trackBefore, string(expression.Kind))
	return true, nil
}

func (t *DemoVisitor) Leave(expression *parser.Sexpression) error {
	t.trackAfter = append(t.trackAfter, string(expression.Kind))
	return nil
}
