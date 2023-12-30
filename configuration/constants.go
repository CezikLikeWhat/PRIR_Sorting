package configuration

import "time"

// ENV
var (
	IsDebugMode   = false
	TimeMeasuring = false
)

// Structs

type Timer struct {
	Start   time.Time
	Elapsed time.Duration
}

// Global variables

var AppTimer = Timer{}
