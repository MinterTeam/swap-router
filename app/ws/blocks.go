package ws

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/v2/blocks"
	"github.com/centrifugal/centrifuge-go"
	log "github.com/sirupsen/logrus"
)

type BlocksChannelHandler struct {
	Subscribers []NewBlockSubscriber
}

type NewBlockSubscriber interface {
	ListenNewBlock(blocks.Resource)
}

func NewBlocksChannelHandler() *BlocksChannelHandler {
	return &BlocksChannelHandler{
		Subscribers: make([]NewBlockSubscriber, 0),
	}
}

func (b *BlocksChannelHandler) AddSubscriber(sub NewBlockSubscriber) {
	b.Subscribers = append(b.Subscribers, sub)
}

func (b *BlocksChannelHandler) OnPublish(sub *centrifuge.Subscription, e centrifuge.PublishEvent) {
	var block blocks.Resource
	if err := json.Unmarshal(e.Data, &block); err != nil {
		log.Errorf("failed to unmarshal block to json: %s", err)
		return
	}

	for _, sub := range b.Subscribers {
		go sub.ListenNewBlock(block)
	}
}
