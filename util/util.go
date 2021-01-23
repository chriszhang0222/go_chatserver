package util

import "fmt"

func GetChannel(userId int, domain string) string{
	str := fmt.Sprintf("%s:user_%d", domain, userId)
	return str
}

func SendRequest(){

}
