/**
 * Copyright 2024 DWANGO Co., Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cache

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dwango/yashiro/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

const (
	keyFileName     = "key"
	keyHashFileName = "keyHash"
)

var defaultCacheBasePath string

type fileCache struct {
	cachePath      string
	cipherBlock    cipher.Block
	expireDuration time.Duration
	filenamePrefix string
}

func newFileCache(cfg config.FileCacheConfig, expireDuration time.Duration, options ...Option) (Cache, error) {
	opts := defaultOpts
	for _, o := range options {
		o(opts)
	}

	cachePath := defaultCacheBasePath
	if len(cfg.CachePath) != 0 {
		cachePath = cfg.CachePath
	}
	filenamePrefix := keyToHex(strings.Join(opts.CacheKeys, "_")) + "_"

	fc := &fileCache{
		cachePath:      cachePath,
		expireDuration: expireDuration,
		filenamePrefix: filenamePrefix,
	}

	// create cache directory
	if err := os.MkdirAll(fc.cachePath, 0777); err != nil {
		return nil, err
	}

	// read or create key
	key, err := fc.readOrCreateKey()
	if err != nil {
		return nil, err
	}

	fc.cipherBlock, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return fc, nil
}

// Load implements Cache.
func (f *fileCache) Load(_ context.Context, key string, decrypt bool) (*string, bool, error) {
	filename := keyToHex(key)

	fInfo, err := f.getFileInfo(filename, false)
	if err != nil {
		// cache file not found
		return nil, false, nil
	}
	// check if cache is expired
	expired := time.Since(fInfo.ModTime().Local()) > f.expireDuration

	cacheByte, err := f.readFile(filename, false)
	if err != nil {
		return nil, false, err
	}
	if !decrypt {
		cache := string(cacheByte)
		return &cache, expired, nil
	}

	valueByte, err := f.decryptCache(cacheByte)
	if err != nil {
		return nil, false, err
	}
	value := string(valueByte)

	return &value, expired, nil
}

// Save implements Cache.
func (f *fileCache) Save(_ context.Context, key string, value *string, encrypt bool) error {
	if value == nil {
		return nil
	}

	filename := keyToHex(key)
	valueByte := []byte(*value)

	if encrypt {
		var err error
		valueByte, err = f.encryptCache(valueByte)
		if err != nil {
			return err
		}
	}

	if err := f.writeToFile(filename, valueByte, false); err != nil {
		return err
	}

	return nil
}

func (f *fileCache) readOrCreateKey() ([]byte, error) {
	var key []byte
	// check key file exists
	if _, err := f.getFileInfo(keyFileName, false); err != nil {
		key = make([]byte, 32)

		// create key file
		if _, err := rand.Read(key); err != nil {
			return nil, err
		}
		if err := f.writeToFile(keyFileName, key, false); err != nil {
			return nil, err
		}

		// hashing key
		keyHash, err := bcrypt.GenerateFromPassword(key, 5)
		if err != nil {
			return nil, err
		}
		if err := f.writeToFile(keyHashFileName, keyHash, true); err != nil {
			return nil, err
		}

		return key, nil
	}

	var err error
	// read key file
	key, err = f.readFile(keyFileName, false)
	if err != nil {
		return nil, err
	}

	// read key hash file
	keyHash, err := f.readFile(keyHashFileName, true)
	if err != nil {
		return nil, err
	}

	// check key is not tampered
	if err := bcrypt.CompareHashAndPassword(keyHash, key); err != nil {
		return nil, err
	}

	return key, nil
}

func (f *fileCache) decryptCache(cipherText []byte) ([]byte, error) {
	if len(cipherText) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(f.cipherBlock, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, nil
}

func (f *fileCache) encryptCache(plainText []byte) ([]byte, error) {
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(f.cipherBlock, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return cipherText, nil
}

func (f fileCache) getFileInfo(filename string, hidden bool) (os.FileInfo, error) {
	filename = f.filenamePrefix + filename
	if hidden {
		filename = "." + filename
	}

	return os.Stat((filepath.Join(f.cachePath, filename)))
}

func (f fileCache) readFile(filename string, hidden bool) ([]byte, error) {
	filename = f.filenamePrefix + filename
	if hidden {
		filename = "." + filename
	}

	return os.ReadFile(filepath.Join(f.cachePath, filename))
}

func (f fileCache) writeToFile(filename string, data []byte, hidden bool) error {
	filename = f.filenamePrefix + filename
	if hidden {
		filename = "." + filename
	}

	file, err := os.Create(filepath.Join(f.cachePath, filename))
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}

func init() {
	const cachePath = "yashiro"

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		defaultCacheBasePath = filepath.Join(os.TempDir(), cachePath, "cache")
		return
	}

	defaultCacheBasePath = filepath.Join(cacheDir, cachePath)
}
