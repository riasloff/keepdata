package logger

import (
	"keepdata/internal/config"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log = zap.NewNop()

func Setup(c *config.Conf) {

	file, err := os.OpenFile(c.Log_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	// Конфигурируем запись в файл
	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(file), zapcore.DebugLevel) // InfoLevel

	// Конфигурируем запись в консоль
	consoleEncoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()) // NewConsoleEncoder
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)

	// Объединяем оба Core с помощью zapcore.Tee
	core := zapcore.NewTee(fileCore, consoleCore)

	// Создаем логгер на основе объединенного core
	logger := zap.New(core)

	Log = logger
}
