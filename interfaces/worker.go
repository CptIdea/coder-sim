package interfaces

import (
	"workSim/internal/entity"
)

type Worker interface {
	Work() (done chan bool)
	Smoke() (done chan bool)
	DrinkCoffee() (done chan bool)
	Sleep() (done chan bool)

	CurrentFeels() entity.Feels
	MaxFeels() entity.Feels
}
