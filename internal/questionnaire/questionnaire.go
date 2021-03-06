package questionnaire

import (
	"errors"
	"fmt"

	tgApi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/Andrianov/psychoHelpBOT/internal/config"
	"github.com/Andrianov/psychoHelpBOT/internal/models"
	"github.com/Andrianov/psychoHelpBOT/internal/storage"
)

type Questionnaire struct {
	cfg     *config.Config
	bot     *tgApi.BotAPI
	storage storage.Storage
}

var index = 24

func New(cfg *config.Config, bot *tgApi.BotAPI, storage storage.Storage) *Questionnaire {
	return &Questionnaire{cfg, bot, storage}
}

func (q *Questionnaire) Intro(update tgApi.Update) error {
	if update.Message == nil {
		return errors.New("message is nil")
	}

	msg := tgApi.NewMessage(update.Message.Chat.ID, startMessage)
	msg.ParseMode = "markdown"
	_, err := q.bot.Send(msg)
	return err
}

func (q *Questionnaire) Cancel(update tgApi.Update) error {
	if update.Message == nil {
		return errors.New("message is nil")
	}

	chatID := update.Message.Chat.ID

	err := q.storage.Delete(chatID)
	if err != nil {
		return err
	}

	msg := tgApi.NewMessage(chatID, cancelMessage)
	_, err = q.bot.Send(msg)
	return err
}

func (q *Questionnaire) Start(update tgApi.Update) error {
	if update.Message == nil {
		return errors.New("message is nil")
	}

	chatID := update.Message.Chat.ID
	userName := update.Message.Chat.UserName

	err := q.storage.Delete(chatID)
	if err != nil {
		return err
	}

	flowSteps := FlowSteps
	if len(userName) == 0 {
		flowSteps = AnonymousFlowSteps
	}

	steps := make([]*models.Step, 0, len(flowSteps))
	for _, step := range flowSteps {
		step := step
		steps = append(steps, &step)
	}

	chat := models.Chat{
		ID:       chatID,
		UserName: userName,
		Flow: &models.Flow{
			Steps: steps,
		},
	}
	err = q.storage.Save(chat)
	if err != nil {
		return err
	}

	return q.next(chat)
}

func (q *Questionnaire) Continue(update tgApi.Update) error {
	var chatID int64
	if update.Message != nil {
		chatID = update.Message.Chat.ID
	} else if update.CallbackQuery != nil && update.CallbackQuery.Message != nil {
		chatID = update.CallbackQuery.Message.Chat.ID
	} else {
		return errors.New("failed to parse chatID")
	}

	chat, err := q.storage.Get(chatID)
	if err != nil {
		if errors.Is(err, storage.ErrChatNotFound) {
			fmt.Println("try to continue chat that not found", chatID, update)
			return nil
		} else {
			return err
		}
	}

	// continue finished chat
	if chat.Flow.IsFinished() {
		return nil
	}

	err = q.saveAnswer(&chat, update)
	if err != nil {
		return err
	}

	if chat.Flow.IsFinished() {
		return q.finish(chat)
	}

	return q.next(chat)
}

func (q *Questionnaire) next(chat models.Chat) error {
	step := chat.Flow.NextStep()
	if step == nil {
		return errors.New("failed to find next step")
	}

	msg := tgApi.NewMessage(chat.ID, step.Question)

	if len(step.Options) != 0 {
		keyboard := tgApi.InlineKeyboardMarkup{}
		for _, option := range step.Options {
			btn := tgApi.NewInlineKeyboardButtonData(option, option)
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tgApi.InlineKeyboardButton{btn})
		}

		msg.ReplyMarkup = keyboard
	}

	_, err := q.bot.Send(msg)
	return err
}

func (q *Questionnaire) saveAnswer(chat *models.Chat, update tgApi.Update) error {
	step := chat.Flow.NextStep()
	if step == nil {
		return errors.New("failed to find next step")
	}

	if len(step.Options) != 0 {
		if update.CallbackQuery != nil {
			step.Answer = update.CallbackQuery.Data
		} else {
			return nil
		}
	} else {
		if update.Message != nil {
			step.Answer = update.Message.Text
		} else {
			return nil
		}
	}

	return q.storage.Save(*chat)
}

func (q *Questionnaire) finish(chat models.Chat) error {
	msg := tgApi.NewMessage(chat.ID, finalMessage)
	_, err := q.bot.Send(msg)
	if err != nil {
		return err
	}

	index++

	text := fmt.Sprintf("%d Новая заявка от @%s !\n", index, chat.UserName)
	for _, step := range chat.Flow.Steps {
		text += fmt.Sprintf("%s: %s\n", step.Name, step.Answer)
	}

	msg = tgApi.NewMessage(q.cfg.MainChatID, text)
	_, err = q.bot.Send(msg)
	if err != nil {
		return err
	}

	msg = tgApi.NewMessage(q.cfg.TechChatID, text)
	_, err = q.bot.Send(msg)
	return err
}
