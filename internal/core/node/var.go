package node

import (
	"sync"

	nodeModel "github.com/bestruirui/bestsub/internal/models/node"
)

var (
	poolMutex   sync.RWMutex
	pool        []nodeModel.Data
	nodeExist   *exist
	nodeProcess *exist

	wgSync     sync.WaitGroup
	wgStatus   bool
	validNodes []nodeModel.Data
	validMutex sync.Mutex

	refreshMutex   sync.Mutex
	subInfoMap     = make(map[uint16]nodeModel.SimpleInfo)
	countryInfoMap = make(map[string]nodeModel.SimpleInfo)
	subAggBuf      = make(map[uint16]*infoSums)
	countryAggBuf  = make(map[string]*infoSums)
)

type infoSums struct {
	sumSpeedUp   uint64
	sumSpeedDown uint64
	sumDelay     uint64
	sumRisk      uint64
	count        uint32
}
