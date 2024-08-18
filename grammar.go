package main

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"github.com/pkg/errors"
	"os"

	bindings "github.com/KyrylR/jack-compiler/parser"
)

type simpleErrorListener struct {
	*antlr.DefaultErrorListener
	errors []string
}

func (l *simpleErrorListener) SyntaxError(_ antlr.Recognizer, _ interface{},
	line, column int, msg string, _ antlr.RecognitionException) {

	errorMsg := fmt.Sprintf("line %d:%d %s", line, column, msg)

	l.errors = append(l.errors, errorMsg)
}

func (l *simpleErrorListener) hasErrors() bool {
	return len(l.errors) > 0
}

func (l *simpleErrorListener) getErrors() []string {
	return l.errors
}

func GetParser(input antlr.CharStream) *bindings.JackParser {
	lexer := bindings.NewJackLexer(input)

	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	return bindings.NewJackParser(stream)
}

// ParseFile parses a single file and returns an error if parsing fails
func ParseFile(filename string) error {
	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return errors.Wrap(err, "error opening file")
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()
	ioStream := antlr.NewIoStream(file)

	parser := GetParser(ioStream)
	parser.GetInterpreter().SetPredictionMode(antlr.PredictionModeSLL)

	parser.RemoveErrorListeners()
	parser.BuildParseTrees = true

	errorListener := &simpleErrorListener{}
	parser.AddErrorListener(errorListener)

	tree := parser.Program()

	if errorListener.hasErrors() {
		return fmt.Errorf("syntax errors encountered in file %s", filename)
	}

	simpleVisitor := NewXMLVisitor(parser)
	simpleVisitor.Visit(tree)

	// save res to file
	res := simpleVisitor.Builder.String()
	outputFileName := filename[:len(filename)-4] + "xml"
	err = os.WriteFile(outputFileName, []byte(res), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return errors.Wrap(err, "error writing file")
	}

	return nil
}
