package internal

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"strings"
)

func CompressMazeData(maze *Maze) (error, string) {
	if jsonData, err := json.Marshal(maze); err == nil {
		var compressionBuffer bytes.Buffer
		compressor := gzip.NewWriter(&compressionBuffer)
		if _, err := compressor.Write(jsonData); err == nil {
			if err := compressor.Close(); err == nil {
				return nil, base64.StdEncoding.EncodeToString(compressionBuffer.Bytes())
			} else {
				return err, ""
			}
		} else {
			return err, ""
		}
	} else {
		return err, ""
	}
}

func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for name, file := range Assets.Files {
		if file.IsDir() || !strings.HasSuffix(name, ".tmpl") {
			continue
		}
		h, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		t, err = t.New(name).Parse(string(h))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
