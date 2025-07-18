package common

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os"
	"path"
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
	sessions            = make([]auth.Session, MaxSessions)
)

// Load 从文件加载会话信息
func init() {
	sessionFile := config.Base().Session.Path
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

	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&sessions); err != nil {
		return
	}
	cleanup()
}

// CloseSession 关闭会话管理器，将会话信息保存到文件
func CloseSession() error {
	cleanup()

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(sessions); err != nil {
		return err
	}
	sessionFile := config.Base().Session.Path

	dir := path.Dir(sessionFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	return os.WriteFile(sessionFile, buf.Bytes(), 0600)
}

func GetOneSession() (uint8, *auth.Session) {
	var oldestIndex uint8 = 0
	var oldestTime uint32 = 0
	found := false
	cleanup()
	for i := range sessions {
		if !sessions[i].IsActive {
			if sessions[i].CreatedAt == 0 {
				return uint8(i), &sessions[i]
			}

			if !found || sessions[i].CreatedAt < oldestTime {
				oldestIndex = uint8(i)
				oldestTime = sessions[i].CreatedAt
				found = true
			}
		}
	}

	if found {
		return oldestIndex, &sessions[oldestIndex]
	}
	return 0, nil
}

// GetSession 获取会话
func GetSession(sessionID uint8) (*auth.Session, error) {
	cleanup()
	if sessionID >= MaxSessions {
		return nil, ErrInvalidSessionID
	}

	if int(sessionID) >= len(sessions) {
		return nil, ErrSessionNotFound
	}
	session := &sessions[sessionID]

	return session, nil
}

// DisableSession 禁用会话
func DisableSession(sessionID uint8) error {
	cleanup()
	if sessionID >= MaxSessions {
		return ErrInvalidSessionID
	}

	sessions[sessionID].IsActive = false
	return nil
}

// DisableAllSession 禁用所有会话
func DisableAllSession() {
	for i := range sessions {
		sessions[i].IsActive = false
	}
}

// GetAllSession 获取所有会话
func GetAllSession() *[]auth.SessionResponse {
	cleanup()

	var sessionResponse []auth.SessionResponse

	for i := range sessions {
		if sessions[i].CreatedAt == 0 {
			continue
		}
		sessionResponse = append(sessionResponse, auth.SessionResponse{
			ID:           uint8(i),
			IsActive:     sessions[i].IsActive,
			ClientIP:     utils.Uint32ToIP(sessions[i].ClientIP),
			UserAgent:    sessions[i].UserAgent,
			ExpiresAt:    time.Unix(int64(sessions[i].ExpiresAt), 0),
			CreatedAt:    time.Unix(int64(sessions[i].CreatedAt), 0),
			LastAccessAt: time.Unix(int64(sessions[i].LastAccessAt), 0),
		})
	}

	return &sessionResponse
}

// cleanup 清理过期会话，返回清理数量
func cleanup() int {

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
