package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var mainKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Make Order 🍽"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("About ❔"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Register cafe [Coming Soon] ⌛"),
	),
)

var orderKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Café Enjothie in UNIST 🇰🇷"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Back 🔙"),
	),
)

var menuFirstKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Smoothies 🥤"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Snacks 🥨"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Submit order 🔘"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Back 🔙"),
	),
)

var menuData = map[string][]string{
	"Smoothies 🥤": []string{
		"Strawberry Shortcake Smoothie 🍓",
		"Apply Smoothie 🍏",
		"Pineapple Smoothie 🍍",
		"Back 🔙",
	},
	"Snacks 🥨": []string{
		"Egg 🥚",
		"Cookie 🍪",
		"Salad 🥗",
		"Back 🔙",
	},
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("MENU"))
	if err != nil {
		log.Printf("Failed to connect with bot. Error: %v\n", err)
		return
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	state := make(map[string]string)
	category := make(map[string]string)
	userCart := make(map[string]Cart)

	orderNumber := 100

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		messageText := update.Message.Text
		currentState := state[update.Message.From.UserName]
		currentCart := userCart[update.Message.From.UserName]

		log.Printf("[%s] %s - state: %s", update.Message.From.UserName, update.Message.Text, currentState)

		if update.InlineQuery != nil {
			fmt.Print("Doing something")
		}

		if currentState == "" || currentState == "main" { // main state
			if messageText == "Make Order 🍽" {
				msgText := "Choose restaurant"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				msg.ReplyMarkup = orderKeyboard
				bot.Send(msg)
				state[update.Message.From.UserName] = "order"
			} else if messageText == "About ❔" {
				msgText := fmt.Sprintln("Smart Menu ordering bot is the new system of simple and fast ordering.\nYou can choose any restaraunt" +
					" from the list and make your orders\n")
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				msg.ReplyMarkup = mainKeyboard
				bot.Send(msg)
				state[update.Message.From.UserName] = "main"
			} else {
				msgText := fmt.Sprintf("Hello, %s!\nWelcome to the Smart Menu ordering bot.\n", update.Message.From.FirstName)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				msg.ReplyMarkup = mainKeyboard
				bot.Send(msg)
				state[update.Message.From.UserName] = "main"
			}

		} else if currentState == "order" {
			if messageText == "Café Enjothie in UNIST 🇰🇷" {
				msgText := fmt.Sprintln("Choose category")
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				msg.ReplyMarkup = menuFirstKeyboard
				bot.Send(msg)
				state[update.Message.From.UserName] = "category"
			} else if messageText == "Back 🔙" {
				msgText := fmt.Sprintf("Hello, %s!\nWelcome to the Smart Menu ordering bot.\n", update.Message.From.FirstName)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				msg.ReplyMarkup = mainKeyboard
				bot.Send(msg)
				state[update.Message.From.UserName] = "main"
			} else {
				msgText := "Choose restaurant"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				msg.ReplyMarkup = orderKeyboard
				bot.Send(msg)
				state[update.Message.From.UserName] = "order"
			}

		} else if currentState == "category" {
			if messageText == "Back 🔙" {
				msgText := "Choose restaurant"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				msg.ReplyMarkup = orderKeyboard
				bot.Send(msg)
				state[update.Message.From.UserName] = "order"
				continue
			} else if menuData[messageText] != nil {
				msgText := "Please choose product"
				var smoothieList [][]tgbotapi.KeyboardButton
				smoothieList = make([][]tgbotapi.KeyboardButton, len(menuData[messageText]))

				for index, drink := range menuData[messageText] {
					smoothieList[index] = make([]tgbotapi.KeyboardButton, 1)
					newButton := tgbotapi.NewKeyboardButton(drink)
					smoothieList[index][0] = newButton
				}
				replyKeyboard := tgbotapi.NewReplyKeyboard(smoothieList...)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				msg.ReplyMarkup = replyKeyboard
				bot.Send(msg)
				state[update.Message.From.UserName] = "product"
				category[update.Message.From.UserName] = messageText
			} else if messageText == "Submit order 🔘" {
				currentCart = Cart{}
				msgText := fmt.Sprintf("Your order has been submitted. Your order number is: %d", orderNumber)
				orderNumber++
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				msg.ReplyMarkup = mainKeyboard
				bot.Send(msg)
				state[update.Message.From.UserName] = "main"
				userCart[update.Message.From.UserName] = currentCart
			} else {
				msgText := fmt.Sprintln("Choose category")
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				msg.ReplyMarkup = menuFirstKeyboard
				bot.Send(msg)
				state[update.Message.From.UserName] = "category"
			}

			// show cart
			var cartList [][]tgbotapi.InlineKeyboardButton
			cartList = make([][]tgbotapi.InlineKeyboardButton, len(currentCart.Products))

			for index, drink := range currentCart.Products {
				cartList[index] = make([]tgbotapi.InlineKeyboardButton, 1)
				newButton := tgbotapi.NewInlineKeyboardButtonData(drink, string(index))
				cartList[index][0] = newButton
			}
			inlineCartKeyboard := tgbotapi.NewInlineKeyboardMarkup(cartList...)
			msgText := fmt.Sprintln("Your cart:")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
			msg.ReplyMarkup = inlineCartKeyboard
			bot.Send(msg)

		} else if currentState == "product" {
			if messageText == "Back 🔙" {
				msgText := fmt.Sprintln("Choose category")
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				msg.ReplyMarkup = menuFirstKeyboard
				bot.Send(msg)
				state[update.Message.From.UserName] = "category"
			} else {
				currentCart.Products = append(currentCart.Products, messageText)
			}
			// show cart

			var cartList [][]tgbotapi.InlineKeyboardButton
			cartList = make([][]tgbotapi.InlineKeyboardButton, len(currentCart.Products))

			for index, drink := range currentCart.Products {
				cartList[index] = make([]tgbotapi.InlineKeyboardButton, 1)
				newButton := tgbotapi.NewInlineKeyboardButtonData(drink, string(index))
				cartList[index][0] = newButton
			}
			inlineCartKeyboard := tgbotapi.NewInlineKeyboardMarkup(cartList...)
			msgText := fmt.Sprintln("Your cart:")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
			msg.ReplyMarkup = inlineCartKeyboard
			bot.Send(msg)

			userCart[update.Message.From.UserName] = currentCart
		} else {
			log.Println("roflan bug")
			state[update.Message.From.UserName] = "main"
		}
	}
}

// Cart cart structure for user
type Cart struct {
	ID       string
	Price    int64
	Products []string
}
