package pluto

type JoinChannel struct {
}

func (p JoinChannel) Process(processable Processable) (Processable, bool) {
	appendable, ok := processable.GetBody().(Appendable)
	if !ok {
		return processable, false
	}

	v, found := appendable["channel"]
	if !found {
		return processable, false
	}

	channel, ok := v.(Channel)
	if !ok {
		return processable, false
	}

	v, found = appendable["identifier"]
	if !found {
		return processable, false
	}

	identifier, ok := v.(Identifier)
	if !ok {
		return processable, false
	}

	v, found = appendable["processor"]
	if !found {
		return processable, false
	}

	processor, ok := v.(Processor)
	if !ok {
		return processable, false
	}

	channel.Join(&DynamicJoinable{identifier, processor})

	return processable, true
}

func (p JoinChannel) GetDescriptor() ProcessorDescriptor {
	return ProcessorDescriptor{
		Name:        "JOIN_CHANNEL",
		Description: "",
		Input:       "channel,identifier,processor",
		Output:      "",
	}
}

type DynamicJoinable struct {
	Identifier
	Processor
}
