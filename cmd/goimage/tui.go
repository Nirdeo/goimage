package main

import (
	"fmt"
	"strings"
)

// Couleurs ANSI
const (
	ColorReset  = "\033[0m"
	ColorBlack  = "\033[30m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	Bold        = "\033[1m"
	Underline   = "\033[4m"
	Italic      = "\033[3m"
)

// Fond de couleur
const (
	BgBlack  = "\033[40m"
	BgRed    = "\033[41m"
	BgGreen  = "\033[42m"
	BgYellow = "\033[43m"
	BgBlue   = "\033[44m"
	BgPurple = "\033[45m"
	BgCyan   = "\033[46m"
	BgWhite  = "\033[47m"
)

// Assurer une largeur minimale
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// Dessine un cadre avec titre
func drawBox(title string, content []string, width int) {
	// Assurer une largeur minimale
	if width < len(title)+4 {
		width = len(title) + 4
	}

	// Ligne supérieure
	fmt.Print(ColorCyan)
	fmt.Print("╭")
	titleStart := (width - len(title)) / 2
	for i := 0; i < width; i++ {
		if i == titleStart-1 {
			fmt.Print("┤ ")
		} else if i == titleStart+len(title) {
			fmt.Print(" ├")
		} else if i >= titleStart && i < titleStart+len(title) {
			fmt.Print(Bold + title[i-titleStart:i-titleStart+1] + ColorReset + ColorCyan)
		} else {
			fmt.Print("─")
		}
	}
	fmt.Println("╮" + ColorReset)

	// Contenu
	for _, line := range content {
		paddedLine := line
		if len(line) < width {
			paddedLine = line + strings.Repeat(" ", width-len(line))
		} else if len(line) > width {
			paddedLine = line[:width-3] + "..."
		}
		fmt.Println(ColorCyan + "│" + ColorReset + paddedLine + ColorCyan + "│" + ColorReset)
	}

	// Ligne inférieure
	fmt.Print(ColorCyan + "╰")
	for i := 0; i < width; i++ {
		fmt.Print("─")
	}
	fmt.Println("╯" + ColorReset)
}

// Dessine un élément de menu
func drawMenuItem(index int, text string, selected bool) string {
	if selected {
		return fmt.Sprintf("%s %s[%d]%s %s%s%s", ColorGreen, Bold, index, ColorReset, BgGreen+ColorBlack, text, ColorReset)
	} else {
		return fmt.Sprintf(" [%d] %s", index, text)
	}
}

// Dessine un menu complet
func drawMenu(title string, items []string, selectedIndex int, width int) {
	menuItems := make([]string, len(items))
	for i, item := range items {
		menuItems[i] = drawMenuItem(i+1, item, i == selectedIndex)
	}
	drawBox(title, menuItems, width)
}

// Affiche un message d'information
func infoMessage(message string) {
	fmt.Println(ColorBlue + "ℹ " + message + ColorReset)
}

// Affiche un message de succès
func successMessage(message string) {
	fmt.Println(ColorGreen + "✓ " + message + ColorReset)
}

// Affiche un message d'erreur
func errorMessage(message string) {
	fmt.Println(ColorRed + "✗ " + message + ColorReset)
}

// Affiche un prompt pour l'utilisateur
func prompt(message string) string {
	fmt.Print(ColorYellow + "? " + message + " " + ColorReset)
	var input string
	fmt.Scanln(&input)
	return input
}

// Dessine une barre de progression
func drawProgressBar(progress float64, width int) {
	barWidth := width - 7 // Espace pour les pourcentages
	completed := int(float64(barWidth) * progress)

	fmt.Print("[")
	for i := 0; i < barWidth; i++ {
		if i < completed {
			fmt.Print(ColorGreen + "█" + ColorReset)
		} else {
			fmt.Print(ColorCyan + "░" + ColorReset)
		}
	}

	percent := int(progress * 100)
	fmt.Printf("] %3d%%\n", percent)
}

// Affiche l'en-tête de l'application
func drawHeader() {
	clearScreen()
	fmt.Println(ColorCyan + Bold + "╭────────────────────────────────────────────────╮" + ColorReset)
	fmt.Println(ColorCyan + Bold + "│            " + ColorPurple + "GoImage - Éditeur d'Images" + ColorCyan + "            │" + ColorReset)
	fmt.Println(ColorCyan + Bold + "╰────────────────────────────────────────────────╯" + ColorReset)
	fmt.Println()
}

// Affiche le pied de page
func drawFooter() {
	fmt.Println()
	fmt.Println(ColorCyan + Bold + "╭────────────────────────────────────────────────╮" + ColorReset)
	fmt.Println(ColorCyan + Bold + "│ " + ColorReset + "Appuyez sur " + ColorYellow + "q" + ColorReset + " pour quitter | " + ColorYellow + "h" + ColorReset + " pour l'aide" + ColorCyan + "        │" + ColorReset)
	fmt.Println(ColorCyan + Bold + "╰────────────────────────────────────────────────╯" + ColorReset)
}
