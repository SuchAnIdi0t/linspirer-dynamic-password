/*
Author: SuchAnIdi0t
Lisence: MIT License

PoC only, use at your own risk.
*/

package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

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
	log.Println(str3)
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
	var swdid, stu_id string
	fmt.Print("swdid:")
	fmt.Scanln(&swdid)
	swdid = strings.ToLower(swdid)
	fmt.Print("stu_id:")
	fmt.Scanln(&stu_id)
	log.Printf("Result: %s", calc(swdid, stu_id)[:6])
}
