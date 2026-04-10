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

var InvalidDecoratorErr = errors.New("invalid decorator")

func generateDecorators[Handle any](member *language.Member, hnd Handle, mappings DecoratorMappings[Handle]) error {
	for decorator := range member.Decorators.Iter() {
		mapping, ok := mappings[decorator.Path.Src]
		if !ok {
			err := fmt.Errorf("%q %w", decorator.Path.Src, &decorator.Path.Pos)
			return errors.Join(InvalidDecoratorErr, err)
		}
		args := []any{}
		for arg := range decorator.Args.Iter() {
			value, err := evaluateConstantExpression(arg)
			if err != nil {
				return err
			}
			args = append(args, value)
		}
		if err := mapping(member, hnd, args); err != nil {
			return fmt.Errorf("%w %w", err, &member.Name.Pos)
		}
	}
	return nil
}

var (
	MissingDecoratorArgsErr = errors.New("missing decorator arguments")
	IncorrectArgTypeErr     = errors.New("incorrect argument type in decorator")
)

const costumeArgs = 1

func costumeAttributeDecorator(attribute *language.Member, target *sb3.TargetHnd, args []any) error {
	if len(args) < costumeArgs {
		err := fmt.Errorf("expected %d got %d", costumeArgs, len(args))
		return errors.Join(MissingDecoratorArgsErr, err)
	}
	path, ok := args[0].(string)
	if !ok {
		err := fmt.Errorf("expected string path %w", &attribute.Name.Pos)
		return errors.Join(IncorrectArgTypeErr, err)
	}
	_, err := target.NewCostume(attribute.Name.Src, path)
	return err
}
