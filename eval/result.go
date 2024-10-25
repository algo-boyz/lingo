package eval

import (
	"bytes"
	"fmt"
	"strings"

	"strconv"

	"github.com/olekukonko/tablewriter"

	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/parser"
	"gitlab.com/gitlab-org/vulnerability-research/foss/lingo/types"
)

// Result represents values computed by the evalation of S-epxressions
type Result interface {
	DeepCopy() Result
	String() string
	Type() types.Type
	Value() interface{}
	SetValue(interface{}) error
	Tidy()
}

// Symbol is used to identify variables or functions
type CharSequenceResult struct{ Val string }

func (r CharSequenceResult) DeepCopy() Result   { return NewCharSequenceResult(r.Val) }
func (r CharSequenceResult) String() string     { return r.Val }
func (r CharSequenceResult) Type() types.Type   { return types.TypeCharSequence }
func (r CharSequenceResult) Tidy()              {}
func (r CharSequenceResult) Value() interface{} { return r.Val }
func (r *CharSequenceResult) SetValue(value interface{}) error {
	strval, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid type for String")
	}
	r.Val = strval
	return nil
}
func NewCharSequenceResult(label string) *CharSequenceResult {
	return &CharSequenceResult{
		label,
	}
}

// Symbol is used to identify variables or functions
type SymbolResult struct{ Val string }

func (r SymbolResult) DeepCopy() Result   { return NewSymbolResult(r.Val) }
func (r SymbolResult) String() string     { return r.Val }
func (r SymbolResult) Type() types.Type   { return types.TypeSymbol }
func (r SymbolResult) Tidy()              {}
func (r SymbolResult) Value() interface{} { return r.Val }
func (r *SymbolResult) SetValue(value interface{}) error {
	strval, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid type for String")
	}
	r.Val = strval
	return nil
}
func NewSymbolResult(label string) *SymbolResult {
	return &SymbolResult{
		label,
	}
}

// IntResult represents an int
type IntResult struct{ Val int }

func (r IntResult) DeepCopy() Result   { return NewIntResult(r.Val) }
func (r IntResult) String() string     { return strconv.Itoa(r.Val) }
func (r IntResult) Type() types.Type   { return types.TypeInt }
func (r IntResult) Tidy()              {}
func (r IntResult) Value() interface{} { return r.Val }
func (r *IntResult) SetValue(value interface{}) error {
	intval, ok := value.(int)
	if !ok {
		return fmt.Errorf("invalid type for Int")
	}
	r.Val = intval
	return nil
}
func NewIntResult(value int) *IntResult {
	return &IntResult{
		value,
	}
}

// StringResult is used for strings
type StringResult struct{ Val string }

func (r StringResult) DeepCopy() Result   { return NewStringResult(r.Val) }
func (r StringResult) String() string     { return r.Val }
func (r StringResult) Type() types.Type   { return types.TypeString }
func (r StringResult) Tidy()              {}
func (r StringResult) Value() interface{} { return r.Val }
func (r *StringResult) SetValue(value interface{}) error {
	strval, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid type for string")
	}
	r.Val = strval
	return nil
}
func NewStringResult(value string) *StringResult { return &StringResult{value} }

// DictResult represents a dictionary, map
type DictResult struct {
	Header []string
	Values map[string][]Result
}

func (r DictResult) DeepCopy() Result {
	dict := NewDictResult(r.Header...)
	for key, values := range r.Values {
		results := []Result{}
		for _, result := range values {
			results = append(results, result.DeepCopy())
		}
		dict.Values[key] = results
	}

	return dict
}

func (r DictResult) String() string {
	writer := bytes.NewBufferString("")
	table := tablewriter.NewWriter(writer)

	var rows [][]string

	for y, h := range r.Header {
		v := r.Values[h]

		if y == 0 {
			// we assume that row len is the same for all entries
			rows = make([][]string, len(v))
		}

		for x, item := range v {
			if len(rows[x]) == 0 {
				rows[x] = make([]string, len(r.Values))
			}
			rows[x][y] = item.String()
		}
	}
	table.AppendBulk(rows)
	table.SetRowLine(true)
	table.SetHeader(r.Header)
	table.Render()

	return writer.String()
}

