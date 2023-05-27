package bots

import (
	"fmt"
	"github.com/amirulandalib/E5SubBot/config"
	"github.com/amirulandalib/E5SubBot/model"
	"github.com/amirulandalib/E5SubBot/pkg/microsoft"
	"github.com/amirulandalib/E5SubBot/service/srv_client"
	"github.com/amirulandalib/E5SubBot/util"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
)

func bBind(m *tb.Message) {
	bot.Send(m.Chat,
		fmt.Sprintf("👉 App registration: [click here for redirection](%s)", microsoft.GetRegURL()),
		tb.ModeMarkdown,
	)

	bot.Send(m.Chat,
		"⚠ Please reply with your `client_id(space)client_secret`",
		&tb.SendOptions{ParseMode: tb.ModeMarkdown,
			ReplyMarkup: &tb.ReplyMarkup{ForceReply: true}},
	)

	UserStatus[m.Chat.ID] = StatusBind1
	UserClientId[m.Chat.ID] = m.Text
}

func bBind1(m *tb.Message) {
	if !m.IsReply() {
		bot.Send(m.Chat, "⚠ Please bind your account by replying to this message")
		return
	}
	tmp := strings.Split(m.Text, " ")
	if len(tmp) != 2 {
		bot.Send(m.Chat, "⚠ Bad format")
		return
	}
	id := tmp[0]
	secret := tmp[1]
	bot.Send(m.Chat,
		fmt.Sprintf("👉 Authorize your account: [click here for redirection](%s)", microsoft.GetAuthURL(id)),
		tb.ModeMarkdown,
	)

	bot.Send(m.Chat,
		"⚠ Please reply with the format eg:---->>> `http://localhost/……...(space) Name or Alias` (for management tag)\n\n example format: `http://localhost/e5sub?code=0.8yruhw8ry847r amirul`",
		&tb.SendOptions{ParseMode: tb.ModeMarkdown,
			ReplyMarkup: &tb.ReplyMarkup{ForceReply: true},
		},
	)
	UserStatus[m.Chat.ID] = StatusBind2
	UserClientId[m.Chat.ID] = id
	UserClientSecret[m.Chat.ID] = secret
}

func bBind2(m *tb.Message) {
	if !m.IsReply() {
		bot.Send(m.Chat, "⚠ Please Bind your account by replying to this message")
		return
	}
	if len(srv_client.GetClients(m.Chat.ID)) == config.BindMaxNum {
		bot.Send(m.Chat, "⚠ You reached the maximum number of bindings!!")
		return
	}
	bot.Send(m.Chat, "Binding The Account...")

	tmp := strings.Split(m.Text, " ")
	if len(tmp) != 2 {
		bot.Send(m.Chat, "😥 Bad format please add a space and a name after the localhost url\n\n example format: `http://localhost/e5sub?code=0.8yruhw8ry847r amirul`")
	}
	code := util.GetURLValue(tmp[0], "code")
	alias := tmp[1]

	id := UserClientId[m.Chat.ID]
	secret := UserClientSecret[m.Chat.ID]

	refresh, err := microsoft.GetTokenWithCode(id, secret, code)
	if err != nil {
		bot.Send(m.Chat, fmt.Sprintf("Unable to get RefreshToken ERROR:%s", err))
		return
	}
	bot.Send(m.Chat, "🎉 Account Token obtained successfully!")

	refresh, info, err := microsoft.GetUserInfo(id, secret, refresh)
	if err != nil {
		bot.Send(m.Chat, fmt.Sprintf("Unable to get user information ERROR:%s", err))
		return
	}
	c := &model.Client{
		TgId:         m.Chat.ID,
		RefreshToken: refresh,
		MsId:         util.Get16MD5Encode(gjson.Get(info, "id").String()),
		Alias:        alias,
		ClientId:     id,
		ClientSecret: secret,
		Other:        "",
	}

	if srv_client.IsExist(c.TgId, c.ClientId) {
		bot.Send(m.Chat, "⚠ The account been bounded already, no need to bind it again")
		return
	}

	bot.Send(m.Chat,
		fmt.Sprintf("ms_id：%s\nuserPrincipalName：%s\ndisplayName：%s",
			c.MsId,
			gjson.Get(info, "userPrincipalName").String(),
			gjson.Get(info, "displayName").String(),
		),
	)

	if err = srv_client.Add(c); err != nil {
		bot.Send(m.Chat, "😥 User information's failed to write to database")
		return
	}

	bot.Send(m.Chat, "✨ The binding is successful!\nNow sit back and relax you dont need to do anything.\nThe Bot will send you account status message every 6 hours..")
	delete(UserStatus, m.Chat.ID)
	delete(UserClientId, m.Chat.ID)
	delete(UserClientSecret, m.Chat.ID)
}

func bUnBind(m *tb.Message) {
	var inlineKeys [][]tb.InlineButton
	clients := srv_client.GetClients(m.Chat.ID)

	for _, u := range clients {
		inlineBtn := tb.InlineButton{
			Unique: "unbind" + strconv.Itoa(u.ID),
			Text:   u.Alias,
			Data:   strconv.Itoa(u.ID),
		}
		bot.Handle(&inlineBtn, bUnBindInlineBtn)
		inlineKeys = append(inlineKeys, []tb.InlineButton{inlineBtn})
	}

	bot.Send(m.Chat,
		fmt.Sprintf("⚠ Select an account to unbind it\n\nCurrent Number of bound accounts: %d/%d", len(srv_client.GetClients(m.Chat.ID)), config.BindMaxNum),
		&tb.ReplyMarkup{InlineKeyboard: inlineKeys},
	)
}
func bUnBindInlineBtn(c *tb.Callback) {
	id, _ := strconv.Atoi(c.Data)
	if err := srv_client.Del(id); err != nil {
		zap.S().Errorw("failed to delete db data",
			"error", err,
			"id", c.Data,
		)
		bot.Send(c.Message.Chat, "⚠ Unbinding failed!")
		return
	}
	bot.Send(c.Message.Chat, "✨ Unbound successfully!")
	bot.Respond(c)
}
