package main

import (
	"fmt"
	"strings"
	"time"
)

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
	Dim         = "\033[2m"
	ColorDim    = "\033[2m"
)

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

const (
	IconInfo     = "ℹ️"
	IconSuccess  = "✅"
	IconError    = "❌"
	IconWarning  = "⚠️"
	IconQuestion = "❓"
	IconFolder   = "📁"
	IconImage    = "🖼️"
	IconTool     = "🛠️"
	IconSave     = "💾"
	IconLoad     = "📥"
	IconEffect   = "✨"
	IconShape    = "🔶"
	IconConvert  = "🔄"
	IconExit     = "🚪"
	IconHelp     = "💡"
	IconTip      = "💡"
	IconKey      = "🔑"
)

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func showWelcomeBanner() {
	clearScreen()
	fmt.Println(ColorCyan + Bold + "╔════════════════════════════════════════════════════════════════════════╗" + ColorReset)
	fmt.Println(ColorCyan + Bold + "║" + ColorReset + "                    " + ColorPurple + Bold + "🎨 Bienvenue dans GoImage 🎨" + ColorReset + "                    " + ColorCyan + Bold + "║" + ColorReset)
	fmt.Println(ColorCyan + Bold + "╠════════════════════════════════════════════════════════════════════════╣" + ColorReset)
	fmt.Println(ColorCyan + Bold + "║" + ColorReset + "                     " + ColorYellow + "Éditeur d'images TUI en Go" + ColorReset + "                     " + ColorCyan + Bold + "║" + ColorReset)
	fmt.Println(ColorCyan + Bold + "║" + ColorReset + "                                                                        " + ColorCyan + Bold + "║" + ColorReset)
	fmt.Println(ColorCyan + Bold + "║" + ColorReset + "  " + ColorGreen + "Première fois ?" + ColorReset + " Voici les étapes recommandées :                  " + ColorCyan + Bold + "║" + ColorReset)
	fmt.Println(ColorCyan + Bold + "║" + ColorReset + "  1. " + IconLoad + " Chargez une image depuis le dossier 'test/' ou autre        " + ColorCyan + Bold + "║" + ColorReset)
	fmt.Println(ColorCyan + Bold + "║" + ColorReset + "  2. " + IconEffect + " Appliquez des effets (négatif, sépia, luminosité...)       " + ColorCyan + Bold + "║" + ColorReset)
	fmt.Println(ColorCyan + Bold + "║" + ColorReset + "  3. " + IconShape + " Dessinez des formes (carré, cercle)                        " + ColorCyan + Bold + "║" + ColorReset)
	fmt.Println(ColorCyan + Bold + "║" + ColorReset + "  4. " + IconSave + " Sauvegardez votre création                                 " + ColorCyan + Bold + "║" + ColorReset)
	fmt.Println(ColorCyan + Bold + "║" + ColorReset + "                                                                        " + ColorCyan + Bold + "║" + ColorReset)
	fmt.Println(ColorCyan + Bold + "║" + ColorReset + "  " + ColorBlue + "Raccourcis utiles:" + ColorReset + " 'h' = aide, 'q' = quitter, chiffres = sélection  " + ColorCyan + Bold + "║" + ColorReset)
	fmt.Println(ColorCyan + Bold + "╚════════════════════════════════════════════════════════════════════════╝" + ColorReset)
	fmt.Println()
	
	fmt.Print(ColorYellow + "Appuyez sur Entrée pour continuer..." + ColorReset)
	fmt.Scanln()
}

