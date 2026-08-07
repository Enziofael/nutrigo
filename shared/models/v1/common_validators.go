package models

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func SanitizeID(id int64) (int64, error) {
	if id > 0 {
		return id, nil
	}
	return 0, fmt.Errorf("Negative id")
}

func SanitizeTgID(tgID int64) (int64, error) {
	if tgID > 0 {
		return tgID, nil
	}
	return 0, fmt.Errorf("Negative tg_id")
}

func SanitizeName(name string) (string, error) {
	clean := strings.Join(strings.Fields(name), " ")
	if clean == "" {
		return "", fmt.Errorf("Name can't be empty")
	}
	if len(clean) > MaxNameLength {
		return "", fmt.Errorf("Name is too long. Max length is %d", MaxNameLength)
	}
	return clean, nil
}

func SanitizeDescription(Desc string) (string, error) {
	clean := strings.TrimSpace(Desc)
	if len(clean) > MaxDescriptionLength {
		return "", fmt.Errorf("Description is too long. Max length is %d", MaxDescriptionLength)
	}
	return clean, nil
}

func SanitizeOrder(order Order) (Order, error) {
	switch strings.ToUpper(string(order)) {
	case "ASC":
		return OrderASC, nil
	case "DESC":
		return OrderDESC, nil
	default:
		return OrderASC, fmt.Errorf("Invalid order: %s", order)
	}
}

func SanitizeOffset(offset int) (int, error) {
	if offset >= 0 {
		return offset, nil
	}
	return 0, fmt.Errorf("Offset can't be negative: %d", offset)
}

func SanitizeLimit(limit int) (int, error) {
	if limit > 0 && limit <= 100 {
		return limit, nil
	}
	return 100, fmt.Errorf("Offset should be > 0 and <= 100: %d", limit)
}

func SanitizeUrl(raw string) (string, error) {
	cleaned := strings.TrimSpace(raw)
	original := cleaned

	if !strings.Contains(original, "://") {
		cleaned = "http://" + cleaned
	}

	parsed, err := url.Parse(cleaned)
	if err != nil {
		return "", errors.New("Invalid link format: " + err.Error())
	}
	if parsed.Host == "" {
		return "", errors.New("Missing host in link")
	}

	if !strings.Contains(original, "://") {
		result := parsed.Host
		if parsed.Path != "" {
			result += parsed.Path
		}
		if parsed.RawQuery != "" {
			result += "?" + parsed.RawQuery
		}
		if parsed.Fragment != "" {
			result += "#" + parsed.Fragment
		}
		return result, nil
	}
	return parsed.String(), nil
}
