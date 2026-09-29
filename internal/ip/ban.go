package ip

import "sync"

type BanList struct {
	mu  sync.RWMutex
	ips map[string]bool
}

func NewBanList() *BanList {
	return &BanList{
		ips: make(map[string]bool),
	}
}

func (b *BanList) Block(ip string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.ips[ip] = true
}

func (b *BanList) Unblock(ip string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.ips, ip)
}

func (b *BanList) Check(ip string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.ips[ip]
}
