// Gemini

package tcp

import (
	"context"
	"dtrat/transport"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Константы для типов пакетов
const (
	pktTypeMsg  byte = 0x01
	pktTypeFile byte = 0x02
)

// Ограничение на размер текстового сообщения (защита от OOM-атак)
const maxMessageSize = 16 * 1024 * 1024 // 16 MB

// tcpTransport — потокобезопасная реализация интерфейса Transport
type tcpTransport struct {
	addr         string
	isServer     bool
	readTimeout  time.Duration
	writeTimeout time.Duration

	listener net.Listener
	conn     net.Conn

	msgChan  chan string
	fileChan chan string

	isStarted bool
	startMu   sync.Mutex

	writeMu   sync.Mutex
	closeOnce sync.Once
	closed    chan struct{}
}

// timeoutConn оборачивает net.Conn для автоматического обновления дедлайнов при каждом I/O событии
type timeoutConn struct {
	net.Conn
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func (c *timeoutConn) Read(b []byte) (int, error) {
	if c.readTimeout > 0 {
		if err := c.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			return 0, err
		}
	}
	return c.Conn.Read(b)
}

func (c *timeoutConn) Write(b []byte) (int, error) {
	if c.writeTimeout > 0 {
		if err := c.SetWriteDeadline(time.Now().Add(c.writeTimeout)); err != nil {
			return 0, err
		}
	}
	return c.Conn.Write(b)
}

// NewTCPClient создает транспорт для клиентской стороны
func NewTCPClient(addr string, readTimeout, writeTimeout time.Duration) transport.Transport {
	return &tcpTransport{
		addr:         addr,
		isServer:     false,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		msgChan:      make(chan string, 100),
		fileChan:     make(chan string, 100),
		closed:       make(chan struct{}),
	}
}

// NewTCPServer создает транспорт для серверной стороны
func NewTCPServer(addr string, readTimeout, writeTimeout time.Duration) transport.Transport {
	return &tcpTransport{
		addr:         addr,
		isServer:     true,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		msgChan:      make(chan string, 100),
		fileChan:     make(chan string, 100),
		closed:       make(chan struct{}),
	}
}

// Start инициализирует соединение (слушает порт или подключается)
func (t *tcpTransport) Start() error {
	t.startMu.Lock()
	defer t.startMu.Unlock()

	if t.isStarted {
		return errors.New("transport already started")
	}

	if t.isServer {
		// Вместо обычного net.Listen используем ListenConfig с включенным KeepAlive
		lc := net.ListenConfig{
			KeepAlive: 15 * time.Second, // ОС будет слать пинги каждые 15 сек, если в канале тишина
		}
		ln, err := lc.Listen(context.Background(), "tcp", t.addr)
		if err != nil {
			return fmt.Errorf("failed to start listener: %w", err)
		}
		t.listener = ln

		rawConn, err := ln.Accept()
		if err != nil {
			t.listener.Close()
			return fmt.Errorf("failed to accept connection: %w", err)
		}

		// Активируем KeepAlive на самом соединении (на всякий случай)
		if tcpConn, ok := rawConn.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(15 * time.Second)
		}

		t.conn = &timeoutConn{Conn: rawConn, readTimeout: t.readTimeout, writeTimeout: t.writeTimeout}
	} else {
		// ... для Клиента:
		dialer := &net.Dialer{
			KeepAlive: 15 * time.Second, // Включает проверку связи со стороны клиента
			Timeout:   10 * time.Second, // Таймаут на само подключение
		}
		rawConn, err := dialer.Dial("tcp", t.addr)
		if err != nil {
			return fmt.Errorf("failed to connect to server: %w", err)
		}

		if tcpConn, ok := rawConn.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(15 * time.Second)
		}

		t.conn = &timeoutConn{Conn: rawConn, readTimeout: t.readTimeout, writeTimeout: t.writeTimeout}
	}

	t.isStarted = true
	go t.readLoop()
	return nil
}

// Send форматирует строку и отправляет текстовый пакет
func (t *tcpTransport) Send(s string, args ...any) error {
	t.startMu.Lock()
	started := t.isStarted
	t.startMu.Unlock()
	if !started {
		return errors.New("transport not started")
	}

	msg := fmt.Sprintf(s, args...)
	payload := []byte(msg)
	payloadLen := uint64(len(payload))

	t.writeMu.Lock()
	defer t.writeMu.Unlock()

	if t.conn == nil {
		return net.ErrClosed
	}

	header := make([]byte, 9)
	header[0] = pktTypeMsg
	binary.BigEndian.PutUint64(header[1:9], payloadLen)

	if _, err := t.conn.Write(header); err != nil {
		return fmt.Errorf("failed to write message header: %w", err)
	}
	if _, err := t.conn.Write(payload); err != nil {
		return fmt.Errorf("failed to write message body: %w", err)
	}

	return nil
}

