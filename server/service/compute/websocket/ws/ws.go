package ws

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CDUwenbojin/websocket"
	"io"
	"sync"
	"time"
	"ygpt/server/service/compute/websocket/protocol"
)

type WsClient struct {
	WsCli   *websocket.Client
	Counter Counter
}

type Counter struct {
	mu sync.Mutex
	ID int
}

func (c *Counter) GetID() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ID++
	if c.ID > 50000 {
		c.ID = 0
	}
	return c.ID
}

func (wc *WsClient) Connect() error {
	return wc.WsCli.Connect()
}

var (
	unconfirmedIDs = make(map[int]struct{}) // 用 map 作为集合，存储未确认的消息 ID
	mu             sync.Mutex               // 保护共享数据的锁
	MaxRetries     = 2
)

type CommonMsg struct {
	protocol.BaseMsg
}

func (wc *WsClient) SendMessageDebug(message interface{}) error {
	fmt.Println("发送消息:", message)
	return wc.WsCli.SendMessage(message)
}

// SendingWithRetry 重发逻辑
func (wc *WsClient) SendingWithRetry(message interface{}) error {
	var msg CommonMsg
	bytes, err := json.Marshal(message)
	err = json.Unmarshal(bytes, &msg)
	if err != nil {
		return err
	}
	msgID := msg.ID
	mu.Lock()
	unconfirmedIDs[msgID] = struct{}{} // 将消息 ID 添加到集合中
	mu.Unlock()

	err = wc.WsCli.SendMessage(message)
	if err != nil {
		return err
	}
	fmt.Printf(fmt.Sprintf("发送消息:%+v", string(bytes)))
	go func() {
		retries := 0
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			mu.Lock()
			if _, exists := unconfirmedIDs[msgID]; !exists || retries >= MaxRetries {
				mu.Unlock()
				return // 消息已确认或达到最大重发次数，退出
			}
			fmt.Println("重发消息:", msgID)
			err := wc.WsCli.SendMessage(message)
			if err != nil {
				mu.Unlock()
				return
			}
			retries++
			mu.Unlock()
		}
	}()
	return nil
}

// 当收到对应的应答报文时，删除对应的消息 ID
func (wc *WsClient) HandleAck(msgID int) {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := unconfirmedIDs[msgID]; exists {
		delete(unconfirmedIDs, msgID) // 从集合中删除已确认的消息 ID
	}
}

// 加密函数（AES-CBC）
func encryptAESCBC(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	copy(ciphertext[:aes.BlockSize], iv)

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

// 解密函数（AES-CBC）
func decryptAESCBC(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return ciphertext, nil
}
