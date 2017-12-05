package getter

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// FileDetector implements Detector to detect file paths.
type FileDetector struct{}

func (d *FileDetector) Detect(src, pwd string) (string, bool, error) {
	fmt.Printf("Megan in file detector\n")
	if len(src) == 0 {
		return "", false, nil
	}
	fmt.Printf("Megan src is %s\n", src)
	if !filepath.IsAbs(src) || ((runtime.GOOS == "windows") && (strings.ToLower(src)[:2] != "c:")) {
		fmt.Printf("Megan src is NOT considered absolute")
		if pwd == "" {
			return "", true, fmt.Errorf(
				"relative paths require a module with a pwd")
		}

		// Stat the pwd to determine if its a symbolic link. If it is,
		// then the pwd becomes the original directory. Otherwise,
		// `filepath.Join` below does some weird stuff.
		//
		// We just ignore if the pwd doesn't exist. That error will be
		// caught later when we try to use the URL.
		if fi, err := os.Lstat(pwd); !os.IsNotExist(err) {
			if err != nil {
				return "", true, err
			}
			if fi.Mode()&os.ModeSymlink != 0 {
				pwd, err = filepath.EvalSymlinks(pwd)
				if err != nil {
					return "", true, err
				}

				// The symlink itself might be a relative path, so we have to
				// resolve this to have a correctly rooted URL.
				pwd, err = filepath.Abs(pwd)
				if err != nil {
					return "", true, err
				}
			}
		}

		fmt.Printf("Megan src is %s\n", src)
		fmt.Printf("Megan pwd is %s\n", pwd)
		src = filepath.Join(pwd, src)
		fmt.Printf("Megan src is %s\n", src)
	}

	return fmtFileURL(src), true, nil
}

func fmtFileURL(path string) string {
	if runtime.GOOS == "windows" {
		// Make sure we're using "/" on Windows. URLs are "/"-based.
		path = filepath.ToSlash(path)
		return fmt.Sprintf("file://%s", path)
	}

	// Make sure that we don't start with "/" since we add that below.
	if path[0] == '/' {
		path = path[1:]
	}
	return fmt.Sprintf("file:///%s", path)
}
