package pubsub_collector

import "sync"

type ChannelMessageMap map[string][]string

type PubSubMessages struct {
	mx *sync.Mutex
	M  ChannelMessageMap
}

func (c *PubSubMessages) Load(key string) ([]string, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()
	val, ok := c.M[key]

	return val, ok
}

func (c *PubSubMessages) Store(channelName, msg string) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.M[channelName] = append(c.M[channelName], msg)
}

func (collector *PubSubMsgCollector) Messages() ChannelMessageMap {
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
