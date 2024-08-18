grammar Jack;

import LexerJack;

program
    :   classDeclaration* EOF
    ;


classDeclaration
    : 'class' className '{' classVarDec* subroutineDec* '}'
    ;

classVarDec
    : classVarType=('static' | 'field') type varName (',' varName)* ';'
    ;

type
    : 'int'
    | 'char'
    | 'boolean'
    | className
    ;

subroutineDec
    : subroutineDecType=('constructor' | 'function' | 'method') ('void' | type) subroutineName '(' parameterList? ')' subroutineBody
    ;

parameterList
    : type varName (',' type varName)*
    ;

subroutineBody
    : '{' varDec* statements '}'
    ;

varDec
    : 'var' type varName (',' varName)* ';'
    ;

className
    : ID
    ;

subroutineName
    : ID
    ;

varName
    : ID
    ;

statements
    : statement*
    ;

statement
    : letStatement
    | ifStatement
    | whileStatement
    | doStatement
    | returnStatement
    ;

letStatement
    : 'let' varName squareBracketExpression? '=' expression ';'
    ;

squareBracketExpression
    : '[' expression ']'
    ;

ifStatement
    : 'if' parExpression curlyStatements elseBlock?
    ;

curlyStatements
    : '{' statements '}'
    ;

parExpression
    : '(' expression ')'
    ;

elseBlock
    : 'else' curlyStatements
    ;

whileStatement
    : 'while' parExpression curlyStatements
    ;

doStatement
    : 'do' subroutineCall ';'
    ;

returnStatement
    : 'return' expression? ';'
    ;

expression
    : term (op term)*
    ;

term
    : integerConstant                       #IntegerConstantTerm
    | stringConstant                        #StringConstantTerm
    | keywordConstant                       #KeywordConstantTerm
    | varName                               #VarNameTerm
    | varName squareBracketExpression       #ArrayTerm
    | subroutineCall                        #SubroutineCallTerm
    | parExpression                         #ParExpressionTerm
    | unaryOp term                          #UnaryOpTerm
    ;

subroutineCall
    : subroutineName parExpressionList                            #SimpleSubroutineCall
    | (className | varName) '.' subroutineName parExpressionList  #NestedSubroutineCall
    ;

expressionList
    : expression (',' expression)*
    ;

parExpressionList
    : '(' expressionList? ')'
    ;

op
    : '+' | '-' | '*' | '/' | '&' | '|' | '<' | '>' | '='
    ;

unaryOp
    : '-' | '~'
    ;

keywordConstant
    : 'true' | 'false' | 'null' | 'this'
    ;

integerConstant
    : NUMBER
    ;

stringConstant
    : STRING
    ;