package compiler

import (
	"fmt"
	"ir"
	"sexp"
)

func compileStmt(cl *Compiler, form sexp.Form) {
	switch form := form.(type) {
	case *sexp.Return:
		compileReturn(cl, form)
	case *sexp.If:
		compileIf(cl, form)
	case *sexp.Block:
		compileBlock(cl, form)
	case *sexp.FormList:
		compileStmtList(cl, form.Forms)
	case *sexp.Bind:
		compileBind(cl, form)
	case *sexp.Rebind:
		compileRebind(cl, form)
	case sexp.CallStmt:
		compileCallStmt(cl, form)
	case *sexp.Panic:
		compilePanic(cl, form)
	case *sexp.While:
		compileWhile(cl, form)

	default:
		panic(fmt.Sprintf("unexpected stmt: %#v\n", form))
	}
}

func compileExpr(cl *Compiler, form sexp.Form) {
	switch form := form.(type) {
	case sexp.Int:
		emit(cl, ir.ConstRef(cl.cvec.InsertInt(form.Val)))
	case sexp.Float:
		emit(cl, ir.ConstRef(cl.cvec.InsertFloat(form.Val)))
	case sexp.String:
		emit(cl, ir.ConstRef(cl.cvec.InsertString(form.Val)))
	case sexp.Symbol:
		emit(cl, ir.ConstRef(cl.cvec.InsertSym(form.Val)))
	case sexp.Bool:
		compileBool(cl, form)
	case sexp.Var:
		compileVar(cl, form)

	case *sexp.NumAdd:
		compileBinOp(cl, ir.NumAdd, form.Args)
	case *sexp.NumSub:
		compileBinOp(cl, ir.NumSub, form.Args)
	case *sexp.NumMul:
		compileBinOp(cl, ir.NumMul, form.Args)
	case *sexp.NumQuo:
		compileBinOp(cl, ir.NumQuo, form.Args)
	case *sexp.NumGt:
		compileBinOp(cl, ir.NumGt, form.Args)
	case *sexp.NumLt:
		compileBinOp(cl, ir.NumLt, form.Args)
	case *sexp.NumEq:
		compileBinOp(cl, ir.NumEq, form.Args)

	case *sexp.Call:
		compileCall(cl, form)
	case *sexp.MultiValueRef:
		compileMultiValueRef(cl, form)

	default:
		panic(fmt.Sprintf("unexpected expr: %#v\n", form))
	}
}