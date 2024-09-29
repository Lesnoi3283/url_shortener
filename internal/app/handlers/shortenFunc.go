package handlers

import "crypto/sha256"

const defaultShortURLLen = 16

func ShortenURL(url []byte) []byte {
	hasher := sha256.New()
	hasher.Write(url)
	//urlShort := fmt.Sprintf("%x", hasher.Sum(nil)) //optimizing:
	sum := hasher.Sum(nil)
	for i, _ := range sum {
		sum[i] = (sum[i] % 26) + 'a' // 'a' - 'z'
	}
	return sum[:defaultShortURLLen]
}
