package main

import (
	bindings "github.com/KyrylR/jack-compiler/parser"
	"github.com/antlr4-go/antlr/v4"
	"strings"
)

type XMLVisitor struct {
	*bindings.BaseJackVisitor
	parser  antlr.Parser
	Builder *strings.Builder
}

func NewXMLVisitor(parser antlr.Parser) *XMLVisitor {
	return &XMLVisitor{
		parser:  parser,
		Builder: &strings.Builder{},
	}
}

func (v *XMLVisitor) Visit(ctx antlr.ParseTree) interface{} {
	return ctx.Accept(v)
}

func (v *XMLVisitor) VisitProgram(ctx *bindings.ProgramContext) interface{} {
	for _, classDec := range ctx.AllClassDeclaration() {
		v.Visit(classDec)
	}

	return nil
}

func (v *XMLVisitor) VisitClassDeclaration(ctx *bindings.ClassDeclarationContext) interface{} {
	v.Builder.WriteString("<class>\n")
	v.Builder.WriteString("<keyword>")
	v.Builder.WriteString(ctx.CLASS().GetText())
	v.Builder.WriteString("</keyword>\n")

	v.Visit(ctx.ClassName())

	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.LC().GetText())
	v.Builder.WriteString("</symbol>\n")

	for _, classVarDec := range ctx.AllClassVarDec() {
		v.Visit(classVarDec)
	}

	for _, subroutineDec := range ctx.AllSubroutineDec() {
		v.Visit(subroutineDec)
	}

	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.RC().GetText())
	v.Builder.WriteString("</symbol>\n")

	v.Builder.WriteString("</class>\n")

	return nil
}

func (v *XMLVisitor) VisitClassVarDec(ctx *bindings.ClassVarDecContext) interface{} {
	v.Builder.WriteString("<classVarDec>\n")
	v.Builder.WriteString("<keyword>")
	v.Builder.WriteString(ctx.GetClassVarType().GetText())
	v.Builder.WriteString("</keyword>\n")
	v.Builder.WriteString("<keyword>")
	v.Builder.WriteString(ctx.Type_().GetText())
	v.Builder.WriteString("</keyword>\n")

	for i, varName := range ctx.AllVarName() {
		v.Visit(varName)
		if i < len(ctx.AllVarName())-1 {
			v.Builder.WriteString("<symbol>")
			v.Builder.WriteString(ctx.COMMA(i).GetText())
			v.Builder.WriteString("</symbol>\n")
		}
	}

	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.SEMICOLON().GetText())
	v.Builder.WriteString("</symbol>\n")

	v.Builder.WriteString("</classVarDec>\n")

	return nil
}

func (v *XMLVisitor) VisitType(ctx *bindings.TypeContext) interface{} {
	v.Builder.WriteString("<keyword>")
	v.Builder.WriteString(ctx.GetText())
	v.Builder.WriteString("</keyword>\n")

	return nil
}

func (v *XMLVisitor) VisitSubroutineDec(ctx *bindings.SubroutineDecContext) interface{} {
	v.Builder.WriteString("<subroutineDec>\n")
	v.Builder.WriteString("<keyword>")
	v.Builder.WriteString(ctx.GetSubroutineDecType().GetText())
	v.Builder.WriteString("</keyword>\n")

	probablyVoid := ctx.VOID()
	probablyType := ctx.Type_()

	if probablyVoid != nil {
		v.Builder.WriteString("<keyword>")
		v.Builder.WriteString(probablyVoid.GetText())
		v.Builder.WriteString("</keyword>\n")
	} else {
		v.Visit(probablyType)
	}

	v.Visit(ctx.SubroutineName())

	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.LP().GetText())
	v.Builder.WriteString("</symbol>\n")

	if ctx.ParameterList() != nil {
		v.Visit(ctx.ParameterList())
	} else {
		v.Builder.WriteString("<parameterList>\n</parameterList>\n")
	}

	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.RP().GetText())
	v.Builder.WriteString("</symbol>\n")

	v.Visit(ctx.SubroutineBody())

	v.Builder.WriteString("</subroutineDec>\n")

	return nil
}

func (v *XMLVisitor) VisitParameterList(ctx *bindings.ParameterListContext) interface{} {
	v.Builder.WriteString("<parameterList>")

	types := ctx.AllType_()
	varNames := ctx.AllVarName()
	commas := ctx.AllCOMMA()

	for i := 0; i < len(types); i++ {
		v.Visit(types[i])
		v.Visit(varNames[i])
		if i < len(commas)-1 {
			v.Builder.WriteString("<symbol>")
			v.Builder.WriteString(commas[i].GetText())
			v.Builder.WriteString("</symbol>\n")
		}
	}

	v.Builder.WriteString("</parameterList>\n")

	return nil
}

