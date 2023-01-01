package nats

import (
	"fmt"
	"log"

	pb "github.com/digital-feather/cryptellation/services/ticks/clients/go/proto"
	"github.com/digital-feather/cryptellation/services/ticks/pkg/models/tick"
	"github.com/nats-io/go-nats"
	"google.golang.org/protobuf/proto"
)

const (
	// SubjectFormat = 'Ticks.<exchange>.<pairSymbol>'
	SubjectFormat = "Ticks.%s.%s"
)

type subWitChan struct {
	sub *nats.Subscription
	ch  chan tick.Tick
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

func (c *Client) Publish(tick tick.Tick) error {
	subject := fmt.Sprintf(SubjectFormat, tick.Exchange, tick.PairSymbol)
	data, err := proto.Marshal(tick.ToProtoBuf())
	if err != nil {
		return err
	}

	return c.natsConn.Publish(subject, data)
}

func (c *Client) Subscribe(symbol string) (<-chan tick.Tick, error) {
	ch := make(chan tick.Tick)

	subject := fmt.Sprintf(SubjectFormat, "*", symbol)
	sub, err := c.natsConn.Subscribe(subject, func(msg *nats.Msg) {
		pbTick := &pb.Tick{}
		if err := proto.Unmarshal(msg.Data, pbTick); err != nil {
			log.Printf("error when receiving tick from %s: %s\n", subject, err)
			return
		}

		t, err := tick.FromProtoBuf(pbTick)
		if err != nil {
			log.Printf("error when decoding protobuf tick from %s: %s\n", subject, err)
			return
		}

		ch <- t
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
