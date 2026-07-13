// ./frontend-bot/internal/actions/mediaMessage.go

// TODO:
// Разнести в отдельные файлы для каждого типа? Потом если что

package actions

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ====================================================================
// PHOTO
// ====================================================================

func SendPhoto(bot *tgbotapi.BotAPI, config tgbotapi.PhotoConfig) {
	if _, err := bot.Send(config); err != nil {
		HandleSendError(bot, config)
	}
}

func ReplacePhoto(bot *tgbotapi.BotAPI, config tgbotapi.EditMessageMediaConfig) {
	if _, err := bot.Request(config); err != nil {
		HandleRequestError(bot, config)
	}
}

func EditPhotoCaption(bot *tgbotapi.BotAPI, config tgbotapi.EditMessageCaptionConfig) {
	if _, err := bot.Request(config); err != nil {
		HandleRequestError(bot, config)
	}
}

// ====================================================================
// VIDEO
// ====================================================================

func SendVideo(bot *tgbotapi.BotAPI, config tgbotapi.VideoConfig) {
	if _, err := bot.Send(config); err != nil {
		HandleSendError(bot, config)
	}
}

func ReplaceVideo(bot *tgbotapi.BotAPI, config tgbotapi.EditMessageMediaConfig) {
	if _, err := bot.Request(config); err != nil {
		HandleRequestError(bot, config)
	}
}

func EditVideoCaption(bot *tgbotapi.BotAPI, config tgbotapi.EditMessageCaptionConfig) {
	if _, err := bot.Request(config); err != nil {
		HandleRequestError(bot, config)
	}
}

// ====================================================================
// DOCUMENT
// ====================================================================

func SendDocument(bot *tgbotapi.BotAPI, config tgbotapi.DocumentConfig) {
	if _, err := bot.Send(config); err != nil {
		HandleSendError(bot, config)
	}
}

func ReplaceDocument(bot *tgbotapi.BotAPI, config tgbotapi.EditMessageMediaConfig) {
	if _, err := bot.Request(config); err != nil {
		HandleRequestError(bot, config)
	}
}

func EditDocumentCaption(bot *tgbotapi.BotAPI, config tgbotapi.EditMessageCaptionConfig) {
	if _, err := bot.Request(config); err != nil {
		HandleRequestError(bot, config)
	}
}

// ====================================================================
// VIDEO MESSAGE
// ====================================================================

func SendVideoMessage(bot *tgbotapi.BotAPI, config tgbotapi.VideoNoteConfig) {
	if _, err := bot.Send(config); err != nil {
		HandleSendError(bot, config)
	}
}

func ReplaceVideoMessage(bot *tgbotapi.BotAPI, config tgbotapi.EditMessageMediaConfig) {
	if _, err := bot.Request(config); err != nil {
		HandleRequestError(bot, config)
	}
}

// Video message has no caption

// ====================================================================
// AUDIO
// ====================================================================

func SendAudio(bot *tgbotapi.BotAPI, config tgbotapi.AudioConfig) {
	if _, err := bot.Send(config); err != nil {
		HandleSendError(bot, config)
	}
}

func ReplaceAudio(bot *tgbotapi.BotAPI, config tgbotapi.EditMessageMediaConfig) {
	if _, err := bot.Request(config); err != nil {
		HandleRequestError(bot, config)
	}
}

func EditAudioCaption(bot *tgbotapi.BotAPI, config tgbotapi.EditMessageCaptionConfig) {
	if _, err := bot.Request(config); err != nil {
		HandleRequestError(bot, config)
	}
}

// ====================================================================
// VOICE
// ====================================================================

func SendVoice(bot *tgbotapi.BotAPI, config tgbotapi.VoiceConfig) {
	if _, err := bot.Send(config); err != nil {
		HandleSendError(bot, config)
	}
}

func ReplaceVoice(bot *tgbotapi.BotAPI, config tgbotapi.EditMessageMediaConfig) {
	if _, err := bot.Request(config); err != nil {
		HandleRequestError(bot, config)
	}
}

// Voice has no caption

// ====================================================================
// LOCATION
// ====================================================================

func SendLocation(bot *tgbotapi.BotAPI, config tgbotapi.LocationConfig) {
	if _, err := bot.Send(config); err != nil {
		HandleSendError(bot, config)
	}
}

// Location can't be edited. Use live location instead

// Location has no caption

// ====================================================================
// LIVE LOCATION
// ====================================================================

/* TODO: SendLiveLocation() Реализация через отправку с каналом обновлений, досрочной остановкой и продлением
Сделать позже в более высокоуровневых слоях */

func UpdateLiveLocation(bot *tgbotapi.BotAPI, config tgbotapi.EditMessageLiveLocationConfig) {
	if _, err := bot.Request(config); err != nil {
		HandleRequestError(bot, config)
	}
}

// Live location has no caption

// ====================================================================
// ANIMATION (GIF)
// ====================================================================

func SendAnimation(bot *tgbotapi.BotAPI, config tgbotapi.AnimationConfig) {
	if _, err := bot.Send(config); err != nil {
		HandleSendError(bot, config)
	}
}

func ReplaceAnimation(bot *tgbotapi.BotAPI, config tgbotapi.EditMessageMediaConfig) {
	if _, err := bot.Request(config); err != nil {
		HandleRequestError(bot, config)
	}
}

func EditAnimationCaption(bot *tgbotapi.BotAPI, config tgbotapi.EditMessageCaptionConfig) {
	if _, err := bot.Request(config); err != nil {
		HandleRequestError(bot, config)
	}
}

// ====================================================================
// STICKER
// ====================================================================

func SendSticker(bot *tgbotapi.BotAPI, config tgbotapi.StickerConfig) {
	if _, err := bot.Send(config); err != nil {
		HandleSendError(bot, config)
	}
}

func ReplaceSticker(bot *tgbotapi.BotAPI, config tgbotapi.EditMessageMediaConfig) {
	if _, err := bot.Request(config); err != nil {
		HandleRequestError(bot, config)
	}
}

// Sticker has no caption

// ====================================================================
// CONTACT
// ====================================================================

func SendContact(bot *tgbotapi.BotAPI, config tgbotapi.ContactConfig) {
	if _, err := bot.Send(config); err != nil {
		HandleSendError(bot, config)
	}
}

// Contact can't be edited

// Contact has no caption

// ====================================================================
// MEDIA GROUP
// ====================================================================

// Either photos, videos, documents, or audio SEPARATELY (!)
func SendMediaGroup(bot *tgbotapi.BotAPI, config tgbotapi.MediaGroupConfig) {
	if _, err := bot.Send(config); err != nil {
		HandleSendError(bot, config)
	}
}

// Should be the same type as the original one
func ReplaceMediaGroup(bot *tgbotapi.BotAPI, config tgbotapi.EditMessageMediaConfig) {
	if _, err := bot.Request(config); err != nil {
		HandleRequestError(bot, config)
	}
}

func EditMediaGroupCaption(bot *tgbotapi.BotAPI, config tgbotapi.EditMessageCaptionConfig) {
	if _, err := bot.Request(config); err != nil {
		HandleRequestError(bot, config)
	}
}