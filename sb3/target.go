package sb3

type TargetHnd struct {
	sb3    *SB3
	target *Target
}

type RegisteredBlock struct {
	id    string
	block *Block
}

func (hnd *TargetHnd) registerBlock(block *Block) *RegisteredBlock {
	registeredBlock := &RegisteredBlock{id: generateId(), block: block}
	hnd.target.Blocks[registeredBlock.id] = block
	return registeredBlock
}

type ProcedureHnd struct {
	target *TargetHnd
	tail   *RegisteredBlock
}

func (hnd *TargetHnd) NewProcedure(proccode string) *ProcedureHnd {
	procedure := &ProcedureHnd{target: hnd}
	prototype := hnd.registerBlock(&Block{
		Opcode: "procedures_prototype",
		Shadow: true,
		Mutation: &Mutation{
			TagName:          "mutation",
			Proccode:         proccode,
			Argumentids:      "[]",
			Argumentnames:    "[]",
			Argumentdefaults: "[]",
		},
	})
	definition := procedure.PushBlock(&Block{
		Opcode: "procedures_definition",
		Inputs: map[string]*Input{
			"custom_block": &Input{shadow: prototype.id},
		},
		TopLevel: true,
	})
	prototype.block.Parent = NonNull(definition.id)
	return procedure
}

func (hnd *ProcedureHnd) PushBlock(block *Block) *RegisteredBlock {
	registeredBlock := hnd.target.registerBlock(block)
	if hnd.tail != nil {
		hnd.tail.block.Next = NonNull(registeredBlock.id)
		registeredBlock.block.Parent = NonNull(hnd.tail.id)
	}
	hnd.tail = registeredBlock
	return registeredBlock
}
