package codegen

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/zuma206/sb3c/language"
	"github.com/zuma206/sb3c/sb3"
)

func Generate(program *language.Program, fileSystem fs.FS) (*sb3.SB3, error) {
	sb3Project := sb3.NewSB3(fileSystem)
	for _, class := range program.Classes {
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
		return sb3Project.NewStage(), nil
	case SpriteClass:
		return sb3Project.NewSprite(class.Name.Src), nil
	default:
		err := fmt.Errorf("%q has an invalid super class %q %w", class.Name.Src, class.Super.Src, &class.Super.Pos)
		return nil, errors.Join(InvalidSuperError, err)
	}
}

func generateMembers(target *sb3.TargetHnd, class *language.Class) error {
	for _, member := range class.Members {
		if err := generateMember(target, member); err != nil {
			return err
		}
	}
	return nil
}
