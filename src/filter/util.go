package filter

import (
	"crypto/md5"
	"crypto/rc4"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func CheckFileIsDirectory(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if fi.IsDir() == false {
		return false, errors.New("target file is not folder")
	}

	return true, nil
}

func GetFileSize(file string) (int64, error) {
	fi, err := os.Stat(file)
	if err != nil {
		return 0, err
	}

	if fi.IsDir() {
		return 0, errors.New(fmt.Sprintf("target file %s is not file", file))
	}

	return fi.Size(), nil
}

func InStringArray(value string, arrays []string) bool {
	for _, v := range arrays {
		if v == value {
			return true
		}
	}

	return false
}

func GetFileMD5sum(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}

	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func HasIntersection(a []string, b []string) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}

	t := strings.Join(b, "%") + "%"
	for _, v := range a {
		if strings.Contains(t, v+"%") {
			return true
		}
	}

	return false
}

func IsTrue(needle interface{}) bool {
	return IsFalse(needle) == false
}

func IsFalse(needle interface{}) bool {
	haystack := []interface{}{
		false,
		0,
		"false",
		"",
	}

	for _, v := range haystack {
		if v == needle {
			return true
		}
	}

	return false
}

func Rc4Decrypt(content []byte, key []byte) ([]byte, error) {
	rc4Cipher, err := rc4.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plainText := make([]byte, len(content))
	rc4Cipher.XORKeyStream(plainText, content)

	return plainText, nil
}
