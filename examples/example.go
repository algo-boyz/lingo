package main

import (
	"fmt"
	"log"
	"strings"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/eval"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/types"
)

var TypeBoolId, TypeBool = types.NewTypeWithProperties("bool", types.Primitive)
var tokBool = parser.HookToken(parser.TokLabel(TypeBool.Name))

// recognize (true) als boolean
type BoolMatcher struct{}

func (i BoolMatcher) Match(s string) parser.TokLabel {
	if strings.ToLower(s) == "true" || strings.ToLower(s) == "false" {
		return tokBool.Label
	}
	return parser.TokUnknown
}
func (i BoolMatcher) Id() string {
	return string(tokBool.Label)
}

func init() {
	parser.HookMatcher(BoolMatcher{})
}

type BoolResult struct{ value bool }

func (r BoolResult) DeepCopy() eval.Result { return NewBoolResult(r.value) }
func (r BoolResult) String() string {
	if r.value {
		return "true"
	}
	return "false"
}
func (r BoolResult) Type() types.Type   { return TypeBool }
func (r BoolResult) Tidy()              {}
func (r BoolResult) Value() interface{} { return r.value }
func (r *BoolResult) SetValue(value interface{}) error {
	boolVal, ok := value.(bool)
	if !ok {
		return fmt.Errorf("invalid type for Bool")
	}
	r.value = boolVal
	return nil
}
func NewBoolResult(value bool) *BoolResult {
	return &BoolResult{
		value,
	}
}

type FunctionBool struct{}

func (f *FunctionBool) Desc() (string, string) {
	return fmt.Sprintf("%s%s %s%s",
			string(parser.TokLeftPar),
			f.Symbol(),
			"x",
			string(parser.TokRightPar)),
		"Converts a boolean symbol to a result [Internal]"
}
func (f *FunctionBool) Symbol() parser.TokLabel {
	return tokBool.Label
}
func (f *FunctionBool) Validate(env *eval.Environment, stack *eval.StackFrame) error {
	if stack.Size() != 1 {
		return eval.TooFewArgs(f.Symbol(), 0, 1)
	}
	if stack.Peek().Type() != types.TypeCharSequence {
		return eval.WrongTypeOfArg(f.Symbol(), 1, stack.Peek())
	}

	boolval := stack.Peek().(*eval.CharSequenceResult)
	boolvals := strings.ToLower(boolval.Val)

	if boolvals != "true" && boolvals != "false" {
		return fmt.Errorf("boolean value should be either true or false")
	}

	return nil
}
func (f *FunctionBool) Evaluate(env *eval.Environment, stack *eval.StackFrame) (eval.Result, error) {
	item := stack.Pop().(*eval.CharSequenceResult)

	boolval := item.Val
	boolval = strings.ToLower(boolval)

	var result bool
	if boolval == "true" {
		result = true
	} else if boolval == "false" {
		result = false
	}

	return NewBoolResult(result), nil
}

func NewFunctionBool() (eval.Function, error) {
	fun := &FunctionBool{}
	return fun, nil
}

type FunctionGt struct{}

func (f *FunctionGt) Desc() (string, string) {
	return fmt.Sprintf("%s%s %s%s",
			string(parser.TokLeftPar),
			f.Symbol(),
			"a b",
			string(parser.TokRightPar)),
		"Cmp if a > b numeric sub-expressions of a b"
}
func (f *FunctionGt) Symbol() parser.TokLabel {
	return parser.TokLabel("gt")
}
func (f *FunctionGt) Validate(env *eval.Environment, stack *eval.StackFrame) error {
	if stack.Empty() {
		return eval.TooFewArgs(f.Symbol(), 0, 1)
	}

	for idx, item := range stack.Items() {
		if item.Type() != types.TypeInt {
			return eval.WrongTypeOfArg(f.Symbol(), idx+1, item)
		}
	}
	return nil
}
func (f *FunctionGt) Evaluate(env *eval.Environment, stack *eval.StackFrame) (eval.Result, error) {
	result := NewBoolResult(false)
	for !stack.Empty() {
		itemA := stack.Pop().(*eval.IntResult)
		itemB := stack.Pop().(*eval.IntResult)
		result.SetValue(itemB.Val > itemA.Val)
	}

	return result, nil
}
func NewFunctionGt() (eval.Function, error) {
	fun := &FunctionGt{}
	parser.HookToken(fun.Symbol())
	return fun, nil
}

type FunctionAnd struct{}

func (f *FunctionAnd) Desc() (string, string) {
	return fmt.Sprintf("%s%s %s%s",
			string(parser.TokLeftPar),
			f.Symbol(),
			"p0 ... pN",
			string(parser.TokRightPar)),
		"apply logical and on sub-expressions"
}
func (f *FunctionAnd) Symbol() parser.TokLabel {
	return parser.TokLabel("and")
}
func (f *FunctionAnd) Validate(env *eval.Environment, stack *eval.StackFrame) error {
	if stack.Empty() {
		return eval.TooFewArgs(f.Symbol(), 0, 1)
	}

	for idx, item := range stack.Items() {
		if item.Type() != TypeBool {
			return eval.WrongTypeOfArg(f.Symbol(), idx+1, item)
		}
	}
	return nil
}
func (f *FunctionAnd) Evaluate(env *eval.Environment, stack *eval.StackFrame) (eval.Result, error) {
	result := true
	for !stack.Empty() {
		item := stack.Pop().(*BoolResult)
		result = result && item.value
	}

	return NewBoolResult(result), nil
}

func NewFunctionAnd() (eval.Function, error) {
	fun := &FunctionAnd{}
	parser.HookToken(fun.Symbol())
	return fun, nil
}

const Prompt = "\033[32m\033[0m > "

func main() {
	// Function to recognize new bool type
	fn, err := NewFunctionBool()
	if err != nil {
		log.Fatalf("failed to create bool function %s:", err.Error())
	}
	err = eval.HookFunction(fn)
	if err != nil {
		log.Fatalf("failed to hook bool function %s:", err.Error())
	}

	// Function to recognize new bool type
	fn, err = NewFunctionAnd()
	if err != nil {
		log.Fatalf("failed to create and function %s:", err.Error())
	}
	err = eval.HookFunction(fn)
	if err != nil {
		log.Fatalf("failed to hook and function %s:", err.Error())
	}

	fn, err = NewFunctionGt()
	if err != nil {
		log.Fatalf("failed to create gt function %s:", err.Error())
	}
	err = eval.HookFunction(fn)
	if err != nil {
		log.Fatalf("failed to hook gt function %s:", err.Error())
	}

	eval.RunLoop()
}
