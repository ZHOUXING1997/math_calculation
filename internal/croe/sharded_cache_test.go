package croe

import (
	"container/list"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/ZHOUXING1997/math_calculation/internal/math_node"
)

func TestShardedLRUCache(t *testing.T) {
	// 创建分片缓存
	cache := NewShardedLRUCache(16)

	// 创建测试节点
	node1 := &math_node.NumberNode{Value: decimal.NewFromInt(1)}
	node2 := &math_node.NumberNode{Value: decimal.NewFromInt(2)}

	// 测试设置和获取
	cache.Set("key1", node1)

	got, ok := cache.Get("key1")
	if !ok {
		t.Errorf("ShardedLRUCache.Get() ok = %v, want %v", ok, true)
	}

	if got.(*math_node.NumberNode).Value.IntPart() != 1 {
		t.Errorf("ShardedLRUCache.Get() = %v, want %v", got.(*math_node.NumberNode).Value.IntPart(), 1)
	}

	// 测试不存在的键
	_, ok = cache.Get("key2")
	if ok {
		t.Errorf("ShardedLRUCache.Get() ok = %v, want %v", ok, false)
	}

	// 测试多个键
	cache.Set("key2", node2)

	got, ok = cache.Get("key2")
	if !ok {
		t.Errorf("ShardedLRUCache.Get() ok = %v, want %v", ok, true)
	}

	if got.(*math_node.NumberNode).Value.IntPart() != 2 {
		t.Errorf("ShardedLRUCache.Get() = %v, want %v", got.(*math_node.NumberNode).Value.IntPart(), 2)
	}

	// 测试更新
	node1Updated := &math_node.NumberNode{Value: decimal.NewFromInt(10)}
	cache.Set("key1", node1Updated)

	got, ok = cache.Get("key1")
	if !ok {
		t.Errorf("ShardedLRUCache.Get() ok = %v, want %v", ok, true)
	}

	if got.(*math_node.NumberNode).Value.IntPart() != 10 {
		t.Errorf("ShardedLRUCache.Get() = %v, want %v", got.(*math_node.NumberNode).Value.IntPart(), 10)
	}
}

func TestShardedLRUCache_RemoveOldest(t *testing.T) {
	// 直接测试内部分片的LRU功能
	// 创建一个容量为1的分片
	shard := &lruCacheShard{
		capacity: 1,
		cache:    make(map[string]math_node.Node),
		lru:      list.New(),
		items:    make(map[string]*list.Element),
	}

	// 创建测试节点
	node1 := &math_node.NumberNode{Value: decimal.NewFromInt(1)}
	node2 := &math_node.NumberNode{Value: decimal.NewFromInt(2)}

	// 添加第一个项
	element1 := shard.lru.PushFront(lruCacheItem{key: "key1", node: node1})
	shard.items["key1"] = element1
	shard.cache["key1"] = node1

	// 创建分片缓存
	cache := &ShardedLRUCache{}

	// 测试removeOldest方法
	cache.removeOldest(shard)

	// 验证分片为空
	if shard.lru.Len() != 0 || len(shard.items) != 0 || len(shard.cache) != 0 {
		t.Errorf("removeOldest() failed, shard not empty: lru.Len()=%d, items.len=%d, cache.len=%d",
			shard.lru.Len(), len(shard.items), len(shard.cache))
	}

	// 再次测试，这次添加两个项
	element1 = shard.lru.PushFront(lruCacheItem{key: "key1", node: node1})
	shard.items["key1"] = element1
	shard.cache["key1"] = node1

	element2 := shard.lru.PushBack(lruCacheItem{key: "key2", node: node2})
	shard.items["key2"] = element2
	shard.cache["key2"] = node2

	// 移除最旧的项（key2）
	cache.removeOldest(shard)

	// 验证key2已被移除
	_, ok1 := shard.items["key1"]
	_, ok2 := shard.items["key2"]
	if !ok1 || ok2 {
		t.Errorf("removeOldest() failed, expected key1 to exist and key2 to be removed, got ok1=%v, ok2=%v", ok1, ok2)
	}

	// 验证分片中只有一个项
	if shard.lru.Len() != 1 || len(shard.items) != 1 || len(shard.cache) != 1 {
		t.Errorf("removeOldest() failed, shard should have 1 item: lru.Len()=%d, items.len=%d, cache.len=%d",
			shard.lru.Len(), len(shard.items), len(shard.cache))
	}
}

func TestResetShardedCache(t *testing.T) {
	// 保存原始缓存
	originalCache := globalShardedCache

	// 设置测试缓存
	testCache := NewShardedLRUCache(10)
	testCache.Set("testKey", &math_node.NumberNode{Value: decimal.NewFromInt(1)})
	globalShardedCache = testCache

	// 重置缓存
	ResetShardedCache()

	// 检查缓存是否被重置
	_, ok := globalShardedCache.Get("testKey")
	if ok {
		t.Errorf("ResetShardedCache() did not reset the cache")
	}

	// 恢复原始缓存
	globalShardedCache = originalCache
}
