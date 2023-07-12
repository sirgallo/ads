package utils

import "crypto/sha256"
import "crypto/sha1"
import "hash/fnv"
import "encoding/hex"


func Sha256Hash(value string) [32]byte {
	data := []byte(value)
	hash := sha256.Sum256(data)
	return hash
}

func Sha256HexHash(value string) string {
	hash := Sha256Hash(value)
	hexStr := hex.EncodeToString(hash[:])
	return hexStr
}

func Sha1Hash(value string) [20]byte {
	data := []byte(value)
	hash := sha1.Sum(data)
	return hash
}

func Sha1HexHash(value string) string {
	hash := Sha1Hash(value)
	hexStr := hex.EncodeToString(hash[:])
	return hexStr
}

func FnvHash(key string) uint32 {
	hash := fnv.New32a()
	data := []byte(key)
	hash.Write(data)
	hashSum := hash.Sum32()
	return hashSum
}