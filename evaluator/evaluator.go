package evaluator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/MarcinKonowalczyk/goruby/ast"
	"github.com/MarcinKonowalczyk/goruby/object"
	"github.com/pkg/errors"
)

type callContext struct {
	object.CallContext
}

func (c *callContext) Eval(node ast.Node, env object.Environment) (object.RubyObject, error) {
	return Eval(node, env)
}

type rubyObjects []object.RubyObject

func (r rubyObjects) Inspect() string {
	toS := make([]string, len(r))
	for i, e := range r {
		toS[i] = e.Inspect()
	}
	return strings.Join(toS, ", ")
}
func (r rubyObjects) Type() object.Type       { return "" }
func (r rubyObjects) Class() object.RubyClass { return nil }

func expandToArrayIfNeeded(obj object.RubyObject) object.RubyObject {
	arr, ok := obj.(rubyObjects)
	if !ok {
		return obj
	}
	return object.NewArray(arr...)
}

// Eval evaluates the given node and traverses recursive over its children
func Eval(node ast.Node, env object.Environment) (object.RubyObject, error) {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.ReturnStatement:
		val, err := Eval(node.ReturnValue, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval of return statement")
		}
		return &object.ReturnValue{Value: val}, nil
	case *ast.BreakStatement:
		var val object.RubyObject
		val, err := Eval(node.Condition, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval of break statement")
		}
		if node.Unless {
			if isTruthy(val) {
				val = object.FALSE
			} else {
				val = object.TRUE
			}
		}
		return &object.BreakValue{Value: val}, nil
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	// Literals
	case (*ast.IntegerLiteral):
		return object.NewInteger(node.Value), nil
	case (*ast.FloatLiteral):
		return object.NewFloat(node.Value), nil
	case (*ast.Boolean):
		return nativeBoolToBooleanObject(node.Value), nil
	case (*ast.Nil):
		return object.NIL, nil
	case (*ast.Self):
		self, _ := env.Get("self")
		return self, nil
	case (*ast.Keyword__FILE__):
		return &object.String{Value: node.Filename}, nil
	case (*ast.InstanceVariable):
		self, _ := env.Get("self")
		selfObj := self.(*object.Self)
		selfAsEnv, ok := selfObj.RubyObject.(object.Environment)
		if !ok {
			return nil, errors.WithStack(
				object.NewSyntaxError(
					fmt.Errorf("instance variable not allowed for %s", selfObj.Name),
				),
			)
		}

		val, ok := selfAsEnv.Get(node.String())
		if !ok {
			return object.NIL, nil
		}
		return val, nil
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.Global:
		val, ok := env.Get(node.Value)
		if !ok {
			return object.NIL, nil
		}
		return val, nil
	case *ast.StringLiteral:
		return &object.String{Value: unescapeStringLiteral(node)}, nil
	case *ast.SymbolLiteral:
		switch value := node.Value.(type) {
		case *ast.Identifier:
			return &object.Symbol{Value: value.Value}, nil
		case *ast.StringLiteral:
			str, err := Eval(value, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval symbol literal string")
			}
			if str, ok := str.(*object.String); ok {
				return &object.Symbol{Value: str.Value}, nil
			}
			panic(errors.WithStack(
				fmt.Errorf("error while parsing SymbolLiteral: expected *object.String, got %T", str),
			))
		default:
			return nil, errors.WithStack(
				object.NewSyntaxError(fmt.Errorf("malformed symbol AST: %T", value)),
			)
		}
	case *ast.FunctionLiteral:
		context, _ := env.Get("self")
		_, inClassOrModule := context.(*object.Self).RubyObject.(object.Environment)
		if node.Receiver != nil {
			rec, err := Eval(node.Receiver, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval function receiver")
			}
			context = rec
			_, recIsEnv := context.(object.Environment)
			if recIsEnv || inClassOrModule {
				inClassOrModule = true
				context = context.Class().(object.RubyClassObject)
			}
		}
		params := make([]*object.FunctionParameter, len(node.Parameters))
		for i, param := range node.Parameters {
			def, err := Eval(param.Default, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval function literal param")
			}
			params[i] = &object.FunctionParameter{Name: param.Name.Value, Default: def}
		}
		body := node.Body
		function := &object.Function{
			Parameters: params,
			Env:        env,
			Body:       body,
		}
		extended := object.AddMethod(context, node.Name.Value, function)
		if node.Receiver != nil && !inClassOrModule {
			envInfo, _ := object.EnvStat(env, context)
			envInfo.Env().Set(node.Receiver.Value, extended)
		}
		return &object.Symbol{Value: node.Name.Value}, nil

	case *ast.ProcedureLiteral:
		// context, _ := env.Get("self")
		// _, inClassOrModule := context.(*object.Self).RubyObject.(object.Environment)
		// if node.Receiver != nil {
		// 	rec, err := Eval(node.Receiver, env)
		// 	if err != nil {
		// 		return nil, errors.WithMessage(err, "eval function receiver")
		// 	}
		// 	context = rec
		// 	_, recIsEnv := context.(object.Environment)
		// 	if recIsEnv || inClassOrModule {
		// 		inClassOrModule = true
		// 		context = context.Class().(object.RubyClassObject)
		// 	}
		// }
		params := make([]*object.FunctionParameter, len(node.Parameters))
		for i, param := range node.Parameters {
			def, err := Eval(param.Default, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval function literal param")
			}
			params[i] = &object.FunctionParameter{Name: param.Name.Value, Default: def}
		}
		body := node.Body
		// function := &object.Function{
		// 	Parameters: params,
		// 	Env:        env,
		// 	Body:       body,
		// }
		// extended := object.AddMethod(context, node.Name.Value, function)
		// if node.Receiver != nil && !inClassOrModule {
		// 	envInfo, _ := object.EnvStat(env, context)
		// 	envInfo.Env().Set(node.Receiver.Value, extended)
		// }
		return &object.Proc{
			Parameters:             params,
			Body:                   body,
			Env:                    env,
			ArgumentCountMandatory: true,
		}, nil

	case *ast.BlockExpression:
		node_params := node.Parameters
		body := node.Body
		params := make([]*object.FunctionParameter, len(node_params))
		for i, param := range node_params {
			def, err := Eval(param.Default, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval function literal param")
			}
			params[i] = &object.FunctionParameter{Name: param.Name.Value, Default: def}
		}
		block := &object.Proc{
			Parameters: params,
			Body:       body,
			Env:        env,
		}
		return block, nil
	case *ast.ArrayLiteral:
		elements, err := evalArrayElements(node.Elements, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval array literal")
		}
		// TODO: If any of the elements is a splat, we need to flatten them
		return &object.Array{Elements: elements}, nil
	case *ast.HashLiteral:
		var hash object.Hash
		for k, v := range node.Map {
			key, err := Eval(k, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval hash key")
			}
			value, err := Eval(v, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval hash value")
			}
			hash.Set(key, value)
		}
		return &hash, nil
	case ast.ExpressionList:
		var objects []object.RubyObject
		for _, e := range node {
			obj, err := Eval(e, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval expression list")
			}
			objects = append(objects, obj)
		}
		return rubyObjects(objects), nil

	// Expressions
	case *ast.Assignment:
		right, err := Eval(node.Right, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval right hand Assignment side")
		}

		switch left := node.Left.(type) {
		case *ast.IndexExpression:
			indexLeft, err := Eval(left.Left, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval left hand Assignment side: eval left side of IndexExpression")
			}
			index, err := Eval(left.Index, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval left hand Assignment side: eval right side of IndexExpression")
			}
			return evalIndexExpressionAssignment(indexLeft, index, expandToArrayIfNeeded(right))
		case *ast.InstanceVariable:
			self, _ := env.Get("self")
			selfObj := self.(*object.Self)
			selfAsEnv, ok := selfObj.RubyObject.(object.Environment)
			if !ok {
				return nil, errors.Wrap(
					object.NewSyntaxError(fmt.Errorf("instance variable not allowed for %s", selfObj.Name)),
					"eval left hand Assignment side",
				)
			}

			right = expandToArrayIfNeeded(right)
			selfAsEnv.Set(left.String(), right)
			return right, nil
		case *ast.Identifier:
			right = expandToArrayIfNeeded(right)
			env.Set(left.Value, right)
			return right, nil
		case *ast.Global:
			right = expandToArrayIfNeeded(right)
			env.SetGlobal(left.Value, right)
			return right, nil
		case ast.ExpressionList:
			values := []object.RubyObject{right}
			if list, ok := right.(rubyObjects); ok {
				values = list
			}
			if len(left) > len(values) {
				// enlarge slice
				for len(values) <= len(left) {
					values = append(values, object.NIL)
				}
			}
			for i, exp := range left {
				if _, ok := exp.(*ast.InstanceVariable); ok {
					self, _ := env.Get("self")
					selfObj := self.(*object.Self)
					selfAsEnv, ok := selfObj.RubyObject.(object.Environment)
					if !ok {
						return nil, errors.Wrap(
							object.NewSyntaxError(fmt.Errorf("instance variable not allowed for %s", selfObj.Name)),
							"eval left hand Assignment side",
						)
					}

					selfAsEnv.Set(exp.String(), values[i])
					continue
				}
				if indexExp, ok := exp.(*ast.IndexExpression); ok {
					indexLeft, err := Eval(indexExp.Left, env)
					if err != nil {
						return nil, errors.WithMessage(err, "eval left hand Assignment side: eval left side of IndexExpression")
					}
					index, err := Eval(indexExp.Index, env)
					if err != nil {
						return nil, errors.WithMessage(err, "eval left hand Assignment side: eval right side of IndexExpression")
					}
					evalIndexExpressionAssignment(indexLeft, index, values[i])
					continue
				}
				env.Set(exp.String(), values[i])
			}
			return expandToArrayIfNeeded(right), nil
		default:
			return nil, errors.WithStack(
				object.NewSyntaxError(fmt.Errorf("assignment not supported to %T", node.Left)),
			)
		}
	case *ast.ModuleExpression:
		module, ok := env.Get(node.Name.Value)
		if !ok {
			module = object.NewModule(node.Name.Value, env)
		}
		moduleEnv := module.(object.Environment)
		moduleEnv.Set("self", &object.Self{RubyObject: module, Name: node.Name.Value})
		bodyReturn, err := Eval(node.Body, moduleEnv)
		if err != nil {
			return nil, errors.WithMessage(err, "eval Module body")
		}
		selfObject, _ := moduleEnv.Get("self")
		self := selfObject.(*object.Self)
		env.Set(node.Name.Value, self.RubyObject)
		return bodyReturn, nil
	case *ast.ClassExpression:
		superClassName := "Object"
		if node.SuperClass != nil {
			superClassName = node.SuperClass.Value
		}
		superClass, ok := env.Get(superClassName)
		if !ok {
			return nil, errors.Wrap(
				object.NewUninitializedConstantNameError(superClassName),
				"eval class superclass",
			)
		}
		class, ok := env.Get(node.Name.Value)
		if !ok {
			class = object.NewClass(node.Name.Value, superClass.(object.RubyClassObject), env)
		}
		classEnv := class.(object.Environment)
		classEnv.Set("self", &object.Self{RubyObject: class, Name: node.Name.Value})
		bodyReturn, err := Eval(node.Body, classEnv)
		if err != nil {
			return nil, errors.WithMessage(err, "eval class body")
		}
		selfObject, _ := classEnv.Get("self")
		self := selfObject.(*object.Self)
		env.Set(node.Name.Value, self.RubyObject)
		return bodyReturn, nil
	case *ast.ContextCallExpression:
		context, err := Eval(node.Context, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval method call receiver")
		}
		if context == nil {
			context, _ = env.Get("self")
		}
		args, err := evalExpressions(node.Arguments, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval method call arguments")
		}
		if node.Block != nil {
			block, err := Eval(node.Block, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval method call block")
			}
			args = append(args, block)
		}
		callContext := &callContext{object.NewCallContext(env, context)}
		return object.Send(callContext, node.Function.Value, args...)
	case *ast.YieldExpression:
		selfObject, _ := env.Get("self")
		self := selfObject.(*object.Self)
		if self.Block == nil {
			return nil, errors.WithStack(object.NewNoBlockGivenLocalJumpError())
		}
		args, err := evalExpressions(node.Arguments, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval yield arguments")
		}
		callContext := &callContext{object.NewCallContext(env, self)}
		return self.Block.Call(callContext, args...)
	case *ast.IndexExpression:
		left, err := Eval(node.Left, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval IndexExpression left side")
		}
		switch left := left.(type) {
		case *object.Proc:
			// special case for procs
			// NOTE: we pass unevaluated index to proc
			return evalProcIndexExpression(env, left, node.Index)
		default:
			index, err := Eval(node.Index, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval IndexExpression index")
			}
			return evalIndexExpression(left, index)
		}
	case *ast.PrefixExpression:
		right, err := Eval(node.Right, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval prefix right side")
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left, err := Eval(node.Left, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval operator left side")
		}

		if node.IsControlExpression() && !node.MustEvaluateRight() && isTruthy(left) {
			return left, nil
		}

		right, err := Eval(node.Right, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval operator right side")
		}
		if node.IsControlExpression() {
			return right, nil
		}
		context := &callContext{object.NewCallContext(env, left)}
		return object.Send(context, node.Operator, right)
	case *ast.ConditionalExpression:
		return evalConditionalExpression(node, env)
	case *ast.ScopedIdentifier:
		self, _ := env.Get("self")
		outer, ok := env.Get(node.Outer.Value)
		if !ok {
			return nil, errors.Wrap(
				object.NewUndefinedLocalVariableOrMethodNameError(self, node.Outer.Value),
				"eval scope outer",
			)
		}
		outerEnv, ok := outer.(object.Environment)
		if !ok {
			return nil, errors.Wrap(
				object.NewUndefinedLocalVariableOrMethodNameError(self, node.Outer.Value),
				"eval scope outer",
			)
		}
		inner, err := Eval(node.Inner, outerEnv)
		if err != nil {
			return nil, errors.WithMessage(err, "eval scope inner")
		}
		return inner, nil
	case *ast.ExceptionHandlingBlock:
		bodyReturn, err := Eval(node.TryBody, env)
		if err == nil {
			return bodyReturn, nil
		}
		return handleException(err, node.Rescues, env)

	case *ast.Comment:
		// ignore comments
		return nil, nil

	case nil:
		return nil, nil

	case *ast.RegexLiteral:
		return &object.Regex{Value: node.Value, Modifiers: node.Modifiers}, nil

	case *ast.RangeLiteral:
		left, err := Eval(node.Left, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval range start")
		}
		right, err := Eval(node.Right, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval range end")
		}
		if left == nil || right == nil {
			return nil, errors.WithStack(
				object.NewSyntaxError(fmt.Errorf("range start or end is nil")),
			)
		}
		if left.Type() != right.Type() {
			return nil, errors.WithStack(
				object.NewSyntaxError(fmt.Errorf("range start and end are not the same type: %s %s", left.Type(), right.Type())),
			)
		}
		if left.Type() != object.INTEGER_OBJ {
			return nil, errors.WithStack(
				object.NewSyntaxError(fmt.Errorf("range start and end are not integers: %s %s", left.Type(), right.Type())),
			)
		}
		leftInt, ok := left.(*object.Integer)
		if !ok {
			return nil, errors.WithStack(
				object.NewSyntaxError(fmt.Errorf("range start is not an integer: %s", left.Type())),
			)
		}
		rightInt, ok := right.(*object.Integer)
		if !ok {
			return nil, errors.WithStack(
				object.NewSyntaxError(fmt.Errorf("range end is not an integer: %s", right.Type())),
			)
		}

		left_int := leftInt.Value
		right_int := rightInt.Value

		if left_int > right_int {
			return nil, errors.WithStack(
				object.NewSyntaxError(fmt.Errorf("range start is greater than end: %d > %d", left_int, right_int)),
			)
		}

		if node.Inclusive {
			right_int++
		}

		elements := make([]object.RubyObject, right_int-left_int)
		for i := left_int; i < right_int; i++ {
			elements[i-left_int] = &object.Integer{Value: i}
		}

		return &object.Array{
			Elements: elements,
		}, nil

	case *ast.Splat:

		val, err := Eval(node.Value, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval splat value")
		}
		if val == nil {
			return nil, errors.WithStack(
				object.NewSyntaxError(fmt.Errorf("splat value is nil")),
			)
		}

		switch val := val.(type) {
		case *object.Array:
			return &object.Array{
				Elements: val.Elements,
			}, nil
		default:
			return &object.Array{
				Elements: []object.RubyObject{val},
			}, nil
		}

	case *ast.LoopExpression:

		return evalLoopExpression(node, env)

	default:
		err := object.NewException("Unknown AST: %T", node)
		return nil, errors.WithStack(err)
	}

}

