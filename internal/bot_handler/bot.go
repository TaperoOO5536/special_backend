package bot_handler

import (
	"log"
	"strconv"

	"github.com/TaperoOO5536/special_backend/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotHandler struct {
	bot *tgbotapi.BotAPI
	adminID int64
}

func NewBotHandler(bot *tgbotapi.BotAPI, adminID int64) BotHandler {
	return BotHandler{bot: bot, adminID: adminID}
}

func (h *BotHandler) HandleOrderMessage(eventType string, userId int, order models.Order) {
	var text string
	switch eventType {
	case "order.create":
		text = "ваш заказ создан, номер вашего заказа: " + string(order.Number)
	case "order.update":
		text = "ваш заказ номер " + string(order.Number) +
						"изменён, новый статус вашего заказа: " + order.Status
	case "order.delete":
		text = "ваш заказ номер " + string(order.Number) + " удалён"
	default:
		log.Fatalf("unknown event type")
	}
	h.sendMessage(userId, text)
}

func (h *BotHandler) HandleUserEventMessage(eventType string, userEvent models.UserEvent, user models.User) {
	var text string
	switch eventType {
	case "userevent.create":
		text = "Новая запись на мероприятие " + userEvent.Event.Title +
	" записался " + user.Nickname + " мест свободно " + 
	strconv.FormatInt(userEvent.Event.OccupiedSeats, 10) + "/" + strconv.FormatInt(userEvent.Event.TotalSeats, 10)
	case "userevent.update":
		text = "Изменена запись на мероприятие " + userEvent.Event.Title +
	" пользователя " + user.Nickname + " мест свободно " + 
	strconv.FormatInt(userEvent.Event.OccupiedSeats, 10) + "/" + strconv.FormatInt(userEvent.Event.TotalSeats, 10)
	case "userevent.delete":
		text = "Отменена запись на мероприятие " + userEvent.Event.Title +
	" пользователя " + user.Nickname + " мест свободно " + 
	strconv.FormatInt(userEvent.Event.OccupiedSeats, 10) + "/" + strconv.FormatInt(userEvent.Event.TotalSeats, 10)
	default:
		log.Fatalf("unknown event type")
	}
	h.sendMessage(int(h.adminID), text)
}

func (h *BotHandler) sendMessage(id int, text string) {
	tgMsg := tgbotapi.NewMessage(int64(id), text)
	h.bot.Send(tgMsg)
}