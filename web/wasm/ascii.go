package main

import (
	"fmt"
	"strings"
	"time"
)

// ASCII block character definitions for contribution levels
const (
	EmptyChar  = ' '  // Empty/Sky: No contributions
	FutureChar = '.'  // Future dates
	LowChar    = '░' // Low level contribution
	MediumChar = '▒' // Medium level contribution
	HighChar   = '▓' // High level contribution
	TopLowChar    = '╻' // Top block with low contributions
	TopMediumChar = '┃' // Top block with medium contributions
	TopHighChar   = '╽' // Top block with high contributions
)

// GenerateASCII creates ASCII art visualization from contribution data.
func GenerateASCII(contributions [][]ContributionDay) (string, error) {
	if len(contributions) == 0 {
		return "", fmt.Errorf("contributions data cannot be empty")
	}

	maxContrib := findMaxContributions(contributions)
	
	var sb strings.Builder
	
	// Process each day (7 rows for 7 days of the week)
	for day := 0; day < 7; day++ {
		for _, week := range contributions {
			if day < len(week) {
				char := getCharForContribution(week[day], maxContrib, day, week)
				sb.WriteRune(char)
			} else {
				sb.WriteRune(EmptyChar)
			}
		}
		sb.WriteRune('\n')
	}
	
	return sb.String(), nil
}

// getCharForContribution determines the ASCII character for a contribution day.
func getCharForContribution(day ContributionDay, maxContrib int, dayIdx int, week []ContributionDay) rune {
	// Check if this is a future date
	if isFutureDate(day.Date) {
		return FutureChar
	}
	
	// No contributions
	if day.ContributionCount == 0 {
		return EmptyChar
	}
	
	// Determine if this is the top block of the week (last day with contributions)
	isTopBlock := isTopOfWeek(dayIdx, week)
	
	// Calculate contribution level (0-3)
	level := getContributionLevel(day.ContributionCount, maxContrib)
	
	// Return appropriate character
	if isTopBlock {
		switch level {
		case 1:
			return TopLowChar
		case 2:
			return TopMediumChar
		case 3:
			return TopHighChar
		default:
			return EmptyChar
		}
	}
	
	switch level {
	case 1:
		return LowChar
	case 2:
		return MediumChar
	case 3:
		return HighChar
	default:
		return EmptyChar
	}
}

// isFutureDate checks if the date is in the future.
func isFutureDate(dateStr string) bool {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}
	return date.After(time.Now())
}

// isTopOfWeek determines if this is the highest day with contributions in the week.
func isTopOfWeek(dayIdx int, week []ContributionDay) bool {
	// Check if any days after this one have contributions
	for i := dayIdx + 1; i < len(week); i++ {
		if week[i].ContributionCount > 0 {
			return false
		}
	}
	return true
}

// getContributionLevel returns a level (0-3) based on contribution count.
func getContributionLevel(count, maxContrib int) int {
	if count == 0 {
		return 0
	}
	if maxContrib == 0 {
		return 1
	}
	
	// Calculate percentage
	percentage := float64(count) / float64(maxContrib)
	
	if percentage > 0.75 {
		return 3 // High
	} else if percentage > 0.5 {
		return 2 // Medium
	} else if percentage > 0 {
		return 1 // Low
	}
	
	return 0
}
