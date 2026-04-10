package codegen

import (
	"errors"
	"fmt"
	"iter"

	"github.com/zuma206/sb3c/language"
	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/sb3"
	"github.com/zuma206/sb3c/utils"
)

func generateMember(target *sb3.TargetHnd, member *language.Member) error {
	switch {
	case member.Value.Method != nil:
		return generateProcedure(target, member)
	case member.Value.Attribute != nil:
		return generateVariable(target, member)
	default:
		panic("malformed class member")
	}
}

var (
	UndefinedMethodErr          = errors.New("undefined method")
	UndefinedMethodDecoratorErr = errors.New("undefined method decorator")
)

func generateProcedure(target *sb3.TargetHnd, method *language.Member) error {
	procedure := target.NewProcedure(method.Name.Src)
	if err := generateProcedureDecorators(method, procedure); err != nil {
		return err
	}
	for call := range method.Value.Method.Calls.Iter() {
		block, err := generateBlock(call)
		if err != nil {
			return err
		}
		procedure.PushBlock(block)
	}
	return nil
}

func generateBlock(call *language.Call) (*sb3.Block, error) {
	mapping, ok := mappings[call.Path.Src]
	if !ok {
		return nil, fmt.Errorf("%w: %q %w", UndefinedMethodErr, call.Path.Src, &call.Path.Pos)
	}
	inputs, err := generateInputs(call.Args, mapping.Inputs)
	if err != nil {
		return nil, fmt.Errorf("%w %w", err, &call.Path.Pos)
	}
	block := &sb3.Block{Opcode: mapping.Opcode, Inputs: inputs}
	return block, nil
}

var NotEnoughArgumentsErr = errors.New("not enough arguments")

var literalTypes = map[*lexer.Type]sb3.LiteralType{
	language.NumberLiteral: sb3.LiteralNumber,
	language.StringLiteral: sb3.LiteralString,
}

func getLiteralType(token *lexer.Token) sb3.LiteralType {
	literalType, ok := literalTypes[token.Type]
	if !ok {
		panic("invalid literal type")
	}
	return literalType
}

func generateInputs(args *utils.List[*lexer.Token], keys []string) (map[string]*sb3.Input, error) {
	inputs := make(map[string]*sb3.Input, len(keys))
	next, stop := iter.Pull(args.Iter())
	defer stop()
	for i, key := range keys {
		arg, ok := next()
		if !ok {
			err := fmt.Errorf("expected %d got %d", len(keys), i)
			return nil, errors.Join(NotEnoughArgumentsErr, err)
		}
		inputs[key] = sb3.LiteralInput(&sb3.Literal{Type: getLiteralType(arg), Value: arg.Src})
	}
	return inputs, nil
}

func generateVariable(target *sb3.TargetHnd, attribute *language.Member) error {
	var initialValue any = ""
	if attribute.Value.Attribute.Initializer != nil {
		var err error
		initialValue, err = evaluateConstantExpression(attribute.Value.Attribute.Initializer)
		if err != nil {
			return err
		}
	}
	target.NewVariable(attribute.Name.Src, initialValue)
	return nil
}