func unescapeStringLiteral(node *ast.StringLiteral) string {
	rep := map[string]string{
		"\\n":  "\n",
		"\\t":  "\t",
		"\\r":  "\r",
		"\\b":  "\b",
		"\\\\": "\\",
	}
	if node.Token.Literal == "\"" {
		rep["\""] = "\""
	} else {
		rep["'"] = "'"
	}
	value := node.Value
	for k, v := range rep {
		value = strings.ReplaceAll(value, k, v)
	}
	return value
}

func evalLoopExpression(node *ast.LoopExpression, env object.Environment) (object.RubyObject, error) {
	condition, err := Eval(node.Condition, env)
	if err != nil {
		return nil, errors.WithMessage(err, "eval loop condition")
	}
	for isTruthy(condition) {
		value, err := evalBlockStatement(node.Block, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval loop body")
		}
		if value != nil {
			switch value := value.(type) {
			case *object.BreakValue:
				if isTruthy(value.Value) {
					return value.Value, nil
				}
			case *object.ReturnValue:
				return value.Value, nil
			default:
				// <shrug> do nothing
			}
		}
		condition, err = Eval(node.Condition, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval loop condition")
		}
	}
	return object.NIL, nil
}

func evalProgram(stmts []ast.Statement, env object.Environment) (object.RubyObject, error) {
	var result object.RubyObject
	var err error
	for _, statement := range stmts {
		if _, ok := statement.(*ast.Comment); ok {
			continue
		}
		result, err = Eval(statement, env)

		if err != nil {
			return nil, errors.WithMessage(err, "eval program statement")
		}

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value, nil
		}

	}
	return result, nil
}

