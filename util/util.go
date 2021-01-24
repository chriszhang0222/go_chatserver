package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func GetChannel(userId int, domain string) string{
	str := fmt.Sprintf("%s:user_%d", domain, userId)
	return str
}

func SendRequest(url string, data map[string]interface{})([]byte, error){
	client := &http.Client{Timeout: 5 * time.Second}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return result, nil
}

func SendRequestWithAuth(url string, data map[string]interface{}, token string) ([]byte, error){
	client := &http.Client{}
	jsonStr, _ := json.Marshal(data)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	request.Header.Add("Authorization", "Bearer " + token)
	response, _ := client.Do(request)
	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	return result, err
}
