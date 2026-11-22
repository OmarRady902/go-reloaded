package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func capitalize(s string) string {
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

func applyTransform(res []string, tag string, n int) {
	for i := 1; i <= n; i++ {
		idx := len(res) - i
		if idx < 0 {
			break
		}
		switch tag {
		case "up":
			res[idx] = strings.ToUpper(res[idx])
		case "low":
			res[idx] = strings.ToLower(res[idx])
		case "cap":
			res[idx] = capitalize(res[idx])
		}
	}
}

func processWords(words []string) []string {
	var res []string
	for i := 0; i < len(words); i++ {
		w := words[i]
		switch w {
		case "(hex)":
			if len(res) > 0 {
				val, err := strconv.ParseInt(res[len(res)-1], 16, 64)
				if err == nil {
					res[len(res)-1] = strconv.FormatInt(val, 10)
				}
			}
		case "(bin)":
			if len(res) > 0 {
				val, err := strconv.ParseInt(res[len(res)-1], 2, 64)
				if err == nil {
					res[len(res)-1] = strconv.FormatInt(val, 10)
				}
			}
		case "(up)":
			if len(res) > 0 {
				res[len(res)-1] = strings.ToUpper(res[len(res)-1])
			}
		case "(low)":
			if len(res) > 0 {
				res[len(res)-1] = strings.ToLower(res[len(res)-1])
			}
		case "(cap)":
			if len(res) > 0 {
				res[len(res)-1] = capitalize(res[len(res)-1])
			}
		case "(up,", "(low,", "(cap,":
			if i+1 < len(words) {
				tag := strings.Trim(w, "(,")
				nStr := strings.TrimSuffix(words[i+1], ")")
				n, err := strconv.Atoi(nStr)
				if err == nil {
					applyTransform(res, tag, n)
				}
				i++
			}
		default:
			res = append(res, w)
		}
	}
	return res
}

func handleArticles(words []string) {
	for i := 0; i < len(words)-1; i++ {
		low := strings.ToLower(words[i])
		if low != "a" && low != "an" {
			continue
		}

		next := strings.ToLower(words[i+1])
		if len(next) == 0 {
			continue
		}

		isVowel := strings.ContainsAny(string(next[0]), "aeiouh")
		if low == "a" && isVowel {
			if words[i] == "A" {
				words[i] = "An"
			} else {
				words[i] = "an"
			}
		} else if low == "an" && !isVowel {
			if words[i] == "An" {
				words[i] = "A"
			} else {
				words[i] = "a"
			}
		}
	}
}

func fixPunctuation(words []string) []string {
	var splitPunct []string
	for _, w := range words {
		currentWord := ""
		for _, r := range w {
			if strings.ContainsRune(".,!?:;", r) {
				if currentWord != "" {
					splitPunct = append(splitPunct, currentWord)
					currentWord = ""
				}
				splitPunct = append(splitPunct, string(r))
			} else {
				currentWord += string(r)
			}
		}
		if currentWord != "" {
			splitPunct = append(splitPunct, currentWord)
		}
	}

	var finalWords []string
	for _, w := range splitPunct {
		if strings.ContainsAny(w, ".,!?:;") && len(finalWords) > 0 {
			finalWords[len(finalWords)-1] += w
		} else {
			finalWords = append(finalWords, w)
		}
	}
	return finalWords
}

func fixQuotes(words []string) []string {
	var quoteFixed []string
	isOpen := false
	for _, word := range words {
		if word == "'" {
			if !isOpen {
				quoteFixed = append(quoteFixed, "'")
				isOpen = true
			} else {
				if len(quoteFixed) > 0 {
					quoteFixed[len(quoteFixed)-1] += "'"
				}
				isOpen = false
			}
		} else {
			if isOpen && len(quoteFixed) > 0 && quoteFixed[len(quoteFixed)-1] == "'" {
				quoteFixed[len(quoteFixed)-1] += word
			} else {
				quoteFixed = append(quoteFixed, word)
			}
		}
	}
	return quoteFixed
}

func fixLine(line string) string {
	if strings.TrimSpace(line) == "" {
		return ""
	}

	words := strings.Fields(line)
	processed := processWords(words)
	handleArticles(processed)
	punctuated := fixPunctuation(processed)
	final := fixQuotes(punctuated)

	return strings.Join(final, " ")
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . <input.txt> <output.txt>")
		os.Exit(1)
	}

	input, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("error of reading file:", err)
		os.Exit(1)
	}

	lines := strings.Split(string(input), "\n")
	var resultLines []string

	for _, line := range lines {
		resultLines = append(resultLines, fixLine(line))
	}

	output := strings.Join(resultLines, "\n")
	err = os.WriteFile(os.Args[2], []byte(output), 0o644)
	if err != nil {
		fmt.Println("error of writing file:", err)
		os.Exit(1)
	}
}
