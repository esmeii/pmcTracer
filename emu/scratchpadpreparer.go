package emu

import (
	"log"

	"gitlab.com/yaotsu/gcn3/insts"
)

// ScratchpadPreparer is the unit that sets the instruction scratchpad
// before the instruction can be emulated.
type ScratchpadPreparer interface {
	// Prepare reads from the register file and write into the instruction
	// scratchpad
	Prepare(instEmuState InstEmuState, wf *Wavefront)

	// Commit write to the register file to reflect the change in the scratchpad
	Commit(instEmuState InstEmuState, wf *Wavefront)
}

// ScratchpadPreparerImpl reads and write registers for the emulator
type ScratchpadPreparerImpl struct {
}

// NewScratchpadPreparerImpl returns a newly created ScratchpadPreparerImpl,
// injecting the dependency of the RegInterface.
func NewScratchpadPreparerImpl() *ScratchpadPreparerImpl {
	p := new(ScratchpadPreparerImpl)
	return p
}

// Prepare read from the register file and sets the scratchpad layout
func (p *ScratchpadPreparerImpl) Prepare(
	instEmuState InstEmuState,
	wf *Wavefront,
) {
	p.clear(instEmuState.Scratchpad())
	inst := instEmuState.Inst()
	switch inst.FormatType {
	case insts.Sop1:
		p.prepareSOP1(instEmuState, wf)
	case insts.Sop2:
		p.prepareSOP2(instEmuState, wf)
	case insts.Sopc:
		p.prepareSOPC(instEmuState, wf)
	case insts.Vop1:
		p.prepareVOP1(instEmuState, wf)
	case insts.Vop2:
		p.prepareVOP2(instEmuState, wf)
	case insts.Vop3:
		p.prepareVOP3(instEmuState, wf)
	case insts.Vopc:
		p.prepareVOPC(instEmuState, wf)
	case insts.Flat:
		p.prepareFlat(instEmuState, wf)
	case insts.Smem:
		p.prepareSMEM(instEmuState, wf)
	case insts.Sopp:
		p.prepareSOPP(instEmuState, wf)
	default:
		log.Panicf("Inst format %s is not supported", inst.Format.FormatName)
	}
}

func (p *ScratchpadPreparerImpl) prepareSOP1(
	instEmuState InstEmuState,
	wf *Wavefront,
) {
	inst := instEmuState.Inst()
	scratchPad := instEmuState.Scratchpad()

	p.readOperand(inst.Src0, wf, 0, scratchPad[0:8])
	copy(scratchPad[24:25], wf.ReadReg(insts.Regs[insts.Scc], 1, 0))
	copy(scratchPad[16:24], wf.ReadReg(insts.Regs[insts.Exec], 1, 0))
}

func (p *ScratchpadPreparerImpl) prepareSOP2(
	instEmuState InstEmuState,
	wf *Wavefront,
) {
	inst := instEmuState.Inst()
	scratchPad := instEmuState.Scratchpad()
	// if inst.Src0.Register != nil && inst.Src1.Register != nil {
	// 	fmt.Println(inst, inst.Src0.Register.Name, inst.Src0.RegCount, inst.Src1.Register.Name)
	// }
	p.readOperand(inst.Src0, wf, 0, scratchPad[0:8])
	p.readOperand(inst.Src1, wf, 0, scratchPad[8:16])
	copy(scratchPad[24:25], wf.ReadReg(insts.Regs[insts.Scc], 1, 0))
}

func (p *ScratchpadPreparerImpl) prepareVOP1(
	instEmuState InstEmuState,
	wf *Wavefront,
) {
	inst := instEmuState.Inst()
	sp := instEmuState.Scratchpad()

	copy(sp[0:8], wf.ReadReg(insts.Regs[insts.Exec], 1, 0))
	copy(sp[520:528], wf.ReadReg(insts.Regs[insts.Vcc], 1, 0))

	offset := 528
	for i := 0; i < 64; i++ {
		p.readOperand(inst.Src0, wf, i, sp[offset:offset+8])
		offset += 8
	}
}

func (p *ScratchpadPreparerImpl) prepareVOP2(
	instEmuState InstEmuState,
	wf *Wavefront,
) {
	inst := instEmuState.Inst()
	sp := instEmuState.Scratchpad()

	copy(sp[0:8], wf.ReadReg(insts.Regs[insts.Exec], 1, 0))
	copy(sp[520:528], wf.ReadReg(insts.Regs[insts.Vcc], 1, 0))

	dstOffset := 8
	src0Offset := 528
	src1Offset := 1040
	for i := 0; i < 64; i++ {
		p.readOperand(inst.Dst, wf, i, sp[dstOffset:dstOffset+8])
		dstOffset += 8
		p.readOperand(inst.Src0, wf, i, sp[src0Offset:src0Offset+8])
		src0Offset += 8
		p.readOperand(inst.Src1, wf, i, sp[src1Offset:src1Offset+8])
		src1Offset += 8
	}
}

