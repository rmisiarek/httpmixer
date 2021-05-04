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
	InformationalCodes = _aggregateCodes(StatusInformational)
	SuccessCodes       = _aggregateCodes(StatusSuccess)
	RedirectionCodes   = _aggregateCodes(StatusRedirection)
	ClientErrorCodes   = _aggregateCodes(StatusClientError)
	ServerErrorCodes   = _aggregateCodes(StatusServerError)
)

type Category int

const (
	InformationalCategory Category = iota
	SuccessCategory
	RedirectionCategory
	ClientErrorCategory
	ServerErrorCategory
	UnknownCategory
)

// func (c Category) String() string {
// 	return [...]string{"Informational", "Success", "Redirection", "Client Error", "Server Error", "Unknown"}[c]
// }

type resolvedCategory struct {
	category    Category
	description string
}

type categoryCache struct {
	codes map[int]resolvedCategory
	sync.Mutex
}

func newCategoryCache() *categoryCache {
	return &categoryCache{
		codes: make(map[int]resolvedCategory),
	}
}

func (c *categoryCache) get(statusCode int) (resolvedCategory, bool) {
	c.Lock()
	defer c.Unlock()

	category, exist := c.codes[statusCode]
	if exist {
		return category, true
	}

	return resolvedCategory{category: UnknownCategory, description: "---"}, false
}

func (c *categoryCache) set(statusCode int, category Category, description string) {
	c.Lock()
	defer c.Unlock()

	c.codes[statusCode] = resolvedCategory{
		category:    category,
		description: description,
	}
}

var cache = newCategoryCache()

func resolveCategory(statusCode int, filter *statusFilter) (resolvedCategory, bool) {
	cat, exist := cache.get(statusCode)
	if exist {
		return cat, true
	}

	if *filter.showAll && _inSlice(InformationalCodes, statusCode) || *filter.onlyInfo && _inSlice(InformationalCodes, statusCode) {
		description := StatusInformational[statusCode]
		cache.set(statusCode, InformationalCategory, description)
		return resolvedCategory{category: InformationalCategory, description: description}, true
	}

	if *filter.showAll && _inSlice(SuccessCodes, statusCode) || *filter.onlySuccess && _inSlice(SuccessCodes, statusCode) {
		description := StatusSuccess[statusCode]
		cache.set(statusCode, SuccessCategory, description)
		return resolvedCategory{category: SuccessCategory, description: description}, true
	}

	if *filter.showAll && _inSlice(RedirectionCodes, statusCode) {
		description := StatusRedirection[statusCode]
		cache.set(statusCode, RedirectionCategory, description)
		return resolvedCategory{category: RedirectionCategory, description: description}, true
	}

	if *filter.showAll && _inSlice(ClientErrorCodes, statusCode) || *filter.onlyClientErr && _inSlice(ClientErrorCodes, statusCode) {
		description := StatusClientError[statusCode]
		cache.set(statusCode, ClientErrorCategory, description)
		return resolvedCategory{category: ClientErrorCategory, description: description}, true
	}

	if *filter.showAll && _inSlice(ServerErrorCodes, statusCode) || *filter.onlyServerErr && _inSlice(ServerErrorCodes, statusCode) {
		description := StatusServerError[statusCode]
		cache.set(statusCode, ServerErrorCategory, description)
		return resolvedCategory{category: ServerErrorCategory, description: description}, true
	}

	if *filter.showAll {
		cache.set(statusCode, UnknownCategory, "---")
		return resolvedCategory{category: UnknownCategory, description: "---"}, false
	}

	return resolvedCategory{category: UnknownCategory, description: "---"}, false
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
