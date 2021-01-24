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

type PersistResponse struct {
	Success bool `json:"success"`
	Message map[string]interface{} `json:"message"`
	Id int `json:"id"`
	Identifier string `json:"identifier"`
	Timestamp string `json:"timestamp"`
	DiscussionDate string `json:"discussion_date"`
	DiscussionTime string `json:"discussion_time"`
	UserIds []int `json:user_ids`
}


func (m *Message) GetUri() string{
	urls := strings.Split(m.Url, "/")
	return urls[1] + "//" + urls[3] + "/api/token/"
}

func (m *Message) GetPersistUri() string{
	urls := strings.Split(m.Url, "/")
	return urls[1] + "//" + urls[3] + "/chat/save_message_to_db/"
}
