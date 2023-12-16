package status

var (
	ErrIgnore  = -2
	ErrUnknown = -1
	None       = 0
)

func Text(code int) string {
	if text, ok := http2text[code]; ok {
		return text
	}
	// add more matching logic here, for instance: gRPC
	return ""
}

func Code(str string) int {
	for code, text := range http2text {
		if text == str {
			return code
		}
	}
	// add more matching logic here, for instance: gRPC
	return ErrUnknown
}
