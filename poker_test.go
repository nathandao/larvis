package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateHand(t *testing.T) {
	testCases := []struct {
		hand        string
		expectError error
	}{
		{hand: "2345Qk", expectError: fmt.Errorf("poker hand must have 5 cards")},
		{hand: "2345", expectError: fmt.Errorf("poker hand must have 5 cards")},
		{hand: "tytjt", expectError: fmt.Errorf("\"Y\" is not a valid card")},
		{hand: "tvb1T", expectError: fmt.Errorf("\"1, B, V\" are not valid cards")},
		{hand: "7TaQK", expectError: nil},
	}

	for _, testCase := range testCases {
		t.Run(testCase.hand, func(t *testing.T) {
			err := validateHand(testCase.hand)
			if testCase.expectError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, testCase.expectError.Error())
			}
		})
	}
}

func TestSanitizeHand(t *testing.T) {
	testCases := []struct {
		hand           string
		expectedResult string
	}{
		{hand: "Ata78", expectedResult: "AAT87"},
		{hand: "99189", expectedResult: "99981"},
		{hand: "99 189 ", expectedResult: "99981"},
		{hand: "QjAkj", expectedResult: "AKQJJ"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.hand, func(t *testing.T) {
			assert.Equal(t, testCase.expectedResult, sanitizeHand(testCase.hand))
		})
	}
}

func TestGetSortedComponents(t *testing.T) {
	testCases := []struct {
		hand           string
		expectedResult []string
	}{
		{hand: "7TAJ2", expectedResult: []string{"A", "J", "T", "7", "2"}},
		{hand: "24AA7", expectedResult: []string{"AA", "7", "4", "2"}},
		{hand: "88AA7", expectedResult: []string{"AA", "88", "7"}},
		{hand: "777AA", expectedResult: []string{"777", "AA"}},
		{hand: "7TA77", expectedResult: []string{"777", "A", "T"}},
		{hand: "7A777", expectedResult: []string{"7777", "A"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.hand, func(t *testing.T) {
			assert.Equal(t, testCase.expectedResult, getSortedComponents(testCase.hand))
		})
	}
}

func TestCalculateResult(t *testing.T) {
	testCases := []struct {
		hands          []string
		expectedResult result
	}{
		// Custom test cases
		{hands: []string{"AAQaQ", "QQAAA"}, expectedResult: tie},
		{hands: []string{"3Q5q2", "5q3Q2"}, expectedResult: tie},

		// Provided test cases
		{hands: []string{"AAAQQ", "QQAAA"}, expectedResult: tie},
		{hands: []string{"53QQ2", "Q53Q2"}, expectedResult: tie},
		{hands: []string{"53888", "88385"}, expectedResult: tie},
		{hands: []string{"QQAAA", "AAAQQ"}, expectedResult: tie},
		{hands: []string{"Q53Q2", "53QQ2"}, expectedResult: tie},
		{hands: []string{"88385", "53888"}, expectedResult: tie},
		{hands: []string{"AAAQQ", "QQQAA"}, expectedResult: hand1Wins},
		{hands: []string{"Q53Q4", "53QQ2"}, expectedResult: hand1Wins},
		{hands: []string{"53888", "88375"}, expectedResult: hand1Wins},
		{hands: []string{"33337", "QQAAA"}, expectedResult: hand1Wins},
		{hands: []string{"22333", "AAA58"}, expectedResult: hand1Wins},
		{hands: []string{"33389", "AAKK4"}, expectedResult: hand1Wins},
		{hands: []string{"44223", "AA892"}, expectedResult: hand1Wins},
		{hands: []string{"22456", "AKQJT"}, expectedResult: hand1Wins},
		{hands: []string{"99977", "77799"}, expectedResult: hand1Wins},
		{hands: []string{"99922", "88866"}, expectedResult: hand1Wins},
		{hands: []string{"9922A", "9922K"}, expectedResult: hand1Wins},
		{hands: []string{"99975", "99965"}, expectedResult: hand1Wins},
		{hands: []string{"99975", "99974"}, expectedResult: hand1Wins},
		{hands: []string{"99752", "99652"}, expectedResult: hand1Wins},
		{hands: []string{"99752", "99742"}, expectedResult: hand1Wins},
		{hands: []string{"99753", "99752"}, expectedResult: hand1Wins},
		{hands: []string{"88822", "QQ777"}, expectedResult: hand1Wins},
		{hands: []string{"99662", "88776"}, expectedResult: hand1Wins},
		{hands: []string{"QQQAA", "AAAQQ"}, expectedResult: hand2Wins},
		{hands: []string{"53QQ2", "Q53Q4"}, expectedResult: hand2Wins},
		{hands: []string{"88375", "53888"}, expectedResult: hand2Wins},
		{hands: []string{"QQAAA", "33337"}, expectedResult: hand2Wins},
		{hands: []string{"AAA58", "22333"}, expectedResult: hand2Wins},
		{hands: []string{"AAKK4", "33389"}, expectedResult: hand2Wins},
		{hands: []string{"AA892", "44223"}, expectedResult: hand2Wins},
		{hands: []string{"AKQJT", "22456"}, expectedResult: hand2Wins},
		{hands: []string{"77799", "99977"}, expectedResult: hand2Wins},
		{hands: []string{"88866", "99922"}, expectedResult: hand2Wins},
		{hands: []string{"9922K", "9922A"}, expectedResult: hand2Wins},
		{hands: []string{"99965", "99975"}, expectedResult: hand2Wins},
		{hands: []string{"99974", "99975"}, expectedResult: hand2Wins},
		{hands: []string{"99652", "99752"}, expectedResult: hand2Wins},
		{hands: []string{"99742", "99752"}, expectedResult: hand2Wins},
		{hands: []string{"99752", "99753"}, expectedResult: hand2Wins},
	}

	for _, testCase := range testCases {
		t.Run(strings.Join(testCase.hands, " "), func(t *testing.T) {
			assert.Equal(t, testCase.expectedResult, calculateResult(testCase.hands[0], testCase.hands[1]))
		})
	}
}
