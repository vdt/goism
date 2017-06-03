package sexpconv

import (
	"exn"
	"go/types"
	"lisp"
	"lisp/function"
	"sexp"
)

func basicTypeZeroValue(typ *types.Basic) sexp.Form {
	switch typ.Kind() {
	case types.String:
		return sexp.Str("")
	case types.Bool:
		return sexp.Bool(false)

	default:
		info := typ.Info()

		if info&types.IsFloat != 0 {
			return sexp.Float(0)
		}
		if info&types.IsInteger != 0 {
			return sexp.Int(0)
		}
	}

	panic(exn.NoImpl("can not provide zero value for %#v", typ))
}

func ZeroValue(typ types.Type) sexp.Form {
	if lisp.Package != nil && types.Identical(typ, lisp.TypSymbol) {
		return sexp.Symbol{Val: "nil"}
	}

	switch typ := typ.(type) {
	case *types.Basic:
		return basicTypeZeroValue(typ)

	case *types.Array:
		zv := ZeroValue(typ.Elem())
		return &sexp.SparseArrayLit{
			Ctor: sexp.NewLispCall(function.MakeVector, sexp.Int(typ.Len()), zv),
			Typ:  typ,
		}

	case *types.Map:
		return nilMap

	case *types.Named:
		if typ, ok := typ.Underlying().(*types.Struct); ok {
			vals := make([]sexp.Form, typ.NumFields())
			for i := range vals {
				vals[i] = ZeroValue(typ.Field(i).Type())
			}
			return &sexp.StructLit{Vals: vals, Typ: typ}
		}
		return basicTypeZeroValue(typ.Underlying().(*types.Basic))
	}

	panic(exn.NoImpl("can not provide zero value for %#v", typ))
}

// Nil values
var (
	nilMap = sexp.Var{
		Name: "Go--nil-map",
		Typ:  types.NewMap(lisp.TypObject, lisp.TypObject),
	}

	nilFunc      = sexp.Var{Name: "Go--nil-function"}
	nilInterface = sexp.Var{Name: "Go--nil-interface"}
	nilSlice     = sexp.Var{Name: "Go--nil-slice"}
)
