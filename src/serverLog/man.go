package serverLog

type Man interface {
	ReOpen() []error
	Close()
}
