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
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"github.com/dwango/yashiro/internal/values"
	"github.com/dwango/yashiro/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

const (
	cacheFileName   = "values"
	keyFileName     = "key"
	keyHashFileName = ".key_hash"
)

var defaultCacheBasePath string

type fileCache struct {
	cacheBasePath string
	cipherBlock   cipher.Block
	expired       bool
}

func newFileCache(cfg config.FileCacheConfig) (Cache, error) {
	fc := &fileCache{
		cacheBasePath: defaultCacheBasePath,
	}

	if len(cfg.CachePath) != 0 {
		fc.cacheBasePath = cfg.CachePath
	}
	// create cache directory
	if err := os.MkdirAll(fc.cacheBasePath, 0777); err != nil {
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

const (
	// 30 days
	expiredDuration time.Duration = 30 * 24 * time.Hour
)

// Load implements Cache.
func (f *fileCache) Load(_ context.Context) (values.Values, bool, error) {
	fInfo, err := f.getFileStat(cacheFileName)
	if err != nil {
		f.expired = true
		return nil, f.expired, nil
	}
	// check if cache is expired
	f.expired = time.Since(fInfo.ModTime().Local()) > expiredDuration

	cacheCipherText, err := f.readFile(cacheFileName)
	if err != nil {
		return nil, false, err
	}

	val, err := f.decryptCache(cacheCipherText)
	if err != nil {
		return nil, false, err
	}

	return val, f.expired, nil
}

// Save implements Cache.
func (f *fileCache) Save(_ context.Context, val values.Values) error {
	if !f.expired {
		return nil
	}

	encryptedCache, err := f.encryptCache(val)
	if err != nil {
		return err
	}

	if err := f.writeToFile(cacheFileName, encryptedCache); err != nil {
		return err
	}

	return nil
}

func (f *fileCache) readOrCreateKey() ([]byte, error) {
	var key []byte
	// check key file exists
	if _, err := f.getFileStat(keyFileName); err != nil {
		key = make([]byte, 32)

		// create key file
		if _, err := rand.Read(key); err != nil {
			return nil, err
		}
		if err := f.writeToFile(keyFileName, key); err != nil {
			return nil, err
		}

		// hashing key
		keyHash, err := bcrypt.GenerateFromPassword(key, 5)
		if err != nil {
			return nil, err
		}
		if err := f.writeToFile(keyHashFileName, keyHash); err != nil {
			return nil, err
		}

		return key, nil
	}

	var err error
	// read key file
	key, err = f.readFile(keyFileName)
	if err != nil {
		return nil, err
	}

	// read key hash file
	keyHash, err := f.readFile(keyHashFileName)
	if err != nil {
		return nil, err
	}

	// check key is not tampered
	if err := bcrypt.CompareHashAndPassword(keyHash, key); err != nil {
		return nil, err
	}

	return key, nil
}

func (f *fileCache) decryptCache(cacheCipherText []byte) (values.Values, error) {
	if len(cacheCipherText) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := cacheCipherText[:aes.BlockSize]
	cacheCipherText = cacheCipherText[aes.BlockSize:]

	cachePlainText := make([]byte, len(cacheCipherText))
	stream := cipher.NewOFB(f.cipherBlock, iv)
	stream.XORKeyStream(cachePlainText, cacheCipherText)

	values := make(values.Values)
	if err := json.Unmarshal(cachePlainText, &values); err != nil {
		return nil, err
	}

	return values, nil
}

func (f *fileCache) encryptCache(values values.Values) ([]byte, error) {
	cacheJSON, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}

	cacheCipherText := make([]byte, aes.BlockSize+len(cacheJSON))
	iv := cacheCipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewOFB(f.cipherBlock, iv)
	stream.XORKeyStream(cacheCipherText[aes.BlockSize:], cacheJSON)

	return cacheCipherText, nil
}

func (f fileCache) getFileStat(filename string) (os.FileInfo, error) {
	return os.Stat(f.cacheBasePath + "/" + filename)
}

func (f fileCache) readFile(filename string) ([]byte, error) {
	return os.ReadFile(f.cacheBasePath + "/" + filename)
}

func (f fileCache) writeToFile(filename string, data []byte) error {
	file, err := os.Create(f.cacheBasePath + "/" + filename)
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
	const cachePath = "/yashiro"

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		defaultCacheBasePath = "/tmp" + cachePath + "/cache"
		return
	}
	defaultCacheBasePath = cacheDir + cachePath
}