func showHelp(context string) {
	clearScreen()
	
	var helpContent []string
	var title string
	
	switch context {
	case "main":
		title = "Aide - Menu Principal"
		helpContent = []string{
			"🎯 NAVIGATION:",
			"• Tapez le numéro (1-6) pour sélectionner une option",
			"• Utilisez 'h' pour afficher cette aide",
			"• Utilisez 'q' pour quitter l'application",
			"",
			"📋 OPTIONS DISPONIBLES:",
			"• [1] Charger une image : Ouvre une image depuis votre disque",
			"• [2] Appliquer un effet : Transforme l'image (négatif, sépia, etc.)",
			"• [3] Dessiner une forme : Ajoute des formes géométriques",
			"• [4] Convertir l'image : Change le format ou redimensionne",
			"• [5] Sauvegarder : Enregistre l'image modifiée",
			"• [6] Quitter : Ferme l'application",
			"",
			"💡 CONSEIL:",
			"Commencez toujours par charger une image (option 1) !",
		}
	case "load":
		title = "Aide - Chargement d'images"
		helpContent = []string{
			"📁 MÉTHODES DE CHARGEMENT:",
			"• Navigation interactive : Parcourez vos dossiers visuellement",
			"• Saisie manuelle : Tapez le chemin direct vers l'image",
			"",
			"🖼️ FORMATS SUPPORTÉS:",
			"• PNG (.png) - Recommandé pour les images avec transparence",
			"• JPEG (.jpg, .jpeg) - Idéal pour les photos",
			"• GIF (.gif) - Pour les images simples",
			"",
			"📂 EXEMPLES DE CHEMINS:",
			"• test/test_image.png",
			"• /chemin/absolu/vers/image.jpg",
			"• ../dossier_parent/image.png",
			"",
			"🔍 NAVIGATION INTERACTIVE:",
			"• Utilisez les numéros pour sélectionner fichiers/dossiers",
			"• '..' remonte au dossier parent",
			"• Seules les images sont affichées",
		}
	case "effects":
		title = "Aide - Effets d'image"
		helpContent = []string{
			"✨ EFFETS DISPONIBLES:",
			"• Négatif : Inverse toutes les couleurs",
			"• Niveaux de gris : Convertit en noir et blanc",
			"• Sépia : Effet vintage brun/doré",
			"• Luminosité : Rend l'image plus claire/sombre",
			"• Contraste : Augmente/diminue les différences de couleur",
			"",
			"⚙️ EFFETS PARAMÉTRABLES:",
			"• Luminosité : 0.5 = sombre, 1.0 = normal, 1.5 = lumineux",
			"• Contraste : 0.5 = faible, 1.0 = normal, 2.0 = fort",
			"",
			"💡 ASTUCE:",
			"Vous pouvez appliquer plusieurs effets successivement !",
		}
	case "shapes":
		title = "Aide - Dessin de formes"
		helpContent = []string{
			"🔶 FORMES DISPONIBLES:",
			"• Carré : Forme rectangulaire remplie",
			"• Cercle : Forme circulaire remplie",
			"",
			"🎨 COULEURS:",
			"• Format RGB : R,G,B (ex: 255,0,0 pour rouge)",
			"• Valeurs entre 0 et 255 pour chaque composante",
			"• Exemples : 0,255,0 (vert), 0,0,255 (bleu)",
			"",
			"📐 POSITIONNEMENT:",
			"• X=0 : bord gauche de l'image",
			"• Y=0 : bord supérieur de l'image",
			"• Vérifiez que la forme reste dans les limites de l'image",
			"",
			"💡 CONSEIL:",
			"Commencez par de petites formes pour tester !",
		}
	default:
		title = "Aide générale"
		helpContent = []string{
			"🚀 UTILISATION GÉNÉRALE:",
			"• Suivez l'ordre logique : Charger → Modifier → Sauvegarder",
			"• Utilisez 'h' à tout moment pour obtenir de l'aide",
			"• Les barres de progression indiquent l'avancement",
			"• Les messages colorés indiquent le statut des opérations",
			"",
			"🎯 BONNES PRATIQUES:",
			"• Travaillez sur des copies de vos images importantes",
			"• Choisissez des noms de fichiers explicites lors de la sauvegarde",
			"• Testez les effets sur de petites images d'abord",
		}
	}
	
	drawBox(title, helpContent, 80)
	fmt.Println()
	fmt.Print(ColorYellow + "Appuyez sur Entrée pour continuer..." + ColorReset)
	fmt.Scanln()
}