func (v *XMLVisitor) VisitSubroutineBody(ctx *bindings.SubroutineBodyContext) interface{} {
	v.Builder.WriteString("<subroutineBody>\n")
	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.LC().GetText())
	v.Builder.WriteString("</symbol>\n")

	for _, varDec := range ctx.AllVarDec() {
		v.Visit(varDec)
	}

	v.Visit(ctx.Statements())

	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.RC().GetText())
	v.Builder.WriteString("</symbol>\n")
	v.Builder.WriteString("</subroutineBody>\n")

	return nil
}

func (v *XMLVisitor) VisitVarDec(ctx *bindings.VarDecContext) interface{} {
	v.Builder.WriteString("<varDec>\n")
	v.Builder.WriteString("<keyword>")
	v.Builder.WriteString(ctx.VAR().GetText())
	v.Builder.WriteString("</keyword>\n")

	v.Visit(ctx.Type_())

	commas := ctx.AllCOMMA()

	for i, varName := range ctx.AllVarName() {
		v.Visit(varName)
		if i < len(commas)-1 {
			v.Builder.WriteString("<symbol>")
			v.Builder.WriteString(commas[i].GetText())
			v.Builder.WriteString("</symbol>\n")
		}
	}

	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.SEMICOLON().GetText())
	v.Builder.WriteString("</symbol>\n")

	v.Builder.WriteString("</varDec>\n")

	return nil
}

func (v *XMLVisitor) VisitClassName(ctx *bindings.ClassNameContext) interface{} {
	v.Builder.WriteString("<identifier>")
	v.Builder.WriteString(ctx.GetText())
	v.Builder.WriteString("</identifier>\n")

	return nil
}

func (v *XMLVisitor) VisitSubroutineName(ctx *bindings.SubroutineNameContext) interface{} {
	v.Builder.WriteString("<identifier>")
	v.Builder.WriteString(ctx.GetText())
	v.Builder.WriteString("</identifier>\n")

	return nil
}

func (v *XMLVisitor) VisitVarName(ctx *bindings.VarNameContext) interface{} {
	v.Builder.WriteString("<identifier>")
	v.Builder.WriteString(ctx.GetText())
	v.Builder.WriteString("</identifier>\n")

	return nil
}

func (v *XMLVisitor) VisitStatements(ctx *bindings.StatementsContext) interface{} {
	v.Builder.WriteString("<statements>\n")

	for _, statement := range ctx.AllStatement() {
		v.Visit(statement)
	}

	v.Builder.WriteString("</statements>\n")

	return nil
}

func (v *XMLVisitor) VisitStatement(ctx *bindings.StatementContext) interface{} {
	switch ctx.GetChild(0).(type) {
	case *bindings.LetStatementContext:
		v.Visit(ctx.LetStatement())
	case *bindings.IfStatementContext:
		v.Visit(ctx.IfStatement())
	case *bindings.WhileStatementContext:
		v.Visit(ctx.WhileStatement())
	case *bindings.DoStatementContext:
		v.Visit(ctx.DoStatement())
	case *bindings.ReturnStatementContext:
		v.Visit(ctx.ReturnStatement())
	default:
		panic("Unknown statement type")
	}

	return nil
}

func (v *XMLVisitor) VisitLetStatement(ctx *bindings.LetStatementContext) interface{} {
	v.Builder.WriteString("<letStatement>\n")

	v.Builder.WriteString("<keyword>")
	v.Builder.WriteString(ctx.LET().GetText())
	v.Builder.WriteString("</keyword>\n")

	v.Visit(ctx.VarName())

	if ctx.SquareBracketExpression() != nil {
		v.Visit(ctx.SquareBracketExpression())
	}

	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.EQ().GetText())
	v.Builder.WriteString("</symbol>\n")

	v.Visit(ctx.Expression())

	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.SEMICOLON().GetText())
	v.Builder.WriteString("</symbol>\n")

	v.Builder.WriteString("</letStatement>\n")

	return nil
}

func (v *XMLVisitor) VisitSquareBracketExpression(ctx *bindings.SquareBracketExpressionContext) interface{} {
	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.LB().GetText())
	v.Builder.WriteString("</symbol>\n")

	v.Visit(ctx.Expression())

	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.RB().GetText())
	v.Builder.WriteString("</symbol>\n")

	return nil
}

