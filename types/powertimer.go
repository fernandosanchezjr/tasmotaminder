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
	Action          *RuleAction      `yaml:"action,omitempty"`
	mtx             sync.Mutex
	started         bool
	startTime       time.Time
}

func (pt *PowerTimer) Evaluate(plug *PlugState, target RuleTarget) {
	mtx := &pt.mtx
	mtx.Lock()
	defer mtx.Unlock()

	if plug.sensorState.POWER1 == TasmotaPowerOFF {
		log.Println("plug state", plug.sensorState.POWER1)
		pt.started = false
		return
	}

	log.Println("powerComparison:", plug.sensor.Energy.Power, pt.PowerComparison, pt.Power, "?")

	if pt.PowerComparison.Compare(plug.sensor.Energy.Power, pt.Power) && !pt.started {
		pt.started = true
		pt.startTime = time.Now()
	} else if !pt.started {
		return
	}

	log.Println("powerComparison succeeded")

	runtime := int(time.Since(pt.startTime).Seconds())

	log.Println(
		"runtimeSeconds check:",
		runtime,
		">=",
		pt.RuntimeSeconds,
		"?",
	)

	if runtime < pt.RuntimeSeconds {
		return
	}

	log.Println("runtimeSeconds succeeded")

	pt.Action.Execute(target)
	pt.started = false
}
