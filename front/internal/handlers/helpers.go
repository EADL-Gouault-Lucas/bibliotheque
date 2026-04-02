package handlers

import "net/url"

// encodeMsg encodes an error message for use in a redirect query parameter.
func encodeMsg(msg string) string {
	return url.QueryEscape(msg)
}
