package common

import (
	"bufio"
	"net/textproto"
	"path"
	"strings"

	"github.com/melonmanchan/go-http-server/statuscodes"
)

func ReadAllHeaders(reader bufio.Reader) []string {
	rdr := textproto.NewReader(&reader)
	output := []string{}

	for {
		str, err := rdr.ReadLine()

		if err != nil || len(str) == 0 {
			break
		}

		output = append(output, str)
	}

	return output
}

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
