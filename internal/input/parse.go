package input

import (
	"encoding/json"
	"fmt"
	"io"
)

func Decode(r io.Reader) (*StatusInput, error) {
	var input StatusInput
	if err := json.NewDecoder(r).Decode(&input); err != nil {
		return nil, fmt.Errorf("decoding stdin: %w", err)
	}
	return &input, nil
}
