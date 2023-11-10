package pluto

type Processor interface {
	ID() string
	// Process
	// The boolean indicates that the next processor can be executed or not.
	Process(Processable) (Processable, bool)
}
