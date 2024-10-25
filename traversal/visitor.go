package traversal

import (
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
)

type ExpressionVisitor interface {
	// Setup
	Before()
	// Callback method that is invoked whenever we enter a node
	Enter(expression *parser.Sexpression) (bool, error)
	// Callback method that is invoked whenever we leave a node
	Leave(expression *parser.Sexpression) error
	// Teardown
	After()
}
