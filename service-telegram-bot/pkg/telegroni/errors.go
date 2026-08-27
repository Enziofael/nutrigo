package telegroni

import "fmt"

type BotError struct {
	Message string
	Err     *BotError
}

func (e *BotError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s:\n%v",
			e.Message, e.Err)
	}
	return e.Message
}

func (e *BotError) Unwrap() error {
	return e.Err
}

func NewBotError(s string) *BotError {
	return &BotError{
		Message: s,
		Err:     nil,
	}
}

func NewBotErrorf(format string, a ...any) *BotError {
	errMessage := fmt.Errorf(format, a...)
	return &BotError{
		Message: errMessage.Error(),
		Err:     nil,
	}
}

func NewBotErrorw(s string, inner *BotError) *BotError {
	return &BotError{
		Message: s,
		Err:     inner,
	}
}
