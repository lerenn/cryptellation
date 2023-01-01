package nats

import (
	"fmt"
	"log"

	pb "github.com/digital-feather/cryptellation/services/backtests/clients/go/proto"
	"github.com/digital-feather/cryptellation/services/backtests/pkg/models/event"
	"github.com/nats-io/go-nats"
	"google.golang.org/protobuf/proto"
)

const (
	// SubjectFormat = 'Backtests.<backtest_id>'
	SubjectFormat = "Backtests.%d"
)

type subWitChan struct {
	sub *nats.Subscription
	ch  chan event.Event
}

type Client struct {
	natsConn      *nats.Conn
	subscriptions []subWitChan
}

func New() (*Client, error) {
	var c Config
	if err := c.Load().Validate(); err != nil {
		return nil, fmt.Errorf("loading nats config: %w", err)
	}

	natsConn, err := nats.Connect(c.URL)
	if err != nil {
		return nil, fmt.Errorf("connecting to nats: %w", err)
	}

	return &Client{
		natsConn:      natsConn,
		subscriptions: make([]subWitChan, 0),
	}, nil
}

func (c *Client) Publish(backtestID uint, event event.Event) error {
	pbEvent, err := event.ToProtoBuf()
	if err != nil {
		return nil
	}

	subject := fmt.Sprintf(SubjectFormat, backtestID)
	data, err := proto.Marshal(pbEvent)
	if err != nil {
		return err
	}

	return c.natsConn.Publish(subject, data)
}

func (c *Client) Subscribe(backtestID uint) (<-chan event.Event, error) {
	ch := make(chan event.Event)

	subject := fmt.Sprintf(SubjectFormat, backtestID)
	sub, err := c.natsConn.Subscribe(subject, func(msg *nats.Msg) {
		pbEvent := &pb.BacktestEvent{}
		if err := proto.Unmarshal(msg.Data, pbEvent); err != nil {
			log.Printf("error when receiving tick from %s: %s\n", subject, err)
			return
		}

		evt, err := event.FromProtoBuf(pbEvent)
		if err != nil {
			log.Printf("error when decoding protobuf tick from %s: %s\n", subject, err)
			return
		}

		ch <- evt
	})
	if err != nil {
		return nil, err
	}

	c.subscriptions = append(c.subscriptions, subWitChan{
		sub: sub,
		ch:  ch,
	})
	return ch, nil
}

func (c *Client) Close() {
	for _, swc := range c.subscriptions {
		if err := swc.sub.Unsubscribe(); err != nil {
			log.Printf("error when unsubscribing from %s: %s\n", swc.sub.Subject, err)
		}
		close(swc.ch)
	}

	c.natsConn.Close()
}
