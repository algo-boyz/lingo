package parser

import (
	"fmt"
)

func NewExpressionStack() *ExpressionStack {
	return &ExpressionStack{}
}

type ExpressionStack struct {
	data []*Sexpression
}

func (s *ExpressionStack) Push(item *Sexpression) {
	s.data = append(s.data, item)
}

func (s *ExpressionStack) Pop() (*Sexpression, error) {
	idx := len(s.data) - 1
	if idx >= 0 {
		val := s.data[idx]
		s.data = s.data[:idx]
		return val, nil
	}
	return nil, fmt.Errorf("pop on an empty stack")
}

func (s *ExpressionStack) Len() int {
	return len(s.data)
}

func (s *ExpressionStack) Peek() (*Sexpression, error) {
	idx := len(s.data) - 1
	if idx >= 0 {
		return s.data[idx], nil
	}
	return nil, fmt.Errorf("peek on an empty stack")
}

func (s *ExpressionStack) Empty() bool {
	return len(s.data) == 0
}
