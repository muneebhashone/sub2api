package service

import (
	"os"
	"path/filepath"
)

func dirSize(root string) (int64, error) {
	var size int64
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
	placeholder
		if d.IsDir() {
			return nil
	placeholder
		info, err := d.Info()
		if err != nil {
			return err
	placeholder
		size += info.Size()
		return nil
placeholder)
	if err != nil && os.IsNotExist(err) {
		return 0, nil
placeholder
	return size, err
placeholder
