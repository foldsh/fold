package watcher

import "time"

type Debouncer struct {
	rate       time.Duration
	onChangeFn func()

	callQueued bool
}

func NewDebouncer(rate time.Duration, onChange func()) *Debouncer {
	return &Debouncer{rate: rate, onChangeFn: onChange}
}

func (d *Debouncer) OnChange() {
	if d.callQueued {
		return
	}
	d.callQueued = true
	go func() {
		select {
		case <-time.After(d.rate):
			d.onChangeFn()
			d.callQueued = false
		}
	}()
}