func (p *ScratchpadPreparerImpl) prepareVOP3(
	instEmuState InstEmuState,
	wf *Wavefront,
) {
	inst := instEmuState.Inst()
	sp := instEmuState.Scratchpad()

	copy(sp[0:8], wf.ReadReg(insts.Regs[insts.Exec], 1, 0))
	copy(sp[520:528], wf.ReadReg(insts.Regs[insts.Vcc], 1, 0))

	src0Offset := 528
	src1Offset := 1040
	src2Offset := 1552
	for i := 0; i < 64; i++ {
		p.readOperand(inst.Src0, wf, i, sp[src0Offset:src0Offset+8])
		src0Offset += 8
		p.readOperand(inst.Src1, wf, i, sp[src1Offset:src1Offset+8])
		src1Offset += 8
		if inst.Src2 != nil {
			p.readOperand(inst.Src2, wf, i, sp[src2Offset:src2Offset+8])
			src2Offset += 8
		}
	}
}

func (p *ScratchpadPreparerImpl) prepareVOPC(
	instEmuState InstEmuState,
	wf *Wavefront,
) {
	inst := instEmuState.Inst()
	sp := instEmuState.Scratchpad()

	src0Offset := 16
	src1Offset := 16 + 64*8
	for i := 0; i < 64; i++ {
		p.readOperand(inst.Src0, wf, i, sp[src0Offset:src0Offset+8])
		src0Offset += 8
		p.readOperand(inst.Src1, wf, i, sp[src1Offset:src1Offset+8])
		src1Offset += 8
	}

	layout := sp.AsVOPC()
	layout.EXEC = wf.Exec
}

func (p *ScratchpadPreparerImpl) prepareFlat(
	instEmuState InstEmuState, wf *Wavefront,
) {
	inst := instEmuState.Inst()
	sp := instEmuState.Scratchpad()

	copy(sp[0:8], wf.ReadReg(insts.Regs[insts.Exec], 1, 0))

	for i := 0; i < 64; i++ {
		p.readOperand(inst.Addr, wf, i, sp[8+i*8:8+i*8+8])
		p.readOperand(inst.Data, wf, i, sp[520+i*16:520+i*16+16])
	}
}

func (p *ScratchpadPreparerImpl) prepareSMEM(
	instEmuState InstEmuState,
	wf *Wavefront,
) {
	inst := instEmuState.Inst()
	scratchpad := instEmuState.Scratchpad()

	if inst.Opcode >= 16 && inst.Opcode <= 26 { // Store instructions
		p.readOperand(inst.Data, wf, 0, scratchpad[0:16])
	}

	p.readOperand(inst.Offset, wf, 0, scratchpad[16:24])
	p.readOperand(inst.Base, wf, 0, scratchpad[24:32])
}

func (p *ScratchpadPreparerImpl) prepareSOPP(
	instEmuState InstEmuState,
	wf *Wavefront,
) {
	inst := instEmuState.Inst()
	scratchPad := instEmuState.Scratchpad()
	layout := scratchPad.AsSOPP()

	layout.PC = wf.PC
	layout.SCC = wf.SCC
	layout.EXEC = wf.Exec
	p.readOperand(inst.SImm16, wf, 0, scratchPad[16:24])
}

func (p *ScratchpadPreparerImpl) prepareSOPC(
	instEmuState InstEmuState,
	wf *Wavefront,
) {
	inst := instEmuState.Inst()
	scratchPad := instEmuState.Scratchpad()

	p.readOperand(inst.Src0, wf, 0, scratchPad[0:8])
	p.readOperand(inst.Src1, wf, 0, scratchPad[8:16])
}

// Commit write to the register file according to the scratchpad layout
func (p *ScratchpadPreparerImpl) Commit(
	instEmuState InstEmuState,
	wf *Wavefront,
) {
	inst := instEmuState.Inst()
	switch inst.FormatType {
	case insts.Sop1:
		p.commitSOP1(instEmuState, wf)
	case insts.Sop2:
		p.commitSOP2(instEmuState, wf)
	case insts.Vop1:
		p.commitVOP1(instEmuState, wf)
	case insts.Vop2:
		p.commitVOP2(instEmuState, wf)
	case insts.Vop3:
		p.commitVOP3A(instEmuState, wf)
	case insts.Vopc:
		p.commitVOPC(instEmuState, wf)
	case insts.Flat:
		p.commitFlat(instEmuState, wf)
	case insts.Smem:
		p.commitSMEM(instEmuState, wf)
	case insts.Sopp:
		p.commitSOPP(instEmuState, wf)
	case insts.Sopc:
		p.commitSOPC(instEmuState, wf)
	default:
		log.Panicf("Inst format %s is not supported", inst.Format.FormatName)
	}
}

