package utils

import (
	"bytes"
	"encoding/json"
	"fmt"

	"math/rand"
	"strconv"

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

	// ['phone' => preg_replace('/[^0-9]+/', '', $request->get('mobilePhoneClient'))],
	send_data := map[string]interface{}{
		"phone": data,
	}

	jdata, _ := json.Marshal(send_data)
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

	code := randomSmsCode()
	tmpl := fmt.Sprintf(
		"{\"data\":[{\"code\":\"%s\"}],\"error\":\"\",\"extra\":\"СМС успешно отправлено\",\"result\":\"1\"}",
		code,
	)
	var d = []byte(tmpl)
	res123 := new(Responce)
	json.Unmarshal(d, &res123)
	return res123, nil
}

func randomSmsCode() string {
	rand.Seed(time.Now().UTC().UnixNano())
	return strconv.Itoa(randInt(1000, 9999))
}

func randInt(min int, max int) int {
    return min + rand.Intn(max-min)
}

