package before

import (
	"database/sql"
	"sync"
	"time"
)

type BatchInserter struct {
	db            *sql.DB
	buffer        []any      // 待插入的数据缓冲区
	mu            sync.Mutex // 保护缓冲区的互斥锁
	batchSize     int
	flushInterval time.Duration
	ticker        *time.Ticker
}

func NewBatchInserter(db *sql.DB, batchSize int, flushInterval time.Duration) *BatchInserter {
	inserter := &BatchInserter{
		db:            db,
		buffer:        make([]any, 0, batchSize),
		batchSize:     batchSize,
		flushInterval: flushInterval,
		ticker:        time.NewTicker(flushInterval),
	}

	// 启动一个goroutine定时刷新缓冲区
	// Start a goroutine to periodically flush the buffer.
	go inserter.flushLoop()

	return inserter
}

// 关闭批量插入器，停止定时器并刷新缓冲区
// Close the batch inserter, stop the ticker, and flush the buffer.
func (b *BatchInserter) Close() {
	b.ticker.Stop()
	// 刷新缓冲区
	// Flush the buffer.
	b.mu.Lock()
	if len(b.buffer) > 0 {
		bufferCopy := make([]any, len(b.buffer))
		copy(bufferCopy, b.buffer)
		b.buffer = b.buffer[:0]
		b.mu.Unlock()
		go b.insertBatch(bufferCopy)
	} else {
		b.mu.Unlock()
	}
}

func (b *BatchInserter) Insert(data any) {
	// 使用互斥锁保护缓冲区，避免并发写入导致数据丢失
	// Lock the mutex to ensure exclusive access to the buffer.
	b.mu.Lock()
	b.buffer = append(b.buffer, data)
	// 如果达到批量大小，立即刷新缓冲区
	// If the buffer size exceeds the batch size, flush the buffer.
	if len(b.buffer) >= b.batchSize {
		// 注意：这里先复制缓冲区，然后清空，再解锁，最后执行插入
		// 这样可以减少锁的持有时间，避免在插入数据库时阻塞其他goroutine写入缓冲区
		bufferCopy := make([]any, len(b.buffer))
		copy(bufferCopy, b.buffer)
		b.buffer = b.buffer[:0] // 清空缓冲区
		b.mu.Unlock()
		// 异步插入，避免阻塞
		go b.insertBatch(bufferCopy)
	} else {
		b.mu.Unlock()
	}
}

// 定时刷新循环
func (b *BatchInserter) flushLoop() {
	for range b.ticker.C {
		// 检查缓冲区是否为空
		// Check if the buffer is empty.
		// 使用互斥锁保护缓冲区，避免数据竞争
		// Lock the mutex to ensure exclusive access to the buffer.
		b.mu.Lock()
		if len(b.buffer) == 0 {
			b.mu.Unlock()
			continue
		}
		bufferCopy := make([]any, len(b.buffer))
		copy(bufferCopy, b.buffer)
		b.buffer = b.buffer[:0]
		b.mu.Unlock()

		go b.insertBatch(bufferCopy)
	}
}

func (b *BatchInserter) insertBatch(batch []interface{}) {
	// 模拟批量插入数据库
	// Simulate batch insert into the database.
	_, err := b.db.Exec("INSERT ...", batch)
	if err != nil {
		// 错误处理
	}
}
