package trade

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/longbridge/openapi-go/internal/util"
	"github.com/longbridge/openapi-go/log"
	"github.com/longbridge/openapi-go/trade/jsontypes"

	tradev1 "github.com/longportapp/openapi-protobufs/gen/go/trade"
	protocol "github.com/longportapp/openapi-protocol/go"
	"github.com/longportapp/openapi-protocol/go/client"
	"github.com/pkg/errors"
)

type core struct {
	client        client.Client
	url           string
	subscriptions []string
	mu            sync.Mutex
}

func newCore(opts *Options) (*core, error) {
	getOTP := func() (string, error) {
		otp, err := opts.httpClient.GetOTPV2(context.Background())
		if err != nil {
			return "", errors.Wrap(err, "failed to get otp")
		}
		return otp, nil
	}

	logger := opts.logger
	logger.SetLevel(opts.logLevel)

	cl := client.New(client.WithLogger(logger))
	err := cl.Dial(context.Background(), opts.tradeURL,
		&protocol.Handshake{
			Version:  1,
			Codec:    protocol.CodecProtobuf,
			Platform: protocol.PlatformOpenapi,
		},
		client.WithAuthTokenGetter(getOTP),
		client.AuthTimeout(opts.lbOpts.AuthTimeout),
		client.DialTimeout(opts.lbOpts.Timeout),
		client.MinGzipSize(opts.lbOpts.MinGzipSize),
		client.ReadBufferSize(opts.lbOpts.ReadBufferSize),
		client.ReadQueueSize(opts.lbOpts.ReadQueueSize),
		client.WriteQueueSize(opts.lbOpts.WriteQueueSize),
	)
	if err != nil {
		return nil, err
	}

	core := &core{client: cl, url: opts.tradeURL}

	core.client.AfterReconnected(func() {
		resubFlag := true
		if err := core.resubscribe(context.Background()); err != nil {
			log.Errorf("failed to do sub, err: %v", err)
			resubFlag = false
		}
		for _, fn := range opts.reconnectCallbacks {
			fn(resubFlag)
		}
	})

	return core, nil
}

func (c *core) SetHandler(f func(*PushEvent)) {
	c.client.Subscribe(uint32(tradev1.Command_CMD_NOTIFY), parseNotifyFunc(f))
}

func (c *core) Subscribe(ctx context.Context, topics []string) (subRes *SubResponse, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.doSubscribe(ctx, topics)
}

func (c *core) doSubscribe(ctx context.Context, topics []string) (subRes *SubResponse, err error) {
	var res *protocol.Packet
	req := &tradev1.Sub{Topics: topics}
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(tradev1.Command_CMD_SUB), Body: req})
	if err != nil {
		return
	}
	var tradeRes tradev1.SubResponse
	if err = res.Unmarshal(&tradeRes); err != nil {
		return
	}
	subRes = &SubResponse{}
	subRes.Current = tradeRes.GetCurrent()
	subRes.Success = tradeRes.GetSuccess()
	subRes.Fail = make([]*SubResponseFail, len(tradeRes.GetFail()))
	for _, f := range tradeRes.GetFail() {
		subRes.Fail = append(subRes.Fail, &SubResponseFail{Topic: f.GetTopic(), Reason: f.GetReason()})
	}
	c.subscriptions = subRes.Current
	return
}

func (c *core) Unsubscribe(ctx context.Context, topics []string) (unsubRes *UnsubResponse, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var res *protocol.Packet
	req := &tradev1.Unsub{Topics: topics}
	res, err = c.client.Do(ctx, &client.Request{Cmd: uint32(tradev1.Command_CMD_UNSUB), Body: req})
	if err != nil {
		return
	}
	var tradeRes tradev1.UnsubResponse
	if err = res.Unmarshal(&tradeRes); err != nil {
		return
	}
	unsubRes = &UnsubResponse{}
	unsubRes.Current = tradeRes.GetCurrent()
	c.subscriptions = tradeRes.GetCurrent()
	return
}

func (c *core) resubscribe(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	res, err := c.doSubscribe(ctx, c.subscriptions)
	if err != nil {
		return errors.Wrap(err, "resubscribe error")
	}
	if len(res.Fail) > 0 {
		log.Errorf("resubscirbe subscription some failed %v", res.Fail)
	}
	return nil
}

func (c *core) Close() error {
	return c.client.Close(nil)
}

func parseNotifyFunc(f func(*PushEvent)) func(*protocol.Packet) {
	return func(packet *protocol.Packet) {
		var notify tradev1.Notification
		if err := packet.Unmarshal(&notify); err != nil {
			log.Errorf("trade context unmarshal notification error:%v", err)
			return
		}
		var data jsontypes.PushEvent
		if err := json.Unmarshal(notify.GetData(), &data); err != nil {
			log.Errorf("trade context json unmarshal push event error:%v", err)
			return
		}
		var event PushEvent
		if err := util.Copy(&event, data); err != nil {
			log.Errorf("trade context copy push event error:%v", err)
			return
		}
		f(&event)
	}
}