func evalExpressions(exps []ast.Expression, env object.Environment) ([]object.RubyObject, error) {
	var result []object.RubyObject

	for _, e := range exps {
		evaluated, err := Eval(e, env)
		if err != nil {
			return nil, err
		}
		result = append(result, evaluated)
	}
	return result, nil
}

func evalArrayElements(elements []ast.Expression, env object.Environment) ([]object.RubyObject, error) {
	var result []object.RubyObject

	for _, e := range elements {
		splat, ok := e.(*ast.Splat)
		if ok {
			// we're a splat! eval the splat and append the elements
			evaluated, err := Eval(splat.Value, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval splat value")
			}
			if evaluated != nil {
				if evaluated.Type() != object.ARRAY_OBJ {
					return nil, errors.WithStack(
						object.NewException("splat value is not an array: %s", evaluated.Type()),
					)
				}
				arrObj, ok := evaluated.(*object.Array)
				if !ok {
					return nil, errors.WithStack(
						object.NewException("splat value is not an array: %s", evaluated.Type()),
					)
				}

				result = append(result, arrObj.Elements...)
			}
		} else {
			// not a splat. just eval the element
			evaluated, err := Eval(e, env)
			if err != nil {
				return nil, errors.WithMessage(err, "eval array elements")
			}
			result = append(result, evaluated)
		}
	}
	return result, nil
}

