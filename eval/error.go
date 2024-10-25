package eval

import (
	"fmt"
	"strings"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
)

func TooManyArgs(lbl parser.TokLabel, provided int, expected int) error {
	return fmt.Errorf(
		ErrorMessage(lbl, "Wrong number of arguments (%d while at least %d arguments are expected)"),
		provided,
		expected)
}

func TooFewArgs(lbl parser.TokLabel, provided int, expected int) error {
	return fmt.Errorf(
		ErrorMessage(lbl, "Wrong number of arguments (%d while at least %d arguments are expected)"),
		provided,
		expected)
}

func WrongNumberOfArgs(lbl parser.TokLabel, provided int, expected int) error {
	return fmt.Errorf(
		ErrorMessage(lbl, "Wrong number of arguments (%d instead of %d)"),
		provided,
		expected)
}

func WrongTypeOfArg(lbl parser.TokLabel, paramIdx int, param Result) error {
	msg := "unsupported type '%s' for param %d"
	return fmt.Errorf(ErrorMessage(lbl, msg), param.Type().Name, paramIdx)
}

func CheckVariableName(lbl parser.TokLabel, paramIdx int, param Result) error {
	// only hidden variables
	if strings.HasPrefix(param.String(), "_") {
		return InvalidVarName(lbl, paramIdx, param)
	}
	return nil
}

func InvalidVarName(lbl parser.TokLabel, paramIdx int, param Result) error {
	msg := "invalid variable name '%s' for param %d"
	return fmt.Errorf(ErrorMessage(lbl, msg), param.Type().Name, paramIdx)
}

func ErrorMessage(lbl parser.TokLabel, msg string) string {
	return fmt.Sprintf("(%s) %s", lbl, msg)
}
