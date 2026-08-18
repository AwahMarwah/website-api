package queue

type InvalidateCachePayload struct {
	Keys []string `json:"keys"`
}
