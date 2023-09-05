package types

import (
	"encoding/json"
	"log"
)

type SensorEnergy struct {
	TotalStartTime *UTCTime `json:"TotalStartTime"`
	Total          float64  `json:"Total"`
	Yesterday      float64  `json:"Yesterday"`
	Today          float64  `json:"Today"`
	Period         int      `json:"Period"`
	Power          int      `json:"Power"`
	ApparentPower  int      `json:"ApparentPower"`
	ReactivePower  int      `json:"ReactivePower"`
	Factor         float64  `json:"Factor"`
	Voltage        int      `json:"Voltage"`
	Current        float64  `json:"Current"`
}

type Sensor struct {
	Time   *UTCTime      `json:"Time"`
	Energy *SensorEnergy `json:"ENERGY"`
}

func NewSensor() *Sensor {
	return &Sensor{Energy: &SensorEnergy{}}
}

func (s *Sensor) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		log.Fatalf("error marshalling Sensor to JSON: %s", err)
	}
	return string(data)
}

func (se *SensorEnergy) String() string {
	data, err := json.Marshal(se)
	if err != nil {
		log.Fatalf("error marshalling SensorEnergy to JSON: %s", err)
	}
	return string(data)
}
