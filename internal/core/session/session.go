package session

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os"
	"path"
	"sync"
	"time"

	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/models/auth"
	"github.com/bestruirui/bestsub/internal/utils"
	"github.com/bestruirui/bestsub/internal/utils/local"
)

const MaxSessions = 10

var (
	ErrSessionPoolFull  = errors.New("session pool is full")
	ErrInvalidSessionID = errors.New("invalid session ID")
	ErrSessionNotFound  = errors.New("session not found")
	sessions            [MaxSessions]auth.Session
	mu                  sync.RWMutex
)

// Load 从文件加载会话信息
func init() {
	mu.Lock()
	defer mu.Unlock()
	sessionFile := config.Get().Session.File
	if _, err := os.Stat(sessionFile); os.IsNotExist(err) {
		return
	}

	data, err := os.ReadFile(sessionFile)
	if err != nil {
		return
	}

	if len(data) == 0 {
		return
	}

	var loadedSessions [MaxSessions]auth.Session
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&loadedSessions); err != nil {
		return
	}

	now := uint32(local.Time().Unix())
	for i := range loadedSessions {
		if loadedSessions[i].IsActive && now > loadedSessions[i].ExpiresAt {
			loadedSessions[i].IsActive = false
		}
		sessions[i] = loadedSessions[i]
	}
}

// Close 关闭会话管理器，将会话信息保存到文件
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(sessions); err != nil {
		return err
	}
	sessionFile := config.Get().Session.File

	dir := path.Dir(sessionFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	return os.WriteFile(sessionFile, buf.Bytes(), 0600)
}

func GetUnActive() (uint8, *auth.Session) {
	mu.Lock()
	defer mu.Unlock()
	for i := range sessions {
		if !sessions[i].IsActive {
			return uint8(i), &sessions[i]
		}
	}
	return 0, nil
}

// Get 获取会话
func Get(sessionID uint8) (*auth.Session, error) {
	if sessionID >= 10 {
		return nil, ErrInvalidSessionID
	}

	mu.Lock()
	defer mu.Unlock()

	session := &sessions[sessionID]
	if !session.IsActive {
		return nil, ErrSessionNotFound
	}

	// 检查是否过期
	now := uint32(local.Time().Unix())
	if now > session.ExpiresAt {
		session.IsActive = false
		return nil, ErrSessionNotFound
	}

	return session, nil
}

// Disable 禁用会话
func Disable(sessionID uint8) error {
	if sessionID >= 10 {
		return ErrInvalidSessionID
	}

	mu.Lock()
	defer mu.Unlock()

	sessions[sessionID].IsActive = false
	return nil
}

// DisableAll 禁用所有会话
func DisableAll() {
	mu.Lock()
	defer mu.Unlock()

	for i := range sessions {
		sessions[i].IsActive = false
	}
}

// GetActive 获取所有活跃会话
func GetActive() *[]auth.SessionResponse {
	mu.RLock()
	defer mu.RUnlock()

	var active []auth.SessionResponse
	now := uint32(local.Time().Unix())

	for i := range sessions {
		if sessions[i].IsActive && now <= sessions[i].ExpiresAt {
			active = append(active, auth.SessionResponse{
				ID:           uint8(i),
				IsActive:     sessions[i].IsActive,
				ClientIP:     utils.Uint32ToIP(sessions[i].ClientIP),
				UserAgent:    sessions[i].UserAgent,
				ExpiresAt:    time.Unix(int64(sessions[i].ExpiresAt), 0),
				CreatedAt:    time.Unix(int64(sessions[i].CreatedAt), 0),
				LastAccessAt: time.Unix(int64(sessions[i].LastAccessAt), 0),
			})
		}
	}

	return &active
}

// Cleanup 清理过期会话，返回清理数量
func Cleanup() int {
	mu.Lock()
	defer mu.Unlock()

	cleaned := 0
	now := uint32(local.Time().Unix())

	for i := range sessions {
		if sessions[i].IsActive && now > sessions[i].ExpiresAt {
			sessions[i].IsActive = false
			cleaned++
		}
	}

	return cleaned
}

// Stats 获取会话统计信息
func Stats() (total, active, expired int) {
	mu.RLock()
	defer mu.RUnlock()

	now := uint32(local.Time().Unix())

	for i := range sessions {
		if sessions[i].IsActive {
			total++
			if now > sessions[i].ExpiresAt {
				expired++
			} else {
				active++
			}
		}
	}

	return total, active, expired
}
