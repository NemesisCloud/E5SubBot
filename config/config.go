package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

const (
	LogBasePath    string = "./log/"
	WelcomeContent string = "Welcome!"
	HelpContent    string = `
	Commands :

	/my View your bound account information

	/bind bind a new account.

	/unbind Unbind the account used to bind previously on this bot.

	/export Export account information as a [.json file] and upload it to telegram.

	/help send this message.
`
)

func Init() {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		zap.S().Fatalw("failed to read config", "error", err)
	}
	BotToken = viper.GetString("bot_token")
	Cron = viper.GetString("cron")
	Socks5 = viper.GetString("socks5")

	viper.SetDefault("errlimit", 5)
	viper.SetDefault("bindmax", 5)
	viper.SetDefault("goroutine", 10)

	BindMaxNum = viper.GetInt("bindmax")
	MaxErrTimes = viper.GetInt("errlimit")
	Notice = viper.GetString("notice")

	MaxGoroutines = viper.GetInt("goroutine")
	Admins = getAdmins()
	DB = viper.GetString("db")
	Table = viper.GetString("table")

	switch DB {
	case "mysql":
		Mysql = mysqlConfig{
			Host:                viper.GetString("mysql.host"),
			Port:                viper.GetInt("mysql.port"),
			User:                viper.GetString("mysql.user"),
			Password:            viper.GetString("mysql.password"),
			DB:                  viper.GetString("mysql.database"),
			SSLMode:             viper.GetString("mysql.ssl_mode"),
			EnabledTLSProtocols: viper.GetString("mysql.enabled_tls_protocols"),
		}
	case "sqlite":
		Sqlite = sqliteConfig{
			DB: viper.GetString("sqlite.db"),
		}
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		MaxGoroutines = viper.GetInt("goroutine")
		BindMaxNum = viper.GetInt("bindmax")
		MaxErrTimes = viper.GetInt("errlimit")
		Notice = viper.GetString("notice")
		Admins = getAdmins()
	})
}
func getAdmins() []int64 {
	var result []int64
	admins := strings.Split(viper.GetString("admin"), ",")
	for _, v := range admins {
		id, _ := strconv.ParseInt(v, 10, 64)
		result = append(result, id)
	}
	return result
}
