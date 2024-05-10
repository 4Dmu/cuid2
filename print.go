package cuid2

import "math/rand"

func createFingerprint() string {
	source := createEntropy(bigLength)
	return source[:bigLength]
}

func createEntropy(length int) string {
	var entropy string

	for len(entropy) < length {
		val := rand.Intn(36)
		entropy = entropy + base36Encode(val)
	}

	return entropy
}
