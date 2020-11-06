package store

import (
	"strconv"

	"github.com/jinzhu/gorm"
)

// Request ...
// type Client struct {
// 	ClID          string `json:"ClID";gorm:"not null;unique"`
// 	DepID         string `json:"DepID";gorm:"not null"`
// 	PhoneNumber   string `json:"PhoneNumber";gorm:"not null"`
// 	PhoneRelevant string `json:"PhoneRelevant"`
// }
type Client struct {
	gorm.Model
	ClientID      string `json:"ClientID";gorm:"not null;unique"`
	DepositID     string `json:"DepositID";gorm:"not null;unique"`
	FullName      string `json:"FullName";gorm:"not null;"`
	Bill          string `json:"Bill";gorm:"not null;"`
	OrgName       string `json:"OrgName";gorm:"not null;"`
	Phone         string `json:"PhoneNumber";`
	PhoneRelevant string `json:"PhoneRelevant";`
}

func SetClient(cl *Client) error {
	return x.Create(cl).Error
}

func GetClientsAll() []*Client {
	var clients []*Client

	x.Raw("SELECT * FROM clients").Scan(&clients)
	return clients
}

func GetClients(page string) []*Client {
	var clients []*Client

	offset, _ := strconv.Atoi(page)
	offset = offset - 1

	x.Raw("SELECT * FROM clients").Offset(offset * 2).Limit(2).Scan(&clients)
	return clients
}

// TODO Delete client

// TODO Update client
