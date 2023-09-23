package types

import (
	"log"
	"sync"
	"time"
)

type PowerTimer struct {
	Power           int              `yaml:"power"`
	PowerComparison *PowerComparison `yaml:"powerComparison,omitempty"`
	RuntimeSeconds  int              `yaml:"runtimeSeconds"`
	Action          *RuleAction      `yaml:"action"`
	mtx             sync.Mutex
	started         bool
	startTime       time.Time
}

func (pt *PowerTimer) Evaluate(plug *PlugState, target RuleTarget) {
	mtx := &pt.mtx
	mtx.Lock()
	defer mtx.Unlock()

	if plug.sensorState.POWER1 == TasmotaPowerOFF {
		pt.started = false
		return
	}

	if pt.PowerComparison.Compare(plug.sensor.Energy.Power, pt.Power) && !pt.started {
		log.Println("powerTimer power check succeeded:", plug.sensor.Energy.Power, pt.PowerComparison, pt.Power)
		pt.started = true
		pt.startTime = time.Now()
	} else if !pt.started {
		return
	}

	runtime := int(time.Since(pt.startTime).Seconds())
	if runtime < pt.RuntimeSeconds {
		return
	}

	log.Println(
		"powerTimer runtime check succeeded:",
		runtime,
		">=",
		pt.RuntimeSeconds,
	)
	pt.Action.Execute(target)
	pt.started = false
}
