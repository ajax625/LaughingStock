package notify

import (
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type TelegramBot struct {
	bot *tgbotapi.BotAPI
	log *zap.Logger
}

func NewTelegramBot(token string, log *zap.Logger) (*TelegramBot, error) {
	if token == "" {
		return nil, nil // Telegram is optional
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TelegramBot{bot: bot, log: log}, nil
}

type SignalPayload struct {
	Symbol     string  `json:"symbol"`
	Direction  string  `json:"direction"`
	Confidence float64 `json:"confidence"`
	Long       *struct {
		Entry  float64 `json:"entry"`
		Target float64 `json:"target"`
		Stop   float64 `json:"stop"`
		ROI    float64 `json:"roi"`
		RR     float64 `json:"rr"`
	} `json:"long,omitempty"`
	Short *struct {
		Entry  float64 `json:"entry"`
		Target float64 `json:"target"`
		Stop   float64 `json:"stop"`
		ROI    float64 `json:"roi"`
		RR     float64 `json:"rr"`
	} `json:"short,omitempty"`
}

// SendSignalAlert sends a signal notification with an Execute inline button.
func (t *TelegramBot) SendSignalAlert(chatID int64, payload []byte) error {
	if t == nil {
		return nil
	}

	var sig SignalPayload
	if err := json.Unmarshal(payload, &sig); err != nil {
		return err
	}

	text := fmt.Sprintf(
		"🚨 *TradeX Signal*\n\nSymbol: `%s`\nDirection: *%s*\nConfidence: %.0f%%",
		sig.Symbol, sig.Direction, sig.Confidence*100,
	)

	if sig.Direction == "LONG" && sig.Long != nil {
		text += fmt.Sprintf(
			"\n\nEntry: `$%.2f`\nTarget: `$%.2f`\nStop: `$%.2f`\nROI: `%.1f%%`\nR/R: `%.1f`",
			sig.Long.Entry, sig.Long.Target, sig.Long.Stop, sig.Long.ROI, sig.Long.RR,
		)
	} else if sig.Direction == "SHORT" && sig.Short != nil {
		text += fmt.Sprintf(
			"\n\nEntry: `$%.2f`\nTarget: `$%.2f`\nStop: `$%.2f`\nROI: `%.1f%%`\nR/R: `%.1f`",
			sig.Short.Entry, sig.Short.Target, sig.Short.Stop, sig.Short.ROI, sig.Short.RR,
		)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Execute Trade",
				fmt.Sprintf("execute:%s:%s", sig.Symbol, sig.Direction),
			),
		),
	)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = keyboard

	_, err := t.bot.Send(msg)
	if err != nil {
		t.log.Error("telegram: send failed", zap.Int64("chat_id", chatID), zap.Error(err))
	}
	return err
}

// SendMessage sends a plain text message to a chat.
func (t *TelegramBot) SendMessage(chatID int64, text string) {
	if t == nil {
		return
	}
	msg := tgbotapi.NewMessage(chatID, text)
	t.bot.Send(msg)
}

// StartCallbackListener polls Telegram for inline button presses and calls onExecute.
// Runs until ctx is cancelled — call in a goroutine.
func (t *TelegramBot) StartCallbackListener(onExecute func(chatID int64, symbol, direction string)) {
	if t == nil {
		return
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := t.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery == nil {
			continue
		}
		var symbol, direction string
		fmt.Sscanf(update.CallbackQuery.Data, "execute:%s:%s", &symbol, &direction)
		if symbol == "" {
			continue
		}
		chatID := update.CallbackQuery.Message.Chat.ID
		t.bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, "Executing..."))
		go onExecute(chatID, symbol, direction)
	}
}
