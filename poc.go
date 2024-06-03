/*
Author: SuchAnIdi0t
Lisence: MIT License

PoC only, use at your own risk.
*/

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var swdid, account, stu_id string

func getStuID(swdid string, acc string) string {
	log.Println("Getting stu_id...")
	postBody := `{"id":"1","!version":"6","jsonrpc":"2.0","is_encrypt":false,"method":"com.linspirer.user.getuserinfo","client_version":"tong_hem_6.00.006.5","params":{"email":"` + acc + `","model":"BZH-W30","swdid":"` + swdid + `"}}`
	res, err := http.Post("https://cloud.linspirer.com:883/public-interface.php", "application/json", strings.NewReader(postBody))
	if err != nil {
		log.Fatal(err)
		return ""
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	log.Println("API responsed with code:", res.StatusCode)

	content := func(body []byte) string {
		block, err := aes.NewCipher([]byte("1191ADF18489D8DA"))
		if err != nil {
			log.Fatal(err)
			return ""
		}
		ciphertext, err := base64.StdEncoding.DecodeString(string(body))
		if err != nil {
			log.Fatal(err)
			return ""
		}

		iv := []byte("5E9B755A8B674394")
		mode := cipher.NewCBCDecrypter(block, iv)
		mode.CryptBlocks(ciphertext, ciphertext)
		plaintext := string(ciphertext)
		return plaintext[:len(plaintext)-int(plaintext[len(plaintext)-1])]
	}(body)
	var data map[string]interface{}
	json.Unmarshal([]byte(content), &data)
	v, ok := data["data"].(string)
	if ok {
		log.Fatal("Error: ", v)
	}

	v2, ok := data["data"].(map[string]interface{})
	if !ok {
		log.Fatal("Unkonwn Error")
	}

	v3, ok := v2["id"].(float64)
	if !ok {
		log.Fatal("Unkonwn Error when getting id")
	}
	v4 := strconv.FormatInt(int64(v3), 10)
	log.Printf("Your stu_id: %s", v4)
	log.Println("You could save this stu_id for future use")
	return v4
}

func ca(str1 string) string {
	str2, _ := strconv.ParseInt(str1, 16, 64)
	v := strconv.Itoa(int(str2))
	if len(v) <= 8 {
		return v
	} else {
		return v[len(v)-8:]
	}
}

func calc(swdid string, stu_id string) string {
	str3 := time.Now().Format("20060102") + swdid + "40E06F51-30D0-D6AD-7F7D-008AD0ADC570" + stu_id
	h := md5.Sum([]byte(str3))
	b := fmt.Sprintf("%x", h)
	log.Println(b)
	if len(b) > 8 {
		return ca(b[len(b)-8:])
	} else {
		return "err"
	}
}

func main() {
	flag.StringVar(&stu_id, "i", "", "stu_id for calculating.If this is given,account will be ignored.")
	flag.StringVar(&account, "a", "", "Account for fetching stu_id. No need to enter if stu_id has already been specified.")
	flag.StringVar(&swdid, "s", "", "SWDID for fetching stu_id and calculating. You should NOT let it empty.")
	flag.Parse()
	if stu_id != "" {
		if swdid == "" {
			fmt.Print("\n*****************\n*\n*  Error: SWDID is empty.\n*\n*****************\n\nUsage:\n\n")
			flag.PrintDefaults()
			return
		}
		swdid = strings.ToLower(swdid)
		log.Printf("Result: %s", calc(swdid, stu_id)[:6])
	} else {
		if account == "" || swdid == "" {
			fmt.Print("\n*****************\n*\n*  Error: SWDID or account is empty.\n*\n*****************\n\nUsage:\n\n")
			flag.PrintDefaults()
			return
		}
		account = strings.ToLower(account)
		log.Printf("Result: %s", calc(swdid, getStuID(swdid, account))[:6])
	}
}
