package internal

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
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