func evalPrefixExpression(operator string, right object.RubyObject) (object.RubyObject, error) {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right), nil
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return nil, errors.WithStack(object.NewException("unknown operator: %s%s", operator, right.Type()))
	}
}

func evalBangOperatorExpression(right object.RubyObject) object.RubyObject {
	switch right {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	case object.NIL:
		return object.TRUE
	default:
		return object.FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.RubyObject) (object.RubyObject, error) {
	switch right := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -right.Value}, nil
	default:
		return nil, errors.WithStack(object.NewException("unknown operator: -%s", right.Type()))
	}
}

func evalConditionalExpression(ce *ast.ConditionalExpression, env object.Environment) (object.RubyObject, error) {
	condition, err := Eval(ce.Condition, env)
	if err != nil {
		return nil, err
	}
	evaluateConsequence := isTruthy(condition)
	if ce.IsNegated() {
		evaluateConsequence = !evaluateConsequence
	}
	if evaluateConsequence {
		return Eval(ce.Consequence, env)
	} else if ce.Alternative != nil {
		return Eval(ce.Alternative, env)
	} else {
		return object.NIL, nil
	}
}

func evalIndexExpressionAssignment(left, index, right object.RubyObject) (object.RubyObject, error) {
	switch target := left.(type) {
	case *object.Array:
		integer, ok := index.(*object.Integer)
		if !ok {
			return nil, errors.Wrap(
				object.NewImplicitConversionTypeError(integer, index),
				"eval array index",
			)
		}
		idx := int(integer.Value)
		if idx >= len(target.Elements) {
			// enlarge slice
			for len(target.Elements) <= idx {
				target.Elements = append(target.Elements, object.NIL)
			}
		}
		target.Elements[idx] = right
		return right, nil
	case *object.Hash:
		target.Set(index, right)
		return right, nil
	default:
		return nil, errors.Wrap(
			object.NewException("assignment target not supported: %s", left.Type()),
			"eval IndexExpression Assignment",
		)
	}
}

