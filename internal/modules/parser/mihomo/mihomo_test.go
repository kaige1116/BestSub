package mihomo

import (
	"fmt"
	"testing"

	"github.com/bestruirui/bestsub/internal/core/nodepool"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

const content = `
proxies:
    - name: TEST_TROJAN
      password: 12345678-12345678-12345678-12345678
      port: 15229
      server: 1.1.1.1
      skip-cert-verify: true
      sni: test.com
      type: trojan
    - alterId: 2
      cipher: auto
      name: TEST_VMESS
      network: ws
      port: 30805
      server: 1.1.1.1
      skip-cert-verify: false
      tls: false
      type: vmess
      uuid: 12345678-12345678-12345678-123456780
      ws-opts:
        headers:
            Host: test.com
        path: /test
    - client-fingerprint: firefox
      name: "TEST_VLESS"
      network: tcp
      port: 34045
      reality-opts:
        public-key: 12345678-12345678-12345678-12345678
        short-id: 12345678
      server: 1.1.1.1
      servername: test.com
      tls: true
      type: vless
      uuid: 12345678-12345678-12345678-12345678
      alpn:
        - h3%2Ch2%2Chttp%2F1.1
      skip-cert-verify: false
      tfo: false
      ws-opts:
        headers:
            Host: test.com
        path: /test
    - {name: TEST_SS, server: test.com, port: 33476, type: ss, cipher: chacha20-ietf-poly1305, password: 12345678-12345678-12345678-12345678, udp: true}
    - {name: TEST_SS_Du, server: test.com, port: 33476, type: ss, cipher: chacha20-ietf-poly1305, password: 12345678-12345678-12345678-12345678, udp: true}
    - name: TEST_VMESS
      type: vmess
      server: 1.1.1.1
      port: 2086
      uuid: 7d92ffc9-02e1-4087-8a46-cc4d76560917
      alterId: 0
      cipher: auto
      udp: true
      tls: false
      network: ws
      servername: test.com
      ws-opts:
        path: /test
        headers:
          Host: test.com
    - name: TEST_VMESS
      type: vmess
      server: 1.1.1.1
      port: 2086
      uuid: 7d92ffc9-02e1-4087-8a46-cc4d76560917
      alterId: 0
      cipher: auto
      udp: true
      tls: false
      network: ws
      servername: test.com
      ws-opts:
        path: /test
        headers:
          Host: test.com
`

func TestParseVLESS(t *testing.T) {
	if err := log.Initialize("debug", "console", ""); err != nil {
		panic(err)
	}
	// 测试第二个节点（带alpn数组）
	content := []byte(content)
	nodeCount, err := Parse(&content, 1)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	fmt.Printf("Node count: %d\n", nodeCount)
	for i := 0; i < nodeCount; i++ {
		nodeData := nodepool.GetNextNode(1)
		if len(nodeData.Config) == 0 {
			t.Fatal("no node")
		}
		fmt.Printf("Node %d config: %s\n", i, string(nodeData.Config))
		fmt.Printf("Node %d info: %+v\n", i, nodeData.Info)
	}
}
