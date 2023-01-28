package ikafka

type Event struct {
	ID        string `json:"id"`
	Payload   []byte `json:"Payload"`
	EventType string `json:"EventType"`
}
