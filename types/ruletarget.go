package types

type RuleTarget interface {
	Off()
	On()
	Reset()
}
