package codegen

import (
	"errors"
	"fmt"
	"iter"
	"strconv"
	"strings"

	"github.com/zuma206/sb3c/language"
	"github.com/zuma206/sb3c/sb3"
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

var (
	UndefinedMethodErr = errors.New("undefined method")
	MissingArgumentErr = errors.New("missing argument")
)

func generateProcedure(target *sb3.TargetHnd, method *language.Member) error {
	procedure := target.NewProcedure(method.Name.Src)
	for call := range method.Value.Method.Calls.Iter() {
		mapping, ok := mappings[call.Path.Src]
		if !ok {
			return fmt.Errorf("%w: %q %w", UndefinedMethodErr, call.Path.Src, &call.Path.Pos)
		}
		block := &sb3.Block{Opcode: mapping.Opcode, Inputs: map[string]*sb3.Input{}}
		next, stop := iter.Pull(call.Args.Iter())
		defer stop()
		for _, input := range mapping.Inputs {
			arg, ok := next()
			if !ok {
				return fmt.Errorf("%w: %q %w", MissingArgumentErr, input, &call.Path.Pos)
			}
			block.Inputs[input] = sb3.LiteralInput(&sb3.Literal{Type: sb3.LiteralNumber, Value: arg.Src})
		}
		procedure.PushBlock(block)
	}
	return nil
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
