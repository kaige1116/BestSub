package jwt

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/bestruirui/bestsub/internal/config"
	timeutil "github.com/bestruirui/bestsub/internal/utils/time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims JWT声明结构
type Claims struct {
	SessionID int64 `json:"session_id"`
	jwt.RegisteredClaims
}

// TokenPair 令牌对
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenHash    string    `json:"token_hash"` // 用于存储到数据库的哈希值
}

// GenerateTokenPair 生成访问令牌和刷新令牌对
func GenerateTokenPair(sessionID int64) (*TokenPair, error) {
	cfg := config.Get()

	// 生成访问令牌
	now := timeutil.Now()
	expiresAt := now.Add(time.Duration(cfg.JWT.ExpiresIn) * time.Second)

	claims := &Claims{
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "bestsub",
			Subject:   fmt.Sprintf("session-%d", sessionID),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// 生成刷新令牌（有效期更长）
	refreshClaims := &Claims{
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)), // 7天
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "bestsub-refresh",
			Subject:   fmt.Sprintf("session-%d", sessionID),
			ID:        uuid.New().String(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	// 生成访问令牌的哈希值用于数据库存储
	tokenHash := HashToken(accessToken)

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		ExpiresAt:    expiresAt,
		TokenHash:    tokenHash,
	}, nil
}

// ValidateToken 验证JWT令牌
func ValidateToken(tokenString string) (*Claims, error) {
	cfg := config.Get()

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ValidateRefreshToken 验证刷新令牌
func ValidateRefreshToken(tokenString string) (*Claims, error) {
	cfg := config.Get()

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// 验证是否为刷新令牌
		if claims.Issuer != "bestsub-refresh" {
			return nil, fmt.Errorf("not a refresh token")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid refresh token")
}

// HashToken 生成令牌的哈希值
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}

// ExtractTokenFromHeader 从Authorization头中提取令牌
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is empty")
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return authHeader[len(bearerPrefix):], nil
}

// IsTokenExpired 检查令牌是否过期
func IsTokenExpired(claims *Claims) bool {
	return time.Now().After(claims.ExpiresAt.Time)
}

// GetTokenRemainingTime 获取令牌剩余有效时间
func GetTokenRemainingTime(claims *Claims) time.Duration {
	return time.Until(claims.ExpiresAt.Time)
}
