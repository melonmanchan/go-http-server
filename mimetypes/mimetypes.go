package mimetypes

import (
	"path"
	"strings"
)

var filetypes = map[string]string{
	".html": "text/html",
	".htm":  "text/html",
	"":      "text/plain",
	".js":   "application/javascript",
	".pdf":  "application/pdf",
	".css":  "text/css",
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".ttf	": "font/ttf",
	".woff":  "font/woff",
	".woff2": "font/woff2",
}

func GetContentType(fileName string) string {
	name := path.Ext(strings.ToLower(fileName))
	return filetypes[name]
}
