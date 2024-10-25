package eval

import (
	"fmt"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
)

func init() {
	parser.HookToken(parser.TokEof)
	parser.HookToken(parser.TokLabel(parser.TokQuote))
	parser.HookToken(parser.TokIllegal)
	parser.HookToken(parser.TokLabel(parser.TokLeftBracket))
	parser.HookToken(parser.TokLabel(parser.TokRightBracket))
	parser.HookToken(parser.TokLabel(parser.TokLeftCurly))
	parser.HookToken(parser.TokLabel(parser.TokRightCurly))
	parser.HookToken(parser.TokLabel(parser.TokComma))
	parser.HookToken(parser.TokUnknown)
	parser.HookToken(parser.TokLabel(parser.TokLeftPar))
	parser.HookToken(parser.TokLabel(parser.TokRightPar))
}

// Function is an interface to add new functionality
type Function interface {
	// Desc returns documentation showing the structure of the command and a
	// description text
	Desc() (string, string)
	// Symbol returns the symbol under which the function is available
	Symbol() parser.TokLabel
	// Validate return an error if there is a mismatch between the parameters
	// on the stack and the function requirements
	Validate(env *Environment, stack *StackFrame) error
	// Evaluate implements the function semantics and returns a Result object
	// that wraps the results of the computation
	Evaluate(env *Environment, stack *StackFrame) (Result, error)
}

var funroot, _ = NewFunctionRoot()
var funcomment, _ = NewFunctionComment()
var funconcat, _ = NewFunctionConcat()
var fundesc, _ = NewFunctionDesc()
var fundef, _ = NewFunctionDef()
var funbind, _ = NewFunctionBindings()
var funresolv, _ = NewFunctionResolve()
var funtidy, _ = NewFunctionTidy()
var funpair, _ = NewFunctionPair()
var fundict, _ = NewFunctionDict()
var funvec, _ = NewFunctionVec()
var funadd, _ = NewFunctionAdd()
var funident, _ = NewFunctionIdentifier()
var funstr, _ = NewFunctionString()
var funint, _ = NewFunctionInt()
var funkw, _ = NewFunctionKeyword()
var funquote, _ = NewFunctionQuote()
var funeval, _ = NewFunctionEval()

var builtins = map[parser.TokLabel]Function{
	funcomment.Symbol(): funcomment,
	funconcat.Symbol():  funconcat,
	fundesc.Symbol():    fundesc,
	funroot.Symbol():    funroot,
	fundef.Symbol():     fundef,
	funbind.Symbol():    funbind,
	funresolv.Symbol():  funresolv,
	funtidy.Symbol():    funtidy,
	funpair.Symbol():    funpair,
	fundict.Symbol():    fundict,
	funvec.Symbol():     funvec,
	funadd.Symbol():     funadd,
	funident.Symbol():   funident,
	funstr.Symbol():     funstr,
	funint.Symbol():     funint,
	funkw.Symbol():      funkw,
	funquote.Symbol():   funquote,
	funeval.Symbol():    funeval,
}

func HookFunction(foo Function) error {
	if _, ok := builtins[foo.Symbol()]; ok {
		return fmt.Errorf("function %s already used", foo.Symbol())
	}
	builtins[foo.Symbol()] = foo
	return nil
}

func IsFunction(tok parser.TokLabel) bool {
	if _, ok := builtins[tok]; ok {
		return true
	}
	return false
}

func Validate(env *Environment, tok parser.TokLabel, stack *StackFrame) error {
	if builtin, ok := builtins[tok]; ok {
		return builtin.Validate(env, stack)
	}
	return fmt.Errorf("Function '%s' not avaiable", tok)
}

func Evaluate(env *Environment, tok parser.TokLabel, stack *StackFrame) (Result, error) {
	if builtin, ok := builtins[tok]; ok {
		return builtin.Evaluate(env, stack)
	}
	return nil, fmt.Errorf("Function '%s' not avaiable", tok)
}
