package codegen

import (
	"errors"
	"fmt"
	"iter"
	"strconv"
	"strings"

	"github.com/zuma206/sb3c/language"
	"github.com/zuma206/sb3c/lexer"
	"github.com/zuma206/sb3c/sb3"
	"github.com/zuma206/sb3c/utils"
)

func Generate(program *language.Program) (*sb3.SB3, error) {
	sb3Project := sb3.NewSB3()
	for class := range program.Classes.Iter() {
		target, err := newTarget(sb3Project, class)
		if err != nil {
			return nil, err
		}
		if err := generateMembers(target, class); err != nil {
			return nil, err
		}
	}
	return sb3Project, nil
}

const (
	StageClass  = "Stage"
	SpriteClass = "Sprite"
)

var InvalidSuperError = errors.New("invalid super")

func newTarget(sb3Project *sb3.SB3, class *language.Class) (*sb3.TargetHnd, error) {
	switch class.Super.Src {
	case StageClass:
		return sb3Project.NewStage()
	default:
		err := fmt.Errorf("%q is an invalid super class %w", class.Name.Src, &class.Name.Pos)
		return nil, errors.Join(InvalidSuperError, err)
	}
}

func generateMembers(target *sb3.TargetHnd, class *language.Class) error {
	for member := range class.Members.Iter() {
		if err := generateMember(target, member); err != nil {
			return err
		}
	}
	return nil
}

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

var UndefinedMethodErr = errors.New("undefined method")

func generateProcedure(target *sb3.TargetHnd, method *language.Member) error {
	procedure := target.NewProcedure(method.Name.Src)
	for call := range method.Value.Method.Calls.Iter() {
		block, err := generateBlock(call)
		if err != nil {
			return err
		}
		procedure.PushBlock(block)
	}
	return nil
}

var MissingArgumentErr = errors.New("missing argument")

func generateBlock(call *language.Call) (*sb3.Block, error) {
	mapping, ok := mappings[call.Path.Src]
	if !ok {
		return nil, fmt.Errorf("%w: %q %w", UndefinedMethodErr, call.Path.Src, &call.Path.Pos)
	}
	inputs := generateInputs(call.Args, mapping.Inputs)
	if len(inputs) < len(mapping.Inputs) {
		err := fmt.Errorf("expected %d arguments, got %d", len(mapping.Inputs), len(inputs))
		return nil, errors.Join(MissingArgumentErr, err)
	}
	block := &sb3.Block{Opcode: mapping.Opcode, Inputs: map[string]*sb3.Input{}}
	return block, nil
}

func generateInputs(args *utils.List[*lexer.Token], keys []string) map[string]*sb3.Input {
	inputs := make(map[string]*sb3.Input, len(keys))
	next, stop := iter.Pull(args.Iter())
	defer stop()
	for _, key := range keys {
		arg, ok := next()
		if !ok {
			break
		}
		inputs[key] = sb3.LiteralInput(&sb3.Literal{Type: sb3.LiteralNumber, Value: arg.Src})
	}
	return inputs
}

func generateVariable(target *sb3.TargetHnd, attribute *language.Member) error {
	initialValue, err := generateVariableInitialValue(attribute.Value.Attribute)
	if err != nil {
		return err
	}
	target.NewVariable(attribute.Name.Src, initialValue)
	return nil
}

var InvalidInitalValueErr = errors.New("invalid attribute initial value")

func generateVariableInitialValue(attribute *language.Attribute) (any, error) {
	if attribute.Initializer == nil {
		return "", nil
	}
	switch attribute.Initializer.Type {
	case language.NumberLiteral:
		return strconv.Atoi(strings.ReplaceAll(attribute.Initializer.Src, "_", ""))
	default:
		return nil, errors.Join(InvalidInitalValueErr, &attribute.Initializer.Pos)
	}
}
