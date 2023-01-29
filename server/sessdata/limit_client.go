package sessdata

import (
	"sync"

	"github.com/bjdgyc/anylink/base"
)

const limitAllKey = "__ALL__"

var (
	limitClient = map[string]int{limitAllKey: 0}
	limitMux    = sync.Mutex{}
)

func LimitClient(user string, close bool) bool {
	limitMux.Lock()
	defer limitMux.Unlock()

	if close {
		releaseClient(user)
		return true
	} else if !CheckLimit(user) {
		return false
	}
	allocClient(user)
	return true
}

// 获取当前用户客户端限额情况
func getCurrUsed(user string) (all, curr int) {
	all = limitClient[limitAllKey]
	curr, ok := limitClient[user]
	if !ok { // 不存在用户
		limitClient[user], curr = 0, 0
	}
	return
}

func CheckLimit(user string) bool {
	_all, c := getCurrUsed(user)
	// 全局判断
	if _all >= base.Cfg.MaxClient {
		return false
	}

	// 超出同一个用户限制
	if c >= base.Cfg.MaxUserClient {
		return false
	}
	return true
}



// 释放客户端
func releaseClient(user string) {
	_all, c := getCurrUsed(user)
	if c == 0 {
		return
	}
	limitClient[user] = c - 1
	limitClient[limitAllKey] = _all - 1
}

// 分配客户端
func allocClient(user string) {
	_all, c := getCurrUsed(user)
	limitClient[user] = c + 1
	limitClient[limitAllKey] = _all + 1
}