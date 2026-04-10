package codegen

import (
	"github.com/zuma206/sb3c/language"
	"github.com/zuma206/sb3c/sb3"
)

type BlockMapping struct {
	Opcode string
	Inputs []string
}

var mappings = map[string]*BlockMapping{
	"this.motion.moveSteps":  {Opcode: "motion_movesteps", Inputs: []string{"STEPS"}},
	"this.looks.say":         {Opcode: "looks_say", Inputs: []string{"MESSAGE"}},
	"this.sound.setVolumeTo": {Opcode: "sound_setvolumeto", Inputs: []string{"VOLUME"}},
}

type DecoratorMapping[Handle any] func(member *language.Member, handle Handle, args []any) error
type DecoratorMappings[Handle any] map[string]DecoratorMapping[Handle]

var procedureDecoratorMappings = DecoratorMappings[*sb3.ProcedureHnd]{
	"whenGreenFlagClicked": procedureCallDecorator("event_whenflagclicked"),
}

var attributeDecoratorMappings = DecoratorMappings[*sb3.TargetHnd]{
	"costume": costumeAttributeDecorator,
}

func procedureCallDecorator(opcode string) DecoratorMapping[*sb3.ProcedureHnd] {
	return func(_ *language.Member, procedure *sb3.ProcedureHnd, _ []any) error {
		blockThread := procedure.Target().NewBlockThread()
		blockThread.PushBlock(&sb3.Block{
			Opcode:   opcode,
			Inputs:   map[string]*sb3.Input{},
			Fields:   struct{}{},
			TopLevel: true,
		})
		blockThread.PushBlock(procedure.Call())
		return nil
	}
}
