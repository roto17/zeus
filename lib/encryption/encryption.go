package encryptions

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/roto17/zeus/lib/config"
)

// EncryptID encrypts an integer ID using AES encryption
func EncryptID(id uint) string {
	// Convert the ID to a byte slice
	idBytes := []byte(fmt.Sprintf("%d", id))

	// Create a new AES cipher with the provided key
	block, err := aes.NewCipher([]byte(config.GetEnv("encryption_key")))
	if err != nil {
		fmt.Printf("%v", err)
	}

	// Generate a new GCM (Galois/Counter Mode) cipher based on AES
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Printf("%v", err)
	}

	// Generate a random nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Printf("%v", err)
	}

	// Encrypt the ID bytes with AES-GCM
	ciphertext := aesGCM.Seal(nonce, nonce, idBytes, nil)

	// Encode the ciphertext to base64 for easier storage or transmission
	return base64.StdEncoding.EncodeToString(ciphertext)
}

// DecryptID decrypts the encrypted string to retrieve the original integer ID
func DecryptID(encryptedID string) uint {
	// Decode the base64-encoded string
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedID)
	if err != nil {
		fmt.Printf("%v", err)
	}

	// Create a new AES cipher with the provided key
	block, err := aes.NewCipher([]byte(config.GetEnv("encryption_key")))
	if err != nil {
		fmt.Printf("%v", err)
	}

	// Generate a new GCM cipher based on AES
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		// return 0, err
		fmt.Printf("%v", err)
	}

	// Split the nonce and the actual ciphertext
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the ciphertext
	idBytes, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		// return 0, err
		fmt.Printf("%v", err)
	}

	// Convert the decrypted bytes back to an integer
	var id uint
	fmt.Sscanf(string(idBytes), "%d", &id)
	return id
}
