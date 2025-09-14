package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type Logger struct {
	*log.Logger
	level Level
}

func New(config Config) *Logger {
	output := config.Output
	if output == nil {
		output = os.Stdout
	}

	return &Logger{
		Logger: log.New(output, "", 0), 
		level:  config.Level,
	}
}

// Debug логирует сообщение на уровне DEBUG
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= LevelDebug {
		l.log(LevelDebug, format, v...)
	}
}

// Info логирует сообщение на уровне INFO
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= LevelInfo {
		l.log(LevelInfo, format, v...)
	}
}

// Warn логирует сообщение на уровне WARN
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= LevelWarn {
		l.log(LevelWarn, format, v...)
	}
}

// Error логирует сообщение на уровне ERROR
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= LevelError {
		l.log(LevelError, format, v...)
	}
}

// Fatal логирует сообщение на уровне FATAL и завершает программу
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.log(LevelFatal, format, v...)
	os.Exit(1)
}

// log внутренний метод для форматирования и вывода логов
func (l *Logger) log(level Level, format string, v ...interface{}) {
	// Получаем информацию о caller'е
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	} else {
		// Оставляем только имя файла, а не полный путь
		file = filepath.Base(file)
	}

	// Форматируем сообщение
	message := fmt.Sprintf(format, v...)
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	
	// Цвета для терминала (опционально)
	var colorCode string
	switch level {
	case LevelDebug:
		colorCode = "\033[36m" // Cyan
	case LevelInfo:
		colorCode = "\033[32m" // Green
	case LevelWarn:
		colorCode = "\033[33m" // Yellow
	case LevelError:
		colorCode = "\033[31m" // Red
	case LevelFatal:
		colorCode = "\033[35m" // Magenta
	}
	resetCode := "\033[0m"

	// Форматированная строка лога
	logEntry := fmt.Sprintf("%s%s [%s] %s:%d - %s%s", 
		colorCode, timestamp, level, file, line, message, resetCode)

	l.Logger.Println(logEntry)
}

// WithFields создает логгер с дополнительными полями (для будущего расширения)
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	// Для простоты возвращаем тот же логгер
	// В реальном проекте можно реализовать контекстное логирование
	return l
}