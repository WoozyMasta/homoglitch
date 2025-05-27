package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

//go:generate rm -f homo*.go
//go:generate go run generate.go glyph.txt

const zwsp = '\u200B' // zero-width space

func main() {
	if len(os.Args) > 1 {
		for _, fname := range os.Args[1:] {
			f, err := os.Open(fname)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				continue
			}
			processInput(bufio.NewScanner(f))
			f.Close()
		}
	} else {
		processInput(bufio.NewScanner(os.Stdin))
	}
}

// injectZeroWidth randomly injects 1–3 zero-width spaces at 2–5 positions in the word
func injectZeroWidth(s string) string {
	runes := []rune(s)
	length := len(runes)

	if length < 3 {
		return s // too short
	}
	count := rand.Intn(4) + 2
	used := map[int]bool{}
	for i := 0; i < count; i++ {
		pos := rand.Intn(length-2) + 1 // not start and end
		if used[pos] {
			continue
		}
		used[pos] = true

		repeat := rand.Intn(3) + 1
		zw := []rune{}
		for j := 0; j < repeat; j++ {
			zw = append(zw, zwsp)
		}

		runes = append(runes[:pos], append(zw, runes[pos:]...)...)
	}

	return string(runes)
}

func obfuscate(text string) string {
	var result []rune
	for _, r := range text {
		if variants, ok := homoglyphMap[r]; ok && len(variants) > 0 {
			result = append(result, variants[rand.Intn(len(variants))])
		} else {
			result = append(result, r)
		}
	}

	words := strings.Fields(string(result))
	for i, word := range words {
		words[i] = injectZeroWidth(word)
	}

	return strings.Join(words, " ")
}

func processInput(scanner *bufio.Scanner) {
	for scanner.Scan() {
		fmt.Println(obfuscate(scanner.Text()))
	}
}
