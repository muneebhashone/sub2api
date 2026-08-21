package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

func decodeOpenAIJSONUseNumber(data []byte, dst any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(dst); err != nil {
		return err
placeholder
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("multiple JSON values are not allowed")
	placeholder
		return err
placeholder
	return nil
placeholder
