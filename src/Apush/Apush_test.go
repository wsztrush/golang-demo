package Apush

import (
	"testing"
	"fmt"
	"encoding/base64"
)

func TestGetQueryString(t *testing.T) {
	GetQueryString()
}

func TestAesEncrypt(t *testing.T) {
	a := AesEncrypt([]byte("token=1508428927981,P_292$0$"), GetKeyByte("YouMustModifyThisToAnOtherString"))
	fmt.Println(base64.StdEncoding.EncodeToString(a))
}