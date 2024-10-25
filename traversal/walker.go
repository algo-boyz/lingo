package traversal

import (
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
)

type ExpressionWalker struct {
	Visitor ExpressionVisitor
}

func NewExpressionWalker(v ExpressionVisitor) *ExpressionWalker {
	return &ExpressionWalker{
		Visitor: v,
	}
}

// Walk along the S-expression tree (DFS, LR)
func (c *ExpressionWalker) Walk(expression *parser.Sexpression) error {
	c.Visitor.Before()
	err := c.walk(expression)
	if err != nil {
		return err
	}
	c.Visitor.After()
	return nil
}

func (c *ExpressionWalker) walk(expression *parser.Sexpression) error {
	enter, err := c.Visitor.Enter(expression)
	if err != nil {
		return err
	}
	subexpressions := expression.SubExpressions
	if enter {
		for _, x := range subexpressions {
			err := c.walk(x)
			if err != nil {
				return err
			}
		}
	}
	err = c.Visitor.Leave(expression)
	if err != nil {
		return err
	}

	return nil
}
