package sim

import (
	"inklings/internal/config"
)

type TravelDirection int

const (
	North TravelDirection = iota
	West
	South
	East
	NumTravelDirections
)

type InklingType int

const (
	RedTrav InklingType = iota
	GreenTrav
	BlueTrav
	NumTravTypes
)

type InklingInfo struct {
	Type InklingType
	Row  int
	Col  int
	Dir  TravelDirection
	Live bool
}

type State struct {
	Cfg *config.Config

	Grid     [][]int
	NumRows  int
	NumCols  int
	Inklings []InklingInfo

	RedLevel   int
	GreenLevel int
	BlueLevel  int

	RedLock   string //what type should this be?
	GreenLock string //what type should this be?
	BlueLock  string //what type should this be?
	GridLock  string //what type should this be? notice a pattern ;)
}

func NewState(cfg *config.Config) *State {
	s := &State{
		Cfg:        cfg,
		NumRows:    cfg.Simulation.DefaultRows,
		NumCols:    cfg.Simulation.DefaultCols,
		RedLevel:   cfg.Ink.InitialRed,
		GreenLevel: cfg.Ink.InitialGreen,
		BlueLevel:  cfg.Ink.InitialBlue,
	}
	// TODO: allocate grid, initialize inklings, etc.
	return s
}

// TODO: implement worker goroutines & producer goroutines.
// func (s *State) RunInkling(id int, wg *sync.WaitGroup) {}
// func (s *State) RunRedProducer(wg *sync.WaitGroup) {}
// func (s *State) RunGreenProducer(wg *sync.WaitGroup) {}
// func (s *State) RunBlueProducer(wg *sync.WaitGroup) {}
