package types

type LingoType uint8

// Function is a generic type interface
type Type struct {
	// Name of the type
	Name string
	// Unique ID
	Id int
	// PropertyMask to store boolean properties
	PropertyMask int32
}

const (
	Primitive = 1 << iota
	Collection
	Symbol
	Numeric
	Expression
)

var tidx = -1

func NewType(name string) (int, Type) {
	tidx += 1
	return tidx, Type{name, tidx, 0}
}

func NewTypeWithProperties(name string, propertyMask int32) (int, Type) {
	tidx += 1
	return tidx, Type{name, tidx, propertyMask}
}

func (t Type) HasProperty(property int32) bool {
	return (t.PropertyMask & property) == property
}

var (
	TypeNilId, TypeNil                   = NewType("nil")
	TypeUnknownId, TypeUnknown           = NewType("nil")
	TypeStringId, TypeString             = NewTypeWithProperties("string", Primitive)
	TypeKeywordId, TypeKeyword           = NewTypeWithProperties("keyword", Symbol)
	TypeIntId, TypeInt                   = NewTypeWithProperties("integer", Primitive)
	TypeSymbolId, TypeSymbol             = NewTypeWithProperties("symbol", Symbol)
	TypeDictionaryId, TypeDictionary     = NewTypeWithProperties("dict", Collection)
	TypeVectorId, TypeVector             = NewTypeWithProperties("vector", Collection)
	TypeCharSequenceId, TypeCharSequence = NewTypeWithProperties("charseq", Primitive)
	TypeSexpressionId, TypeSexpression   = NewTypeWithProperties("sexp", Expression)
)
