package ws

func (h *Hub) isNicknameTaken(nickname string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for conn := range h.clients {
		if conn.nickname == nickname {
			return true
		}
	}

	return false
}
