package codegen

import (
	"errors"
	"fmt"

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
	for declaration := range program.Declarations.Iter() {
		super := string(declaration.Super.Src)
		switch super {
		case StageClass:
			if err := cg.generateStage(declaration); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%w: %q", InvalidSuperError, super)
		}
	}
	return nil
}

func (cg *CodeGenerator) generateStage(class *language.ClassDeclaration) error {
	_, err := cg.sb3.NewStage(string(class.Name.Src))
	if err != nil {
		return fmt.Errorf("%w %w", err, &class.Super.Pos)
	}
	return nil
}
