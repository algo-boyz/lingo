package macro

import (
	"log"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/traversal"
)

// MacroVisitor is the interface for macro transformations
type MacroVisitor interface {
	// Name of the macro
	Identifier() string
	// Transformation result
	Result() (*parser.Sexpression, error)
	// Called before the transformation
	Before()
	// Called when entering the nodes of the source S-expression
	Enter(expression *parser.Sexpression) (bool, error)
	// Called when leaving the nodes of the source S-expression
	Leave(expression *parser.Sexpression) error
	// Called after the transformation
	After()
}

var macros = map[string]MacroVisitor{}

func init() {
	HookMacro(NewThreadMacro(Last))
	HookMacro(NewThreadMacro(First))
	HookMacro(NewQuoteMacro())
}

func HookMacro(visitor MacroVisitor) {
	tokLabel := parser.TokLabel(visitor.Identifier())
	parser.HookToken(tokLabel)
	if _, ok := macros[visitor.Identifier()]; ok {
		log.Fatalf("macro with id %s already registered", visitor.Identifier())
	}
	macros[visitor.Identifier()] = visitor
}

func ApplyMacros(expression *parser.Sexpression) (*parser.Sexpression, error) {
	rewritten := expression
	for _, m := range macros {
		safe := rewritten.DeepCopy()
		visitor := traversal.NewExpressionWalker(m)
		err := visitor.Walk(rewritten)
		if err != nil {
			return nil, err
		}
		rewritten, err = m.Result()
		if err != nil {
			rewritten = safe
		}
	}
	return rewritten, nil
}
