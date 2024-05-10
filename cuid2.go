package cuid2

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
)

const (
	defaultLength   = 24
	bigLength       = 32
	initialCountMax = 476782367
)

var defaultCounter = createCounter(rand.Intn(initialCountMax))
var defaultFingerprint = createFingerprint()
var defaultAlphabet = createAlphabet()

type GenOpts struct {
	length      *int
	fingerprint *string
	alphabet    []rune
	counter     *func() int
}

type Gen struct {
	fingerprint string
	alphabet    []rune
	counter     Counter
	length      int
}

func New(opts GenOpts) Gen {
	var alphabet []rune
	var fingerprint string
	var length int
	var counter func() int

	if opts.alphabet == nil {
		alphabet = defaultAlphabet
	} else {
		alphabet = opts.alphabet
	}

	if opts.fingerprint == nil {
		fingerprint = defaultFingerprint
	} else {
		fingerprint = *opts.fingerprint
	}

	if opts.length == nil {
		length = defaultLength
	} else {
		length = *opts.length
	}

	if opts.counter == nil {
		counter = defaultCounter
	} else {
		counter = *opts.counter
	}

	return Gen{alphabet: alphabet, fingerprint: fingerprint, counter: counter, length: length}
}

func (gen Gen) Cuid2() string {
	firstLetter := string(randomLetter(gen.alphabet))
	time := base36Encode(time.Now().UnixMilli())
	count := base36Encode(gen.counter())
	salt := createEntropy(gen.length)
	hashInput := time + salt + count + gen.fingerprint

	return strings.ToLower(firstLetter + hash(hashInput)[1:gen.length])
}

func IsCuid(id string, minLength, maxLength int) bool {
	length := len(id)

	match, _ := regexp.MatchString("^[0-9a-z]+$", id)

	if length >= minLength && length <= maxLength && match {
		return true
	}

	return false
}
