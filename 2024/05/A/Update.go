package main

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
)

type Update struct {
	pages        []int
	orderChecked bool
	validOrder   bool
}

// init parses input line to populate Update
func (u *Update) init(inputLine string) {
	u.orderChecked = false
	u.validOrder = true

	pagesString := strings.Split(inputLine, ",")
	for _, page := range pagesString {
		p, err := strconv.Atoi(page)
		if err != nil {
			log.Fatalf("Failed to parse page %q.\n", page)
		}
		u.pages = append(u.pages, p)
	}
}

func (u *Update) string() (msg string) {
	for _, page := range u.pages {
		msg += fmt.Sprintf("%d,", page)
	}
	msg = strings.TrimSuffix(msg, ",")
	return msg
}

func (u *Update) validate(rules []*Rule) bool {
	for _, rule := range rules {
		// Broke a rule
		if u.orderChecked && !u.validOrder {
			continue
		}

		// Check if rule applies
		if slices.Contains(u.pages, rule.before) && slices.Contains(u.pages, rule.after) {
			beforeIndex := slices.Index(u.pages, rule.before)
			afterIndex := slices.Index(u.pages, rule.after)

			verbosef(
				"\tRule %q applies.\n\t\tbefore %d; after %d\n\t\tcurrent valid: %t\n\t\tnew valid: %t\n",
				rule.string(),
				beforeIndex,
				afterIndex,
				u.validOrder,
				beforeIndex < afterIndex,
			)
			u.validOrder = beforeIndex < afterIndex
			u.orderChecked = true
		}
	}

	return u.validOrder
}
