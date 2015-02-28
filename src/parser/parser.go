package parser

import (
	"strings"
)


func GetContentType(request string) string {

	if strings.HasSuffix(request, ".html") || strings.HasSuffix(request, ".html/") {
		return "text/html"
	} else if strings.HasSuffix(request, ".css") || strings.HasSuffix(request, ".css/") {
		return "text/css"
	} else if strings.HasSuffix(request, ".js") || strings.HasSuffix(request, ".js/") {
		return "text/javascript"
	} else if strings.HasSuffix(request, ".jpg") || strings.HasSuffix(request, ".jpg/") {
		return "image/jpeg"
	} else if strings.HasSuffix(request, ".jpeg") || strings.HasSuffix(request, ".jpeg/") {
		return "image/jpeg"
	} else if strings.HasSuffix(request, ".png") || strings.HasSuffix(request, ".png/") {
		return "image/png"
	} else if strings.HasSuffix(request, ".gif") || strings.HasSuffix(request, ".gif/") {
		return "image/gif"
	} else if strings.HasSuffix(request, ".swf") || strings.HasSuffix(request, ".swf/") {
		return "application/x-shockwave-flash"
	}
	return "false"
}