func (v *XMLVisitor) VisitIfStatement(ctx *bindings.IfStatementContext) interface{} {
	v.Builder.WriteString("<ifStatement>\n")

	v.Builder.WriteString("<keyword>")
	v.Builder.WriteString(ctx.IF().GetText())
	v.Builder.WriteString("</keyword>\n")

	v.Visit(ctx.ParExpression())
	v.Visit(ctx.CurlyStatements())

	if ctx.ElseBlock() != nil {
		v.Visit(ctx.ElseBlock())
	}

	v.Builder.WriteString("</ifStatement>\n")

	return nil
}

func (v *XMLVisitor) VisitCurlyStatements(ctx *bindings.CurlyStatementsContext) interface{} {
	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.LC().GetText())
	v.Builder.WriteString("</symbol>\n")

	v.Visit(ctx.Statements())

	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.RC().GetText())
	v.Builder.WriteString("</symbol>\n")

	return nil
}

func (v *XMLVisitor) VisitParExpression(ctx *bindings.ParExpressionContext) interface{} {
	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.LP().GetText())
	v.Builder.WriteString("</symbol>\n")

	v.Visit(ctx.Expression())

	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.RP().GetText())
	v.Builder.WriteString("</symbol>\n")

	return nil
}

func (v *XMLVisitor) VisitElseBlock(ctx *bindings.ElseBlockContext) interface{} {
	v.Builder.WriteString("<keyword>")
	v.Builder.WriteString(ctx.ELSE().GetText())
	v.Builder.WriteString("</keyword>\n")

	v.Visit(ctx.CurlyStatements())

	return nil
}

func (v *XMLVisitor) VisitWhileStatement(ctx *bindings.WhileStatementContext) interface{} {
	v.Builder.WriteString("<whileStatement>\n")

	v.Builder.WriteString("<keyword>")
	v.Builder.WriteString(ctx.WHILE().GetText())
	v.Builder.WriteString("</keyword>\n")

	v.Visit(ctx.ParExpression())
	v.Visit(ctx.CurlyStatements())

	v.Builder.WriteString("</whileStatement>\n")

	return nil
}

func (v *XMLVisitor) VisitDoStatement(ctx *bindings.DoStatementContext) interface{} {
	v.Builder.WriteString("<doStatement>\n")

	v.Builder.WriteString("<keyword>")
	v.Builder.WriteString(ctx.DO().GetText())
	v.Builder.WriteString("</keyword>\n")

	v.Visit(ctx.SubroutineCall())

	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.SEMICOLON().GetText())
	v.Builder.WriteString("</symbol>\n")

	v.Builder.WriteString("</doStatement>\n")

	return nil
}

func (v *XMLVisitor) VisitReturnStatement(ctx *bindings.ReturnStatementContext) interface{} {
	v.Builder.WriteString("<returnStatement>\n")

	v.Builder.WriteString("<keyword>")
	v.Builder.WriteString(ctx.RETURN().GetText())
	v.Builder.WriteString("</keyword>\n")

	if ctx.Expression() != nil {
		v.Visit(ctx.Expression())
	}

	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.SEMICOLON().GetText())
	v.Builder.WriteString("</symbol>\n")

	v.Builder.WriteString("</returnStatement>\n")

	return nil
}

func (v *XMLVisitor) VisitExpression(ctx *bindings.ExpressionContext) interface{} {
	v.Builder.WriteString("<expression>\n")

	for i, term := range ctx.AllTerm() {
		v.Visit(term)
		if i < len(ctx.AllOp()) {
			v.Visit(ctx.AllOp()[i])
		}
	}

	v.Builder.WriteString("</expression>\n")

	return nil
}

func (v *XMLVisitor) WrapTerm(ctx antlr.ParserRuleContext) interface{} {
	v.Builder.WriteString("<term>\n")
	v.Visit(ctx)
	v.Builder.WriteString("</term>\n")

	return nil
}

func (v *XMLVisitor) VisitIntegerConstantTerm(ctx *bindings.IntegerConstantTermContext) interface{} {
	return v.WrapTerm(ctx.IntegerConstant())
}

func (v *XMLVisitor) VisitStringConstantTerm(ctx *bindings.StringConstantTermContext) interface{} {
	return v.WrapTerm(ctx.StringConstant())
}

func (v *XMLVisitor) VisitKeywordConstantTerm(ctx *bindings.KeywordConstantTermContext) interface{} {
	return v.WrapTerm(ctx.KeywordConstant())
}

func (v *XMLVisitor) VisitVarNameTerm(ctx *bindings.VarNameTermContext) interface{} {
	return v.WrapTerm(ctx.VarName())
}

