package main

import (
	"sync"
)

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
	InformationalCodes = intKeysToSlice(StatusInformational)
	SuccessCodes       = intKeysToSlice(StatusSuccess)
	RedirectionCodes   = intKeysToSlice(StatusRedirection)
	ClientErrorCodes   = intKeysToSlice(StatusClientError)
	ServerErrorCodes   = intKeysToSlice(StatusServerError)
)

const (
	InformationalCategory = "1"
	SuccessCategory       = "2"
	RedirectionCategory   = "3"
	ClientErrorCategory   = "4"
	ServerErrorCategory   = "5"
	UnknownCategory       = "0"
)

type categoryCache struct {
	codes map[int]string
	sync.Mutex
}

func newCategoryCache() *categoryCache {
	return &categoryCache{
		codes: make(map[int]string),
	}
}

func (c *categoryCache) get(statusCode int) (string, bool) {
	c.Lock()
	defer c.Unlock()

	category, exist := c.codes[statusCode]
	if exist {
		return category, true
	}

	return "", false
}

func (c *categoryCache) set(statusCode int, description string) {
	c.Lock()
	defer c.Unlock()

	c.codes[statusCode] = description
}

var cache = newCategoryCache()

func resolveCodeDescription(statusCode int, filter *statusFilter) string {
	cat, exist := cache.get(statusCode)
	if exist {
		return cat
	}

	if filter.showAll && intSliceContains(InformationalCodes, statusCode) ||
		filter.onlyInfo && intSliceContains(InformationalCodes, statusCode) {
		description := StatusInformational[statusCode]
		cache.set(statusCode, description)
		return description
	}

	if filter.showAll && intSliceContains(SuccessCodes, statusCode) ||
		filter.onlySuccess && intSliceContains(SuccessCodes, statusCode) {
		description := StatusSuccess[statusCode]
		cache.set(statusCode, description)
		return description
	}

	if filter.showAll && intSliceContains(RedirectionCodes, statusCode) {
		description := StatusRedirection[statusCode]
		cache.set(statusCode, description)
		return description
	}

	if filter.showAll && intSliceContains(ClientErrorCodes, statusCode) ||
		filter.onlyClientErr && intSliceContains(ClientErrorCodes, statusCode) {
		description := StatusClientError[statusCode]
		cache.set(statusCode, description)
		return description
	}

	if filter.showAll && intSliceContains(ServerErrorCodes, statusCode) ||
		filter.onlyServerErr && intSliceContains(ServerErrorCodes, statusCode) {
		description := StatusServerError[statusCode]
		cache.set(statusCode, description)
		return description
	}

	if filter.showAll {
		cache.set(statusCode, "N/A")
	}

	return "N/A"
}
