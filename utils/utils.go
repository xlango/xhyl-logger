package utils

import (
	"bytes"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func AppendByte(buffer *bytes.Buffer, b byte) {
	str := strconv.FormatInt(int64(b), 16)
	if len(str) < 2 {
		buffer.WriteString("0")
	}
	buffer.WriteString(str)
}

func GetCurrentExeDir() string {
	dirPath := ""
	file, _ := exec.LookPath(os.Args[0])
	filePath, _ := filepath.Abs(file)
	pathSeparator := string(os.PathSeparator)
	if runtime.GOOS == "windows" {
		lastPos := strings.LastIndex(filePath, pathSeparator)
		if lastPos >= 0 {
			dirPath = Substr(filePath, 0, lastPos)
		}
	} else {
		dirPath = path.Dir(filePath)
	}
	return dirPath
}

func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

func GenerateUuid() int64 {
	rand.Seed(int64(time.Now().Nanosecond()))
	return rand.Int63()
}
