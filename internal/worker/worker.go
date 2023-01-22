package worker

import (
	"time"

	"workSim/interfaces"
	"workSim/internal/entity"
)

type worker struct {
	maxFeels entity.Feels
	curFeels entity.Feels

	speed int

	isBusy bool
}

func NewWorker(maxFeels entity.Feels, speed int) interfaces.Worker {
	return &worker{
		maxFeels: maxFeels, curFeels: entity.Feels{
			Energy:     maxFeels.Energy,
			NeedCoffee: 0,
			NeedSmoke:  0,
		}, speed: speed,
	}
}

func (w *worker) Work() (done chan bool) {
	if w.isBusy {
		return nil
	}

	w.isBusy = true
	w.curFeels.Energy -= 2
	w.curFeels.NeedSmoke += 1
	w.curFeels.NeedCoffee += 1
	done = make(chan bool)

	go func() {
		time.Sleep(time.Second * 20 / time.Duration(w.speed))
		w.isBusy = false
		done <- !w.isDead()
	}()

	return
}

func (w *worker) Smoke() (done chan bool) {
	if w.isBusy {
		return nil
	}

	w.isBusy = true
	w.curFeels.NeedSmoke = 0
	w.curFeels.Energy -= 1
	w.curFeels.NeedCoffee += 1

	done = make(chan bool)

	go func() {
		time.Sleep(time.Second * 5 / time.Duration(w.speed))
		w.isBusy = false
		done <- !w.isDead()
	}()

	return
}

func (w *worker) DrinkCoffee() (done chan bool) {
	if w.isBusy {
		return nil
	}

	w.isBusy = true
	w.curFeels.NeedCoffee = 0
	w.curFeels.Energy -= 1
	w.curFeels.NeedSmoke += 1

	done = make(chan bool)

	go func() {
		time.Sleep(time.Second * 5 / time.Duration(w.speed))
		w.isBusy = false
		done <- !w.isDead()
	}()

	return
}

func (w *worker) Sleep() (done chan bool) {
	if w.isBusy {
		return nil
	}

	w.isBusy = true
	w.curFeels.NeedCoffee = 0
	w.curFeels.Energy = w.maxFeels.Energy
	w.curFeels.NeedSmoke += 3

	done = make(chan bool)

	go func() {
		time.Sleep(time.Second * 100 / time.Duration(w.speed))
		w.isBusy = false
		done <- !w.isDead()
	}()

	return
}

func (w *worker) CurrentFeels() entity.Feels {
	return w.curFeels
}

func (w *worker) MaxFeels() entity.Feels {
	return w.maxFeels
}

func (w *worker) isDead() bool {
	return w.curFeels.Energy < 0 || w.curFeels.NeedCoffee > w.maxFeels.NeedCoffee || w.curFeels.NeedSmoke > w.maxFeels.NeedSmoke
}
