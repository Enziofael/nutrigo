package telegroni

import "fmt"

type BotError struct {
	Message string
	Err     error
}

func (e *BotError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("BOT error: message=%s, error=%v",
			e.Message, e.Err)
	}
	return fmt.Sprintf("BOT error: message=%s",
		e.Message)
}

func (e *BotError) Unwrap() error {
	return e.Err
}
