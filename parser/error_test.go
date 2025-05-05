package parser

import (
	"testing"

	"github.com/MarcinKonowalczyk/goruby/token"
	"github.com/MarcinKonowalczyk/goruby/utils"
	"github.com/pkg/errors"
)

func TestIsEOF(t *testing.T) {
	t.Run("Errors with unexpected token EOF", func(t *testing.T) {
		err := &UnexpectedTokenError{
			ExpectedTokens: []token.Type{token.IDENT},
			ActualToken:    token.EOF,
		}
		utils.Assert(t, IsEOFError(NewErrors("", err)), "Expected an EOF error, got %T:%q\n", err, err)
	})
	t.Run("Errors with unexpected token not EOF", func(t *testing.T) {
		err := &UnexpectedTokenError{
			ExpectedTokens: []token.Type{token.IDENT},
			ActualToken:    token.NEWLINE,
		}
		utils.Assert(t, !IsEOFError(NewErrors("", err)), "Expected no EOF error, got %T:%q\n", err, err)
	})
	t.Run("unexpected token EOF", func(t *testing.T) {
		err := &UnexpectedTokenError{
			ExpectedTokens: []token.Type{token.IDENT},
			ActualToken:    token.EOF,
		}
		utils.Assert(t, IsEOFError(err), "Expected an EOF error, got %T:%q\n", err, err)
	})
	t.Run("unexpected token not EOF", func(t *testing.T) {
		err := &UnexpectedTokenError{
			ExpectedTokens: []token.Type{token.IDENT},
			ActualToken:    token.NEWLINE,
		}
		utils.Assert(t, !IsEOFError(err), "Expected no EOF error, got %T:%q\n", err, err)
	})
	t.Run("unexpected token EOF wrapped error", func(t *testing.T) {
		err := &UnexpectedTokenError{
			ExpectedTokens: []token.Type{token.IDENT},
			ActualToken:    token.EOF,
		}
		wrapped := errors.Wrap(err, "some error")
		utils.Assert(t, IsEOFError(wrapped), "Expected an EOF error, got %T:%q\n", wrapped, wrapped)
	})
	t.Run("unexpected token not EOF wrapped error", func(t *testing.T) {
		err := &UnexpectedTokenError{
			ExpectedTokens: []token.Type{token.IDENT},
			ActualToken:    token.NEWLINE,
		}
		wrapped := errors.Wrap(err, "some error")
		utils.Assert(t, !IsEOFError(wrapped), "Expected no EOF error, got %T:%q\n", wrapped, wrapped)
	})
	t.Run("random error", func(t *testing.T) {
		err := errors.New("some error")
		utils.Assert(t, !IsEOFError(err), "Expected no EOF error, got %T:%q\n", err, err)
	})
}

func TestIsEOFInsteadOfNewline(t *testing.T) {
	t.Run("Errors with unexpected token EOF, expected token NEWLINE", func(t *testing.T) {
		err := &UnexpectedTokenError{
			ExpectedTokens: []token.Type{token.NEWLINE},
			ActualToken:    token.EOF,
		}
		utils.Assert(t, IsEOFInsteadOfNewlineError(NewErrors("", err)), "Expected an EOF NEWLINE error, got %T:%q\n", err, err)
	})
	t.Run("unexpected token EOF, expected token NEWLINE", func(t *testing.T) {
		err := &UnexpectedTokenError{
			ExpectedTokens: []token.Type{token.NEWLINE},
			ActualToken:    token.EOF,
		}
		utils.Assert(t, IsEOFInsteadOfNewlineError(err), "Expected an EOF NEWLINE error, got %T:%q\n", err, err)
	})
	t.Run("Errors with unexpected token not EOF", func(t *testing.T) {
		err := &UnexpectedTokenError{
			ExpectedTokens: []token.Type{token.IDENT},
			ActualToken:    token.NEWLINE,
		}
		utils.Assert(t, !IsEOFInsteadOfNewlineError(NewErrors("", err)), "Expected no EOF NEWLINE error, got %T:%q\n", err, err)
	})
	t.Run("unexpected token not EOF", func(t *testing.T) {
		err := &UnexpectedTokenError{
			ExpectedTokens: []token.Type{token.IDENT},
			ActualToken:    token.NEWLINE,
		}
		utils.Assert(t, !IsEOFInsteadOfNewlineError(err), "Expected no EOF NEWLINE error, got %T:%q\n", err, err)
	})
	t.Run("unexpected token EOF expected NEWLINE wrapped error", func(t *testing.T) {
		err := &UnexpectedTokenError{
			ExpectedTokens: []token.Type{token.NEWLINE},
			ActualToken:    token.EOF,
		}
		wrapped := errors.Wrap(err, "some error")
		utils.Assert(t, IsEOFInsteadOfNewlineError(wrapped), "Expected an EOF NEWLINE error, got %T:%q\n", wrapped, wrapped)
	})
	t.Run("unexpected token EOF expected not NEWLINE wrapped error", func(t *testing.T) {
		err := &UnexpectedTokenError{
			ExpectedTokens: []token.Type{token.IDENT},
			ActualToken:    token.EOF,
		}
		wrapped := errors.Wrap(err, "some error")
		utils.Assert(t, !IsEOFInsteadOfNewlineError(wrapped), "Expected no EOF NEWLINE error, got %T:%q\n", wrapped, wrapped)
	})
	t.Run("unexpected token not EOF wrapped error", func(t *testing.T) {
		err := &UnexpectedTokenError{
			ExpectedTokens: []token.Type{token.IDENT},
			ActualToken:    token.NEWLINE,
		}
		wrapped := errors.Wrap(err, "some error")
		utils.Assert(t, !IsEOFInsteadOfNewlineError(wrapped), "Expected no EOF NEWLINE error, got %T:%q\n", wrapped, wrapped)
	})
	t.Run("random error", func(t *testing.T) {
		err := errors.New("some error")
		utils.Assert(t, !IsEOFInsteadOfNewlineError(err), "Expected no EOF NEWLINE error, got %T:%q\n", err, err)
	})
}
