package main

type Configuration struct {
	ServerUrl                 string
	EventTopicName            string
	HeartbeatTopicName        string
	DeviceId                  string
	HeartbeatFrequencySeconds uint64
}
type IotEvent struct {
	DeviceId         string `json:"device_id"`
	Timestamp        string `json:"timestamp"`
	EventDescription string `json:"event_description,omitempty"`
	EventType        string `json:"event_type"`
}
