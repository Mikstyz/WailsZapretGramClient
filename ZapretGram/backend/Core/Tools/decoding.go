package Tools

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"

	//"encoding/base64"
	"encoding/json"
	"errors"
	//"strings"
)

type Pubkey struct {
	pubkey *string
}

func NewKey(key string) *Pubkey {
	return &Pubkey{pubkey: &key}
}

func (pubkey *Pubkey) SetKey(k string) {
	pubkey.pubkey = &k
}

// Шифруем структуру ent и возвращаем []byte
func (pubkey *Pubkey) EncPublicKey(ent interface{}) ([]byte, error) {
	// Проверяем ключ
	if pubkey == nil || pubkey.pubkey == nil {
		return nil, errors.New("ключ не задан")
	}

	// Сериализуем структуру в JSON
	plaintext, err := json.Marshal(ent)
	if err != nil {
		return nil, err
	}

	// Тот же ключ, что и в расшифровке
	key := sha256.Sum256([]byte(*pubkey.pubkey))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Генерируем случайный nonce (ВАЖНО!)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Шифруем JSON → ciphertext
	// nonce добавляется в начало ciphertext (так ждёт твоя DecPublicKey)
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

// Расшифровка []byte в структуру out
func (pubkey *Pubkey) DecPublicKey(enc []byte, out interface{}) error {
	if pubkey == nil || pubkey.pubkey == nil {
		return errors.New("ключ не задан")
	}

	key := sha256.Sum256([]byte(*pubkey.pubkey))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	if len(enc) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, ct := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return err
	}

	return json.Unmarshal(plaintext, out)
}

func (pubkey *Pubkey) MyKey() string {
	return *pubkey.pubkey
}
