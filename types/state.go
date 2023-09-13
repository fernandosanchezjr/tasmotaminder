package types

import (
	"encoding/json"
	"log"
)

type Wifi struct {
	AP        int    `json:"AP"`
	SSId      string `json:"SSId"`
	BSSId     string `json:"BSSId"`
	Channel   int    `json:"Channel"`
	Mode      string `json:"Mode"`
	RSSI      int    `json:"RSSI"`
	Signal    int    `json:"Signal"`
	LinkCount int    `json:"LinkCount"`
	Downtime  string `json:"Downtime"`
}

type State struct {
	Time      string  `json:"Time"`
	Uptime    string  `json:"Uptime"`
	UptimeSec uint64  `json:"UptimeSec"`
	Heap      uint64  `json:"Heap"`
	SleepMode string  `json:"SleepMode"`
	Sleep     uint64  `json:"Sleep"`
	LoadAvg   float64 `json:"LoadAvg"`
	MqttCount uint64  `json:"MqttCount"`
	POWER1    string  `json:"POWER1"`
	Wifi      *Wifi   `json:"Wifi"`
}

func NewState() *State {
	return &State{Wifi: &Wifi{}}
}

func (s *State) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		log.Fatalf("error marshalling State to JSON: %s", err)
	}
	return string(data)
}

func (w *Wifi) String() string {
	data, err := json.Marshal(w)
	if err != nil {
		log.Fatalf("error marshalling Sensor to JSON: %s", err)
	}
	return string(data)
}

func (s *State) Unmarshal(payload []byte) {
	if err := json.Unmarshal(payload, s); err != nil {
		log.Fatalf("error unmarshalling state: %s", err)
	}
}
