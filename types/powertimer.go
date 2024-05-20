package types

import (
	"fmt"
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

func (pt *PowerTimer) Evaluate(state *State, plug *PlugState, rule *PlugRule, target RuleTarget) {
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
	go pt.act(state, rule, target)
}

func (pt *PowerTimer) act(state *State, rule *PlugRule, target RuleTarget) {
	defer pt.release()

	pt.Action.Execute(target)
	if rule.Notify {
		title := fmt.Sprintf("%s %s", rule.DeviceName(), pt.Action)
		body := fmt.Sprintf("Executed from rule:\n%s", rule)
		err := state.notify(title, "tasmotaminder,power", body)
		if err != nil {
			log.Printf("Error sending notification: %v", err)
		}
	}
}

func (pt *PowerTimer) release() {
	pt.mtx.Lock()
	defer pt.mtx.Unlock()

	pt.acting = false
	pt.started = false
}
