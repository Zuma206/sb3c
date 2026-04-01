package codegen

import (
	"errors"
	"fmt"
	"iter"

	"github.com/zuma206/sb3c/language"
	"github.com/zuma206/sb3c/sb3"
)

type CodeGenerator struct {
	sb3 *sb3.SB3
}

func Generate(program *language.Program) (*sb3.SB3, error) {
	sb3Project := sb3.NewSB3()
	generator := &CodeGenerator{sb3: sb3Project}
	err := generator.generate(program)
	return sb3Project, err
}

var InvalidSuperError = errors.New("invalid super")

const (
	StageClass  = "Stage"
	SpriteClass = "Sprite"
)

func (cg *CodeGenerator) generate(program *language.Program) error {
	for declaration := range program.Classes.Iter() {
		switch declaration.Super.Src {
		case StageClass:
			if err := cg.generateStage(declaration); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%w: %q", InvalidSuperError, declaration.Super.Src)
		}
	}
	return nil
}

func (cg *CodeGenerator) generateStage(class *language.Class) error {
	stage, err := cg.sb3.NewStage()
	if err != nil {
		return fmt.Errorf("%w %w", err, &class.Super.Pos)
	}
	for method := range class.Members.Iter() {
		if err := cg.generateProcedure(stage, method); err != nil {
			return err
		}
	}
	return nil
}

var (
	UndefinedMethodErr = errors.New("undefined method")
	MissingArgumentErr = errors.New("missing argument")
)

func (cg *CodeGenerator) generateProcedure(target *sb3.TargetHnd, method *language.Method) error {
	procedure := target.NewProcedure(method.Name.Src)
	for call := range method.Calls.Iter() {
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
