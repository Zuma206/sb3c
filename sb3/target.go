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

type BlockThread struct {
	target *TargetHnd
	tail   *RegisteredBlock
}

func (hnd *TargetHnd) NewBlockThread() *BlockThread {
	return &BlockThread{target: hnd}
}

type ProcedureHnd struct {
	BlockThread
	proccode string
}

func (hnd *TargetHnd) NewProcedure(proccode string) *ProcedureHnd {
	procedure := &ProcedureHnd{
		BlockThread: BlockThread{
			target: hnd,
		},
		proccode: proccode,
	}
	prototype := hnd.registerBlock(&Block{
		Opcode: "procedures_prototype",
		Shadow: true,
		Mutation: &Mutation{
			TagName:          "mutation",
			Proccode:         proccode,
			Argumentids:      "[]",
			Argumentnames:    "[]",
			Argumentdefaults: "[]",
			Children:         []struct{}{},
		},
		Inputs: map[string]*Input{},
	})
	definition := procedure.PushBlock(&Block{
		Opcode: "procedures_definition",
		Inputs: map[string]*Input{
			"custom_block": {shadow: prototype.id},
		},
		TopLevel: true,
	})
	prototype.block.Parent = NonNull(definition.id)
	return procedure
}

func (hnd *ProcedureHnd) Target() *TargetHnd {
	return hnd.target
}

func (hnd *TargetHnd) Project() *SB3 {
	return hnd.sb3
}

func (hnd *TargetHnd) NewVariable(name string, value any) string {
	id := generateId() + "-" + name
	hnd.target.Variables[id] = &Variable{Name: name, Value: value}
	return id
}

func (blockThread *BlockThread) PushBlock(block *Block) *RegisteredBlock {
	registeredBlock := blockThread.target.registerBlock(block)
	if blockThread.tail != nil {
		blockThread.tail.block.Next = NonNull(registeredBlock.id)
		registeredBlock.block.Parent = NonNull(blockThread.tail.id)
	}
	blockThread.tail = registeredBlock
	return registeredBlock
}

func (hnd *ProcedureHnd) Call() *Block {
	return &Block{
		Opcode: "procedures_call",
		Inputs: map[string]*Input{},
		Fields: struct{}{},
		Mutation: &Mutation{
			TagName:     "mutation",
			Children:    []struct{}{},
			Proccode:    hnd.proccode,
			Argumentids: "[]",
		},
	}
}
