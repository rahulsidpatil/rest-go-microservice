package util

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/varstr/uaparser"
)

var _hostName = getHost()

// Palindrome ... check if msg is a palindrome
func Palindrome(msg string) (palindrome bool) {
	var revMsg string
	for _, v := range msg {
		revMsg = string(v) + revMsg
	}
	if revMsg == msg {
		palindrome = true
	}
	return
}

// RequestFrom ... records the source of request.
func RequestFrom(tags map[string]string, time time.Time) {
	name := "requestFrom"
	name = addTagsToName(name, tags)
	if os.Getenv("STATS") == "on" {
		fmt.Printf("%v at: %v\n", name, time.String())
	}
}

// RecordLatency records a handler latency.
func RecordLatency(tags map[string]string, d time.Duration) {
	name := "handler.latency"
	name = addTagsToName(name, tags)
	if os.Getenv("STATS") == "on" {
		fmt.Printf("RecordLatency: %v = %v\n", name, d)
	}
}

// GetStatsTags ...
func GetStatsTags(r *http.Request) map[string]string {
	userBrowser, userOS := parseUserAgent(r.UserAgent())
	stats := map[string]string{
		"browser":  userBrowser,
		"os":       userOS,
		"endpoint": filepath.Base(r.URL.Path),
	}
	if _hostName != "" {
		stats["host"] = _hostName
	}
	return stats
}

func getHost() string {
	host, err := os.Hostname()
	if err != nil {
		return ""
	}
	return host
}

func addTagsToName(name string, tags map[string]string) string {
	// The format we want is: host.endpoint.os.browser
	// if there's no host tag, then we don't use it.
	keyOrder := make([]string, 0, 4)
	if _, ok := tags["host"]; ok {
		keyOrder = append(keyOrder, "host")
	}
	keyOrder = append(keyOrder, "endpoint", "os", "browser")

	buf := &bytes.Buffer{}
	buf.WriteString(name)
	for _, k := range keyOrder {
		buf.WriteByte('.')

		v, ok := tags[k]
		if !ok || v == "" {
			buf.WriteString("no-")
			buf.WriteString(k)
			continue
		}

		writeClean(buf, v)
	}

	return buf.String()
}

func parseUserAgent(uaString string) (browser, os string) {
	ua := uaparser.Parse(uaString)

	if ua.Browser != nil {
		browser = ua.Browser.Name
	}
	if ua.OS != nil {
		os = ua.OS.Name
	}

	return browser, os
}

// writeClean cleans value (e.g. replaces special characters with '-') and
// writes out the cleaned value to buf.
func writeClean(buf *bytes.Buffer, value string) {
	for i := 0; i < len(value); i++ {
		switch c := value[i]; c {
		case '{', '}', '/', '\\', ':', ' ', '\t', '.':
			buf.WriteByte('-')
		default:
			buf.WriteByte(c)
		}
	}
}
