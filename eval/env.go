package eval

import (
	"strings"
)

type Environment struct {
	Bindings map[string]Result
}

func NewEnvironment() *Environment {
	return &Environment{
		Bindings: map[string]Result{},
	}
}

func (env *Environment) Tidy() {
	for _, val := range env.Bindings {
		val.Tidy()
	}
	env.Bindings = map[string]Result{}
}

func (env *Environment) GC() {
	for sym, val := range env.Bindings {
		if !strings.HasPrefix(sym, "_v") {
			continue
		}
		val.Tidy()
		delete(env.Bindings, sym)
	}
}

func (env *Environment) DeepCopy() *Environment {
	nenv := NewEnvironment()
	for name, val := range env.Bindings {
		nenv.Bind(name, val.DeepCopy())
	}
	return nenv
}

type Stack struct {
	items []StackFrame
}

type StackFrame struct {
	items []Result
}

func (env *Environment) NewStack() *Stack {
	return &Stack{
		items: []StackFrame{},
	}
}

func (s *Stack) PushStackFrame() *StackFrame {
	sf := StackFrame{}
	s.items = append(s.items, sf)
	return &sf
}

func (s *Stack) Items() []StackFrame {
	return s.items
}

func (s *Stack) PopStackFrame() *StackFrame {
	if len(s.items) == 0 {
		return nil
	}
	lastItem := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return &lastItem
}

func (stack *Stack) PeekStackFrame() *StackFrame {
	if len(stack.items) == 0 {
		return nil
	}
	return &stack.items[len(stack.items)-1]
}

func (s *Stack) Empty() bool {
	return len(s.items) == 0
}

func (env *Environment) Bind(variable string, value Result) {
	env.Bindings[variable] = value
}

func (stack *StackFrame) Size() int {
	return len(stack.items)
}

func (stack *StackFrame) Push(item Result) {
	stack.items = append(stack.items, item)
}

func (stack *StackFrame) GetArgument(index int) Result {
	if len(stack.items) < index {
		return nil
	}
	return stack.items[index]
}

func (stack *StackFrame) Append(other *StackFrame) {
	stack.items = append(stack.items, other.items...)
}

func (stack *StackFrame) Pop() Result {
	if len(stack.items) == 0 {
		return nil
	}
	lastItem := stack.items[len(stack.items)-1]
	stack.items = stack.items[:len(stack.items)-1]
	return lastItem
}

func (s *StackFrame) Items() []Result {
	return s.items
}

func (stack *StackFrame) Empty() bool {
	return len(stack.items) == 0
}

func (stack *StackFrame) Reset() {
	stack.items = []Result{}
}

func (stack *StackFrame) Peek() Result {
	if len(stack.items) == 0 {
		return nil
	}
	return stack.items[len(stack.items)-1]
}
