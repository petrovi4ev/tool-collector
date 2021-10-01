package socket_responses_collector

import "sync"

// ResponsesMap [request][response, response, response, response]
type ResponsesMap map[string][]string

type ResponseMessages struct {
	mx *sync.Mutex
	M  ResponsesMap
}

func (c *ResponseMessages) Load(key string) ([]string, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()
	val, ok := c.M[key]

	return val, ok
}

func (c *ResponseMessages) Store(channelRequest, response string) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.M[channelRequest] = append(c.M[channelRequest], response)
}

func (collector *ResponsesCollector) Messages() ResponsesMap {
	collector.messages.mx.Lock()
	result := make(map[string][]string)
	for k, v := range collector.messages.M {
		messages := make([]string, len(v))
		copy(messages, v)
		result[k] = messages
	}
	collector.messages.mx.Unlock()

	return result
}
