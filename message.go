package event

func CreateErrorMessage(code int64, description string, meta MessageMetadata) Message {
	meta.Key = "event_error"
	return Message{
		Metadata: meta,
		Body:     []byte(description),
		Extra:    "",
	}
}