func drawBox(title string, content []string, width int) {
	if width < len(title)+4 {
		width = len(title) + 4
	}

	fmt.Print(ColorCyan + Bold)
	fmt.Print("╭")
	titlePadding := width - len(title) - 4
	leftPadding := titlePadding / 2
	rightPadding := titlePadding - leftPadding
	
	for i := 0; i < leftPadding; i++ {
		fmt.Print("─")
	}
	fmt.Print("┤ " + ColorPurple + title + ColorCyan + " ├")
	for i := 0; i < rightPadding; i++ {
		fmt.Print("─")
	}
	fmt.Println("╮" + ColorReset)

	for _, line := range content {
		paddedLine := line
		if len(line) < width {
			paddedLine = line + strings.Repeat(" ", width-len(line))
		} else if len(line) > width {
			paddedLine = line[:width-3] + "..."
		}
		fmt.Println(ColorCyan + Bold + "│" + ColorReset + " " + paddedLine + " " + ColorCyan + Bold + "│" + ColorReset)
	}

	fmt.Print(ColorCyan + Bold + "╰")
	for i := 0; i < width; i++ {
		fmt.Print("─")
	}
	fmt.Println("╯" + ColorReset)
}

func drawMenuItem(index int, icon string, text string, shortcut string, selected bool) string {
	shortcutText := ""
	if shortcut != "" {
		shortcutText = ColorDim + " [" + shortcut + "]" + ColorReset
	}
	
	if selected {
		return fmt.Sprintf("%s %s[%d]%s %s %s%s%s%s", 
			ColorGreen, Bold, index, ColorReset, icon, BgGreen+ColorBlack, text, ColorReset, shortcutText)
	} else {
		return fmt.Sprintf(" %s[%d]%s %s %s%s", 
			ColorCyan, index, ColorReset, icon, text, shortcutText)
	}
}

func drawMenu(title string, items []string, icons []string, shortcuts []string, selectedIndex int, width int) {
	menuItems := make([]string, len(items))
	for i, item := range items {
		icon := ""
		shortcut := ""
		if i < len(icons) {
			icon = icons[i]
		}
		if i < len(shortcuts) {
			shortcut = shortcuts[i]
		}
		menuItems[i] = drawMenuItem(i+1, icon, item, shortcut, i == selectedIndex)
	}
	drawBox(title, menuItems, width)
}

func infoMessage(message string) {
	fmt.Println(ColorBlue + IconInfo + " " + message + ColorReset)
}

func successMessage(message string) {
	fmt.Println(ColorGreen + IconSuccess + " " + message + ColorReset)
}

func errorMessage(message string) {
	fmt.Println(ColorRed + IconError + " " + Bold + "ERREUR:" + ColorReset + " " + message)
}

func errorMessageWithTip(message, tip string) {
	fmt.Println(ColorRed + IconError + " " + Bold + "ERREUR:" + ColorReset + " " + message)
	fmt.Println(ColorYellow + IconTip + " " + Italic + "Conseil: " + tip + ColorReset)
}

func warningMessage(message string) {
	fmt.Println(ColorYellow + IconWarning + " " + Bold + "ATTENTION:" + ColorReset + " " + message)
}

func promptWithValidation(message string, validOptions []string) string {
	optionsStr := ""
	if len(validOptions) > 0 {
		optionsStr = " (" + strings.Join(validOptions, "/") + ")"
	}
	
	fmt.Print(ColorYellow + IconQuestion + " " + message + optionsStr + ColorReset + " ")
	var input string
	fmt.Scanln(&input)
	return input
}

func prompt(message string) string {
	fmt.Print(ColorYellow + IconQuestion + " " + message + ColorReset + " ")
	var input string
	fmt.Scanln(&input)
	return input
}

func drawProgressBarAnimated(progress float64, width int, message string) {
	barWidth := width - 10 // Espace pour les pourcentages et crochets
	completed := int(float64(barWidth) * progress)
	percent := int(progress * 100)
	
	animChars := []string{"▰", "▱", "▲", "▼"}
	animChar := animChars[percent%len(animChars)]
	
	fmt.Print(ColorBlue + message + ColorReset + " [")
	for i := 0; i < barWidth; i++ {
		if i < completed {
			fmt.Print(ColorGreen + "█" + ColorReset)
		} else if i == completed {
			fmt.Print(ColorYellow + animChar + ColorReset)
		} else {
			fmt.Print(ColorCyan + "░" + ColorReset)
		}
	}
	
	fmt.Printf("] %s%3d%%%s", ColorGreen, percent, ColorReset)
	
	if progress > 0.1 && progress < 0.9 {
		fmt.Print(ColorDim + " ⏱️" + ColorReset)
	}
	
	fmt.Print("\n")
}

