package signatures

import (
	"crypto/md5"
	"fmt"
	"strings"
)

// Sign returns the signatures for the given message.
func Sign(content string, secret string) string {
	data := []byte(secret + content + secret)
	result := fmt.Sprintf("%x", md5.Sum(data))
	return strings.ToUpper(result)
}
