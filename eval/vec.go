package eval

import (
	"fmt"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/types"
)

type FunctionVec struct{}

func (f *FunctionVec) Desc() (string, string) {
	return fmt.Sprintf("%c%s %s%c",
			parser.TokLeftPar,
			f.Symbol(),
			"(p0) ... (pN)",
			parser.TokRightPar),
		"Generate a vector from (string|int|keyword) subexpressions pX [Internal]"
}

func (f *FunctionVec) Symbol() parser.TokLabel {
	return parser.TokLabel(parser.TokVector)
}

func (f *FunctionVec) Validate(env *Environment, stack *StackFrame) error {
	if stack.Empty() {
		return TooFewArgs(f.Symbol(), 0, 1)
	}
	for idx, item := range stack.items {
		switch item.Type().Id {
		case types.TypeKeywordId, types.TypeStringId, types.TypeIntId:
			continue
		default:
			return WrongTypeOfArg(f.Symbol(), idx+1, item)
		}
	}
	return nil
}

func (f *FunctionVec) Evaluate(env *Environment, stack *StackFrame) (Result, error) {
	result := NewVecResult()

	for !stack.Empty() {
		item := stack.Pop()
		if err := result.PrependResult(item); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func NewFunctionVec() (Function, error) {
	fun := &FunctionVec{}
	parser.HookToken(fun.Symbol())
	return fun, nil
}
