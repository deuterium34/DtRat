package transport

import (
	"dtrat/config"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Tgbot struct {
	bot     *tgbotapi.BotAPI
	updates *tgbotapi.UpdatesChannel
	userID  int64

	cfg    config.Config
	sendMu sync.Mutex
	waitMu sync.Mutex
}

func NewTgBot(cfg config.Config) (Transport, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Transport.Telegram.Token)
	if err != nil {
		return nil, fmt.Errorf("tgbotapi.NewBotAPI: %w", err)
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	upd, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		return nil, fmt.Errorf("GetUpdatesChan: %w", err)
	}

	b := &Tgbot{
		bot:     bot,
		updates: &upd,
		userID:  int64(cfg.Transport.Telegram.UserID),
		cfg:     cfg,
	}

	return b, nil
}

func (b *Tgbot) Close() error {
	return b.Send("Соединение закрыто.")
}

func (b *Tgbot) Start() error {
	return nil
}

func (b *Tgbot) Send(s string, args ...any) error {
	b.sendMu.Lock()
	defer b.sendMu.Unlock()

	msg := tgbotapi.NewMessage(b.userID, fmt.Sprintf(s, args...))
	_, err := b.bot.Send(msg)
	return err
}

func (b *Tgbot) SendFile(file string) error {
	b.sendMu.Lock()
	defer b.sendMu.Unlock()

	document := tgbotapi.NewDocumentUpload(b.userID, file)
	_, err := b.bot.Send(document)
	return err
}

func (b *Tgbot) Wait() (string, error) {
	b.waitMu.Lock()
	defer b.waitMu.Unlock()

	for update := range *b.updates {
		if update.Message == nil || update.Message.From.ID != int(b.userID) {
			continue
		}

		return update.Message.Text, nil
	}

	return "", nil
}

func (b *Tgbot) WaitFile() (string, error) {
	b.waitMu.Lock()
	defer b.waitMu.Unlock()

	for update := range *b.updates {
		if update.Message == nil || update.Message.From.ID != int(b.userID) {
			continue
		}

		var fileID string
		var fileName string

		if update.Message.Document != nil {
			fileID = update.Message.Document.FileID
			fileName = update.Message.Document.FileName
		} else if update.Message.Photo != nil && len(*update.Message.Photo) > 0 {
			photos := *update.Message.Photo
			photo := photos[len(photos)-1]
			fileID = photo.FileID
			fileName = fmt.Sprintf("photo_%s.jpg", fileID)
		}

		if fileID == "" {
			continue
		}

		fileURL, err := b.bot.GetFileDirectURL(fileID)
		if err != nil {
			return "", fmt.Errorf("ошибка получения URL файла: %w", err)
		}

		downloadDir := "downloads"
		if err := os.MkdirAll(downloadDir, os.ModePerm); err != nil {
			return "", fmt.Errorf("ошибка создания папки: %w", err)
		}

		localPath := filepath.Join(downloadDir, fileName)

		if err := b.downloadFile(fileURL, localPath); err != nil {
			return "", err
		}

		return localPath, nil
	}

	return "", nil
}

func (b *Tgbot) downloadFile(url string, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
