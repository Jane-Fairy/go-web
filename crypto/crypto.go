package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
)

func demo() {
	//import "crypto/sha256"
	h := sha256.New()
	io.WriteString(h, "His money is twice tainted: 'taint yours and 'taint mine.")
	fmt.Printf("% x", h.Sum(nil))

	//import "crypto/sha1"
	b := sha1.New()
	io.WriteString(h, "His money is twice tainted: 'taint yours and 'taint mine.")
	fmt.Printf("% x", b.Sum(nil))

	//import "crypto/md5"
	c := md5.New()
	io.WriteString(h, "需要加密的密码")
	fmt.Printf("%x", c.Sum(nil))
}
