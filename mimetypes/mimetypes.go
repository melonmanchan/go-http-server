package mimetypes

import "path"

var filetypes = map[string]string{
	".html": "text/html",
	".htm":  "text/html",
	".js":   "application/javascript",
	".css":  "text/css",
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	"":      "text/plain",
}

func GetContentType(fileName string) string {
	name := path.Ext(fileName)
	return filetypes[name]
}
