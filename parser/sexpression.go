package parser

import (
	"fmt"
	"strings"
)

type Sexpression struct {
	Kind           TokLabel
	SubExpressions []*Sexpression
	Value          string
}

func (s Sexpression) Len() int {
	return len(s.SubExpressions)
}

func (s *Sexpression) IsAtomic() bool {
	return len(s.SubExpressions) == 0
}

func (s *Sexpression) Append(sub *Sexpression) bool {
	s.SubExpressions = append(s.SubExpressions, sub)
	return true
}

// ExpressionString produces a "clean" parseable expressoin string
func (s Sexpression) ExpressionString() string {
	return s.stringWithIndentation(0, true)
}

func (s Sexpression) String() string {
	return s.stringWithIndentation(0, false)
}

func (s Sexpression) DeepCopy() *Sexpression {
	if s.IsAtomic() {
		return NewAtomicSexpression(s.Kind, s.Value)
	}
	copy := NewSexpression(s.Kind, []*Sexpression{})
	for _, sub := range s.SubExpressions {
		copy.SubExpressions = append(copy.SubExpressions, sub.DeepCopy())
	}
	return copy
}

func (s Sexpression) stringWithIndentation(level int, valuesOnly bool) string {
	var sb strings.Builder

	if s.IsAtomic() {
		var sb strings.Builder
		if s.Kind == TokString {
			if valuesOnly {
				sb.WriteString(fmt.Sprintf("\"%s\" ", s.Value))
			} else {
				sb.WriteString(strings.Repeat(" ", level))
				sb.WriteString(fmt.Sprintf("(%s \"%s\")\n", s.Kind, s.Value))
			}
		} else {
			if valuesOnly {
				sb.WriteString(fmt.Sprintf("%s ", s.Value))
			} else {
				sb.WriteString(strings.Repeat(" ", level))
				sb.WriteString(fmt.Sprintf("(%s %s)\n", s.Kind, s.Value))
			}
		}
		return sb.String()
	}

	if valuesOnly {
		sb.WriteString(fmt.Sprintf("(%s %s", s.Kind, s.Value))
	} else {
		sb.WriteString(strings.Repeat(" ", level))
		sb.WriteString(fmt.Sprintf("(%s:%s\n", s.Kind, s.Value))
	}

	for _, sub := range s.SubExpressions {
		sb.WriteString(sub.stringWithIndentation(level+1, valuesOnly))
	}

	if valuesOnly {
		sb.WriteString(") ")
	} else {
		sb.WriteString(strings.Repeat(" ", level))
		sb.WriteString(")\n")
	}

	return sb.String()
}

func NewSexpression(label TokLabel, subexpression []*Sexpression) *Sexpression {
	return &Sexpression{
		label,
		subexpression,
		"",
	}
}

func NewAtomicSexpression(label TokLabel, value string) *Sexpression {
	return &Sexpression{
		label,
		[]*Sexpression{},
		value,
	}
}
