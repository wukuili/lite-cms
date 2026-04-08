package counter

import (
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

// ViewCounter 浏览量批量计数器
// 在内存中累积浏览量，定期批量刷新到数据库，减少DB写入压力
type ViewCounter struct {
	mu      sync.Mutex
	counts  map[uint]int
	db      *gorm.DB
	stopCh  chan struct{}
}

// New 创建浏览量计数器
func New(db *gorm.DB) *ViewCounter {
	vc := &ViewCounter{
		counts: make(map[uint]int),
		db:     db,
		stopCh: make(chan struct{}),
	}
	go vc.flushLoop()
	return vc
}

// Increment 增加某篇文章的浏览量（仅在内存中累积）
func (vc *ViewCounter) Increment(articleID uint) {
	vc.mu.Lock()
	vc.counts[articleID]++
	vc.mu.Unlock()
}

// flushLoop 定期将累积的浏览量刷新到数据库
func (vc *ViewCounter) flushLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			vc.Flush()
		case <-vc.stopCh:
			vc.Flush() // 关闭前最后刷新一次
			return
		}
	}
}

// Flush 将当前累积的浏览量写入数据库
func (vc *ViewCounter) Flush() {
	vc.mu.Lock()
	if len(vc.counts) == 0 {
		vc.mu.Unlock()
		return
	}
	// 取走当前计数，释放锁后再操作DB
	counts := vc.counts
	vc.counts = make(map[uint]int)
	vc.mu.Unlock()

	for id, count := range counts {
		if err := vc.db.Exec(
			"UPDATE articles SET view_count = view_count + ? WHERE id = ?",
			count, id,
		).Error; err != nil {
			log.Printf("批量刷新浏览量失败 id=%d count=%d: %v", id, count, err)
			// 失败的计数放回去，下次再试
			vc.mu.Lock()
			vc.counts[id] += count
			vc.mu.Unlock()
		}
	}
}

// Stop 停止计数器
func (vc *ViewCounter) Stop() {
	close(vc.stopCh)
}
