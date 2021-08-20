package logger

type Logger interface {
	Info(message string)
	InfoWithMeta(message string, meta map[string]interface{})
}
