package model

import "strings"

type Message struct{
	User int `json:"user_id"`
	Company int `json:"company_id"`
	SubDomain string `json:"sub_domain"`
	Url string `json:"domain_uri"`
	Roomid int `json:"room_id"`
	MessageBody string `json:"body"`
	Identifier string `json:"identifier"`
	Type string `json:"type"`
}

type Token struct{
	Access string `json:"access"`
	Refresh string `json:"refresh"`
}
func (m *Message) GetUri() string{
	urls := strings.Split(m.Url, "/")
	return urls[1] + "//" + urls[3] + "/api/token/"
}
