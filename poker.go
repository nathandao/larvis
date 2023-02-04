package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/exp/slices"
)

type result string
type comboType string

const (
	tie       result = "Tie"
	hand1Wins result = "Hand 1"
	hand2Wins result = "Hand 2"

	fullHouse   comboType = "FullHouse"
	fourOfAKind comboType = "FourOfAKind"
	twoPairs    comboType = "TwoPairs"
	triple      comboType = "Triple"
	onePair     comboType = "OnePair"
	highCards   comboType = "HighCards"

	// Regex to match 1 or more repeating cards. Assuming cards are uppercase and sorted.
	comboRegex = `2{1,}|3{1,}|4{1,}|5{1,}|6{1,}|7{1,}|8{1,}|9{1,}|T{1,}|J{1,}|Q{1,}|K{1,}|A{1,}`
)

var (
	cardRank  = []string{"2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K", "A"}
	comboRank = []comboType{highCards, onePair, twoPairs, triple, fullHouse, fourOfAKind}
)

// getComboRank returns the index of the combo type in the comboRank slice. Higher index means a stronger combo.
func getComboRank(combo comboType) int {
	return slices.Index(comboRank, combo)
}

// getCardRank returns the index of the index of the card int the cardRank slice. Higher index means a stronger card.
func getCardRank(card string) int {
	return slices.Index(cardRank, card)
}

// validateHand returns any error related to the poker hand input.
func validateHand(hand string) error {
	// Make sure to work with only uppercase chars
	hand = sanitizeHand(hand)

	if len(hand) != 5 {
		return fmt.Errorf("poker hand must have 5 cards")
	}

	invalidCards := regexp.MustCompile(`[^2-9TJQKA]`).FindAllString(hand, -1)
	if len(invalidCards) > 0 {
		// Get the unique set of invalid cards
		sort.Strings(invalidCards)
		uniqueInvalidCards := slices.Compact(invalidCards)

		if len(uniqueInvalidCards) == 1 {
			return fmt.Errorf("\"%s\" is not a valid card", invalidCards[0])
		}

		return fmt.Errorf("\"%s\" are not valid cards", strings.Join(uniqueInvalidCards, ", "))
	}

	return nil
}

// sanitizeHand returns the poker hand as uppercased, and sorted from higher to lower ranked cards.
func sanitizeHand(hand string) string {
	hand = strings.ReplaceAll(strings.ToUpper(hand), " ", "")
	cards := strings.Split(hand, "")
	sort.Slice(cards, func(i, j int) bool {
		return getCardRank(cards[i]) > getCardRank(cards[j])
	})
	return strings.Join(cards, "")
}

// getSortedComponents returns a sorted slice of the poker hand components.
// Each component is either a string with of same card, or a single card.
// The slice of components is sorted descendingly, first by repeated card length, and then by card rank
// if the repeat lengths are equal.
// Examples: []string{"AA", "44", "6"} or []string{"222", "AA"} or []string{"A", "Q", "T", "4", "2"}
func getSortedComponents(hand string) []string {
	components := regexp.MustCompile(comboRegex).FindAllString(hand, -1)

	sort.Slice(components, func(i, j int) bool {
		if len(components[i]) != len(components[j]) {
			return len(components[i]) > len(components[j])
		}

		return getCardRank(components[i][:1]) > getCardRank(components[j][:1])
	})

	return components
}

// getComboType returns the hand's combyType
func getComboType(hand string) comboType {
	components := getSortedComponents(hand)
	comboCount := map[int]int{4: 0, 3: 0, 2: 0, 1: 0}

	for _, component := range components {
		comboCount[len(component)] += 1
	}

	switch {
	case comboCount[4] == 1:
		return fourOfAKind
	case comboCount[3] == 1 && comboCount[2] == 1:
		return fullHouse
	case comboCount[3] == 1:
		return triple
	case comboCount[2] == 2:
		return twoPairs
	case comboCount[2] == 1:
		return onePair
	default:
		return highCards
	}
}

// calculateResult returns the result of the winner or if it's a tie.
func calculateResult(firstHand string, secondHand string) result {
	firstHand = sanitizeHand(firstHand)
	secondHand = sanitizeHand(secondHand)

	firstHandComboRank := getComboRank(getComboType(firstHand))
	secondHandComboRank := getComboRank(getComboType(secondHand))

	if firstHandComboRank > secondHandComboRank {
		return hand1Wins
	}

	if firstHandComboRank < secondHandComboRank {
		return hand2Wins
	}

	// At this point, both hands should have the same combo type.
	// We'll need to compare the "sorted components" to find the winner
	firstHandComponents := getSortedComponents(firstHand)
	secondHandComponents := getSortedComponents(secondHand)

	// Since both hands have the same combo type, and the components slice are aready sorted by rank, the winner is then
	// decided by iterating through the ranked component slice, and whoever has the stronger component first is the winner.
	for i, component := range firstHandComponents {
		// Both components are strings of repeated cards of the same length (or a single card), compare them using the first character.
		component1Rank := slices.Index(cardRank, component[:1])
		component2Rank := slices.Index(cardRank, secondHandComponents[i][:1])

		if component1Rank == component2Rank {
			continue
		}

		if component1Rank > component2Rank {
			return hand1Wins
		}

		return hand2Wins
	}

	// It's a tie of all sorted components are equal
	return tie
}
