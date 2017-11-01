package common

import (
	"path"
	"strings"

	"github.com/melonmanchan/go-http-server/statuscodes"
)

func sanitizeQueryParameter(url string) string {
	return strings.Split(url, "?")[0]
}

func GetPathFromHeader(header string) (string, *statuscodes.HTTPStatus) {
	paths := strings.Split(header, " ")

	if paths[0] != "GET" {
		return "", &statuscodes.MethodNotAllowed
	}

	return sanitizeQueryParameter(paths[1]), nil
}

func SafePath(reqPath string, basePath string) string {
	return "./" + path.Join(basePath, strings.Replace(reqPath, "..", "", -1))
}
