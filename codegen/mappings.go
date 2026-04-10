package codegen

import "github.com/zuma206/sb3c/sb3"

type BlockMapping struct {
	Opcode string
	Inputs []string
}

var mappings = map[string]*BlockMapping{
	"this.motion.moveSteps":  {Opcode: "motion_movesteps", Inputs: []string{"STEPS"}},
	"this.looks.say":         {Opcode: "looks_say", Inputs: []string{"MESSAGE"}},
	"this.sound.setVolumeTo": {Opcode: "sound_setvolumeto", Inputs: []string{"VOLUME"}},
}

type ProcedureDecoratorMapping func(procedure *sb3.ProcedureHnd)

var procedureDecoratorMappings = map[string]ProcedureDecoratorMapping{
	"whenGreenFlagClicked": procedureCallDecorator("event_whenflagclicked"),
}

func procedureCallDecorator(opcode string) ProcedureDecoratorMapping {
	return func(procedure *sb3.ProcedureHnd) {
		blockThread := procedure.Target().NewBlockThread()
		blockThread.PushBlock(&sb3.Block{
			Opcode:   opcode,
			Inputs:   map[string]*sb3.Input{},
			Fields:   struct{}{},
			TopLevel: true,
		})
		blockThread.PushBlock(procedure.Call())
	}
}
