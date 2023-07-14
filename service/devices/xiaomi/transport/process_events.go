package transport

import "log"

type EventConsumeFunc func(data string)

func (t *Transport) RegisterHeartBeatConsumer(sid string, consumeFunc EventConsumeFunc) {
	t.heartBeatConsumersMutex.Lock()
	defer t.heartBeatConsumersMutex.Unlock()

	t.heartBeatConsumers[sid] = consumeFunc
}

func (t *Transport) RegisterReportConsumer(sid string, consumeFunc EventConsumeFunc) {
	t.reportConsumersMutex.Lock()
	defer t.reportConsumersMutex.Unlock()

	t.reportConsumers[sid] = consumeFunc
}

func (t *Transport) processIncomingEvent(msg *message) {
	switch msg.Cmd {
	case eventHeartbeat:
		func() {
			t.heartBeatConsumersMutex.Lock()
			defer t.heartBeatConsumersMutex.Unlock()

			if consumerFunc, fnd := t.heartBeatConsumers[msg.Sid]; fnd {
				go consumerFunc(msg.Data)
			}
		}()
	case eventReport:
		func() {
			t.reportConsumersMutex.Lock()
			defer t.reportConsumersMutex.Unlock()

			if consumerFunc, fnd := t.reportConsumers[msg.Sid]; fnd {
				go consumerFunc(msg.Data)
			}
		}()
	default:
		log.Printf("Unknown event: %+v", msg)
	}
}
