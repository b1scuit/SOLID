package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/b1scuit/solid/rdf/lexer/lexertoken"
)

type LexerOption func(*Lexer)

func WithInput(i string) LexerOption {
	return func(l *Lexer) {
		l.Input = i
	}
}

func WihInitalState(lf LexFn) LexerOption {
	return func(l *Lexer) {
		l.State = lf
	}
}

type Lexer struct {
	Input  string
	Tokens chan lexertoken.Token
	State  LexFn

	Start int
	Pos   int
	Width int
}

func New(opts ...LexerOption) (*Lexer, error) {
	l := &Lexer{
		Tokens: make(chan lexertoken.Token),
	}

	for _, f := range opts {
		f(l)
	}

	return l, nil
}

/*
Backup to the beginning of the last read token.
*/
func (l *Lexer) Backup() {
	l.Pos -= l.Width
}

/*
Returns a slice of the current input from the current lexer start position
to the current position.
*/
func (l *Lexer) CurrentInput() string {
	return l.Input[l.Start:l.Pos]
}

/*
Decrement the position
*/
func (l *Lexer) Dec() {
	l.Pos--
}

/*
Puts a token onto the token channel. The value of l token is
read from the input based on the current lexer position.
*/
func (l *Lexer) Emit(tokenType lexertoken.TokenType) {
	l.Tokens <- lexertoken.Token{Type: tokenType, Value: strings.TrimSpace(l.Input[l.Start:l.Pos])}
	l.Start = l.Pos
}

/*
Returns a token with error information.
*/
func (l *Lexer) Errorf(format string, args ...interface{}) LexFn {
	l.Tokens <- lexertoken.Token{
		Type:  lexertoken.TOKEN_ERROR,
		Value: fmt.Sprintf(format, args...),
	}

	return nil
}

/*
Ignores the current token by setting the lexer's start
position to the current reading position.
*/
func (l *Lexer) Ignore() {
	l.Start = l.Pos
}

/*
Increment the position
*/
func (l *Lexer) Inc() {
	l.Pos++
	if l.Pos >= utf8.RuneCountInString(l.Input) {
		l.Emit(lexertoken.TOKEN_EOF)
	}
}

/*
Return a slice of the input from the current lexer position
to the end of the input string.
*/
func (l *Lexer) InputToEnd() string {
	return l.Input[l.Pos:]
}

/*
Returns the true/false if the lexer is at the end of the
input stream.
*/
func (l *Lexer) IsEOF() bool {
	return l.Pos >= len(l.Input)
}

/*
Returns true/false if then next character is whitespace
*/
func (l *Lexer) IsWhitespace() bool {
	ch, _ := utf8.DecodeRuneInString(l.Input[l.Pos:])
	return unicode.IsSpace(ch)
}

/*
Reads the next rune (character) from the input stream
and advances the lexer position.
*/
func (l *Lexer) Next() rune {
	if l.Pos >= utf8.RuneCountInString(l.Input) {
		l.Width = 0
		return lexertoken.EOF
	}

	result, width := utf8.DecodeRuneInString(l.Input[l.Pos:])

	l.Width = width
	l.Pos += l.Width
	return result
}

/*
Return the next token from the channel
*/

func (l *Lexer) NextToken() chan lexertoken.Token {
	return l.Tokens
}

/*
func (l *Lexer) NextToken() (t lexertoken.Token) {

	defer func() {
		if r := recover(); r != nil {
			t = lexertoken.Token{
				Type: lexertoken.TOKEN_EOF,
			}
		}
	}()

	for {
		select {
		case token := <-l.Tokens:
			return token
		default:
			l.State = l.State(l)
		}
	}
}*/

/*
Returns the next rune in the stream, then puts the lexer
position back. Basically reads the next rune without consuming
it.
*/
func (l *Lexer) Peek() rune {
	rune := l.Next()
	l.Backup()
	return rune
}

/*
Starts the lexical analysis and feeding tokens into the
token channel.
*/
func (l *Lexer) Run() {
	state := l.State
	for state := state; state != nil; {
		state = state(l)
	}

	l.Shutdown()
}

/*
Shuts down the token stream
*/
func (l *Lexer) Shutdown() {
	close(l.Tokens)
}

/*
Skips whitespace until we get something meaningful.
*/
func (l *Lexer) SkipWhitespace() {
	for {
		ch := l.Next()

		if !unicode.IsSpace(ch) {
			l.Dec()
			break
		}

		if ch == lexertoken.EOF {
			l.Emit(lexertoken.TOKEN_EOF)
			break
		}
	}
}
