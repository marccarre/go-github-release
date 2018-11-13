package validate

import "os"

// Files validates the existence and readability of the provided files, or
// returns the relevant error.
func Files(files []string) error {
	for _, file := range files {
		if err := File(file); err != nil {
			return err
		}
	}
	return nil
}

// File validates the existence and readability of the provided file, or
// returns the relevant error.
func File(file string) error {
	_, err := os.Stat(file)
	if err != nil {
		return err
	}
	return nil
}
