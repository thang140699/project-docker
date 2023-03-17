package ultilities

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func LoadEnvFromFile(config interface{}, configPrefix, envPath string) (err error) {
	godotenv.Load(envPath)
	err = envconfig.Process(configPrefix, config)
	return
}
func LoadEnvFromDir(config interface{}, configPrefix, dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	filePaths := make([]string, 0)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		filePaths = append(filePaths, filepath.Join(dir, f.Name()))
	}
	if err := godotenv.Load(filePaths...); err != nil {
		return err
	}
	return envconfig.Process(configPrefix, config)
}
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}
func StringInArray(str string, arr []string) bool {
	if len(arr) == 0 {
		return false
	}
	for _, val := range arr {
		if strings.TrimSpace(str) == strings.TrimSpace(val) {
			return true
		}
	}
	return false
}
func GetQuery(req *http.Request, key string) (string, bool) {
	if values := req.URL.Query().Get(key); len(values) > 0 {
		return values, true
	}
	return "", false
}
func IntInArray(key int, arr []int) bool {
	low := 0
	high := len(arr) - 1
	for low <= high {
		medium := (low + high) / 2
		if arr[medium] < key {
			low = medium + 1
		} else {
			high = medium - 1
		}
	}
	if low == len(arr) || arr[low] != key {
		return false
	}
	return true
}