func evalIndexExpression(left, index object.RubyObject) (object.RubyObject, error) {
	switch target := left.(type) {
	case *object.Array:
		return evalArrayIndexExpression(target, index), nil
	case *object.Hash:
		return evalHashIndexExpression(target, index), nil
	case *object.String:
		return evalStringIndexExpression(target, index), nil
	default:
		var left_type string = string(left.Type())
		if left_type == "" {
			left_type = fmt.Sprintf("%T", left)
		}
		return nil, errors.WithStack(object.NewException("index operator not supported: %s", left_type))
	}
}

func evalProcIndexExpression(env object.Environment, target *object.Proc, index ast.Expression) (object.RubyObject, error) {
	switch index.(type) {
	case *ast.Splat:
		// evaluate the splat literal
		literal := index.(*ast.Splat)
		evaluated, err := Eval(literal.Value, env)
		if err != nil {
			return nil, errors.WithMessage(err, "eval splat value")
		}
		if evaluated == nil {
			return nil, errors.WithStack(
				object.NewException("splat value is nil"),
			)
		}
		if evaluated.Type() != object.ARRAY_OBJ {
			return nil, errors.WithStack(
				object.NewException("splat value is not an array: %s", evaluated.Type()),
			)
		}
		arrObj, ok := evaluated.(*object.Array)
		if !ok {
			return nil, errors.WithStack(
				object.NewException("splat value is not an array: %s", evaluated.Type()),
			)
		}

		// call the proc with the splat arguments
		args := make([]object.RubyObject, len(arrObj.Elements))
		for i, e := range arrObj.Elements {
			if e == nil {
				args[i] = object.NIL
			} else {
				args[i] = e
			}
		}
		printable_args := make([]string, len(args))
		for i, e := range args {
			if e == nil {
				printable_args[i] = "nil"
			} else {
				printable_args[i] = e.Inspect()
			}
		}
		callContext := &callContext{object.NewCallContext(env, nil)}
		value, err := target.Call(callContext, args...)
		return value, err
	default:
		// not implemented yet
		return nil, errors.WithStack(
			object.NewException("proc index operator not supported: %s", index),
		)
	}
}
func evalArrayIndexExpression(arrayObject *object.Array, index object.RubyObject) object.RubyObject {
	idx := index.(*object.Integer).Value
	maxNegative := -int64(len(arrayObject.Elements))
	maxPositive := maxNegative*-1 - 1
	if maxPositive < 0 {
		return object.NIL
	}

	if idx > 0 && idx > maxPositive {
		return object.NIL
	}
	if idx < 0 && idx < maxNegative {
		return object.NIL
	}
	if idx < 0 {
		return arrayObject.Elements[len(arrayObject.Elements)+int(idx)]
	}
	return arrayObject.Elements[idx]
}

