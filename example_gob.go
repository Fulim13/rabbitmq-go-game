// The encoding/gob package in Go provides a more efficient binary serialization format. While not human-readable, it offers faster encoding and decoding than JSON, making it suitable for performance-critical applications.

package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type GameLog struct {
	Event   string
	Details string
}

func encode(gameLog GameLog) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(gameLog)
	return buffer.Bytes(), err
}

func decode(data []byte) (GameLog, error) {
	var gameLog GameLog
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&gameLog)
	return gameLog, err
}

func main() {
	original := GameLog{Event: "Start", Details: "Game started"}
	encodedData, _ := encode(original)
	decoded, _ := decode(encodedData)
	fmt.Println(decoded)
}
