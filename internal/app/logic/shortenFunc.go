package logic

import "crypto/sha256"

const defaultShortURLLen = 16

//todo: make it not exportable.

// ShortenURL generates a short version of given slice of bytes.
// It takes SHA256 sum of given bytes, translates sum into letters (from 'a' to 'z'), cuts it and returns.
func ShortenURL(url []byte) []byte {
	hasher := sha256.New()
	hasher.Write(url)
	//urlShort := fmt.Sprintf("%x", hasher.Sum(nil)) //optimizing:
	sum := hasher.Sum(nil)
	for i := range sum {
		sum[i] = (sum[i] % 26) + 'a' // from 'a' to 'z'
	}
	return sum[:defaultShortURLLen]
}
