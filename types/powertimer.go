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
	acting          bool
	startTime       time.Time
}

func (pt *PowerTimer) Evaluate(plug *PlugState, target RuleTarget) {
	mtx := &pt.mtx
	mtx.Lock()
	defer mtx.Unlock()

	if pt.acting {
		log.Println("plug performing action, waiting...")
		return
	}

	if plug.sensorState.POWER1 == TasmotaPowerOFF {
		log.Println("plug state", plug.sensorState.POWER1)

		go pt.release()
		return
	}

	log.Println("powerComparison:", plug.sensor.Energy.Power, pt.PowerComparison, pt.Power, "?")

	if pt.PowerComparison.Compare(plug.sensor.Energy.Power, pt.Power) && !pt.started {
		pt.started = true
		pt.startTime = time.Now()
	} else if !pt.started {
		return
	}

	log.Println("powerComparison condition met")

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

	log.Println("runtimeSeconds condition met")

	pt.acting = true
	go pt.act(target)
}

func (pt *PowerTimer) act(target RuleTarget) {
	defer pt.release()

	pt.Action.Execute(target)
}

func (pt *PowerTimer) release() {
	pt.mtx.Lock()
	defer pt.mtx.Lock()

	pt.acting = false
	pt.started = false
}
