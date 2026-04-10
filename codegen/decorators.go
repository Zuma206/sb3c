package codegen

import (
	"errors"
	"fmt"

	"github.com/zuma206/sb3c/language"
	"github.com/zuma206/sb3c/sb3"
)

func generateProcedureDecorators(method *language.Member, procedure *sb3.ProcedureHnd) error {
	return generateDecorators(method, procedure, procedureDecoratorMappings)
}
func generateAttributeDecorators(method *language.Member, target *sb3.TargetHnd) error {
	return generateDecorators(method, target, attributeDecoratorMappings)
}

func generateDecorators[Handle any](member *language.Member, hnd Handle, mappings DecoratorMappings[Handle]) error {
	for decorator := range member.Decorators.Iter() {
		mapping, ok := mappings[decorator.Path.Src]
		if !ok {
			err := fmt.Errorf("%q %w", decorator.Path.Src, &decorator.Path.Pos)
			return errors.Join(UndefinedMethodDecoratorErr, err)
		}
		if err := mapping(member, hnd); err != nil {
			return fmt.Errorf("%w %w", err, &member.Name.Pos)
		}
	}
	return nil
}
