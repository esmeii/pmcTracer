package pagemigrationcontroller

import "gitlab.com/akita/akita/v3/sim"

// A Builder can build pagemigrationControllers
type Builder struct {
	engine         sim.Engine
	freq           sim.Freq
	numReqPerCycle int
	pageSize       uint64

	remotePort   sim.Port
	ctrlPort     sim.Port
	localMemPort sim.Port
}

// MakeBuilder returns a Builder
func MakeBuilder() Builder {
	return Builder{
		freq:           1 * sim.GHz,
		numReqPerCycle: 4,
		pageSize:       4096,
	}
}

// WithEngine sets the engine that the pagemigrationControllers to use
func (b Builder) WithEngine(engine sim.Engine) Builder {
	b.engine = engine
	return b
}

// WithFreq sets the freq the pagemigrationControllers use
func (b Builder) WithFreq(freq sim.Freq) Builder {
	b.freq = freq
	return b
}

// WithPageSize sets the page size that the pagemigrationController works with.
func (b Builder) WithPageSize(n uint64) Builder {
	b.pageSize = n
	return b
}

// WithNumReqPerCycle sets the number of requests per cycle can be processed by
// a pagemigrationController
func (b Builder) WithNumReqPerCycle(n int) Builder {
	b.numReqPerCycle = n
	return b
}

// Build creates a new pagemigrationController
func (b Builder) Build(name string) *PageMigrationController {
	pagemigrationController := &PageMigrationController{}
	pagemigrationController.TickingComponent =
		sim.NewTickingComponent(name, b.engine, b.freq, pagemigrationController)
	//remotePort   sim.Port
	//	ctrlPort     sim.Port
	//	localMemPort sim.Port
	pagemigrationController.remotePort = b.remotePort
	pagemigrationController.ctrlPort = b.ctrlPort
	pagemigrationController.localMemPort = b.localMemPort

	b.createPorts(name, pagemigrationController)

	//pagemigrationController.reset()

	return pagemigrationController
}

func (b Builder) createPorts(name string, pagemigrationController *PageMigrationController) {
	pagemigrationController.remotePort = sim.NewLimitNumMsgPort(pagemigrationController, b.numReqPerCycle,
		name+".remotePort")
	pagemigrationController.AddPort("remote", pagemigrationController.remotePort)

	pagemigrationController.ctrlPort = sim.NewLimitNumMsgPort(pagemigrationController, b.numReqPerCycle,
		name+".ctrlPort")
	pagemigrationController.AddPort("ctrl", pagemigrationController.ctrlPort)

	pagemigrationController.localMemPort = sim.NewLimitNumMsgPort(pagemigrationController, 1,
		name+".localMemPort")
	pagemigrationController.AddPort("local", pagemigrationController.localMemPort)
}
