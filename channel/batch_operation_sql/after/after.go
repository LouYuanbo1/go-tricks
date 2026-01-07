package after

import (
	"database/sql"
	"log"
	"time"
)

type ChannelBatchInserter struct {
	db *sql.DB
	// 使用 chan 代替缓冲区，避免显式锁的使用(chan中内置了锁,是并发安全的)
	// Use a channel (chan) instead of a buffer to avoid the use of explicit locks
	// (chan has a built-in lock, making it concurrency-safe).
	dataChan      chan any
	batchSize     int
	flushInterval time.Duration
	ticker        *time.Ticker
}

func NewChannelBatchInserter(db *sql.DB, batchSize int, flushInterval time.Duration) *ChannelBatchInserter {
	inserter := &ChannelBatchInserter{
		db:            db,
		dataChan:      make(chan any, batchSize), // 缓冲通道
		batchSize:     batchSize,
		flushInterval: flushInterval,
		ticker:        time.NewTicker(flushInterval),
	}

	// 单个消费者 goroutine，避免并发问题
	// Start a single consumer goroutine to process the batch.
	go inserter.processBatch()

	return inserter
}

func (c *ChannelBatchInserter) Close() {
	// 关闭通道，触发 processBatch 中的 for 循环退出
	// Close the channel to trigger the exit of the for loop in processBatch.
	close(c.dataChan)
	c.ticker.Stop()
}

func (c *ChannelBatchInserter) Insert(data any) {
	// 并发安全：多个 goroutine 可以同时发送
	// Concurrent-safe: multiple goroutines can send data concurrently.
	select {
	case c.dataChan <- data:
	default:
		log.Println("Channel full, data dropped")
	}
}

func (c *ChannelBatchInserter) processBatch() {
	batch := make([]any, 0, c.batchSize)

	for {
		// 使用 select 语句监听通道和定时器
		// Use select statement to listen to the channel and the ticker.
		select {
		case data := <-c.dataChan:
			batch = append(batch, data)
			if len(batch) >= c.batchSize {
				c.flush(batch)
				batch = batch[:0] // 清空
			}

		case <-c.ticker.C:
			if len(batch) > 0 {
				c.flush(batch)
				batch = batch[:0]
			}
		}
	}
}

func (c *ChannelBatchInserter) flush(batch []any) {
	// 模拟批量插入数据库
	// Simulate batch insert into the database.
	_, err := c.db.Exec("INSERT ...", batch)
	if err != nil {
		// 错误处理
	}
}
