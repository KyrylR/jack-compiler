lexer grammar LexerJack;

CLASS: 'class' ;
CONSTRUCTOR: 'constructor' ;
FUNCTION: 'function' ;
METHOD: 'method' ;
FIELD: 'field' ;
STATIC: 'static' ;
VAR: 'var' ;
INT: 'int' ;
CHAR: 'char' ;
BOOLEAN: 'boolean' ;
VOID: 'void' ;
TRUE: 'true' ;
FALSE: 'false' ;
NULL: 'null' ;
THIS: 'this' ;
LET: 'let' ;
DO: 'do' ;
IF: 'if' ;
ELSE: 'else' ;
WHILE: 'while' ;
RETURN: 'return' ;

// Symbols
LP: '(' ;
RP: ')' ;
LB: '[' ;
RB: ']' ;
LC: '{' ;
RC: '}' ;
DOT: '.' ;
COMMA: ',' ;
SEMICOLON: ';' ;
PLUS: '+' ;
MINUS: '-' ;
MULTIPLY: '*' ;
DIVIDE: '/' ;
AND: '&' ;
OR: '|' ;
LT: '<' ;
GT: '>' ;
EQ: '=' ;
BNOT: '~' ;

ID          :   LETTER (LETTER|DIGIT)*;
fragment
LETTER      :   [a-zA-Z\u0080-\u00FF_] ;

NUMBER: DIGIT+ ;                                      // match integers
fragment
DIGIT: [0-9] ;                                        // match single digit

STRING      :   '"' (ESC|.)*? '"' ;
fragment ESC: '\\' [btnrf"\\] ;

COMMENT
    : '/*' .*? '*/'    -> channel(HIDDEN)             // match anything between /* and */
    ;

LINE_COMMENT
    : '//' ~[\r\n]* '\r'? '\n' -> channel(HIDDEN)
    ;

WS  : [ \r\t\u000C\n]+ -> channel(HIDDEN)
    ;