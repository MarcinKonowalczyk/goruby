package parser

import (
	"fmt"
	"testing"

	"github.com/MarcinKonowalczyk/assert"
	"github.com/MarcinKonowalczyk/goruby/token"
)

// isEOFInsteadOfNewlineError returns true if err represents an unexpectedTokenError with its
// actual token type set to token.EOF and if its expected token types includes
// token.NEWLINE.
//
// It returns false for any other error.
func isEOFInsteadOfNewlineError(err error) bool {
	if !IsEOFError(err) {
		return false
	}

	if errors, ok := err.(*Errors); ok {
		for _, e := range errors.Errors {
			if isEOFInsteadOfNewlineError(e) {
				return true
			}
		}
	}

	if tokenErr, ok := err.(*UnexpectedTokenError); ok {
		for _, expectedToken := range tokenErr.ExpectedTokens {
			if expectedToken == token.NEWLINE {
				return true
			}
		}
	}

	return false
}

func TestIsEOF(t *testing.T) {
	t.Run("Errors with unexpected token EOF", func(t *testing.T) {
		err := &UnexpectedTokenError{
			ExpectedTokens: []token.Type{token.IDENT},
			ActualToken:    token.EOF,
		}
		assert.That(t, IsEOFError(NewErrors("", err)), "Expected an EOF error, got %T:%q\n", err, err)
	})
	t.Run("Errors with unexpected token not EOF", func(t *testing.T) {
		err := &UnexpectedTokenError{
			ExpectedTokens: []token.Type{token.IDENT},
			ActualToken:    token.NEWLINE,
		}
		errs := NewErrors("", err)
		assert.That(t, !IsEOFError(errs), "Expected no EOF error, got %T:%q\n", err, err)
	})
	t.Run("unexpected token EOF", func(t *testing.T) {
		err := &UnexpectedTokenError{
			ExpectedTokens: []token.Type{token.IDENT},
			ActualToken:    token.EOF,
		}
		assert.That(t, IsEOFError(err), "Expected an EOF error, got %T:%q\n", err, err)
	})
	t.Run("unexpected token not EOF", func(t *testing.T) {
		err := &UnexpectedTokenError{
			ExpectedTokens: []token.Type{token.IDENT},
			ActualToken:    token.NEWLINE,
		}
		assert.That(t, !IsEOFError(err), "Expected no EOF error, got %T:%q\n", err, err)
	})
	t.Run("random error", func(t *testing.T) {
		err := fmt.Errorf("some error")
		assert.That(t, !IsEOFError(err), "Expected no EOF error, got %T:%q\n", err, err)
	})
}

func TestIsEOFInsteadOfNewline(t *testing.T) {
	t.Run("Errors with unexpected token EOF, expected token NEWLINE", func(t *testing.T) {
		err := &UnexpectedTokenError{
			ExpectedTokens: []token.Type{token.NEWLINE},
			ActualToken:    token.EOF,
		}
		assert.That(t, isEOFInsteadOfNewlineError(NewErrors("", err)), "Expected an EOF NEWLINE error, got %T:%q\n", err, err)
	})
	t.Run("unexpected token EOF, expected token NEWLINE", func(t *testing.T) {
		err := &UnexpectedTokenError{
			ExpectedTokens: []token.Type{token.NEWLINE},
			ActualToken:    token.EOF,
		}
		assert.That(t, isEOFInsteadOfNewlineError(err), "Expected an EOF NEWLINE error, got %T:%q\n", err, err)
	})
	t.Run("Errors with unexpected token not EOF", func(t *testing.T) {
		err := &UnexpectedTokenError{
			ExpectedTokens: []token.Type{token.IDENT},
			ActualToken:    token.NEWLINE,
		}
		assert.That(t, !isEOFInsteadOfNewlineError(NewErrors("", err)), "Expected no EOF NEWLINE error, got %T:%q\n", err, err)
	})
	t.Run("unexpected token not EOF", func(t *testing.T) {
		err := &UnexpectedTokenError{
			ExpectedTokens: []token.Type{token.IDENT},
			ActualToken:    token.NEWLINE,
		}
		assert.That(t, !isEOFInsteadOfNewlineError(err), "Expected no EOF NEWLINE error, got %T:%q\n", err, err)
	})
	t.Run("random error", func(t *testing.T) {
		err := fmt.Errorf("some error")
		assert.That(t, !isEOFInsteadOfNewlineError(err), "Expected no EOF NEWLINE error, got %T:%q\n", err, err)
	})
}
