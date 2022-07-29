package messagebroker

// Producer
type Publisher interface {
	Start() error
	Stop() error
	Publish(key, body []byte, logBody string) error
}