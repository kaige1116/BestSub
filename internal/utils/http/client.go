package utilshttp

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/bestruirui/bestsub/internal/config"
	"github.com/bestruirui/bestsub/internal/models/system"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

// HTTP客户端池，使用统一池和共享传输层
type clientPool struct {
	clientPool       sync.Pool // 统一的客户端池
	pooledClientPool sync.Pool // PooledClient对象池
	transportPool    sync.Pool // 代理传输层池

	// 共享的传输层实例（http.Transport是线程安全的）
	sharedTransport *http.Transport
}

// Client 池化的HTTP客户端，包含释放方法
type Client struct {
	*http.Client
	pool *clientPool
}

var globalPool *clientPool

// 连接池配置常量
const (
	maxIdleConns        = 100
	maxIdleConnsPerHost = 20
	maxConnsPerHost     = 50
	idleConnTimeout     = 90 * time.Second
	tlsHandshakeTimeout = 10 * time.Second
	clientTimeout       = 30 * time.Second
)

// init 初始化HTTP客户端连接池
func init() {
	sharedTransport := createTransport()

	globalPool = &clientPool{
		sharedTransport: sharedTransport,
	}

	// 初始化HTTP客户端池
	globalPool.clientPool.New = func() any {
		return &http.Client{
			Timeout:   clientTimeout,
			Transport: sharedTransport, // 所有客户端共享同一个传输层
		}
	}

	// 初始化PooledClient对象池
	globalPool.pooledClientPool.New = func() any {
		return &Client{
			pool: globalPool,
		}
	}

	// 初始化代理传输层池
	globalPool.transportPool.New = func() any {
		return createTransport()
	}

	log.Debug("HTTP客户端连接池初始化完成")
}

// createTransport 创建优化的HTTP传输层，支持连接池
func createTransport() *http.Transport {
	return &http.Transport{
		MaxIdleConns:        maxIdleConns,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
		MaxConnsPerHost:     maxConnsPerHost,
		IdleConnTimeout:     idleConnTimeout,
		TLSHandshakeTimeout: tlsHandshakeTimeout,
		DisableKeepAlives:   false,
		DisableCompression:  false,
	}
}

// Release 将HTTP客户端放回池中
func (pc *Client) Release() {
	if pc.Client != nil {
		// 如果是代理传输层，需要放回传输层池
		if transport, ok := pc.Client.Transport.(*http.Transport); ok && transport != pc.pool.sharedTransport {
			pc.pool.transportPool.Put(transport)
		}

		// 重置传输层为共享传输层
		pc.Client.Transport = pc.pool.sharedTransport
		pc.pool.clientPool.Put(pc.Client)

		// 重置PooledClient并放回池中
		pc.Client = nil
		pc.pool.pooledClientPool.Put(pc)
	}
}

// Direct 从池中获取一个直连的HTTP客户端
func Direct() *Client {
	client := globalPool.clientPool.Get().(*http.Client)
	pooledClient := globalPool.pooledClientPool.Get().(*Client)
	pooledClient.Client = client
	return pooledClient
}

// Proxy 从池中获取一个使用代理的HTTP客户端
func Proxy() *Client {
	// 获取代理配置
	proxyConfig, err := config.Proxy()
	if err != nil {
		log.Errorf("Failed to get proxy config: %v", err)
		// 返回直连客户端作为降级方案
		return Direct()
	}

	// 如果代理未启用，返回直连客户端
	if !proxyConfig.Enable {
		log.Warn("Proxy is not enabled, returning direct client")
		return Direct()
	}

	// 验证代理配置
	if proxyConfig.Host == "" || proxyConfig.Port == 0 {
		log.Warn("Proxy is enabled but host or port is empty, using direct connection")
		return Direct()
	}

	// 构建代理URL
	proxyURL := buildProxyURL(proxyConfig)

	// 从池中获取客户端和PooledClient
	client := globalPool.clientPool.Get().(*http.Client)
	pooledClient := globalPool.pooledClientPool.Get().(*Client)

	// 从池中获取代理传输层并配置
	proxyTransport := globalPool.transportPool.Get().(*http.Transport)
	proxyTransport.Proxy = http.ProxyURL(proxyURL)

	// 更新客户端的传输层
	client.Transport = proxyTransport
	pooledClient.Client = client

	log.Debugf("Using proxy: %s://%s:%d", proxyConfig.Type, proxyConfig.Host, proxyConfig.Port)
	return pooledClient
}

// Transport 从池中获取传输层并允许自定义配置
func Transport(configFunc func(*http.Transport)) *Client {
	if configFunc == nil {
		log.Warn("config function is nil, using direct connection")
		return Direct()
	}

	// 从池中获取客户端、PooledClient和传输层
	client := globalPool.clientPool.Get().(*http.Client)
	pooledClient := globalPool.pooledClientPool.Get().(*Client)
	transport := globalPool.transportPool.Get().(*http.Transport)

	// 让用户配置传输层
	configFunc(transport)

	client.Transport = transport
	pooledClient.Client = client

	return pooledClient
}

// buildProxyURL 构建代理URL
func buildProxyURL(config *system.ProxyConfig) *url.URL {
	// 默认代理类型为http
	proxyType := config.Type
	if proxyType == "" {
		proxyType = "http"
	}

	// 构建代理地址
	proxyAddr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	// 构建完整的代理URL
	var proxyURL *url.URL
	var err error

	if config.Username != "" && config.Password != "" {
		// 带认证的代理
		proxyURL, err = url.Parse(fmt.Sprintf("%s://%s:%s@%s",
			proxyType,
			url.QueryEscape(config.Username),
			url.QueryEscape(config.Password),
			proxyAddr))
	} else {
		// 无认证的代理
		proxyURL, err = url.Parse(fmt.Sprintf("%s://%s", proxyType, proxyAddr))
	}

	if err != nil {
		log.Errorf("failed to parse proxy URL: %v", err)
		return nil
	}

	return proxyURL
}
