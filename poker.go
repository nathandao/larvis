package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/exp/slices"
)

type result string

const (
	tie       result = "Tie"
	hand1Wins result = "Hand 1"
	hand2Wins result = "Hand 2"
)

var cardRank = []string{"2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K", "A"}

// getCardRank returns the index of the index of the card int the cardRank slice. Higher index means a stronger card.
func getCardRank(card string) int {
	return slices.Index(cardRank, card)
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

// validateHand returns any error related to the poker hand input.
func validateHand(hand string) error {
	// Make sure to work with only uppercase chars
	hand = sanitizeHand(hand)

	if len(hand) != 5 {
		return fmt.Errorf("poker hand must have 5 cards")
	}

	invalidCards := regexp.MustCompile(`[^2-9TJQKA]`).FindAllString(hand, -1)
	if len(invalidCards) > 0 {
		sort.Strings(invalidCards)
		uniqueInvalidCards := slices.Compact(invalidCards)

		if len(uniqueInvalidCards) == 1 {
			return fmt.Errorf("\"%s\" is not a valid card", invalidCards[0])
		}

		return fmt.Errorf("\"%s\" are not valid cards", strings.Join(uniqueInvalidCards, ", "))
	}

	return nil
}

// getSortedComponents returns a sorted slice of the poker hand components.
// Each component is either a string of combo cards, or a single card.
// The slice of components is sorted descendingly, first by combo card length, and
// then by card rank if the combo lengths are equal.
// Examples: []string{"AA", "44", "6"} or []string{"222", "AA"} or []string{"A", "Q", "T", "4", "2"}
func getSortedComponents(hand string) []string {
	// Make sure cards are uppercased and sorted
	hand = sanitizeHand(hand)

	// Regex to match 1 or more repeating cards. Assuming cards are uppercase and sorted.
	comboRegex := `2{1,}|3{1,}|4{1,}|5{1,}|6{1,}|7{1,}|8{1,}|9{1,}|T{1,}|J{1,}|Q{1,}|K{1,}|A{1,}`
	components := regexp.MustCompile(comboRegex).FindAllString(hand, -1)

	sort.Slice(components, func(i, j int) bool {
		if len(components[i]) != len(components[j]) {
			return len(components[i]) > len(components[j])
		}

		return getCardRank(components[i][:1]) > getCardRank(components[j][:1])
	})

	return components
}

// calculateResult returns the result of the winner or if it's a tie.
func calculateResult(firstHand string, secondHand string) result {
	firstHandComponents := getSortedComponents(firstHand)
	secondHandComponents := getSortedComponents(secondHand)

	// Below is a list of all possible combinations and their respective
	// "sorted components". Each row is a set of identical cards.
	//
	//              4 of    full           two    one   high
	// Combo        a kind  house  triple  pairs  pair  card
	//=========================================================
	// Components   ||||    |||    |||     ||      ||     |
	//              |       ||     |       ||      |      |
	//         	               |       |       |      |
	//                                             |      |
	//                                                    |
	//
	// From this table, it's easy to see that there are 2 cases when comparing 2 hands:
	//
	// - 2 hands have different amount of components: then the hand with least row
	//   of components wins.
	//
	// - 2 hands have the same number of rows: Since component row within each hand
	//   are already sorted from stronger to weaker by `getSortedComponents`, all we
	//   need to do is compare rows from the same index with each other. The stronger
	//   row is one that is longer (has more idential cards), or in case of equal
	//   lengths, card rank can be compared. The iteration ends immediately once the
	//   a winning row is found, since all later rows can't affect the result.
	//   If no winning row is found when all rows were iterated, then it's a tie.

	// Compare by component length
	if len(firstHandComponents) < len(secondHandComponents) {
		return hand1Wins
	}

	if len(firstHandComponents) > len(secondHandComponents) {
		return hand2Wins
	}

	// If number of rows are equal, compare by rows
	for i, row1 := range firstHandComponents {
		row2 := secondHandComponents[i]

		// Check by row length
		if len(row1) > len(row2) {
			return hand1Wins
		}

		if len(row1) < len(row2) {
			return hand2Wins
		}

		// If both rows have the same number of cards, rank them by card.
		row1CardRank := slices.Index(cardRank, row1[:1])
		row2CardRank := slices.Index(cardRank, row2[:1])

		if row1CardRank == row2CardRank {
			continue
		}

		if row1CardRank > row2CardRank {
			return hand1Wins
		}

		return hand2Wins
	}

	// It's a tie if all components are equal
	return tie
}
