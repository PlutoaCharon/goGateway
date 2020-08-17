package public

import (
	"crypto/sha256"
	"fmt"
)

// 获取加盐密码
func GenSaltPassword(passport string, salt string) string {
	sOb1 := sha256.New()
	sOb1.Write([]byte(passport))
	rs1 := fmt.Sprintf("%x", sOb1.Sum(nil))
	sOb2 := sha256.New()
	sOb2.Write([]byte(rs1 + salt))
	rs2 := fmt.Sprintf("%x", sOb2.Sum(nil))
	return rs2
}