func drawProgressBar(progress float64, width int) {
	drawProgressBarAnimated(progress, width, "Progression")
}

func drawHeader() {
	clearScreen()
	fmt.Println(ColorCyan + Bold + "╭──────────────────────────────────────────────────────────────────────╮" + ColorReset)
	fmt.Println(ColorCyan + Bold + "│" + ColorReset + "                    " + ColorPurple + Bold + "🎨 GoImage - Éditeur d'Images 🎨" + ColorReset + "                    " + ColorCyan + Bold + "│" + ColorReset)
	fmt.Println(ColorCyan + Bold + "│" + ColorReset + "                      " + ColorYellow + "Interface TUI - Version 1.0" + ColorReset + "                      " + ColorCyan + Bold + "│" + ColorReset)
	fmt.Println(ColorCyan + Bold + "╰──────────────────────────────────────────────────────────────────────╯" + ColorReset)
	fmt.Println()
}

func drawFooter() {
	fmt.Println()
	fmt.Println(ColorCyan + Bold + "╭──────────────────────────────────────────────────────────────────────╮" + ColorReset)
	fmt.Println(ColorCyan + Bold + "│" + ColorReset + " " + ColorGreen + IconKey + " Raccourcis:" + ColorReset + " " + ColorYellow + "q" + ColorReset + "=quitter " + ColorYellow + "h" + ColorReset + "=aide " + ColorYellow + "1-6" + ColorReset + "=sélection " + ColorYellow + "Entrée" + ColorReset + "=confirmer" + ColorCyan + "     │" + ColorReset)
	fmt.Println(ColorCyan + Bold + "│" + ColorReset + " " + ColorBlue + IconTip + " Astuce:" + ColorReset + " Suivez l'ordre logique: Charger → Modifier → Sauvegarder" + ColorCyan + "    │" + ColorReset)
	fmt.Println(ColorCyan + Bold + "╰──────────────────────────────────────────────────────────────────────╯" + ColorReset)
}

func drawStatusBar(currentImage string, hasImage bool) {
	status := "Aucune image chargée"
	icon := IconError
	color := ColorRed
	
	if hasImage {
		status = "Image: " + currentImage
		icon = IconSuccess
		color = ColorGreen
	}
	
	fmt.Println(ColorCyan + "╭─ " + ColorYellow + "STATUT" + ColorCyan + " ─────────────────────────────────────────────────────────────╮" + ColorReset)
	fmt.Println(ColorCyan + "│" + ColorReset + " " + color + icon + " " + status + strings.Repeat(" ", 62-len(status)) + ColorCyan + "│" + ColorReset)
	fmt.Println(ColorCyan + "╰────────────────────────────────────────────────────────────────╯" + ColorReset)
	fmt.Println()
}

func showNotification(message string, duration time.Duration, isSuccess bool) {
	icon := IconSuccess
	color := ColorGreen
	
	if !isSuccess {
		icon = IconError
		color = ColorRed
	}
	
	fmt.Println(color + icon + " " + message + ColorReset)
	time.Sleep(duration)
}

func confirmAction(message string) bool {
	fmt.Print(ColorYellow + IconQuestion + " " + message + " (o/N) " + ColorReset)
	var input string
	fmt.Scanln(&input)
	return strings.ToLower(input) == "o" || strings.ToLower(input) == "oui"
}

func displayImageInfo(width, height int, format string) {
	fmt.Println(ColorBlue + "╭─ " + ColorYellow + "INFORMATIONS IMAGE" + ColorBlue + " ─────────────────────────────────────────╮" + ColorReset)
	fmt.Println(ColorBlue + "│" + ColorReset + " " + IconImage + " Dimensions: " + ColorGreen + fmt.Sprintf("%d × %d pixels", width, height) + ColorReset + strings.Repeat(" ", 35-len(fmt.Sprintf("%d × %d pixels", width, height))) + ColorBlue + "│" + ColorReset)
	if format != "" {
		fmt.Println(ColorBlue + "│" + ColorReset + " " + IconTool + " Format: " + ColorGreen + format + ColorReset + strings.Repeat(" ", 47-len(format)) + ColorBlue + "│" + ColorReset)
	}
	fmt.Println(ColorBlue + "╰──────────────────────────────────────────────────────────────────╯" + ColorReset)
	fmt.Println()
}