// SendFile стримит файл порциями, не загружая его целиком в память
func (t *tcpTransport) SendFile(filePath string) error {
	t.startMu.Lock()
	started := t.isStarted
	t.startMu.Unlock()
	if !started {
		return errors.New("transport not started")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	filename := filepath.Base(filePath)
	fnBytes := []byte(filename)
	fnLen := len(fnBytes)

	if fnLen > 65535 {
		return errors.New("filename is too long")
	}

	// Структура полезной нагрузки: [2 байта длина имени] + [имя файла] + [содержимое]
	payloadLen := uint64(2 + fnLen + int(stat.Size()))

	t.writeMu.Lock()
	defer t.writeMu.Unlock()

	if t.conn == nil {
		return net.ErrClosed
	}

	header := make([]byte, 9)
	header[0] = pktTypeFile
	binary.BigEndian.PutUint64(header[1:9], payloadLen)

	if _, err := t.conn.Write(header); err != nil {
		return fmt.Errorf("failed to write file header: %w", err)
	}

	fnLenBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(fnLenBuf, uint16(fnLen))
	if _, err := t.conn.Write(fnLenBuf); err != nil {
		return fmt.Errorf("failed to write filename length: %w", err)
	}

	if _, err := t.conn.Write(fnBytes); err != nil {
		return fmt.Errorf("failed to write filename: %w", err)
	}

	if _, err := io.Copy(t.conn, file); err != nil {
		return fmt.Errorf("failed to stream file data: %w", err)
	}

	return nil
}

// Wait блокирует горутину до получения текстового сообщения
func (t *tcpTransport) Wait() (string, error) {
	t.startMu.Lock()
	started := t.isStarted
	t.startMu.Unlock()
	if !started {
		return "", errors.New("transport not started")
	}

	select {
	case msg, ok := <-t.msgChan:
		if !ok {
			return "", net.ErrClosed
		}
		return msg, nil
	case <-t.closed:
		return "", net.ErrClosed
	}
}

// WaitFile блокирует горутину до получения файла и возвращает путь к нему
func (t *tcpTransport) WaitFile() (string, error) {
	t.startMu.Lock()
	started := t.isStarted
	t.startMu.Unlock()
	if !started {
		return "", errors.New("transport not started")
	}

	select {
	case path, ok := <-t.fileChan:
		if !ok {
			return "", net.ErrClosed
		}
		return path, nil
	case <-t.closed:
		return "", net.ErrClosed
	}
}

// Close корректно завершает работу, закрывает сокеты и каналы без паник
func (t *tcpTransport) Close() error {
	var err error
	t.closeOnce.Do(func() {
		close(t.closed)

		t.writeMu.Lock()
		if t.conn != nil {
			err = t.conn.Close()
		}
		t.writeMu.Unlock()

		if t.listener != nil {
			if lErr := t.listener.Close(); lErr != nil && err == nil {
				err = lErr
			}
		}
	})
	return err
}

// Внутренний цикл чтения пакетов из сокета
func (t *tcpTransport) readLoop() {
	defer func() {
		t.Close()
		close(t.msgChan)
		close(t.fileChan)
	}()

	headerBuf := make([]byte, 9)
	fnLenBuf := make([]byte, 2)

	for {
		_, err := io.ReadFull(t.conn, headerBuf)
		if err != nil {
			return // Соединение закрыто, таймаут или EOF
		}

		pktType := headerBuf[0]
		payloadLen := binary.BigEndian.Uint64(headerBuf[1:9])

		switch pktType {
		case pktTypeMsg:
			if payloadLen > maxMessageSize {
				return // Защита памяти от аномально больших пакетов
			}
			msgBuf := make([]byte, payloadLen)
			if _, err := io.ReadFull(t.conn, msgBuf); err != nil {
				return
			}
			select {
			case t.msgChan <- string(msgBuf):
			case <-t.closed:
				return
			}

		case pktTypeFile:
			if _, err := io.ReadFull(t.conn, fnLenBuf); err != nil {
				return
			}
			fnLen := binary.BigEndian.Uint16(fnLenBuf)

			fnBuf := make([]byte, fnLen)
			if _, err := io.ReadFull(t.conn, fnBuf); err != nil {
				return
			}
			// Безопасное извлечение имени (Path Traversal Protection)
			filename := filepath.Base(string(fnBuf))

			tmpFile, err := os.CreateTemp("", "transport_*_"+filename)
			if err != nil {
				return
			}

			fileContentSize := int64(payloadLen - 2 - uint64(fnLen))
			// Потоковое сохранение данных на диск без аллокации RAM
			_, err = io.CopyN(tmpFile, t.conn, fileContentSize)
			tmpFile.Close()
			if err != nil {
				os.Remove(tmpFile.Name())
				return
			}

			select {
			case t.fileChan <- tmpFile.Name():
			case <-t.closed:
				os.Remove(tmpFile.Name())
				return
			}

		default:
			return // Неизвестный тип пакета — десинхронизация потока, закрываем сокет
		}
	}
}
