package main

// Prepared based on https://httpstatuses.com/

var StatusInformational = map[int]string{
	100: "Continue",
	101: "Switching Protocols",
	102: "Processing",
}

var StatusSuccess = map[int]string{
	200: "OK",
	201: "Created",
	202: "Accepted",
	203: "Non-authoritative Information",
	204: "No Content",
	205: "Reset Content",
	206: "Partial Content",
	207: "Multi-Status",
	208: "Already Reported",
	226: "IM Used",
}

var StatusRedirection = map[int]string{
	300: "Multiple Choices",
	301: "Moved Permanently",
	302: "Found",
	303: "See Other",
	304: "Not Modified",
	305: "Use Proxy",
	307: "Temporary Redirect",
	308: "Permanent Redirect",
}

var StatusClientError = map[int]string{
	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	406: "Not Acceptable",
	407: "Proxy Authentication Required",
	408: "Request Timeout",
	409: "Conflict",
	410: "Gone",
	411: "Length Required",
	412: "Precondition Failed",
	413: "Payload Too Large",
	414: "Request-URI Too Long",
	415: "Unsupported Media Type",
	416: "Requested Range Not Satisfiable",
	417: "Expectation Failed",
	418: "I'm a teapot",
	421: "Misdirected Request",
	422: "Unprocessable Entity",
	423: "Locked",
	424: "Failed Dependency",
	426: "Upgrade Required",
	428: "Precondition Required",
	429: "Too Many Requests",
	431: "Request Header Fields Too Large",
	444: "Connection Closed Without Response",
	451: "Unavailable For Legal Reasons",
	499: "Client Closed Request",
}

var StatusServerError = map[int]string{
	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
	505: "HTTP Version Not Supported",
	506: "Variant Also Negotiates",
	507: "Insufficient Storage",
	508: "Loop Detected",
	510: "Not Extended",
	511: "Network Authentication Required",
	599: "Network Connect Timeout Error",
}

var (
	InformationalCodes = _aggregateCodes(StatusInformational)
	SuccessCodes       = _aggregateCodes(StatusSuccess)
	ClientErrorCodes   = _aggregateCodes(StatusClientError)
	ServerErrorCodes   = _aggregateCodes(StatusServerError)
)

type Category int

const (
	InformationalCategory Category = iota
	SuccessCategory
	ClientErrorCategory
	ServerErrorCategory
	UnknownCategory
)

func whichCategory(statusCode int) Category {

	// TODO: first check in cache

	if _inSlice(InformationalCodes, statusCode) {
		return InformationalCategory
	}

	if _inSlice(SuccessCodes, statusCode) {
		return SuccessCategory
	}

	if _inSlice(ClientErrorCodes, statusCode) {
		return ClientErrorCategory
	}

	if _inSlice(ServerErrorCodes, statusCode) {
		return ServerErrorCategory
	}

	return UnknownCategory
}

func _aggregateCodes(m map[int]string) []int {
	keys := make([]int, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}

func _inSlice(s []int, v int) bool {
	for _, a := range s {
		if a == v {
			return true
		}
	}
	return false
}
