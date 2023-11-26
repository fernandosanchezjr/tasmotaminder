package types

type PowerSchedule struct {
	Cron   string      `yaml:"cron"`
	Action *RuleAction `yaml:"action,omitempty"`
}
