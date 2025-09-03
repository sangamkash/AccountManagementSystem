package log_helper

const (
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorOrange  = "\033[38;5;208m" // 256-color orange
	colorReset   = "\033[0m"
)

func LogServiceInitializing(s string) string {
	return colorYellow + "=====initialized ::" + s + " ::starting======" + colorReset
}
func LogServiceInitialized(s string) string {
	return colorGreen + "=====initialized ::" + s + " ::successfully======" + colorReset
}
func LogServiceFailToStarted(s string) string {
	return colorRed + "=====fail to start ::" + s + " :: fail!!!======" + colorReset
}

func LogPanic(s string) string {
	return colorRed + "=====!!!!!Panic ::" + s + " :: Panic!!!!!======" + colorReset
}

func LogValidator(s string) string {
	return colorYellow + "Validator::" + s + colorReset
}
