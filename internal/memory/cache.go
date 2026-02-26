package memory

import (
	"container/list"
	"sync"
	"time"
)

// LRUCache LRU 缓存实现
type LRUCache struct {
	maxSize      int
	maxAge       time.Duration
	mu           sync.RWMutex
	cache        map[string]*list.Element
	lruList      *list.List
	accessTimes  map[string]time.Time
}

// cacheEntry 缓存条目
type cacheEntry struct {
	key      string
	value    string
	fileSize int64
}

// NewLRUCache 创建新的 LRU 缓存
func NewLRUCache(maxSize int, maxAge time.Duration) *LRUCache {
	lc := &LRUCache{
		maxSize:     maxSize,
		maxAge:      maxAge,
		cache:       make(map[string]*list.Element),
		lruList:     list.New(),
		accessTimes: make(map[string]time.Time),
	}

	// 启动清理协程
	go lc.cleanup()

	return lc
}

// Get 获取缓存值
func (lc *LRUCache) Get(key string) (string, bool) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	elem, ok := lc.cache[key]
	if !ok {
		return "", false
	}

	// 检查是否过期
	if lc.maxAge > 0 {
		if accessTime, exists := lc.accessTimes[key]; exists {
			if time.Since(accessTime) > lc.maxAge {
				// 过期，删除
				lc.removeElement(elem)
				return "", false
			}
		}
	}

	// 更新访问时间
	lc.accessTimes[key] = time.Now()

	// 移动到链表头部（最近使用）
	lc.lruList.MoveToFront(elem)

	entry := elem.Value.(*cacheEntry)
	return entry.value, true
}

// Set 设置缓存值
func (lc *LRUCache) Set(key string, value string, fileSize int64) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	// 如果已存在，更新并移动到头部
	if elem, ok := lc.cache[key]; ok {
		entry := elem.Value.(*cacheEntry)
		entry.value = value
		entry.fileSize = fileSize
		lc.accessTimes[key] = time.Now()
		lc.lruList.MoveToFront(elem)
		return
	}

	// 如果超过最大大小，删除最久未使用的
	for lc.lruList.Len() >= lc.maxSize {
		lc.removeOldest()
	}

	// 添加新条目
	entry := &cacheEntry{
		key:      key,
		value:    value,
		fileSize: fileSize,
	}
	elem := lc.lruList.PushFront(entry)
	lc.cache[key] = elem
	lc.accessTimes[key] = time.Now()
}

// Remove 删除缓存项
func (lc *LRUCache) Remove(key string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if elem, ok := lc.cache[key]; ok {
		lc.removeElement(elem)
	}
}

// removeElement 删除元素
func (lc *LRUCache) removeElement(elem *list.Element) {
	entry := elem.Value.(*cacheEntry)
	delete(lc.cache, entry.key)
	delete(lc.accessTimes, entry.key)
	lc.lruList.Remove(elem)
}

// removeOldest 删除最久未使用的项
func (lc *LRUCache) removeOldest() {
	back := lc.lruList.Back()
	if back != nil {
		lc.removeElement(back)
	}
}

// cleanup 定期清理过期项
func (lc *LRUCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		lc.mu.Lock()
		now := time.Now()
		for key, accessTime := range lc.accessTimes {
			if lc.maxAge > 0 && now.Sub(accessTime) > lc.maxAge {
				if elem, ok := lc.cache[key]; ok {
					lc.removeElement(elem)
				}
			}
		}
		lc.mu.Unlock()
	}
}

// Size 返回当前缓存大小
func (lc *LRUCache) Size() int {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	return lc.lruList.Len()
}

// Clear 清空缓存
func (lc *LRUCache) Clear() {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.cache = make(map[string]*list.Element)
	lc.lruList = list.New()
	lc.accessTimes = make(map[string]time.Time)
}
