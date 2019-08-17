package dwlog

import (
	"fmt"
	"os"
	"time"
	"path"
	"sort"
	"strings"
	"io/ioutil"
	"path/filepath"
)

type FileWriter struct {
	Name string
	MaxCount int
	file *os.File
	nextDayTime int64
}

func (w *FileWriter) Write(p []byte) (n int, err error) {
	now := time.Now()

	if w.file == nil {
		w.openFile(&now)
	} else if now.Unix() >= w.nextDayTime {
		w.file.Close()
		w.openFile(&now)
	}

	return w.file.Write(p)
}

func (w *FileWriter) openFile(now *time.Time) (err error) {
	name := fmt.Sprintf("%s.%s", w.Name, now.Format("20060102"))

	// remove symbol link if exist
	os.Remove(w.Name)

	w.file, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// create symbol
	err = os.Symlink(path.Base(name), w.Name)
	if err != nil {
		fmt.Println(err)
		return err
	}

	year, month, day := now.Date()
	w.nextDayTime = time.Date(year, month, day+1, 0, 0, 0, 0, now.Location()).Unix()

	if w.MaxCount > 0 {
		go w.cleanFiles()
	}

	return nil
}

func (w *FileWriter) cleanFiles() {
	dir := path.Dir(w.Name)

	fileList, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	prefix := path.Base(w.Name) + "."

	var matches []string
	for _, f := range fileList {
		if !f.IsDir() && strings.HasPrefix(f.Name(), prefix) {
			matches = append(matches, f.Name())
		}
	}

	if len(matches) > w.MaxCount {
		sort.Sort(sort.Reverse(sort.StringSlice(matches)))

		for _, f := range matches[w.MaxCount:] {
			file := filepath.Join(dir, f)
			os.Remove(file)
		}
	}
}