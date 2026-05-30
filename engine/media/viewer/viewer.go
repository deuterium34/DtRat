/*
Пакет для вывода видео и изображений на экран.
Релизация под Windows с установленным Edge.
*/
package viewer

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Viewer управляет запущенным процессом полноэкранного просмотра.
type Viewer struct {
	cmd      *exec.Cmd
	tempPath string
}

// tpl — минималистичный HTML-шаблон для отображения медиа без интерфейса.
const tpl = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <style>
        body {
            margin: 0;
            padding: 0;
            background-color: black;
            overflow: hidden; /* Убирает скроллбары */
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            width: 100vw;
            cursor: none; /* Прячет курсор мыши */
        }
        video, img {
            max-width: 100%;
            max-height: 100%;
            object-fit: contain; /* Сохраняет пропорции без обрезки */
        }
    </style>
</head>
<body>
    {{if .IsVideo}}
    <video src="{{.URL}}" autoplay loop></video>
    {{else}}
    <img src="{{.URL}}">
    {{end}}
</body>
</html>`

type templateData struct {
	URL     template.URL
	IsVideo bool
}

// Show выводит файл (изображение или видео) на экран в полноэкранном режиме.
// Возвращает объект Viewer, с помощью которого можно остановить показ.
func Show(filePath string) (*Viewer, error) {
	// 1. Проверки файла
	if filePath == "" {
		return nil, errors.New("путь к файлу не может быть пустым")
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения абсолютного пути: %w", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("файл не найден: %s", absPath)
	}

	// 2. Определение типа файла
	ext := strings.ToLower(filepath.Ext(absPath))
	isVideo := false
	switch ext {
	case ".mp4", ".webm", ".ogg", ".mov":
		isVideo = true
	case ".jpg", ".jpeg", ".png", ".webp", ".gif", ".bmp":
		isVideo = false
	default:
		return nil, fmt.Errorf("неподдерживаемый формат файла: %s", ext)
	}

	// 3. Формирование безопасного file:// URL
	fileURL := "file:///" + filepath.ToSlash(absPath)

	// 4. Создание временного HTML-файла-обертки
	tempFile, err := os.CreateTemp("", "media-*.html")
	if err != nil {
		return nil, fmt.Errorf("ошибка создания временного файла: %w", err)
	}
	tempFileName := tempFile.Name()

	t, err := template.New("media").Parse(tpl)
	if err != nil {
		tempFile.Close()
		os.Remove(tempFileName)
		return nil, fmt.Errorf("ошибка парсинга шаблона: %w", err)
	}

	err = t.Execute(tempFile, templateData{
		URL:     template.URL(fileURL),
		IsVideo: isVideo,
	})
	tempFile.Close()
	if err != nil {
		os.Remove(tempFileName)
		return nil, fmt.Errorf("ошибка записи HTML шаблона: %w", err)
	}

	// 5. Поиск Microsoft Edge в стандартных директориях Windows
	edgePaths := []string{
		`C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`,
		`C:\Program Files\Microsoft\Edge\Application\msedge.exe`,
	}

	var edgeExe string
	for _, p := range edgePaths {
		if _, err := os.Stat(p); err == nil {
			edgeExe = p
			break
		}
	}

	if edgeExe == "" {
		os.Remove(tempFileName)
		return nil, errors.New("браузер Microsoft Edge не найден (необходим для работы)")
	}

	// 6. Запуск процесса в режиме киоска
	htmlURL := "file:///" + filepath.ToSlash(tempFileName)

	// Флаг --kiosk делает окно полноэкранным без возможности выхода стандартными средствами
	// Флаг --autoplay-policy=no-user-gesture-required позволяет видео со звуком играть сразу
	cmd := exec.Command(edgeExe,
		"--kiosk", htmlURL,
		"--edge-kiosk-type=fullscreen",
		"--autoplay-policy=no-user-gesture-required",
		"--disable-infobars",
	)

	if err := cmd.Start(); err != nil {
		os.Remove(tempFileName)
		return nil, fmt.Errorf("ошибка запуска плеера: %w", err)
	}

	return &Viewer{
		cmd:      cmd,
		tempPath: tempFileName,
	}, nil
}

// Close закрывает полноэкранное окно и очищает временные файлы.
func (v *Viewer) Close() error {
	var errs []error

	if v.cmd != nil && v.cmd.Process != nil {
		if err := v.cmd.Process.Kill(); err != nil {
			errs = append(errs, fmt.Errorf("ошибка завершения процесса: %w", err))
		}
		// Дожидаемся фактического завершения, чтобы не было зомби-процессов
		_ = v.cmd.Wait()
	}

	if v.tempPath != "" {
		if err := os.Remove(v.tempPath); err != nil && !os.IsNotExist(err) {
			errs = append(errs, fmt.Errorf("ошибка удаления временного файла: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("ошибки при закрытии: %v", errs)
	}
	return nil
}
