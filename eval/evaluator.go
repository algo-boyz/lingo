package eval

import (
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
)

func NewEvaluator(env *Environment) Evaluator {
	stack := env.NewStack()
	sf := stack.PushStackFrame()
	return Evaluator{
		false,
		stack,
		sf,
		env,
	}
}

type Evaluator struct {
	isBinding  bool
	stack      *Stack
	stackFrame *StackFrame
	env        *Environment
}

func (t *Evaluator) Result() Result {
	if t.stack.Empty() {
		return NewEmptyResult()
	}

	return t.stack.PeekStackFrame().Peek()
}

func (t *Evaluator) Before() {}
func (t *Evaluator) After() {
}

func (t *Evaluator) Enter(expression *parser.Sexpression) (bool, error) {
	if expression.Kind == fundef.Symbol() {
		t.isBinding = true
	}

	if IsFunction(expression.Kind) {
		t.stack.PushStackFrame()
		t.stackFrame = t.stack.PeekStackFrame()

		if expression.Kind == funquote.Symbol() {
			res := NewSexpressionResult(expression.DeepCopy())
			t.stackFrame.Push(res)
			//t.evaluate(expression.Kind)

			return false, nil
		}
	}

	return true, nil
}

func (t *Evaluator) evaluate(kind parser.TokLabel) (Result, error) {
	err := Validate(t.env, kind, t.stackFrame)
	if err != nil {
		return nil, err
	}

	res, err := Evaluate(t.env, kind, t.stackFrame)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (t *Evaluator) Leave(expression *parser.Sexpression) error {
	kind := expression.Kind
	if expression.IsAtomic() && expression.Value != "" {
		// for atomic expressions, we put the literal onto the stack first
		t.stackFrame.Push(NewCharSequenceResult(expression.Value))
	}

	// resolve the variable
	if !t.isBinding && kind == parser.TokIdentifier {
		kind = funresolv.Symbol()
	}

	if expression.Kind == fundef.Symbol() {
		t.isBinding = false
	}

	res, err := t.evaluate(kind)
	if err != nil {
		return err
	}

	t.stackFrame.Push(res)

	if IsFunction(expression.Kind) {
		last := t.stack.PopStackFrame()
		t.stackFrame = t.stack.PeekStackFrame()
		t.stackFrame.Append(last)
	}

	return nil
}
