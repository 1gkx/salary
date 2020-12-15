package utils

import (
	"bytes"
	"encoding/json"
	"fmt"

	"net/http"
	"time"

	"github.com/1gkx/salary/internal/conf"
)

// Responce ...
type Responce struct {
	Data   []map[string]interface{} `json:"data"`
	Error  string                   `json:"error"`
	Extra  string                   `json:"extra"`
	Result string                   `json:"result"`
}

// Post ...
func Post(data interface{}, method string) (*Responce, error) {

	if !conf.Prod() {
		return returnTestCode()
	}

	jdata, _ := json.Marshal(data)
	req, err := http.NewRequest(
		"POST",
		conf.Cfg.Gateway,
		bytes.NewReader([]byte(jdata)),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Action", method)
	req.Header.Set("HOOK", "Y")
	// req.Header.Set("TEST", "Y")

	client := &http.Client{Timeout: 300 * time.Second}
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	res := new(Responce)
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (res *Responce) GetSmsCode() string {
	return fmt.Sprintf("%s", res.Data[0]["code"])
}

func (res *Responce) GetExpiredSmsCode() string {
	return time.Now().Add(1 * time.Minute).Format(time.RFC3339)
}

func returnTestCode() (*Responce, error) {
	var d = []byte(`{"data":[{"code":"1234"}],"error":"","extra":"СМС успешно отправлено","result":"1"}`)
	res123 := new(Responce)
	json.Unmarshal(d, &res123)
	return res123, nil
}
