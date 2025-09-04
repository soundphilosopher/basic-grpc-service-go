// Package talk implements an ELIZA-like chatbot that provides conversational responses
// using pattern matching and reflection techniques.
package talk

import (
	"fmt"
	"math/rand"
	"strings"
)

// preprocess normalizes user input by converting to lowercase and removing
// leading/trailing whitespace and punctuation.
func preprocess(input string) string {
	return strings.Trim(strings.ToLower(strings.TrimSpace(input)), `.!?'"`)
}

// randomElementFrom returns a randomly selected element from the provided slice.
func randomElementFrom(list []string) string {
	return list[rand.Intn(len(list))] //nolint:gosec
}

// reflect converts personal pronouns in a text fragment from first person to
// second person and vice versa, creating a conversational reflection effect.
func reflect(fragment string) string {
	words := strings.Fields(fragment)
	for i, word := range words {
		if reflectedWord, ok := reflectedWords[word]; ok {
			words[i] = reflectedWord
		}
	}
	return strings.Join(words, " ")
}

// lookupResponse searches for a matching regex pattern in the input and returns
// an appropriate response. If a pattern matches and contains %s placeholders,
// the captured groups are reflected and inserted into the response template.
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

// Reply processes user input and returns a conversational response along with
// a boolean indicating whether the conversation should end. Returns true for
// goodbye phrases, false otherwise.
func Reply(input string) (string, bool) {
	input = preprocess(input)
	if _, ok := goodbyeInputSet[input]; ok {
		return randomElementFrom(goodbyeResponses), true
	}

	return lookupResponse(input), false
}

// GetIntroResponses generates an introduction sequence for a new conversation,
// including personalized greetings, a random fact about ELIZA, and an opening question.
func GetIntroResponses(name string) []string {
	intros := make([]string, 0, len(introResponses)+2)
	for _, n := range introResponses {
		intros = append(intros, fmt.Sprintf(n, name))
	}

	intros = append(intros, randomElementFrom(facts))
	intros = append(intros, "How are you feeling today?")
	return intros
}
