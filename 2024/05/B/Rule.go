package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Rule struct {
	before int
	after  int
}

func (r *Rule) init(inputLine string) {
	var err error
	stringsIn := strings.Split(inputLine, "|")

	r.before, err = strconv.Atoi(stringsIn[0])
	if err != nil {
		log.Fatalf("Failed to parse rule before value %q", stringsIn[0])
	}

	r.after, err = strconv.Atoi(stringsIn[1])
	if err != nil {
		log.Fatalf("Failed to parse rule after value %q", stringsIn[1])
	}
}

func (r *Rule) string() string {
	return fmt.Sprintf("%d before %d", r.before, r.after)
}