func evalHashIndexExpression(hash *object.Hash, index object.RubyObject) object.RubyObject {
	result, ok := hash.Get(index)
	if !ok {
		return object.NIL
	}
	return result
}

func evalStringIndexExpression(stringObject *object.String, index object.RubyObject) object.RubyObject {
	idx := index.(*object.Integer).Value
	maxNegative := -int64(len(stringObject.Value))
	maxPositive := maxNegative*-1 - 1
	if maxPositive < 0 {
		return object.NIL
	}

	if idx > 0 && idx > maxPositive {
		return object.NIL
	}
	if idx < 0 && idx < maxNegative {
		return object.NIL
	}
	if idx < 0 {
		return &object.String{Value: string(stringObject.Value[len(stringObject.Value)+int(idx)])}
	}
	return &object.String{Value: string(stringObject.Value[idx])}
}

func evalBlockStatement(block *ast.BlockStatement, env object.Environment) (object.RubyObject, error) {
	var result object.RubyObject
	var err error
	for _, statement := range block.Statements {
		result, err = Eval(statement, env)
		if err != nil {
			return nil, err
		}
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ {
				return result, nil
			} else if rt == object.BREAK_VALUE_OBJ {
				if isTruthy(result.(*object.BreakValue).Value) {
					return result, nil
				} else {
				}
			}
		}
	}
	if result == nil {
		return object.NIL, nil
	}
	return result, nil
}

