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

// SubscribeFilterlog support getting logs from `From` block to `To` block and
// auto reconnect if network disconnected.
func (cs *ChainSubscrier) SubscribeFilterlogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) error {
	// checkChan := make(chan types.Log)

	// Support from `From` block to latest block.
	logs, err := cs.c.FilterLogs(ctx, q)
	if err != nil {
		return err
	}

	checkChan := make(chan types.Log, len(logs))

	for _, l := range logs {
		checkChan <- l
	}

	resubscribeFunc := func() (ethereum.Subscription, error) {
		return cs.c.SubscribeFilterLogs(ctx, q, checkChan)
	}

	return cs.subscribeFilterlog(ctx, resubscribeFunc, q, checkChan, ch)
}

func (cs *ChainSubscrier) subscribeFilterlog(ctx context.Context, fn resubscribeFunc, query ethereum.FilterQuery, checkChan <-chan types.Log, resultChan chan<- types.Log) error {
	// Pipeline: ethclient subscribe --> checkChan(validate log and get missing log) --> resultChan --> user

	// Report whether the comming log has seen.
	hasSeen := func(lastLog, commingLog types.Log) bool {
		if lastLog.BlockNumber > commingLog.BlockNumber {
			return true
		} else if lastLog.BlockNumber == commingLog.BlockNumber {
			if lastLog.TxIndex > commingLog.TxIndex {
				return true
			} else if lastLog.TxIndex == commingLog.TxIndex &&
				lastLog.Index >= commingLog.Index {
				return true
			}
		}

		return false
	}

	// The goroutine for geting missing log and sending log to result channel.
	go func() {
		var lastLog *types.Log
		for {
			select {
			case commingLog := <-checkChan:
				if lastLog != nil {
					if hasSeen(*lastLog, commingLog) {
						log.Warn("Duplicate logs", "block", commingLog.BlockNumber, "tx", commingLog.TxHash.Hex(),
							"txIndex", commingLog.TxIndex, "index", commingLog.Index)
						continue
					} else {
						// Lost some logs between lastLog and commingLog if the network disconnected.
						// Retrieve potentially missing log and make sure not duplicate.

						// TODO: There are many duplicate logs, and optimize here in future.
						start, end := lastLog.BlockNumber, commingLog.BlockNumber
						for start <= end {
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
								l := l
								if hasSeen(*lastLog, l) {
									log.Debug("Duplicate logs", "block", l.BlockNumber, "tx", l.TxHash.Hex(),
										"txIndex", l.TxIndex, "index", l.Index, "last", *lastLog)
									continue
								}
								lastLog = &l
								resultChan <- l
							}

							start = end + 1
						}
					}
				} else {
					lastLog = &commingLog
					resultChan <- commingLog
				}
			case <-ctx.Done():
				log.Debug("SubscribeFilterlog exit...")
				return
			}
		}
	}()

	// The goroutine to subscribe filter log and send log to check channel.
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
	// The goroutine for geting missing header and sending header to result channel.
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
						// Ignore duplicate
						continue
					} else {
						// Get missing headers
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

	// The goroutine to subscribe new header and send header to check channel.
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