func (r DictResult) Type() types.Type   { return types.TypeDictionary }
func (r DictResult) Tidy()              {}
func (r DictResult) Value() interface{} { return r.Values }
func (r *DictResult) AddPair(k string, v ...Result) error {
	if _, ok := r.Values[k]; !ok {
		return fmt.Errorf("key '%s' not present", k)
	}

	r.Values[k] = append(r.Values[k], v...)
	return nil
}
func (r *DictResult) Merge(other *DictResult) *DictResult {
	hmap := map[string]string{}
	headers := []string{}

	for _, h := range append(r.Header, other.Header...) {
		if _, ok := hmap[h]; ok {
			continue
		}
		headers = append(headers, h)
		hmap[h] = h
	}

	dict := NewDictResult(headers...)

	for _, h := range headers {
		if val, ok := r.Values[h]; ok {
			dict.Values[h] = val
		}
		if val, ok := other.Values[h]; ok {
			dict.Values[h] = val
		}
	}

	return dict
}
func (r *DictResult) SetValue(value interface{}) error {
	dictval, ok := value.(DictResult)
	if !ok {
		return fmt.Errorf("invalid type for dictionary")
	}
	r.Header = dictval.Header
	r.Values = dictval.Values
	return nil
}
func NewDictResult(head ...string) *DictResult {
	d := map[string][]Result{}
	headers := []string{}

	for _, h := range head {
		d[h] = []Result{}
		headers = append(headers, h)
	}

	return &DictResult{
		headers,
		d,
	}
}

func NewVecResult() *VecResult {
	return &VecResult{
		[]Result{},
	}
}

// VecResult represents a vector
type VecResult struct {
	Data []Result
}

func (r VecResult) DeepCopy() Result {
	vec := NewVecResult()
	for _, item := range r.Data {
		vec.Data = append(vec.Data, item.DeepCopy())
	}
	return vec
}

func (r VecResult) String() string {
	var buffer strings.Builder

	for idx, item := range r.Data {
		if idx > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(item.String())
	}

	return buffer.String()
}
func (r VecResult) Type() types.Type   { return types.TypeVector }
func (r VecResult) Tidy()              {}
func (r VecResult) Value() interface{} { return r.Data }
func (r *VecResult) AppendResult(value Result) error {
	r.Data = append(r.Data, value)
	return nil
}
func (r *VecResult) PrependResult(value Result) error {
	r.Data = append([]Result{value}, r.Data...)
	return nil
}
func (r *VecResult) SetValue(value interface{}) error {
	data, ok := value.([]Result)
	if !ok {
		return fmt.Errorf("invalid type for vector")
	}
	r.Data = data
	return nil
}

func NewEmptyResult() *EmptyResult { return &EmptyResult{} }

// EmptyResult represents the empty set a.k.a. as null or nil
type EmptyResult struct{}

func (r EmptyResult) DeepCopy() Result   { return NewEmptyResult() }
func (r EmptyResult) String() string     { return "nil" }
func (r EmptyResult) Type() types.Type   { return types.TypeNil }
func (r EmptyResult) Tidy()              {}
func (r EmptyResult) Value() interface{} { return nil }
func (r *EmptyResult) SetValue(value interface{}) error {
	return nil
}

// KeywordResult is used for dictionary/hash keys
type KeyWordResult struct{ Val string }

func (r KeyWordResult) DeepCopy() Result   { return NewKeyWordResult(r.Val) }
func (r KeyWordResult) String() string     { return r.Val }
func (r KeyWordResult) Type() types.Type   { return types.TypeKeyword }
func (r KeyWordResult) Tidy()              {}
func (r KeyWordResult) Value() interface{} { return r.Val }
func (r *KeyWordResult) SetValue(value interface{}) error {
	strval, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid type for String")
	}
	r.Val = strval
	return nil
}

func NewKeyWordResult(value string) *KeyWordResult { return &KeyWordResult{value} }

// SexpressoinResult is used for dictionary/hash keys
type SexpressionResult struct{ Exp *parser.Sexpression }

func (r SexpressionResult) DeepCopy() Result   { return NewSexpressionResult(r.Exp.DeepCopy()) }
func (r SexpressionResult) String() string     { return r.Exp.String() }
func (r SexpressionResult) Type() types.Type   { return types.TypeSexpression }
func (r SexpressionResult) Tidy()              {}
func (r SexpressionResult) Value() interface{} { return r.Exp }
func (r *SexpressionResult) SetValue(value interface{}) error {
	exp, ok := value.(*parser.Sexpression)
	if !ok {
		return fmt.Errorf("invalid type for String")
	}
	r.Exp = exp
	return nil
}

func NewSexpressionResult(expression *parser.Sexpression) *SexpressionResult {
	return &SexpressionResult{expression}
}