func evalIdentifier(node *ast.Identifier, env object.Environment) (object.RubyObject, error) {
	val, ok := env.Get(node.Value)
	if ok {
		return val, nil
	}

	if node.IsConstant() {
		return nil, errors.Wrap(
			object.NewUninitializedConstantNameError(node.Value),
			"eval identifier",
		)
	}

	self, _ := env.Get("self")
	context := &callContext{object.NewCallContext(env, self)}
	val, err := object.Send(context, node.Value)
	if err != nil {
		return nil, errors.Wrap(
			object.NewUndefinedLocalVariableOrMethodNameError(self, node.Value),
			"eval ident as method call",
		)
	}
	return val, nil
}

func unwrapReturnValue(obj object.RubyObject) object.RubyObject {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func handleException(err error, rescues []*ast.RescueBlock, env object.Environment) (object.RubyObject, error) {
	if err != nil && len(rescues) == 0 {
		return nil, err
	}
	errorObject := err.(object.RubyObject)
	errClass := errorObject.Class().Name()
	rescueEnv := object.WithScopedLocalVariables(env)

	var catchAll *ast.RescueBlock
	for _, r := range rescues {
		if len(r.ExceptionClasses) == 0 {
			catchAll = r
			continue
		}
		if r.Exception != nil {
			rescueEnv.Set(r.Exception.Value, errorObject)
		}
		for _, cl := range r.ExceptionClasses {
			if cl.Value == errClass {
				rescueRet, err := Eval(r.Body, rescueEnv)
				return rescueRet, err
			}
		}
	}

	if catchAll != nil {
		ancestors := getAncestors(errorObject)
		sort.Strings(ancestors)
		if sort.SearchStrings(ancestors, "StandardError") >= len(ancestors) {
			return nil, err
		}

		if catchAll.Exception != nil {
			rescueEnv.Set(catchAll.Exception.Value, errorObject)
		}
		rescueRet, err := Eval(catchAll.Body, rescueEnv)
		return rescueRet, err
	}

	return nil, err
}

func getAncestors(obj object.RubyObject) []string {
	class := obj.Class()
	if c, ok := obj.(object.RubyClass); ok {
		class = c
	}
	var ancestors []string
	ancestors = append(ancestors, class.Name())

	superClass := class.SuperClass()
	if superClass != nil {
		superAncestors := getAncestors(superClass.(object.RubyClassObject))
		ancestors = append(ancestors, superAncestors...)
	}
	return ancestors
}

func isTruthy(obj object.RubyObject) bool {
	switch obj {
	case object.NIL:
		return false
	case object.TRUE:
		return true
	case object.FALSE:
		return false
	default:
		switch obj := obj.(type) {
		case *object.Boolean:
			return obj.Value
		case *object.Integer:
			return obj.Value != 0
		case *object.Float:
			return obj.Value != 0.0
		case *object.String:
			return obj.Value != ""
		case *object.Array:
			return len(obj.Elements) > 0
		case *object.Hash:
			return obj.Len() > 0
		default:
			return true
		}
	}
}

// IsError returns true if the given RubyObject is an object.Error or an
// object.Exception (or any subclass of object.Exception)
func IsError(obj object.RubyObject) bool {
	if obj != nil {
		return obj.Type() == object.EXCEPTION_OBJ
	}
	return false
}

func nativeBoolToBooleanObject(input bool) object.RubyObject {
	if input {
		return object.TRUE
	}
	return object.FALSE
}
