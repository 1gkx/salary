package utils

var (
	SmsSuccess = []byte(`{"data":[{"code":"1234"}],"error":"","extra":"СМС успешно отправлено","result":"1"}`)
	SmsError   = []byte(`{"data":"","error":"Send SMS Error","extra":"Ошибка отправки смс","result":"0"}`)

	ClientSuccess = []byte(`{"data":"","error":"","extra":"Операция по утверждению карты успешно завершена","result":"1"}`)
	ClientError1  = []byte(`{"data":"","error":"","extra":"Ошибка при подключении Faktura у клиента","result":"0"}`)
	ClientError2  = []byte(`{"data":"","error":"","extra":"Ошибка при подключении ПримСоцЛайн. У клиента обнаружено несколько договоров","result":"0"}`)
	ClientError3  = []byte(`{"data":"","error":"","extra":"Ошибка при подключении ПримСоцЛайн, напишите vld-webdev@pskb.com","result":"0"}`)
	ClientError4  = []byte(`{"data":"","error":"","extra":"Клиент не прошел процедуру упрощенной идентификации","result":"0"}`)
)
