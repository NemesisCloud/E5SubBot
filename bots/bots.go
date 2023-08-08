package bots

import (
	"fmt"
	"github.com/amirulandalib/E5SubBot/config"
	"github.com/amirulandalib/E5SubBot/db"
	"github.com/amirulandalib/E5SubBot/logger"
	"go.uber.org/zap"
	"golang.org/x/net/proxy"
	tb "gopkg.in/tucnak/telebot.v2"
	"net/http"
	"time"
)

var bot *tb.Bot

const (
	logo = `
		____    ____   ____     __     ___       __    _____  __
		/ __/___/ __/  / __/_ __/ /    / _ )___  / /_  / __/ |/ /
	   / _//___/__ \  _\ \/ // / _ \  / _  / _ \/ __/ / _//    / 
	  /___/   /____/ /___/\_,_/_.__/ /____/\___/\__/ /___/_/|_/  
				Translated to English by @CaptainLightyear
	`
)


func Start() {
	var err error
	fmt.Print(logo)

	config.Init()
	logger.Init()
	db.Init()
	InitTask()

	poller := &tb.LongPoller{Timeout: 15 * time.Second}
	midPoller := tb.NewMiddlewarePoller(poller, func(upd *tb.Update) bool {
		if upd.Message == nil {
			return true
		}
		if !upd.Message.Private() {
			return false
		}
		return true
	})
	setting := tb.Settings{
		Token:  config.BotToken,
		Poller: midPoller,
	}

	if config.Socks5 != "" {
		dialer, err := proxy.SOCKS5("tcp", config.Socks5, nil, proxy.Direct)
		if err != nil {
			zap.S().Fatalw("failed to get dialer",
				"error", err, "proxy", config.Socks5)
		}
		transport := &http.Transport{}
		transport.Dial = dialer.Dial
		setting.Client = &http.Client{Transport: transport}
	}

	bot, err = tb.NewBot(setting)
	if err != nil {
		zap.S().Fatalw("failed to create bot", "error", err)
	}
	fmt.Printf("Bot: %d %s\n", bot.Me.ID, bot.Me.Username)

	makeHandlers()
	fmt.Println("Bot Started successfully")
	fmt.Println("------------")
	bot.Start()
}

func makeHandlers() {
	// 所有用户
	bot.Handle("/start", bStart)
	bot.Handle("/my", bMy)
	bot.Handle("/bind", bBind)
	bot.Handle("/unbind", bUnBind)
	bot.Handle("/export", bExport)
	bot.Handle("/help", bHelp)
	bot.Handle(tb.OnText, bOnText)
	// 管理员
	bot.Handle("/task", bTask)
	bot.Handle("/log", bLog)
}
