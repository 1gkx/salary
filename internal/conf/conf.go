package conf

import (
	"github.com/go-ini/ini"
)

type database struct {
	Driver   string `json:"driver";ini:"driver"`
	Host     string `json:"host";ini:"host"`
	Port     string `json:"port";ini:"port"`
	User     string `json:"user";ini:"user"`
	Password string `json:"pass";ini:"pass"`
	Path     string `json:"path";ini:"path"`
}

type mail struct {
	Host string `json:"host";ini:"host"`
	Port string `json:"port";ini:"port"`
}

type Config struct {
	Prod bool `json:"prod";ini:"prod"`
	Database database          `json:"database";ini:"database"`
	Mail     mail              `json:"mail";ini:"mail"`
	Gateway  string            `json:"gateway";ini:"gateway"`
	Methods  map[string]string `json:"methods";ini:"methods"`
}

var Cfg *Config
var fileConfig *ini.File

func Read() error {

	var err error

	fileConfig, err = ini.InsensitiveLoad("conf/app.ini")
	if err != nil {
		return err
	}

	Cfg = new(Config)

	Cfg.Prod = fileConfig.Section("").Key("prod").MustBool(false)

	Cfg.Database.Driver = fileConfig.Section("database").Key("driver").In(
		"sqlite3", []string{"sqlite3", "postgres", "mysql"},
	)
	Cfg.Database.Host = fileConfig.Section("database").Key("host").String()
	Cfg.Database.Port = fileConfig.Section("database").Key("port").String()
	Cfg.Database.User = fileConfig.Section("database").Key("user").String()
	Cfg.Database.Password = fileConfig.Section("database").Key("password").String()
	Cfg.Database.Path = fileConfig.Section("database").Key("path").String()

	Cfg.Mail.Host = fileConfig.Section("mail").Key("host").String()
	Cfg.Mail.Port = fileConfig.Section("mail").Key("port").String()

	Cfg.Gateway = fileConfig.Section("").Key("gateway").String()

	Cfg.Methods = map[string]string{
		"SMS":     fileConfig.Section("methods").Key("sms").String(),
		"APPROVE": fileConfig.Section("methods").Key("approve").String(),
	}

	return nil

}

func Save() {
	fileConfig.Section("database").Key("driver").SetValue(Cfg.Database.Driver)
	fileConfig.Section("database").Key("host").SetValue(Cfg.Database.Host)
	fileConfig.Section("database").Key("port").SetValue(Cfg.Database.Port)
	fileConfig.Section("database").Key("user").SetValue(Cfg.Database.User)
	fileConfig.Section("database").Key("password").SetValue(Cfg.Database.Password)
	fileConfig.Section("database").Key("path").SetValue(Cfg.Database.Path)

	fileConfig.Section("mail").Key("host").SetValue(Cfg.Mail.Host)
	fileConfig.Section("mail").Key("port").SetValue(Cfg.Mail.Port)

	fileConfig.Section("").Key("gateway").SetValue(Cfg.Gateway)

	fileConfig.SaveTo("conf/app.ini")
}

func Prod() bool {
	return Cfg.Prod
}
