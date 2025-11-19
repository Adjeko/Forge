package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

// ASCII Art constants
const (
	top    = "█▀▀▀▀▀▀▀ ▄▀▀▀▀▀▄ █▀▀▀▀▀█ ▄▀▀▀▀▀▀ █▀▀▀▀▀▀▀"
	middle = "█▀▀▀▀▀▀  █     █ █▀▀▀▀█▀ █   ▀▀█ █▀▀▀▀▀▀ "
	bottom = "▀         ▀▀▀▀▀  ▀     ▀  ▀▀▀▀▀▀ ▀▀▀▀▀▀▀▀"
)

const (
	leftPattern   = "////"
	fillerPattern = "////"
	prefix        = "SEW"
	version       = "v1.0.0" // Placeholder version
)

var (
	red   = lipgloss.Color("#FF0000")
	white = lipgloss.Color("#FFFFFF")
)

// RenderHeader renders the header component with the given width
func RenderHeader(width int) string {
	if width == 0 {
		return ""
	}

	lines := []string{
		buildLine(buildVersionLine(), "bold red", width),
		buildLine(applyGradient(top, red, white), "", width),
		buildLine(applyGradient(middle, red, white), "", width),
		buildLine(applyGradient(bottom, red, white), "", width),
		buildLine("/////////////////////////////////////////", "red", width),
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func buildVersionLine() string {
	totalWidth := utf8.RuneCountInString(top)
	availableForVersion := totalWidth - len(prefix)
	if availableForVersion <= 0 {
		return prefix
	}

	versionDisplay := version
	if len(version) > availableForVersion {
		versionDisplay = version[:availableForVersion]
	}

	spacingLen := totalWidth - (len(prefix) + len(versionDisplay))
	if spacingLen < 0 {
		spacingLen = 0
	}

	return prefix + strings.Repeat(" ", spacingLen) + versionDisplay
}

func buildLine(content string, colorStyle string, width int) string {
	// Left pattern
	left := lipgloss.NewStyle().Foreground(red).Render(leftPattern)

	// Main content
	var main string
	if colorStyle != "" {
		style := lipgloss.NewStyle()
		if strings.Contains(colorStyle, "bold") {
			style = style.Bold(true)
		}
		if strings.Contains(colorStyle, "red") {
			style = style.Foreground(red)
		}
		main = style.Render(content)
	} else {
		main = content
	}

	// Calculate visible length (stripping ANSI codes)
	visibleLen := lipgloss.Width(left) + 1 + lipgloss.Width(main) + 1 // +1 for spaces
	remaining := width - visibleLen

	if remaining <= 0 {
		return fmt.Sprintf("%s %s", left, main)
	}

	// Filler
	var fillerBuilder strings.Builder
	currentLen := 0
	for currentLen < remaining {
		next := fillerPattern
		if currentLen+len(next) > remaining {
			next = next[:remaining-currentLen]
		}
		fillerBuilder.WriteString(next)
		currentLen += len(next)
	}
	filler := lipgloss.NewStyle().Foreground(red).Render(fillerBuilder.String())

	return fmt.Sprintf("%s %s %s", left, main, filler)
}

func applyGradient(text string, start, end lipgloss.Color) string {
	startHex := string(start)
	endHex := string(end)

	sC, _ := colorful.Hex(startHex)
	eC, _ := colorful.Hex(endHex)

	runes := []rune(text)
	gradientChars := 0
	for _, r := range runes {
		if r != ' ' {
			gradientChars++
		}
	}

	if gradientChars <= 1 {
		return text
	}

	var sb strings.Builder
	processed := 0
	for _, r := range runes {
		if r == ' ' {
			sb.WriteRune(r)
			continue
		}

		t := float64(processed) / float64(gradientChars-1)

		// Linear RGB Interpolation to match C# implementation
		rVal := sC.R + (eC.R-sC.R)*t
		gVal := sC.G + (eC.G-sC.G)*t
		bVal := sC.B + (eC.B-sC.B)*t

		c := colorful.Color{R: rVal, G: gVal, B: bVal}.Hex()

		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(c)).Render(string(r)))
		processed++
	}

	return sb.String()
}