func (p *ScratchpadPreparerImpl) commitSOP1(
	instEmuState InstEmuState,
	wf *Wavefront,
) {
	inst := instEmuState.Inst()
	scratchpad := instEmuState.Scratchpad()

	p.writeOperand(inst.Dst, wf, 0, scratchpad[8:16])
	wf.WriteReg(insts.Regs[insts.Exec], 1, 0, scratchpad[16:24])
	wf.WriteReg(insts.Regs[insts.Scc], 1, 0, scratchpad[24:25])
}

func (p *ScratchpadPreparerImpl) commitSOP2(
	instEmuState InstEmuState,
	wf *Wavefront,
) {
	inst := instEmuState.Inst()
	scratchpad := instEmuState.Scratchpad()
	p.writeOperand(inst.Dst, wf, 0, scratchpad[16:24])
	wf.WriteReg(insts.Regs[insts.Scc], 1, 0, scratchpad[24:25])
}

func (p *ScratchpadPreparerImpl) commitVOP1(
	instEmuState InstEmuState,
	wf *Wavefront,
) {
	inst := instEmuState.Inst()
	scratchpad := instEmuState.Scratchpad()

	wf.WriteReg(insts.Regs[insts.Vcc], 1, 0, scratchpad[520:528])

	for i := 63; i >= 0; i-- {
		offset := 8 + i*8
		p.writeOperand(inst.Dst, wf, i, scratchpad[offset:offset+8])
	}
}

func (p *ScratchpadPreparerImpl) commitVOP2(
	instEmuState InstEmuState,
	wf *Wavefront,
) {
	inst := instEmuState.Inst()
	scratchpad := instEmuState.Scratchpad()

	wf.WriteReg(insts.Regs[insts.Vcc], 1, 0, scratchpad[520:528])

	for i := 63; i >= 0; i-- {
		offset := 8 + i*8
		p.writeOperand(inst.Dst, wf, i, scratchpad[offset:offset+8])
	}
}

func (p *ScratchpadPreparerImpl) commitVOP3A(
	instEmuState InstEmuState,
	wf *Wavefront,
) {
	inst := instEmuState.Inst()
	sp := instEmuState.Scratchpad()

	wf.WriteReg(insts.Regs[insts.Vcc], 1, 0, sp[520:528])

	for i := 63; i >= 0; i-- {
		offset := 8 + i*8
		p.writeOperand(inst.Dst, wf, i, sp[offset:offset+8])
	}
}

func (p *ScratchpadPreparerImpl) commitVOPC(
	instEmuState InstEmuState,
	wf *Wavefront,
) {
	sp := instEmuState.Scratchpad().AsVOPC()
	wf.VCC = sp.VCC
	wf.Exec = sp.EXEC
}

func (p *ScratchpadPreparerImpl) commitFlat(
	instEmuState InstEmuState,
	wf *Wavefront,
) {
	inst := instEmuState.Inst()
	scratchpad := instEmuState.Scratchpad()

	if inst.Opcode < 24 || inst.Opcode > 31 { // Skip store instructions
		for i := 0; i < 64; i++ {
			p.writeOperand(inst.Dst, wf, i, scratchpad[1544+i*16:1544+i*16+16])
		}
	}
}

func (p *ScratchpadPreparerImpl) commitSMEM(
	instEmuState InstEmuState,
	wf *Wavefront,
) {
	inst := instEmuState.Inst()
	scratchpad := instEmuState.Scratchpad()

	if inst.Opcode <= 12 { // Load instructions
		p.writeOperand(inst.Data, wf, 0, scratchpad[32:96])
	}
}

func (p *ScratchpadPreparerImpl) commitSOPC(
	instEmuState InstEmuState,
	wf *Wavefront,
) {
	scratchpad := instEmuState.Scratchpad()
	wf.SCC = scratchpad.AsSOPC().SCC
}

func (p *ScratchpadPreparerImpl) commitSOPP(
	instEmuState InstEmuState,
	wf *Wavefront,
) {
	scratchpad := instEmuState.Scratchpad()
	wf.PC = scratchpad.AsSOPP().PC
}

func (p *ScratchpadPreparerImpl) readOperand(
	operand *insts.Operand,
	wf *Wavefront,
	laneID int,
	buf []byte,
) {
	switch operand.OperandType {
	case insts.RegOperand:
		copy(buf, wf.ReadReg(operand.Register, operand.RegCount, laneID))
	case insts.IntOperand:
		copy(buf, insts.Uint64ToBytes(uint64(operand.IntValue)))
	case insts.LiteralConstant:
		copy(buf, insts.Uint32ToBytes(operand.LiteralConstant))
	default:
		log.Panicf("Operand %s is not supported", operand.String())
	}
}

func (p *ScratchpadPreparerImpl) writeOperand(
	operand *insts.Operand,
	wf *Wavefront,
	laneID int,
	buf []byte,
) {
	if operand.OperandType != insts.RegOperand {
		log.Panic("Can only write into reg operand")
	}

	wf.WriteReg(operand.Register, operand.RegCount, laneID, buf)
}

func (p *ScratchpadPreparerImpl) clear(buf []byte) {
	for i := 0; i < len(buf); i++ {
		buf[i] = 0
	}
}
