package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// جلب التوكن من إعدادات فيرسل
	botToken := os.Getenv("BOT_TOKEN")
	devID := "8590415901"

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return
	}

	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		return
	}

	if update.Message != nil {
		chatID := update.Message.Chat.ID

		// التفاعل مع أمر البداية
		if update.Message.IsCommand() && update.Message.Command() == "start" {
			msg := tgbotapi.NewMessage(chatID, "✅ البوت يعمل الآن على فيرسل!\nأرسل جهة الاتصال للبدء.")
			bot.Send(msg)
		}

		// إشعار المطور (أنت) بأي رسالة تصل للبوت
		notification := tgbotapi.NewMessage(8590415901, fmt.Sprintf("رسالة من %d إلى المطور %s", chatID, devID))
		bot.Send(notification)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
