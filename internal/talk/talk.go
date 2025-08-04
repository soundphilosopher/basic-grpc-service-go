package talk

import (
	"fmt"
	"math/rand"
	"strings"
)

func preprocess(input string) string {
	return strings.Trim(strings.ToLower(strings.TrimSpace(input)), `.!?'"`)
}

func randomElementFrom(list []string) string {
	return list[rand.Intn(len(list))] //nolint:gosec
}

func reflect(fragment string) string {
	words := strings.Fields(fragment)
	for i, word := range words {
		if reflectedWord, ok := reflectedWords[word]; ok {
			words[i] = reflectedWord
		}
	}
	return strings.Join(words, " ")
}

func lookupResponse(input string) string {
	for re, responses := range requestInputRegexToResponseOptions {
		matches := re.FindStringSubmatch(input)
		if len(matches) < 1 {
			continue
		}
		response := randomElementFrom(responses)
		// If the response has an entry point, reflect the input phrase (so "I"
		// becomes "you").
		if !strings.Contains(response, "%s") {
			return response
		}
		if len(matches) > 1 {
			fragment := reflect(matches[1])
			response = fmt.Sprintf(response, fragment)
			return response
		}
	}
	return randomElementFrom(defaultResponses)
}

func Reply(input string) (string, bool) {
	input = preprocess(input)
	if _, ok := goodbyeInputSet[input]; ok {
		return randomElementFrom(goodbyeResponses), true
	}

	return lookupResponse(input), false
}

func GetIntroResponses(name string) []string {
	intros := make([]string, 0, len(introResponses)+2)
	for _, n := range introResponses {
		intros = append(intros, fmt.Sprintf(n, name))
	}

	intros = append(intros, randomElementFrom(facts))
	intros = append(intros, "How are you feeling today?")
	return intros
}
