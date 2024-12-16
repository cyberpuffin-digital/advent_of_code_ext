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
	wasInvalid   bool
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
		u.wasInvalid = false
	}
}

func (u *Update) string() (msg string) {
	for _, page := range u.pages {
		msg += fmt.Sprintf("%d,", page)
	}
	msg = strings.TrimSuffix(msg, ",")
	return msg
}

func (u *Update) validate(pages []int, rules []*Rule) bool {
	u.validOrder = true

	for ruleIndex, rule := range rules {
		// Skip rules that don't apply
		if !slices.Contains(pages, rule.before) || !slices.Contains(pages, rule.after) {
			continue
		}

		beforeIndex := slices.Index(pages, rule.before)
		afterIndex := slices.Index(pages, rule.after)

		if beforeIndex < afterIndex {
			verbosef("\t%q passes rule %d (%q).\n", u.string(), ruleIndex, rule.string())
			continue
		}

		u.wasInvalid = true
		verbosef("\t%q fails rule %d (%q).\n", u.string(), ruleIndex, rule.string())

		swap := pages[beforeIndex]
		pages[beforeIndex] = pages[afterIndex]
		pages[afterIndex] = swap

		u.validOrder = u.validate(pages, rules)
	}

	u.pages = pages
	return u.validOrder
}
