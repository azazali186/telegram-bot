package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	token := "6452255074:AAF9mEzQD7yt0zSbXiTEiZ3NxRE4n6N6KoQ"
	if token == "" {
		log.Fatal("Telegram token not provided!")
	}

	setBotCommands(token)

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the database
	db, err := sql.Open("mysql", "root:Aj189628@@tcp(localhost:3306)/telegramData")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				handleMessageCommand(update.Message, bot)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command. Type /start to see available commands.")
				bot.Send(msg)
			}
		} else if update.CallbackQuery != nil {
			insertMessageData(update.CallbackQuery, db)
			handleCallbackQuery(update.CallbackQuery, bot)
		}
	}
}

func insertMessageData(callback *tgbotapi.CallbackQuery, db *sql.DB) {
	// fmt.Printf("Hello World\n\n\n")
	username := callback.Message.From.UserName
	messageID := callback.Message.MessageID
	command := callback.Data
	chatID := callback.Message.Chat.ID
	location := ""

	if callback.Message.Location != nil {
		location = fmt.Sprintf("%f;%f", callback.Message.Location.Latitude, callback.Message.Location.Longitude)
	}

	_, err := db.Exec(`INSERT INTO messages (username, message_id, chat_id, user_location, command) VALUES (?, ?, ?, ?, ?)`, username, messageID, chatID, location, command)
	if err != nil {
		log.Println("Failed to insert message into database:", err)
	}
}

func setBotCommands(token string) {
	type BotCommand struct {
		Command     string `json:"command"`
		Description string `json:"description"`
	}

	data := struct {
		Commands []BotCommand `json:"commands"`
	}{
		Commands: []BotCommand{
			{Command: "start", Description: "Start the bot and get the command list."},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	url := "https://api.telegram.org/bot" + token + "/setMyCommands"
	_, err = http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
}

func handleMessageCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	switch message.Command() {
	case "start":

		desc := "Hello " + message.From.UserName

		// Add more rows as needed...
		commandsKeyboard := getCommandList()

		msg := tgbotapi.NewMessage(message.Chat.ID, desc)
		msg.ReplyMarkup = commandsKeyboard
		bot.Send(msg)
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Unknown command. Type /start to see available commands.")
		bot.Send(msg)
	}
}

func getCommandList() tgbotapi.InlineKeyboardMarkup {
	commandsKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Product", "product"),
			tgbotapi.NewInlineKeyboardButtonData("Balance", "balance"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Wallet", "wallet"),
			tgbotapi.NewInlineKeyboardButtonData("Withdraw", "withdraw"),
		),

		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Deposit", "deposit"),
			tgbotapi.NewInlineKeyboardButtonData("Referral", "referral"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Network", "network"),
			tgbotapi.NewInlineKeyboardButtonData("Delete", "delete"),
		),
	)
	return commandsKeyboard
}

func getWalletInfo() tgbotapi.InlineKeyboardMarkup {
	commandsKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("BTC", "btc"),
			tgbotapi.NewInlineKeyboardButtonData("ETH", "eth"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("USDT", "usdt"),
			tgbotapi.NewInlineKeyboardButtonData("BNB", "bnb"),
		),

		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("TRON", "tron"),
			tgbotapi.NewInlineKeyboardButtonData("Back to main", "depositeBack"),
		),
	)
	return commandsKeyboard
}

func handleCallbackQuery(callback *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI) {
	commandsKeyboard := getCommandList()
	switch callback.Data {
	case "product":
		bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "You chose Product!"))
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Thanks for choose "+callback.Data)
		msg.ReplyMarkup = commandsKeyboard
		bot.Send(msg)
	case "balance":
		bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "You chose Balance!"))
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Thanks for choose "+callback.Data)
		msg.ReplyMarkup = commandsKeyboard
		bot.Send(msg)
	case "wallet":
		bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "You chose Wallet!"))
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Thanks for choose "+callback.Data)
		msg.ReplyMarkup = commandsKeyboard
		bot.Send(msg)
	case "withdraw":
		bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "You chose Withdraw!"))
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Thanks for choose "+callback.Data)
		msg.ReplyMarkup = commandsKeyboard
		bot.Send(msg)
	case "referral":
		bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "You chose Referral!"))
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Thanks for choose "+callback.Data)
		msg.ReplyMarkup = commandsKeyboard
		bot.Send(msg)
	case "network":
		bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "You chose Network!"))
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Thanks for choose "+callback.Data)
		msg.ReplyMarkup = commandsKeyboard
		bot.Send(msg)
	case "delete":
		deleteMsgConfig := tgbotapi.DeleteMessageConfig{
			ChatID:    callback.Message.Chat.ID,
			MessageID: callback.Message.MessageID,
		}
		_, err := bot.DeleteMessage(deleteMsgConfig)
		if err != nil {
			log.Println("Failed to delete message:", err)
		}
		bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "You chose Delete!"))
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Type /start to see available commands.")
		bot.Send(msg)
	case "deposit":
		bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "You chose Deposit!"))
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Thanks for choose "+callback.Data)
		msg.ReplyMarkup = getWalletInfo()
		bot.Send(msg)
	case "btc":
		bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "You chose BitCoin!"))
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Thanks for choose "+callback.Data)
		msg.ReplyMarkup = getWalletInfo()
		bot.Send(msg)
	case "eth":
		bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "You chose Ether!"))
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Thanks for choose "+callback.Data)
		msg.ReplyMarkup = getWalletInfo()
		bot.Send(msg)
	case "bnb":
		bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "You chose BNB!"))
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Thanks for choose "+callback.Data)
		msg.ReplyMarkup = getWalletInfo()
		bot.Send(msg)
	case "tron":
		bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "You chose TRON!"))
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Thanks for choose "+callback.Data)
		msg.ReplyMarkup = getWalletInfo()
		bot.Send(msg)
	case "usdt":
		bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "You chose USDT!"))
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Thanks for choose "+callback.Data)
		msg.ReplyMarkup = getWalletInfo()
		bot.Send(msg)
	case "depositeBack":
		desc := "Hello " + callback.Message.From.UserName
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, desc)
		msg.ReplyMarkup = commandsKeyboard
		bot.Send(msg)
	default:
		bot.AnswerCallbackQuery(tgbotapi.NewCallback(callback.ID, "Unknown command!"))
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Thanks for choose "+callback.Data)
		msg.ReplyMarkup = commandsKeyboard
		bot.Send(msg)
	}
}
