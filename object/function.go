package object

import (
	"strings"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/object/call"
	"github.com/MarcinKonowalczyk/goruby/object/env"
	"github.com/MarcinKonowalczyk/goruby/object/hash"
	"github.com/MarcinKonowalczyk/goruby/object/ruby"
	"github.com/MarcinKonowalczyk/goruby/trace"
)

// ReturnValue represents a wrapper object for a return statement. It is no
// real Ruby object and only used within the interpreter evaluation
type ReturnValue struct {
	Value ruby.Object
}

func (rv *ReturnValue) Inspect() string   { return rv.Value.Inspect() }
func (rv *ReturnValue) Class() ruby.Class { return rv.Value.Class() }
func (rv *ReturnValue) HashKey() hash.Key { return rv.Value.HashKey() }

var (
	_ ruby.Object = &ReturnValue{}
)

// BreakValue represents a wrapper object for a break statement. It is no
// real Ruby object and only used within the interpreter evaluation
type BreakValue struct {
	Value ruby.Object
}

func (bv *BreakValue) Inspect() string   { return bv.Value.Inspect() }
func (bv *BreakValue) Class() ruby.Class { return bv.Value.Class() }
func (bv *BreakValue) HashKey() hash.Key { return bv.Value.HashKey() }

var (
	_ ruby.Object = &BreakValue{}
)

// FunctionParameters represents a list of function parameters.
type functionParameters []*FunctionParameter

func (f functionParameters) defaultParamCount() int {
	count := 0
	for _, p := range f {
		if p.Default != nil {
			count++
		}
	}
	return count
}

func (f functionParameters) separateDefaultParams() ([]*FunctionParameter, []*FunctionParameter) {
	mandatory, defaults := make([]*FunctionParameter, 0), make([]*FunctionParameter, 0)
	for _, p := range f {
		if p.Default != nil {
			defaults = append(defaults, p)
		} else {
			mandatory = append(mandatory, p)
		}
	}
	return mandatory, defaults
}

// FunctionParameter represents a parameter within a function
type FunctionParameter struct {
	Name    string
	Default ruby.Object
	Splat   bool
}

func (f *FunctionParameter) String() string {
	var out strings.Builder
	if f.Splat {
		out.WriteString("*")
	}
	out.WriteString(f.Name)
	if f.Default != nil {
		out.WriteString(" = ")
		out.WriteString(f.Default.Inspect())
	}
	return out.String()
}

// A Function represents a user defined function. It is no real Ruby object.
type Function struct {
	Name       string
	Parameters []*FunctionParameter
	Body       *ast.BlockStatement
	Env        env.Environment[ruby.Object]
}

// String returns the function literal
func (f *Function) String() string {
	var out strings.Builder
	out.WriteString("{")
	if len(f.Parameters) != 0 {
		args := []string{}
		for _, a := range f.Parameters {
			args = append(args, a.String())
		}
		out.WriteString("|")
		out.WriteString(strings.Join(args, ", "))
		out.WriteString("|")
	}
	out.WriteString("")
	out.WriteString(f.Body.String())
	out.WriteString("}")
	return out.String()
}

// Call implements the RubyMethod interface. It evaluates f.Body and returns its result
func (f *Function) Call(ctx call.Context[ruby.Object], args ...ruby.Object) (ruby.Object, error) {
	defer trace.TraceCtx(ctx, "Function.Call")()
	trace.MessageCtx(ctx, f.Name)
	trace.MessageCtx(ctx, f.String())
	// TODO: Handle tail splats
	if len(f.Parameters) == 1 && f.Parameters[0].Splat {
		// Only one splat parameter.
		args_arr := NewArray(args...)
		extendedEnv := env.NewEnclosedEnvironment(f.Env)
		extendedEnv.Set(f.Parameters[0].Name, args_arr)
		evaluated, err := ctx.Eval(f.Body, extendedEnv)
		if err != nil {
			return nil, err
		}
		return f.unwrapReturnValue(evaluated), nil

	} else {
		// normal evaluation
		defaultParams := functionParameters(f.Parameters).defaultParamCount()
		if len(args) < len(f.Parameters)-defaultParams || len(args) > len(f.Parameters) {
			return nil, NewWrongNumberOfArgumentsError(len(f.Parameters), len(args))
		}
		params, err := f.populateParameters(args)
		if err != nil {
			return nil, err
		}
		extendedEnv := env.NewEnclosedEnvironment(f.Env)
		for _, v := range params {
			extendedEnv.Set(v.name, v.value)
		}
		evaluated, err := ctx.Eval(f.Body, extendedEnv)
		if err != nil {
			return nil, err
		}
		return f.unwrapReturnValue(evaluated), nil
	}
}

type populatedParameter struct {
	name  string
	value ruby.Object
}

func (f *Function) populateParameters(args []ruby.Object) ([]populatedParameter, error) {
	if len(args) > len(f.Parameters) {
		return nil, NewWrongNumberOfArgumentsError(len(f.Parameters), len(args))
	}
	params := make([]populatedParameter, 0, len(f.Parameters))

	mandatory, defaults := functionParameters(f.Parameters).separateDefaultParams()

	if len(args) < len(mandatory)-len(defaults) || len(args) > len(f.Parameters) {
		return nil, NewWrongNumberOfArgumentsError(len(f.Parameters), len(args))
	}

	if len(args) == len(f.Parameters) {
		for paramIdx, param := range f.Parameters {
			params = append(params, populatedParameter{
				name:  param.Name,
				value: args[paramIdx],
			})
		}
		return params, nil
	}

	parameters := append(mandatory, defaults...)

	for paramIdx, param := range parameters {
		if paramIdx >= len(args) {
			// params[param.Name] = param.Default
			params = append(params, populatedParameter{
				name:  param.Name,
				value: param.Default,
			})
			continue
		}
		// params[param.Name] = args[paramIdx]
		params = append(params, populatedParameter{
			name:  param.Name,
			value: args[paramIdx],
		})
	}
	return params, nil
}

func (f *Function) unwrapReturnValue(obj ruby.Object) ruby.Object {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

var (
	_ ruby.Method = &Function{}
)
