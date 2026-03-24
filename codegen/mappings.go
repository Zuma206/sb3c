package codegen

type BlockMapping struct {
	Opcode string
	Inputs []string
}

var mappings = map[string]*BlockMapping{
	"this.motion.moveSteps": {Opcode: "motion_movesteps", Inputs: []string{"STEPS"}},
}
