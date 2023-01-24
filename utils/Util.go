package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

/**
 * Return absolute path to the file
 */
func Abs(path string) string {
	e, err := os.Executable()
	if err != nil {
		panic(err)
	}
	eio := filepath.Dir(e)
	file := filepath.Join(eio, path)
	_, err = os.Stat(file)
	if err != nil {
		panic(err)
	}
	return file
}

/**
 * Log and print the error
 */
func LogError(err interface{}) string {
	str := fmt.Sprintf("Error: %v\n", err)
	fmt.Fprintf(os.Stderr, str)
	log.Output(3, str)
	return ""
}

/**
 * Quickly and easily generate individual or bulk sets of universally unique identifiers (UUIDs).
 * See: constants.Env.SECRET_TOKEN
 */
func UUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return LogError(err)
	}
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

/**
 * Returns the value to which the specified key is mapped, or defaultValue if this map contains no mapping for the key.
 * @param key  - the key whose associated value is to be returned
 * @param defVal - the default mapping of the key
 * @return - the value to which the specified key is mapped, or defVal if this map contains no mapping for the key or NIL of default defVal is omitted
 */
func Get(source map[string]interface{}, key string, defVal ...interface{}) interface{} {
	val, found := source[key]
	if found {
		return val
	}
	if len(defVal) > 0 {
		return defVal[0]
	}
	return nil
}
