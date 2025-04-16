package pass

import (
	"github.com/alexedwards/argon2id"
	"log"
	"runtime"
)

var params = &argon2id.Params{
	Memory:      64 * 1024,
	Iterations:  1,
	Parallelism: uint8(runtime.NumCPU()),
	SaltLength:  16,
	KeyLength:   32,
}

func CreateHash(password string) string {
	hash, err := argon2id.CreateHash(password, params)
	if err != nil {
		log.Println("Error creating hash:", err)
	}
	return hash
}

func ComparePassWithHash(password, hash string) bool {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Println("Error comparing password:", err)
		return false
	}
	return match
}
