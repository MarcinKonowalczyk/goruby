goruby
======

[![Go Report Card](https://goreportcard.com/badge/github.com/MarcinKonowalczyk/goruby)](https://goreportcard.com/report/github.com/MarcinKonowalczyk/goruby)

An implementation of *A SUBSET* of Ruby in Go. The goal of this implementation is to be able to run [Pyramid-Scheme](https://github.com/ConorOBrien-Foxx/Pyramid-Scheme).


## REPL
There is a basic REPL within `cmd/girb`. It supports multiline expressions and all syntax elements the language supports yet.

To run it ad hoc run `go run cmd/girb/main.go` and exit the REPL with CTRL-D.

## Command
To run the command as one off run `go run cmd/goruby/main.go`.