package scheduler

import (
	"context"
	"log"
	"time"
)

// SimpleCron 简单的定时任务调度器（不依赖外部库）
type SimpleCron struct {
	tickers []*time.Ticker
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewSimpleCron 创建简单的定时器
func NewSimpleCron(ctx context.Context) *SimpleCron {
	ctx, cancel := context.WithCancel(ctx)
	return &SimpleCron{
		tickers: []*time.Ticker{},
		ctx:     ctx,
		cancel:  cancel,
	}
}

// AddDailyTask 添加每日任务（在指定时间执行）
func (sc *SimpleCron) AddDailyTask(hour, minute int, task func()) error {
	now := time.Now()
	target := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
	
	// 如果目标时间已过，设置为明天
	if target.Before(now) {
		target = target.Add(24 * time.Hour)
	}

	// 计算首次延迟
	delay := target.Sub(now)

	go func() {
		// 等待首次执行时间
		select {
		case <-time.After(delay):
			task()
		case <-sc.ctx.Done():
			return
		}

		// 之后每24小时执行一次
		ticker := time.NewTicker(24 * time.Hour)
		sc.tickers = append(sc.tickers, ticker)

		for {
			select {
			case <-ticker.C:
				task()
			case <-sc.ctx.Done():
				return
			}
		}
	}()

	return nil
}

// Stop 停止所有定时任务
func (sc *SimpleCron) Stop() {
	sc.cancel()
	for _, ticker := range sc.tickers {
		ticker.Stop()
	}
	log.Println("SimpleCron stopped")
}
