package health

import (
	"fmt"
	"log"
	"runtime/debug"
	"time"
)

// ProbeType defines probe type of info.
type ProbeType string

// ProbeType enums.
const (
	ProbeNone  ProbeType = "none"
	ProbeAlive ProbeType = "alive"
	ProbeReady ProbeType = "ready"
)

// Status is status, don't use it.
type Status string

const (
	statInit    Status = "init"
	statRunning Status = "running"
	statPause   Status = "pause"
	statExited  Status = "exited"
)

// Info defines info.
type Info struct {
	// for json marshalling only.
	Expired   bool
	Vars      map[string]any
	Status    Status
	Age       string
	ProbeType ProbeType

	taskName string
	lastTime time.Time
	interval time.Duration
}

// NewInfo creates info, name must be unique.
func NewInfo(name string, d time.Duration, pt ProbeType) *Info {
	if _, ok := infos.Load(name); ok {
		log.Panicln("name is duplicated:", name)
	}

	return &Info{
		taskName:  name,
		Status:    statInit,
		interval:  d,
		ProbeType: pt,
	}
}

// UpdateVars updates info.
func (i *Info) UpdateVars(vars map[string]any) {
	i.lastTime = time.Now()
	i.Vars = vars
	infos.Store(i.taskName, *i)
}

// Up marks this task status to up.
func (i *Info) Up() {
	i.lastTime = time.Now()
	i.Status = statRunning
	infos.Store(i.taskName, *i)
}

// Down marks this task status to down and is unrecoverable.
func (i *Info) Down() {
	i.lastTime = time.Now()
	i.Status = statExited
	infos.Store(i.taskName, *i)

	log.Printf("task(%v) is exiting", i.taskName)

	// Panic handling.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered: ", r)
			debug.PrintStack()
		}
	}()
}

// Pause marks this task status to pause which is under self healing.
func (i *Info) Pause() {
	i.lastTime = time.Now()
	i.Status = statPause
	infos.Store(i.taskName, *i)
}
