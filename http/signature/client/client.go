package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/codec"
	"github.com/zeromicro/zero-examples/http/signature/internal"
)

var crypt = flag.Bool("crypt", false, "encrypt body or not")

func hs256(key []byte, body string) string {
	h := hmac.New(sha256.New, key)
	io.WriteString(h, body)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func main() {
	flag.Parse()

	var err error
	body := "admin"
	if *crypt {
		bodyBytes, err := codec.EcbEncrypt(internal.Key, []byte(body))
		if err != nil {
			log.Fatal(err)
		}
		body = base64.StdEncoding.EncodeToString(bodyBytes)
	}

	r, err := http.NewRequest(http.MethodPost, "http://localhost:3333/a/b?c=first&d=second", strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	timestamp := time.Now().Unix()
	sha := sha256.New()
	sha.Write([]byte(body))
	bodySign := fmt.Sprintf("%x", sha.Sum(nil))
	contentOfSign := strings.Join([]string{
		strconv.FormatInt(timestamp, 10),
		http.MethodPost,
		r.URL.Path,
		r.URL.RawQuery,
		bodySign,
	}, "\n")
	sign := hs256(internal.Key, contentOfSign)
	var mode string
	if *crypt {
		mode = "1"
	} else {
		mode = "0"
	}
	content := strings.Join([]string{
		"version=v1",
		"type=" + mode,
		fmt.Sprintf("key=%s", base64.StdEncoding.EncodeToString(internal.Key)),
		"time=" + strconv.FormatInt(timestamp, 10),
	}, "; ")

	encrypter, err := codec.NewRsaEncrypter(internal.PubKey)
	if err != nil {
		log.Fatal(err)
	}

	output, err := encrypter.Encrypt([]byte(content))
	if err != nil {
		log.Fatal(err)
	}

	encryptedContent := base64.StdEncoding.EncodeToString(output)
	r.Header.Set("X-Content-Security", strings.Join([]string{
		fmt.Sprintf("key=%s", internal.Fingerprint),
		"secret=" + encryptedContent,
		"signature=" + sign,
	}, "; "))
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
	io.Copy(os.Stdout, resp.Body)
}

/*
待编码字符：
"admin"

对应的字节数组：
[97,100,109,105,110]

对应的二进制表示：
[01100001,01100100,01101101,01101001,01101110]

重新按照6bit进行分组(最后一组不足6bit,在后面补0)：
[011000,010110,010001,101101,011010,010110,1110 00]

转换成整数表示：
[24,22,17,45,26,22,56]

查表得到对应字符(不足4byte,补'=')：
['Y','W','R','t','a','W','4','=']

可见字符表(大写字母 + 小写字母 + 数字 + '+' + '/')：
*/
