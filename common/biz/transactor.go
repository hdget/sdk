package biz

import "sync"

type transactor interface {
	Ref(tx any)      // 引用, 因为可能嵌套NewTransactor, 这里加入计数
	Unref()          // 解引用, 如果计数为0，表示需要Finalize，并Destroy tx
	GetTx() any      // 获取Tx
	ReachRoot() bool // 是否到达最外层
}

type safeTxImpl struct {
	mu    sync.RWMutex
	tx    any
	count int
}

func newTransactor() transactor {
	return &safeTxImpl{
		mu: sync.RWMutex{},
	}
}

func (t *safeTxImpl) Ref(tx any) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.tx == nil {
		t.tx = tx
	}
	t.count++
}

func (t *safeTxImpl) Unref() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.count > 0 {
		t.count--
	}

	// 如果计数为0后销毁tx
	if t.count == 0 {
		t.tx = nil
	}
}

func (t *safeTxImpl) ReachRoot() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	// 如果为1则表示到达最外层，需要Commit Or Rollback
	return t.count == 1
}

func (t *safeTxImpl) GetTx() any {
	t.mu.RLock() // 使用读锁，允许多个 Get 调用并发执行
	defer t.mu.RUnlock()
	return t.tx
}
