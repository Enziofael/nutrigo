package consts

const (
	ANSI_RESET = "\033[0m"

	ANSI_BOLD      = "\033[1m"
	ANSI_DIM       = "\033[2m"
	ANSI_ITALIC    = "\033[3m"
	ANSI_UNDERLINE = "\033[4m"

	ANSI_BLINK_SLOW = "\033[5m"
	ANSI_BLINK_FAST = "\033[6m"

	ANSI_REVERSE = "\033[7m" // Inversion

	ANSI_HIDDEN = "\033[8m"

	ANSI_STRIKETHROUGH = "\033[9m"

	ANSI_DOUBLE_UNDERLINE = "\033[21m"

	// Style cancel
	ANSI_NORMAL_INTENSITY = "\033[22m" // Bold и Dim
	ANSI_NO_ITALIC        = "\033[23m"
	ANSI_NO_UNDERLINE     = "\033[24m"
	ANSI_NO_BLINK         = "\033[25m"
	ANSI_NO_REVERSE       = "\033[27m"
	ANSI_NO_HIDDEN        = "\033[28m"
	ANSI_NO_STRIKETHROUGH = "\033[29m"
	ANSI_NO_OVERLINE      = "\033[55m"
	ANSI_OVERLINE         = "\033[53m"
)

const (
	ANSI_BLACK   = "\033[30m"
	ANSI_RED     = "\033[31m"
	ANSI_GREEN   = "\033[32m"
	ANSI_YELLOW  = "\033[33m"
	ANSI_BLUE    = "\033[34m"
	ANSI_MAGENTA = "\033[35m"
	ANSI_CYAN    = "\033[36m"
	ANSI_WHITE   = "\033[37m"

	ANSI_DEFAULT = "\033[39m"
)

const (
	ANSI_BRIGHT_BLACK   = "\033[90m"
	ANSI_BRIGHT_RED     = "\033[91m"
	ANSI_BRIGHT_GREEN   = "\033[92m"
	ANSI_BRIGHT_YELLOW  = "\033[93m"
	ANSI_BRIGHT_BLUE    = "\033[94m"
	ANSI_BRIGHT_MAGENTA = "\033[95m"
	ANSI_BRIGHT_CYAN    = "\033[96m"
	ANSI_BRIGHT_WHITE   = "\033[97m"
)

const (
	ANSI_BG_BLACK   = "\033[40m"
	ANSI_BG_RED     = "\033[41m"
	ANSI_BG_GREEN   = "\033[42m"
	ANSI_BG_YELLOW  = "\033[43m"
	ANSI_BG_BLUE    = "\033[44m"
	ANSI_BG_MAGENTA = "\033[45m"
	ANSI_BG_CYAN    = "\033[46m"
	ANSI_BG_WHITE   = "\033[47m"

	ANSI_BG_DEFAULT = "\033[49m"
)

const (
	ANSI_BG_BRIGHT_BLACK   = "\033[100m"
	ANSI_BG_BRIGHT_RED     = "\033[101m"
	ANSI_BG_BRIGHT_GREEN   = "\033[102m"
	ANSI_BG_BRIGHT_YELLOW  = "\033[103m"
	ANSI_BG_BRIGHT_BLUE    = "\033[104m"
	ANSI_BG_BRIGHT_MAGENTA = "\033[105m"
	ANSI_BG_BRIGHT_CYAN    = "\033[106m"
	ANSI_BG_BRIGHT_WHITE   = "\033[107m"
)

func ANSIColor256(color int) string {
	return "\033[38;5;" + itoa(color) + "m"
}

func ANSIBgColor256(color int) string {
	return "\033[48;5;" + itoa(color) + "m"
}

const (
	ANSI_256_ORANGE = 208
	ANSI_256_PINK   = 205
	ANSI_256_PURPLE = 129
	ANSI_256_TEAL   = 80
	ANSI_256_LIME   = 118
	ANSI_256_GRAY_1 = 244
	ANSI_256_GRAY_2 = 250
	ANSI_256_GRAY_3 = 255
)

func ANSIColorRGB(r, g, b int) string {
	return "\033[38;2;" + itoa(r) + ";" + itoa(g) + ";" + itoa(b) + "m"
}

// ANSIBgColorRGB возвращает код для установки цвета фона по RGB
func ANSIBgColorRGB(r, g, b int) string {
	return "\033[48;2;" + itoa(r) + ";" + itoa(g) + ";" + itoa(b) + "m"
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	if i < 0 {
		// не используется, но для полноты
		return "-" + itoa(-i)
	}
	var buf [3]byte // максимум 3 цифры для 255
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(buf[pos:])
}

