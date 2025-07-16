package utils

import (
	"errors"
	"regexp"
	"unicode/utf8"
)

func ValidateAdTitle(title string) error {
	if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 100 {
		return errors.New("title must be 3-100 characters")
	}
	return nil
}

func ValidateAdDescription(desc string) error {
	if utf8.RuneCountInString(desc) < 50 || utf8.RuneCountInString(desc) > 1000 {
		return errors.New("description must be 50-1000 characters")
	}
	return nil
}

func ValidateAdPrice(price float64) error {
	if price <= 0 || price > 10000000 {
		return errors.New("price must be > 0 and < 10_000_000")
	}
	return nil
}

func ValidateAdImageURL(url string) error {
	matched, _ := regexp.MatchString(`^https?://.+\.(jpg|jpeg|png|gif)$`, url)
	if !matched {
		return errors.New("image_url must be a valid URL ending with .jpg, .jpeg, .png, .gif")
	}
	return nil
}
