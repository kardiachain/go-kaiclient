package hexutil

import (
	"fmt"
	"testing"
)

func TestSign_Verify(t *testing.T) {
	isValid, addr := VerifySign("0x1b5a79e827562665dcf6ff6ee449494f502b9f2d6c87d7b7bd275589f37981337919cca8b3ebe720fc700720ab10031354abdbca6cd910b92af5553523c6bc351c", "Long Dep Trai")

	fmt.Println("IsValid", isValid)
	fmt.Println("Addr", addr)
}
