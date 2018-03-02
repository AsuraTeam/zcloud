package os

import gos "os"
import "path/filepath"

// Delete a tree of empty directories. Returns non-nil error if
// an empty directory cannot be deleted. Does not return an error
// if there are non-empty directories. This function can be used
// to prune empty directories from a tree.
func RemoveEmpty(path string) error {
	var dirs []string

	return filepath.Walk(path, func(path string, info gos.FileInfo, err error) error {
		if info.Mode().IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})

	for i := len(dirs) - 1; i >= 0; i-- {
		err := gos.Remove(dirs[i])
		if err != nil && !IsNotEmpty(err) {
			return err
		}
	}

	return nil
}
