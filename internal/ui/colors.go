package ui

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"
)

const (
	checkMark = "✓"
	cross     = "✗"
	arrow     = "→"
	notice    = "⚠"
)

func Success(msg string) string {
	return colorGreen + checkMark + " " + msg + colorReset
}

func Error(msg string) string {
	return colorRed + cross + " " + msg + colorReset
}

func Working(msg string) string {
	return colorBlue + arrow + " " + msg + colorReset
}

func Warn(msg string) string {
	return colorYellow + notice + " " + msg + colorReset
}