func (v *XMLVisitor) VisitArrayTerm(ctx *bindings.ArrayTermContext) interface{} {
	v.Builder.WriteString("<term>\n")

	v.Visit(ctx.VarName())
	v.Visit(ctx.SquareBracketExpression())

	v.Builder.WriteString("</term>\n")

	return nil
}

func (v *XMLVisitor) VisitSubroutineCallTerm(ctx *bindings.SubroutineCallTermContext) interface{} {
	return v.WrapTerm(ctx.SubroutineCall())
}

func (v *XMLVisitor) VisitParExpressionTerm(ctx *bindings.ParExpressionTermContext) interface{} {
	return v.WrapTerm(ctx.ParExpression())
}

func (v *XMLVisitor) VisitUnaryOpTerm(ctx *bindings.UnaryOpTermContext) interface{} {
	v.Builder.WriteString("<term>\n")

	v.Visit(ctx.UnaryOp())
	v.Visit(ctx.Term())

	v.Builder.WriteString("</term>\n")

	return nil
}

func (v *XMLVisitor) VisitSimpleSubroutineCall(ctx *bindings.SimpleSubroutineCallContext) interface{} {
	v.Visit(ctx.SubroutineName())
	v.Visit(ctx.ParExpressionList())

	return nil
}

func (v *XMLVisitor) VisitNestedSubroutineCall(ctx *bindings.NestedSubroutineCallContext) interface{} {
	possibleClassName := ctx.ClassName()
	possibleVarName := ctx.VarName()

	if possibleClassName != nil {
		v.Visit(possibleClassName)
	} else {
		v.Visit(possibleVarName)
	}

	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.DOT().GetText())
	v.Builder.WriteString("</symbol>\n")

	v.Visit(ctx.SubroutineName())
	v.Visit(ctx.ParExpressionList())

	return nil
}

func (v *XMLVisitor) VisitExpressionList(ctx *bindings.ExpressionListContext) interface{} {
	v.Builder.WriteString("<expressionList>\n")

	for i, expression := range ctx.AllExpression() {
		v.Visit(expression)
		if i < len(ctx.AllExpression())-1 {
			v.Builder.WriteString("<symbol>")
			v.Builder.WriteString(ctx.COMMA(i).GetText())
			v.Builder.WriteString("</symbol>\n")
		}
	}

	v.Builder.WriteString("</expressionList>\n")

	return nil
}

func (v *XMLVisitor) VisitParExpressionList(ctx *bindings.ParExpressionListContext) interface{} {
	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.LP().GetText())
	v.Builder.WriteString("</symbol>\n")

	if ctx.ExpressionList() != nil {
		v.Visit(ctx.ExpressionList())
	} else {
		v.Builder.WriteString("<expressionList>\n</expressionList>\n")
	}

	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.RP().GetText())
	v.Builder.WriteString("</symbol>\n")

	return nil
}

func (v *XMLVisitor) VisitOp(ctx *bindings.OpContext) interface{} {
	v.Builder.WriteString("<symbol>")

	switch ctx.GetText() {
	case "<":
		v.Builder.WriteString("&lt;")
	case ">":
		v.Builder.WriteString("&gt;")
	case "&":
		v.Builder.WriteString("&amp;")
	case "\"":
		v.Builder.WriteString("&quot;")
	default:
		v.Builder.WriteString(ctx.GetText())
	}

	v.Builder.WriteString("</symbol>\n")

	return nil
}

func (v *XMLVisitor) VisitUnaryOp(ctx *bindings.UnaryOpContext) interface{} {
	v.Builder.WriteString("<symbol>")
	v.Builder.WriteString(ctx.GetText())
	v.Builder.WriteString("</symbol>\n")

	return nil
}

func (v *XMLVisitor) VisitKeywordConstant(ctx *bindings.KeywordConstantContext) interface{} {
	v.Builder.WriteString("<keyword>")
	v.Builder.WriteString(ctx.GetText())
	v.Builder.WriteString("</keyword>\n")

	return nil
}

func (v *XMLVisitor) VisitIntegerConstant(ctx *bindings.IntegerConstantContext) interface{} {
	v.Builder.WriteString("<integerConstant>")
	v.Builder.WriteString(ctx.GetText())
	v.Builder.WriteString("</integerConstant>\n")

	return nil
}

func (v *XMLVisitor) VisitStringConstant(ctx *bindings.StringConstantContext) interface{} {
	v.Builder.WriteString("<stringConstant>")

	text := ctx.GetText()
	text = text[1 : len(text)-1]

	v.Builder.WriteString(text)

	v.Builder.WriteString("</stringConstant>\n")

	return nil
}
