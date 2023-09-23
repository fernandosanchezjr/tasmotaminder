package types

import "time"

const (
	TasmotaSensorTopic          = "tele/+/SENSOR"
	TasmotaStateTopic           = "tele/+/STATE"
	TasmotaCommandTopic         = "cmnd/+/POWER"
	TasmotaSensorSuffix         = "SENSOR"
	TasmotaStateSuffix          = "STATE"
	TasmotaPowerOFF             = "OFF"
	TasmotaPowerON              = "ON"
	TasmotaCommandWildcard      = "+"
	TasmotaDefaultResetDuration = time.Second * 5
)
