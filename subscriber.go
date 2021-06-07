package ethclient

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
)

var (
	reconnectInterval = 2 * time.Second
)

var _ Subscriber = (*ChainSubscrier)(nil)

// ChainSubscrier implements Subscriber interface
type ChainSubscrier struct {
	c *ethclient.Client
}

// NewChainSubscriber .
func NewChainSubscriber(c *ethclient.Client) (*ChainSubscrier, error) {
	return &ChainSubscrier{c}, nil
}

// SubscribeFilterlog .
func (cs *ChainSubscrier) SubscribeFilterlogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) error {
	checkChan := make(chan types.Log)
	resubscribeFunc := func() (ethereum.Subscription, error) {
		return cs.c.SubscribeFilterLogs(ctx, q, checkChan)
	}

	return cs.subscribeFilterlog(ctx, resubscribeFunc, q, checkChan, ch)
}

func (cs *ChainSubscrier) subscribeFilterlog(ctx context.Context, fn resubscribeFunc, query ethereum.FilterQuery, checkChan <-chan types.Log, resultChan chan<- types.Log) error {
	// pipeline: ethclient subscribe --> checkChan(validate log and get missing log) --> resultChan --> user

	// the goroutine for geting missing log and sending log to result channel.
	go func() {
		var lastLog *types.Log
		for {
			select {
			case result := <-checkChan:
				if lastLog != nil {
					if lastLog.BlockNumber > result.BlockNumber {
						// ignore duplicate
						continue
					} else {
						// retrieve potentially missing log
						start, end := lastLog.BlockNumber+1, result.BlockNumber
						for start < end {
							query.FromBlock = big.NewInt(int64(start))
							vlog, err := cs.c.FilterLogs(ctx, query)
							if err != nil {
								if err == context.Canceled || err == context.DeadlineExceeded {
									log.Debug("SubscribeFilterlog Filterlog exit...")
									return
								}

								log.Warn("Client subscribeFilterlog filterlog", "err", err)
								time.Sleep(reconnectInterval)
								continue
							}

							if len(vlog) != 0 {
								log.Debug("Client got missing log", "from", start, "to", end)
							}

							for _, l := range vlog {
								resultChan <- l
							}

							start = end
						}

					}
				}
				lastLog = &result
				resultChan <- result
			case <-ctx.Done():
				log.Debug("SubscribeFilterlog exit...")
				return
			}
		}
	}()

	// the goroutine to subscribe filter log and send log to check channel.
	go func() {
		for {
			log.Debug("Client resubscribe log...")

			sub, err := fn()
			switch {
			case err == context.Canceled || err == context.DeadlineExceeded:
				log.Debug("SubscribeFilterlog exit...")
				return
			case err != nil:
				log.Warn("Client resubscribelogFunc  err: ", err)
				time.Sleep(reconnectInterval)
				continue
			}

			select {
			case err := <-sub.Err():
				log.Warn("Client subscribe log err: ", err)
				sub.Unsubscribe()
				time.Sleep(reconnectInterval)
			case <-ctx.Done():
				log.Debug("SubscribeFilterlog exit...")
				return
			}
		}
	}()

	return nil
}

// SubscribeNewHead .
func (cs *ChainSubscrier) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) error {
	checkChan := make(chan *types.Header)
	resubscribeFunc := func() (ethereum.Subscription, error) {
		return cs.c.SubscribeNewHead(ctx, checkChan)
	}

	return cs.subscribeNewHead(ctx, resubscribeFunc, checkChan, ch)
}

// subscribeNewHead subscribes new header and auto reconnect if the connection lost.
func (cs *ChainSubscrier) subscribeNewHead(ctx context.Context, fn resubscribeFunc, checkChan <-chan *types.Header, resultChan chan<- *types.Header) error {
	// the goroutine for geting missing header and sending header to result channel.
	go func() {
		var lastHeader *types.Header
		for {
			select {
			case <-ctx.Done():
				log.Debug("SubscribeNewHead exit...")
				return
			case result := <-checkChan:
				if lastHeader != nil {
					if lastHeader.Number.Cmp(result.Number) >= 0 {
						// ignore duplicate
						continue
					} else {
						// get missing headers
						start, end := new(big.Int).Add(lastHeader.Number, big.NewInt(1)), result.Number
						for start.Cmp(end) < 0 {
							header, err := cs.c.HeaderByNumber(ctx, start)
							switch err {
							case context.DeadlineExceeded, context.Canceled:
								log.Debug("SubscribeNewHead HeaderByNumber exit...")
								return
							case ethereum.NotFound:
								log.Warn("Client subscribeNewHead err: header not found")
								time.Sleep(reconnectInterval)
								continue
							case nil:
								log.Debug("Client get missing header", "number", start)
								start.Add(start, big.NewInt(1))
								resultChan <- header
							default: // ! nil
								log.Warn("Client subscribeNewHead", "err", err)
								time.Sleep(reconnectInterval)
								continue
							}
						}
					}
				}
				lastHeader = result
				resultChan <- result
			}
		}
	}()

	// the goroutine to subscribe new header and send header to check channel.
	go func() {
		for {
			log.Debug("Client resubscribe...")
			sub, err := fn()
			if err != nil {
				if err == context.Canceled || err == context.DeadlineExceeded {
					log.Debug("SubscribeNewHead exit...")
					return
				}
				log.Warn("ChainClient resubscribeHeadFunc", "err", err)
				time.Sleep(reconnectInterval)
				continue
			}

			select {
			case err := <-sub.Err():
				log.Warn("ChainClient subscribe head", "err", err)
				sub.Unsubscribe()
				time.Sleep(reconnectInterval)
			}
		}
	}()

	return nil
}
