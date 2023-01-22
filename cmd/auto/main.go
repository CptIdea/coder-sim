package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"workSim/interfaces"
	"workSim/internal/entity"
	"workSim/internal/worker"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var users map[int64]interfaces.Worker

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New("5923996581:AAGzjP9t8JvbGgdUY1fgE46UXY-5J9QTIZA", opts...)
	if err != nil {
		panic(err)
	}

	users = make(map[int64]interfaces.Worker)

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	text := update.Message.Text
	fmt.Println(update.Message.Chat.ID, text)

	w, ok := users[update.Message.Chat.ID]
	if !ok || text == "/start" {
		w = worker.NewWorker(
			entity.Feels{
				Energy:     10,
				NeedCoffee: 5,
				NeedSmoke:  5,
			},
			10,
		)
		users[update.Message.Chat.ID] = w
	}

	switch text {
	case "/start":
		sendStats(ctx, b, w, update)
	case "/work":
		done := w.Work()
		if done == nil {
			sendBusyMessage(ctx, b, w, update)
			return
		}

		sendWaitMessage(ctx, b, w, update)
		go func() {
			if <-done {
				sendStats(ctx, b, w, update)
			} else {
				sendDeadMessage(ctx, b, w, update)
				delete(users, update.Message.Chat.ID)
			}
		}()
	case "/smoke":
		done := w.Smoke()
		if done == nil {
			sendBusyMessage(ctx, b, w, update)
			return
		}

		sendWaitMessage(ctx, b, w, update)
		go func() {
			if <-done {
				sendStats(ctx, b, w, update)
			} else {
				sendDeadMessage(ctx, b, w, update)
				delete(users, update.Message.Chat.ID)
			}
		}()
	case "/drink_coffee":
		done := w.DrinkCoffee()
		if done == nil {
			sendBusyMessage(ctx, b, w, update)
			return
		}

		sendWaitMessage(ctx, b, w, update)
		go func() {
			if <-done {
				sendStats(ctx, b, w, update)
			} else {
				sendDeadMessage(ctx, b, w, update)
				delete(users, update.Message.Chat.ID)
			}
		}()

	case "/sleep":
		done := w.Sleep()
		if done == nil {
			sendBusyMessage(ctx, b, w, update)
			return
		}

		sendWaitMessage(ctx, b, w, update)
		go func() {
			if <-done {
				sendStats(ctx, b, w, update)
			} else {
				sendDeadMessage(ctx, b, w, update)
				delete(users, update.Message.Chat.ID)
			}
		}()
	}
}

func sendStats(ctx context.Context, b *bot.Bot, w interfaces.Worker, update *models.Update) {
	b.SendMessage(
		ctx, &bot.SendMessageParams{
			Text: fmt.Sprintf(
				"```\n%s```", PrettyJson(
					map[string]interface{}{
						"current_feels": w.CurrentFeels(),
						"max_feels":     w.MaxFeels(),
					},
				),
			),
			ChatID:    update.Message.Chat.ID,
			ParseMode: models.ParseModeMarkdown,
		},
	)
}

func sendDeadMessage(ctx context.Context, b *bot.Bot, w interfaces.Worker, update *models.Update) {
	b.SendMessage(
		ctx, &bot.SendMessageParams{
			Text: fmt.Sprintf(
				"```\n%s```", PrettyJson(
					map[string]interface{}{
						"current_feels": w.CurrentFeels(),
						"status":        "dead",
					},
				),
			),
			ChatID:    update.Message.Chat.ID,
			ParseMode: models.ParseModeMarkdown,
		},
	)
}

func sendWaitMessage(ctx context.Context, b *bot.Bot, w interfaces.Worker, update *models.Update) {
	b.SendMessage(
		ctx, &bot.SendMessageParams{
			Text: fmt.Sprintf(
				"```\n%s```", PrettyJson(
					map[string]interface{}{
						"status": "doing something",
					},
				),
			),
			ChatID:    update.Message.Chat.ID,
			ParseMode: models.ParseModeMarkdown,
		},
	)
}

func sendBusyMessage(ctx context.Context, b *bot.Bot, w interfaces.Worker, update *models.Update) {
	b.SendMessage(
		ctx, &bot.SendMessageParams{
			Text: fmt.Sprintf(
				"```\n%s```", PrettyJson(
					map[string]interface{}{
						"status": "busy",
					},
				),
			),
			ChatID:    update.Message.Chat.ID,
			ParseMode: models.ParseModeMarkdown,
		},
	)
}

func PrettyJson(body map[string]interface{}) string {
	bytes, err := json.MarshalIndent(body, "", "\t")
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}
