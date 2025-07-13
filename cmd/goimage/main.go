package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	// Support des formats d'image standards
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	
	// Import des effets depuis le package
	"github.com/nirdeo/goimage/pkg/effects"
)

// Variable globale pour déterminer si c'est la première utilisation
var isFirstTime = true

func main() {
	StartTUI()
}

func StartTUI() {
	var img image.Image
	var err error
	var currentFilePath string
	var imageFormat string

	if isFirstTime {
		showWelcomeBanner()
		isFirstTime = false
	}

	menuItems := []string{
		"Charger une image",
		"Appliquer un effet",
		"Dessiner une forme",
		"Convertir l'image",
		"Sauvegarder l'image",
		"Quitter",
	}

	menuIcons := []string{
		IconLoad,
		IconEffect,
		IconShape,
		IconConvert,
		IconSave,
		IconExit,
	}

	menuShortcuts := []string{
		"Ctrl+O",
		"Ctrl+E",
		"Ctrl+D",
		"Ctrl+C",
		"Ctrl+S",
		"Ctrl+Q",
	}

	for {
		drawHeader()
		
		drawStatusBar(currentFilePath, img != nil)
		
		drawMenu("Menu Principal", menuItems, menuIcons, menuShortcuts, -1, 70)
		drawFooter()

		choice := promptWithValidation("Choisissez une option", []string{"1", "2", "3", "4", "5", "6", "h", "q"})

		if choice == "h" || choice == "H" {
			showHelp("main")
			continue
		}

		switch choice {
		case "1":
			if choice == "h" {
				showHelp("load")
				continue
			}
			
			clearScreen()
			drawBox("Charger une image", []string{
				"Choisissez une méthode de sélection d'image:",
				"",
				"1. " + IconFolder + " Navigation interactive dans les fichiers (Recommandé)",
				"2. " + IconTool + " Saisie manuelle du chemin",
				"",
				"💡 Astuce: La navigation interactive vous montre uniquement les images supportées",
			}, 80)
			
			methodChoice := promptWithValidation("Méthode de chargement", []string{"1", "2", "h"})
			
			if methodChoice == "h" {
				showHelp("load")
				continue
			}
			
			var filePath string
			if methodChoice == "1" {
				filePath, err = navigateToFile()
			} else if methodChoice == "2" {
				filePath = readUserInput("Chemin de l'image à charger (ex: test/test_image.png)")
			} else {
				warningMessage("Option invalide, retour au menu principal")
				time.Sleep(1 * time.Second)
				continue
			}
			
			if err != nil {
				if err.Error() != "navigation annulée" {
					errorMessageWithTip(fmt.Sprintf("Erreur lors de la sélection: %v", err), "Vérifiez que le fichier existe et que vous avez les permissions")
				}
				time.Sleep(2 * time.Second)
				continue
			}
			
			img, currentFilePath, imageFormat, err = loadImageFromPathEnhanced(filePath)
			if err != nil {
				errorMessageWithTip(fmt.Sprintf("Erreur lors du chargement: %v", err), "Vérifiez que le fichier est une image valide (PNG, JPEG, GIF)")
				time.Sleep(2 * time.Second)
			} else if img == nil {
				errorMessage("L'image n'a pas pu être chargée correctement")
				time.Sleep(2 * time.Second)
			} else {
				bounds := img.Bounds()
				successMessage(fmt.Sprintf("Image chargée avec succès!"))
				displayImageInfo(bounds.Dx(), bounds.Dy(), imageFormat)
				time.Sleep(2 * time.Second)
			}
			
		case "2":
			if img == nil {
				errorMessageWithTip("Veuillez d'abord charger une image", "Utilisez l'option 1 pour charger une image")
				time.Sleep(2 * time.Second)
				continue
			}
			
			if choice == "h" {
				showHelp("effects")
				continue
			}
			
			img = applyEffectEnhanced(img)
			
		case "3":
			if img == nil {
				errorMessageWithTip("Veuillez d'abord charger une image", "Utilisez l'option 1 pour charger une image")
				time.Sleep(2 * time.Second)
				continue
			}
			
			if choice == "h" {
				showHelp("shapes")
				continue
			}
			
			img = drawShapeEnhanced(img)
			
		case "4":
			if img == nil {
				errorMessageWithTip("Veuillez d'abord charger une image", "Utilisez l'option 1 pour charger une image")
				time.Sleep(2 * time.Second)
				continue
			}
			
			modifiedImg, err := convertImageEnhanced(img)
			if err != nil {
				errorMessageWithTip(fmt.Sprintf("Erreur lors de la conversion: %v", err), "Vérifiez le format de sortie et les permissions d'écriture")
				time.Sleep(2 * time.Second)
			} else if modifiedImg != nil {
				// Mettre à jour l'image principale si elle a été modifiée
				img = modifiedImg
				successMessage("Image convertie avec succès!")
				time.Sleep(1 * time.Second)
			}
			
		case "5":
			if img == nil {
				errorMessageWithTip("Veuillez d'abord charger une image", "Utilisez l'option 1 pour charger une image")
				time.Sleep(2 * time.Second)
				continue
			}
			err := saveImageEnhanced(img)
			if err != nil {
				errorMessageWithTip(fmt.Sprintf("Erreur lors de la sauvegarde: %v", err), "Vérifiez le chemin et les permissions d'écriture")
				time.Sleep(2 * time.Second)
			}
			
		case "6", "q", "Q":
			if img != nil {
				if confirmAction("Vous avez une image en cours d'édition. Quitter quand même ?") {
					clearScreen()
					successMessage("Merci d'avoir utilisé GoImage. À bientôt! 👋")
					return
				}
			} else {
				clearScreen()
				successMessage("Merci d'avoir utilisé GoImage. À bientôt! 👋")
				return
			}
			
		default:
			warningMessage("Option invalide. Utilisez les numéros 1-6, 'h' pour l'aide, ou 'q' pour quitter")
			time.Sleep(1 * time.Second)
		}
	}
}

// readUserInput lit une entrée utilisateur et renvoie la chaîne de caractères
func readUserInput(promptMessage string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(ColorYellow + IconQuestion + " " + promptMessage + ColorReset + " ")
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// loadImageFromPathEnhanced charge une image avec un meilleur feedback UX
func loadImageFromPathEnhanced(filePath string) (image.Image, string, string, error) {
	clearScreen()
	drawBox("Chargement de l'image", []string{
		"Fichier sélectionné: " + filePath,
		"Formats supportés: PNG, JPEG, GIF",
		"",
		"Vérification du fichier...",
	}, 80)
	fmt.Println()

	// Vérification que le fichier existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, "", "", fmt.Errorf("le fichier n'existe pas: %s", filePath)
	}

	// Simulation de progression avec des étapes claires
	steps := []string{
		"Vérification du fichier...",
		"Ouverture du fichier...",
		"Analyse du format...",
		"Décodage de l'image...",
		"Finalisation...",
	}

	for i, step := range steps {
		drawProgressBarAnimated(float64(i)/float64(len(steps)), 50, step)
		time.Sleep(100 * time.Millisecond)
		if i < len(steps)-1 {
			fmt.Print("\033[1A\r") // Remonte d'une ligne
		}
	}

	// Ouverture du fichier
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", "", fmt.Errorf("impossible d'ouvrir le fichier: %v", err)
	}
	defer file.Close()

	// Tentative de décodage de l'image
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", "", fmt.Errorf("impossible de décoder l'image (format non supporté?): %v", err)
	}

	drawProgressBarAnimated(1.0, 50, "Chargement terminé")
	fmt.Println()

	if img == nil {
		return nil, "", "", fmt.Errorf("l'image a été décodée comme nil")
	}

	return img, filePath, format, nil
}

// readMetadata lit et affiche les métadonnées basiques d'une image
func readMetadata(img image.Image) {
	clearScreen()

	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Déterminer le type de couleur
	colorModel := "Inconnu"
	switch img.ColorModel() {
	case color.RGBAModel:
		colorModel = "RGBA (8 bits par canal)"
	case color.RGBA64Model:
		colorModel = "RGBA (16 bits par canal)"
	case color.NRGBAModel:
		colorModel = "NRGBA (alpha non-prémultiplié)"
	case color.GrayModel:
		colorModel = "Niveaux de gris (8 bits)"
	case color.Gray16Model:
		colorModel = "Niveaux de gris (16 bits)"
	case color.CMYKModel:
		colorModel = "CMYK"
	}

	// Récupérer quelques échantillons de couleurs
	topLeft := img.At(bounds.Min.X, bounds.Min.Y)
	topRight := img.At(bounds.Max.X-1, bounds.Min.Y)
	bottomLeft := img.At(bounds.Min.X, bounds.Max.Y-1)
	bottomRight := img.At(bounds.Max.X-1, bounds.Max.Y-1)
	center := img.At((bounds.Min.X+bounds.Max.X)/2, (bounds.Min.Y+bounds.Max.Y)/2)

	drawBox("Métadonnées de l'image", []string{
		fmt.Sprintf("Dimensions: %d × %d pixels", width, height),
		fmt.Sprintf("Format de couleur: %s", colorModel),
		fmt.Sprintf("Rectangle: (%d,%d) à (%d,%d)", bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y),
		"",
		"Échantillons de couleurs:",
		fmt.Sprintf("- Haut gauche: %v", topLeft),
		fmt.Sprintf("- Haut droite: %v", topRight),
		fmt.Sprintf("- Bas gauche: %v", bottomLeft),
		fmt.Sprintf("- Bas droite: %v", bottomRight),
		fmt.Sprintf("- Centre: %v", center),
	}, 70)

	readUserInput("Appuyez sur Entrée pour continuer")
	return
}

// applyEffectEnhanced applique un effet avec une meilleure UX
func applyEffectEnhanced(img image.Image) image.Image {
	clearScreen()
	effectItems := []string{
		"Négatif",
		"Niveaux de gris", 
		"Sépia",
		"Luminosité",
		"Contraste",
		"Aide",
		"Retour",
	}

	effectIcons := []string{
		"🔄", "⚫", "🟤", "☀️", "🔆", IconHelp, "↩️",
	}

	drawBox("Appliquer un effet", []string{
		"Choisissez un effet à appliquer à l'image actuelle",
		"",
		"💡 Astuce: Certains effets sont paramétrables",
		"⚠️ L'effet remplacera l'image actuelle",
	}, 80)
	fmt.Println()

	drawMenu("Effets disponibles", effectItems, effectIcons, []string{}, -1, 70)

	choice := promptWithValidation("Choisissez un effet", []string{"1", "2", "3", "4", "5", "6", "7", "h"})

	if choice == "h" || choice == "6" {
		showHelp("effects")
		return img
	}

	if choice == "7" {
		return img
	}

	var effect effects.Effect

	switch choice {
	case "1":
		effect = &effects.NegativeEffect{}
	case "2":
		effect = &effects.GrayscaleEffect{}
	case "3":
		effect = &effects.SepiaEffect{}
	case "4":
		// Effet luminosité avec aide contextuelle
		infoMessage("Ajustement de la luminosité")
		fmt.Println("💡 Valeurs recommandées:")
		fmt.Println("  • 0.5 = Image plus sombre")
		fmt.Println("  • 1.0 = Luminosité normale")
		fmt.Println("  • 1.5 = Image plus lumineuse")
		fmt.Println()
		
		factorStr := readUserInput("Facteur de luminosité (0.1 à 3.0)")
		factor, err := strconv.ParseFloat(factorStr, 64)
		if err != nil || factor <= 0 || factor > 3.0 {
			warningMessage("Valeur invalide, utilisation de la valeur par défaut (1.0)")
			factor = 1.0
		}
		effect = &effects.BrightnessEffect{Factor: factor}
	case "5":
		// Effet contraste avec aide contextuelle
		infoMessage("Ajustement du contraste")
		fmt.Println("💡 Valeurs recommandées:")
		fmt.Println("  • 0.5 = Contraste faible")
		fmt.Println("  • 1.0 = Contraste normal")
		fmt.Println("  • 2.0 = Contraste fort")
		fmt.Println()
		
		factorStr := readUserInput("Facteur de contraste (0.1 à 3.0)")
		factor, err := strconv.ParseFloat(factorStr, 64)
		if err != nil || factor <= 0 || factor > 3.0 {
			warningMessage("Valeur invalide, utilisation de la valeur par défaut (1.0)")
			factor = 1.0
		}
		effect = &effects.ContrastEffect{Factor: factor}
	default:
		warningMessage("Option invalide, retour au menu principal")
		time.Sleep(1 * time.Second)
		return img
	}

	clearScreen()
	drawBox("Application de l'effet", []string{
		"Effet sélectionné: " + effect.Name(),
		"Description: " + effect.Description(),
		"",
		"⏳ Traitement en cours...",
	}, 80)
	fmt.Println()

	// Simulation de progression pour l'effet
	for i := 0; i <= 100; i += 5 {
		drawProgressBarAnimated(float64(i)/100.0, 50, "Application de l'effet")
		time.Sleep(30 * time.Millisecond)
		if i < 100 {
			fmt.Print("\033[1A\r") // Remonte d'une ligne
		}
	}

	modifiedImg := effect.Apply(img)
	
	successMessage("Effet appliqué avec succès!")
	infoMessage("Vous pouvez maintenant appliquer d'autres effets ou sauvegarder l'image")
	time.Sleep(2 * time.Second)
	
	return modifiedImg
}

// drawShapeEnhanced dessine une forme avec une meilleure UX
func drawShapeEnhanced(img image.Image) image.Image {
	clearScreen()
	shapeItems := []string{
		"Carré",
		"Cercle",
		"Aide",
		"Retour",
	}

	shapeIcons := []string{
		"⬜", "⭕", IconHelp, "↩️",
	}

	bounds := img.Bounds()
	drawBox("Dessiner une forme", []string{
		"Choisissez une forme à dessiner sur l'image",
		"",
		fmt.Sprintf("🖼️ Dimensions de l'image: %d × %d pixels", bounds.Dx(), bounds.Dy()),
		"💡 Astuce: Vérifiez que la forme reste dans les limites",
	}, 80)
	fmt.Println()

	drawMenu("Formes disponibles", shapeItems, shapeIcons, []string{}, -1, 70)

	choice := promptWithValidation("Choisissez une forme", []string{"1", "2", "3", "4", "h"})

	if choice == "h" || choice == "3" {
		showHelp("shapes")
		return img
	}

	if choice == "4" {
		return img
	}

	// Définition de la couleur avec aide
	fmt.Println()
	infoMessage("Configuration de la couleur")
	fmt.Println("💡 Exemples de couleurs:")
	fmt.Println("  • 255,0,0 = Rouge")
	fmt.Println("  • 0,255,0 = Vert")
	fmt.Println("  • 0,0,255 = Bleu")
	fmt.Println("  • 255,255,0 = Jaune")
	fmt.Println("  • 255,0,255 = Magenta")
	fmt.Println("  • 0,255,255 = Cyan")
	fmt.Println("  • 0,0,0 = Noir")
	fmt.Println("  • 255,255,255 = Blanc")
	fmt.Println()

	colorInput := readUserInput("Couleur au format R,G,B (ex: 255,0,0)")
	colorParts := strings.Split(colorInput, ",")
	if len(colorParts) != 3 {
		warningMessage("Format de couleur invalide, utilisation du rouge par défaut")
		colorParts = []string{"255", "0", "0"}
	}

	r, err1 := strconv.Atoi(strings.TrimSpace(colorParts[0]))
	g, err2 := strconv.Atoi(strings.TrimSpace(colorParts[1]))
	b, err3 := strconv.Atoi(strings.TrimSpace(colorParts[2]))
	
	if err1 != nil || err2 != nil || err3 != nil || r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		warningMessage("Valeurs de couleur invalides, utilisation du rouge par défaut")
		r, g, b = 255, 0, 0
	}
	
	shapeColor := color.RGBA{uint8(r), uint8(g), uint8(b), 255}

	switch choice {
	case "1": // Carré
		clearScreen()
		drawBox("Paramètres du carré", []string{
			"Définissez la position et la taille du carré",
			fmt.Sprintf("🖼️ Dimensions de l'image: %d × %d pixels", bounds.Dx(), bounds.Dy()),
			fmt.Sprintf("🎨 Couleur sélectionnée: RGB(%d,%d,%d)", r, g, b),
			"",
			"💡 Conseils:",
			"  • X=0, Y=0 = coin supérieur gauche",
			"  • Vérifiez que X+Taille ≤ largeur image",
			"  • Vérifiez que Y+Taille ≤ hauteur image",
		}, 80)
		fmt.Println()
		
		// Conseils pour le positionnement du carré
		fmt.Println("💡 Suggestions de positionnement:")
		fmt.Printf("  • Coin supérieur gauche: X=0, Y=0\n")
		fmt.Printf("  • Centre de l'image: X=%d, Y=%d\n", bounds.Dx()/2, bounds.Dy()/2)
		fmt.Printf("  • Coin supérieur droit: X=%d, Y=0\n", bounds.Dx()-50)
		fmt.Printf("  • Coin inférieur gauche: X=0, Y=%d\n", bounds.Dy()-50)
		fmt.Printf("  • Taille recommandée: %d-%d pixels\n", bounds.Dx()/10, bounds.Dx()/5)
		fmt.Println()

		xStr := readUserInput("Position X (0 à gauche)")
		yStr := readUserInput("Position Y (0 en haut)")
		sizeStr := readUserInput("Taille en pixels")

		x, err1 := strconv.Atoi(xStr)
		y, err2 := strconv.Atoi(yStr)
		size, err3 := strconv.Atoi(sizeStr)

		if err1 != nil || err2 != nil || err3 != nil || x < 0 || y < 0 || size <= 0 {
			errorMessageWithTip("Valeurs invalides", "Utilisez des nombres positifs")
			time.Sleep(2 * time.Second)
			return img
		}

		if x+size > bounds.Dx() || y+size > bounds.Dy() {
			errorMessageWithTip("Le carré dépasse les limites de l'image", fmt.Sprintf("Réduisez la taille ou ajustez la position (image: %dx%d)", bounds.Dx(), bounds.Dy()))
			time.Sleep(2 * time.Second)
			return img
		}

		clearScreen()
		drawBox("Dessin du carré", []string{
			fmt.Sprintf("Position: (%d, %d)", x, y),
			fmt.Sprintf("Taille: %d pixels", size),
			fmt.Sprintf("Couleur: RGB(%d,%d,%d)", r, g, b),
		}, 80)
		fmt.Println()

		// Simulation de progression
		for i := 0; i <= 100; i += 10 {
			drawProgressBarAnimated(float64(i)/100.0, 50, "Dessin du carré")
			time.Sleep(50 * time.Millisecond)
			if i < 100 {
				fmt.Print("\033[1A\r") // Remonte d'une ligne
			}
		}

		squareEffect := &SquareEffect{
			X:     x,
			Y:     y,
			Size:  size,
			Color: shapeColor,
		}
		modifiedImg := squareEffect.Apply(img)
		successMessage("Carré dessiné avec succès!")
		time.Sleep(2 * time.Second)
		return modifiedImg

	case "2": // Cercle
		clearScreen()
		drawBox("Paramètres du cercle", []string{
			"Définissez la position et le rayon du cercle",
			fmt.Sprintf("🖼️ Dimensions de l'image: %d × %d pixels", bounds.Dx(), bounds.Dy()),
			fmt.Sprintf("🎨 Couleur sélectionnée: RGB(%d,%d,%d)", r, g, b),
			"",
			"💡 Conseils:",
			"  • Centre X, Y = position du centre du cercle",
			"  • Vérifiez que le cercle reste dans l'image",
			"  • Rayon max recommandé: " + fmt.Sprintf("%d", min(bounds.Dx(), bounds.Dy())/2),
		}, 80)
		fmt.Println()
		
		// Conseils pour le positionnement du cercle
		fmt.Println("💡 Suggestions de positionnement:")
		fmt.Printf("  • Centre de l'image: X=%d, Y=%d\n", bounds.Dx()/2, bounds.Dy()/2)
		fmt.Printf("  • Quart supérieur gauche: X=%d, Y=%d\n", bounds.Dx()/4, bounds.Dy()/4)
		fmt.Printf("  • Quart supérieur droit: X=%d, Y=%d\n", 3*bounds.Dx()/4, bounds.Dy()/4)
		fmt.Printf("  • Quart inférieur gauche: X=%d, Y=%d\n", bounds.Dx()/4, 3*bounds.Dy()/4)
		fmt.Printf("  • Rayon petit: %d pixels\n", bounds.Dx()/20)
		fmt.Printf("  • Rayon moyen: %d pixels\n", bounds.Dx()/10)
		fmt.Printf("  • Rayon grand: %d pixels\n", bounds.Dx()/6)
		fmt.Println()

		xStr := readUserInput("Centre X")
		yStr := readUserInput("Centre Y") 
		radiusStr := readUserInput("Rayon en pixels")

		x, err1 := strconv.Atoi(xStr)
		y, err2 := strconv.Atoi(yStr)
		radius, err3 := strconv.Atoi(radiusStr)

		if err1 != nil || err2 != nil || err3 != nil || x < 0 || y < 0 || radius <= 0 {
			errorMessageWithTip("Valeurs invalides", "Utilisez des nombres positifs")
			time.Sleep(2 * time.Second)
			return img
		}

		// Vérifier que le centre du cercle est dans l'image et qu'au moins une partie du cercle sera visible
		if x < 0 || x >= bounds.Dx() || y < 0 || y >= bounds.Dy() {
			errorMessageWithTip("Le centre du cercle doit être dans l'image", fmt.Sprintf("Utilisez des coordonnées entre 0 et %dx%d", bounds.Dx()-1, bounds.Dy()-1))
			time.Sleep(2 * time.Second)
			return img
		}
		
		// Avertir si le cercle dépasse largement les limites (mais permettre quand même)
		if x+radius < 0 || x-radius >= bounds.Dx() || y+radius < 0 || y-radius >= bounds.Dy() {
			errorMessageWithTip("Le cercle est entièrement en dehors de l'image", fmt.Sprintf("Ajustez le centre ou le rayon (image: %dx%d)", bounds.Dx(), bounds.Dy()))
			time.Sleep(2 * time.Second)
			return img
		}

		clearScreen()
		drawBox("Dessin du cercle", []string{
			fmt.Sprintf("Centre: (%d, %d)", x, y),
			fmt.Sprintf("Rayon: %d pixels", radius),
			fmt.Sprintf("Couleur: RGB(%d,%d,%d)", r, g, b),
		}, 80)
		fmt.Println()

		// Simulation de progression
		for i := 0; i <= 100; i += 10 {
			drawProgressBarAnimated(float64(i)/100.0, 50, "Dessin du cercle")
			time.Sleep(50 * time.Millisecond)
			if i < 100 {
				fmt.Print("\033[1A\r") // Remonte d'une ligne
			}
		}

		circleEffect := &CircleEffect{
			CenterX: x,
			CenterY: y,
			Radius:  radius,
			Color:   shapeColor,
		}
		modifiedImg := circleEffect.Apply(img)
		successMessage("Cercle dessiné avec succès!")
		time.Sleep(2 * time.Second)
		return modifiedImg

	default:
		warningMessage("Option invalide, retour au menu principal")
		time.Sleep(1 * time.Second)
		return img
	}
}

// resizeImage redimensionne l'image selon les dimensions spécifiées
func resizeImage(img image.Image, newWidth, newHeight int) image.Image {
	// Si les deux dimensions sont 0, on ne fait rien
	if newWidth == 0 && newHeight == 0 {
		return img
	}

	// Obtenir les dimensions actuelles
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Calculer les dimensions manquantes si l'utilisateur n'en a spécifié qu'une
	if newWidth == 0 {
		// Préserver le ratio si seulement la hauteur est spécifiée
		newWidth = int(float64(width) * float64(newHeight) / float64(height))
	} else if newHeight == 0 {
		// Préserver le ratio si seulement la largeur est spécifiée
		newHeight = int(float64(height) * float64(newWidth) / float64(width))
	}

	// Créer la nouvelle image redimensionnée
	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Algorithme de redimensionnement par interpolation au plus proche voisin
	// Simple mais moins qualitatif que la méthode bilinéaire
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// Calculer les coordonnées correspondantes dans l'image source
			srcX := bounds.Min.X + x*width/newWidth
			srcY := bounds.Min.Y + y*height/newHeight

			// Copier la couleur du pixel source au pixel de destination
			newImg.Set(x, y, img.At(srcX, srcY))
		}
	}

	return newImg
}

// convertImage permet à l'utilisateur de convertir une image dans un autre format
func convertImage(img image.Image) error {
	clearScreen()
	formatItems := []string{
		"PNG",
		"JPEG (qualité standard)",
		"JPEG (haute qualité)",
		"JPEG (qualité personnalisée)",
		"GIF",
		"Redimensionner l'image",
		"Afficher les métadonnées",
		"Retour",
	}

	drawBox("Conversion d'image", []string{
		"Choisissez le format vers lequel vous souhaitez convertir l'image",
		"ou une autre opération",
	}, 70)
	fmt.Println()

	drawMenu("Options disponibles", formatItems, []string{}, []string{}, -1, 50)

	choice := readUserInput("Choisissez une option")

	if choice == "8" || choice == "0" {
		return nil
	}

	// Option pour afficher les métadonnées
	if choice == "7" {
		readMetadata(img)
		return nil
	}

	// Option pour redimensionner
	if choice == "6" {
		clearScreen()
		drawBox("Redimensionnement d'image", []string{
			"Spécifiez les nouvelles dimensions de l'image",
			"(entrez 0 pour une dimension pour conserver le ratio)",
		}, 70)
		fmt.Println()

		// Obtenir les dimensions actuelles
		bounds := img.Bounds()
		width := bounds.Max.X - bounds.Min.X
		height := bounds.Max.Y - bounds.Min.Y

		infoMessage(fmt.Sprintf("Dimensions actuelles: %d × %d pixels", width, height))
		
		// Conseils pour le redimensionnement
		fmt.Println("💡 Suggestions de redimensionnement:")
		fmt.Printf("  • Réduire de moitié: %d × %d\n", width/2, height/2)
		fmt.Printf("  • Doubler la taille: %d × %d\n", width*2, height*2)
		fmt.Printf("  • Format HD (16:9): 1920 × 1080\n")
		fmt.Printf("  • Format carré: %d × %d\n", min(width, height), min(width, height))
		fmt.Printf("  • Largeur fixe 800px: 800 × %d\n", (800*height)/width)
		fmt.Printf("  • Hauteur fixe 600px: %d × 600\n", (600*width)/height)
		fmt.Println("  • Entrez 0 pour une dimension pour conserver le ratio")
		fmt.Println()

		widthStr := readUserInput("Nouvelle largeur (pixels)")
		heightStr := readUserInput("Nouvelle hauteur (pixels)")

		// Convertir les entrées en nombres
		newWidth, err1 := strconv.Atoi(widthStr)
		newHeight, err2 := strconv.Atoi(heightStr)

		if err1 != nil || err2 != nil || (newWidth < 0) || (newHeight < 0) {
			errorMessage("Dimensions invalides")
			time.Sleep(1 * time.Second)
			return fmt.Errorf("dimensions invalides")
		}

		clearScreen()
		infoMessage("Redimensionnement en cours...")

		// Simulation de progression
		for i := 0; i <= 100; i += 2 {
			drawProgressBar(float64(i)/100.0, 40)
			time.Sleep(15 * time.Millisecond)
			if i < 100 {
				fmt.Print("\033[1A\r") // Remonte d'une ligne
			}
		}

		// Redimensionner l'image
		resizedImg := resizeImage(img, newWidth, newHeight)

		// Mettre à jour l'image pour les traitements ultérieurs
		*(&img) = resizedImg

		successMessage(fmt.Sprintf("Image redimensionnée avec succès: %d × %d pixels",
			resizedImg.Bounds().Max.X-resizedImg.Bounds().Min.X,
			resizedImg.Bounds().Max.Y-resizedImg.Bounds().Min.Y))

		time.Sleep(1 * time.Second)
		return nil
	}

	outputPath := readUserInput("Chemin du fichier de sortie (avec extension)")

	// Vérifier et ajouter l'extension si nécessaire
	switch choice {
	case "1": // PNG
		if !strings.HasSuffix(strings.ToLower(outputPath), ".png") {
			outputPath += ".png"
		}

		clearScreen()
		infoMessage("Conversion en PNG...")

		// Simulation de progression
		for i := 0; i <= 100; i += 5 {
			drawProgressBar(float64(i)/100.0, 40)
			time.Sleep(50 * time.Millisecond)
			if i < 100 {
				fmt.Print("\033[1A\r") // Remonte d'une ligne
			}
		}

		file, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer file.Close()

		err = png.Encode(file, img)
		if err != nil {
			return err
		}

	case "2": // JPEG qualité standard
		if !strings.HasSuffix(strings.ToLower(outputPath), ".jpg") &&
			!strings.HasSuffix(strings.ToLower(outputPath), ".jpeg") {
			outputPath += ".jpg"
		}

		clearScreen()
		infoMessage("Conversion en JPEG (qualité standard)...")

		// Simulation de progression
		for i := 0; i <= 100; i += 5 {
			drawProgressBar(float64(i)/100.0, 40)
			time.Sleep(30 * time.Millisecond)
			if i < 100 {
				fmt.Print("\033[1A\r") // Remonte d'une ligne
			}
		}

		file, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer file.Close()

		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 75})
		if err != nil {
			return err
		}

	case "3": // JPEG haute qualité
		if !strings.HasSuffix(strings.ToLower(outputPath), ".jpg") &&
			!strings.HasSuffix(strings.ToLower(outputPath), ".jpeg") {
			outputPath += ".jpg"
		}

		clearScreen()
		infoMessage("Conversion en JPEG (haute qualité)...")

		// Simulation de progression
		for i := 0; i <= 100; i += 5 {
			drawProgressBar(float64(i)/100.0, 40)
			time.Sleep(40 * time.Millisecond)
			if i < 100 {
				fmt.Print("\033[1A\r") // Remonte d'une ligne
			}
		}

		file, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer file.Close()

		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 95})
		if err != nil {
			return err
		}

	case "4": // JPEG qualité personnalisée
		if !strings.HasSuffix(strings.ToLower(outputPath), ".jpg") &&
			!strings.HasSuffix(strings.ToLower(outputPath), ".jpeg") {
			outputPath += ".jpg"
		}

		fmt.Println("💡 Guide des qualités JPEG:")
		fmt.Println("  • 10-30 = Très faible (petite taille, qualité médiocre)")
		fmt.Println("  • 30-50 = Faible (pour aperçus)")
		fmt.Println("  • 50-70 = Moyenne (usage web)")
		fmt.Println("  • 70-90 = Bonne (recommandé pour la plupart des usages)")
		fmt.Println("  • 90-100 = Excellente (gros fichiers, qualité maximale)")
		fmt.Println()
		
		qualityStr := readUserInput("Qualité (1-100)")
		quality, err := strconv.Atoi(qualityStr)
		if err != nil || quality < 1 || quality > 100 {
			errorMessage("Valeur invalide, utilisation de la qualité par défaut (75)")
			quality = 75
		}

		clearScreen()
		infoMessage(fmt.Sprintf("Conversion en JPEG (qualité %d)...", quality))

		for i := 0; i <= 100; i += 5 {
			drawProgressBar(float64(i)/100.0, 40)
			time.Sleep(35 * time.Millisecond)
			if i < 100 {
				fmt.Print("\033[1A\r")
			}
		}

		file, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer file.Close()

		err = jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
		if err != nil {
			return err
		}

	case "5":
		if !strings.HasSuffix(strings.ToLower(outputPath), ".gif") {
			outputPath += ".gif"
		}

		clearScreen()
		infoMessage("Conversion en GIF...")

		for i := 0; i <= 100; i += 5 {
			drawProgressBar(float64(i)/100.0, 40)
			time.Sleep(45 * time.Millisecond)
			if i < 100 {
				fmt.Print("\033[1A\r") 
			}
		}

		file, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer file.Close()

		// création palette couleurs
		var palette color.Palette
		if q, ok := img.(image.PalettedImage); ok {
			palette = q.ColorModel().(color.Palette)
		} else {
			palette = color.Palette{
				color.RGBA{0, 0, 0, 255},       // Noir
				color.RGBA{255, 255, 255, 255}, // Blanc
				color.RGBA{255, 0, 0, 255},     // Rouge
				color.RGBA{0, 255, 0, 255},     // Vert
				color.RGBA{0, 0, 255, 255},     // Bleu
				color.RGBA{255, 255, 0, 255},   // Jaune
				color.RGBA{255, 0, 255, 255},   // Magenta
				color.RGBA{0, 255, 255, 255},   // Cyan
			}
		}

		// conversion
		bounds := img.Bounds()
		palettedImg := image.NewPaletted(bounds, palette)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				palettedImg.Set(x, y, img.At(x, y))
			}
		}

		err = gif.Encode(file, palettedImg, &gif.Options{
			NumColors: len(palette),
			Quantizer: nil,
			Drawer:    nil,
		})

		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("option de conversion invalide")
	}

	successMessage(fmt.Sprintf("Image convertie avec succès: %s", outputPath))
	time.Sleep(1 * time.Second)
	return nil
}



type SquareEffect struct {
	X, Y, Size int
	Color      color.Color
}

func (s *SquareEffect) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	// copie
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			result.Set(x, y, img.At(x, y))
		}
	}

	// carré
	for y := s.Y; y < s.Y+s.Size && y < bounds.Max.Y; y++ {
		for x := s.X; x < s.X+s.Size && x < bounds.Max.X; x++ {
			if x >= bounds.Min.X && y >= bounds.Min.Y {
				result.Set(x, y, s.Color)
			}
		}
	}
	return result
}

type CircleEffect struct {
	CenterX, CenterY, Radius int
	Color                    color.Color
}

func (c *CircleEffect) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	// copie
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			result.Set(x, y, img.At(x, y))
		}
	}

	// cercle
	for y := c.CenterY - c.Radius; y <= c.CenterY+c.Radius; y++ {
		for x := c.CenterX - c.Radius; x <= c.CenterX+c.Radius; x++ {
			dx := x - c.CenterX
			dy := y - c.CenterY
			if dx*dx+dy*dy <= c.Radius*c.Radius {
				if x >= bounds.Min.X && x < bounds.Max.X && y >= bounds.Min.Y && y < bounds.Max.Y {
					result.Set(x, y, c.Color)
				}
			}
		}
	}
	return result
}



func saveImageEnhanced(img image.Image) error {
	clearScreen()
	drawBox("Sauvegarder l'image", []string{
		"Entrez le chemin où sauvegarder l'image modifiée",
		"",
		"📁 Extensions supportées: .png, .jpg, .jpeg",
		"💡 Astuce: Utilisez des noms explicites (ex: image_effet_sepia.png)",
		"⚠️ Attention: Un fichier existant sera écrasé",
	}, 80)
	fmt.Println()
	
	fmt.Println("💡 Suggestions de noms de fichiers:")
	fmt.Println("  • image_modifiee.png")
	fmt.Println("  • sortie/effet_sepia.jpg")
	fmt.Println("  • images/final_lumineux.png")
	fmt.Println("  • test_cercle_rouge.png")
	fmt.Println("  • backup/original_copie.jpg")
	fmt.Println()

	filePath := readUserInput("Chemin de sauvegarde (ex: output/mon_image.png)")

	ext := strings.ToLower(filepath.Ext(filePath))

	if ext == "" {
		warningMessage("Extension manquante, ajout automatique de .png")
		filePath += ".png"
		ext = ".png"
	}

	if _, err := os.Stat(filePath); err == nil {
		if !confirmAction("Le fichier existe déjà. L'écraser ?") {
			infoMessage("Sauvegarde annulée")
			return nil
		}
	}

	dir := filepath.Dir(filePath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("impossible de créer le dossier: %v", err)
		}
	}

	clearScreen()
	drawBox("Sauvegarde en cours", []string{
		"Fichier de destination: " + filePath,
		"Format: " + strings.ToUpper(ext[1:]),
		"",
		"⏳ Traitement en cours...",
	}, 80)
	fmt.Println()

	steps := []string{
		"Préparation du fichier...",
		"Encodage de l'image...",
		"Écriture sur disque...",
		"Finalisation...",
	}

	for i, step := range steps {
		drawProgressBarAnimated(float64(i)/float64(len(steps)), 50, step)
		time.Sleep(200 * time.Millisecond)
		if i < len(steps)-1 {
			fmt.Print("\033[1A\r")
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("impossible de créer le fichier: %v", err)
	}
	defer file.Close()

	switch ext {
	case ".png":
		err = png.Encode(file, img)
	case ".jpg", ".jpeg":
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	default:
		return fmt.Errorf("format non supporté: %s (utilisez .png, .jpg ou .jpeg)", ext)
	}

	if err != nil {
		return fmt.Errorf("erreur lors de l'encodage: %v", err)
	}

	drawProgressBarAnimated(1.0, 50, "Sauvegarde terminée")
	fmt.Println()

	successMessage(fmt.Sprintf("Image sauvegardée avec succès: %s", filePath))
	
	if stat, err := os.Stat(filePath); err == nil {
		infoMessage(fmt.Sprintf("Taille du fichier: %s", formatFileSize(stat.Size())))
	}
	
	time.Sleep(2 * time.Second)
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func convertImageEnhanced(img image.Image) (image.Image, error) {
	clearScreen()
	formatItems := []string{
		"PNG",
		"JPEG (qualité standard)",
		"JPEG (haute qualité)",
		"JPEG (qualité personnalisée)",
		"GIF",
		"Redimensionner l'image",
		"Afficher les métadonnées",
		"Retour",
	}

	drawBox("Conversion d'image", []string{
		"Choisissez le format vers lequel vous souhaitez convertir l'image",
		"ou une autre opération",
	}, 70)
	fmt.Println()

	drawMenu("Options disponibles", formatItems, []string{}, []string{}, -1, 50)

	choice := readUserInput("Choisissez une option")

	if choice == "8" || choice == "0" {
		return img, nil
	}

	if choice == "7" {
		readMetadata(img)
		return img, nil
	}

	if choice == "6" {
		clearScreen()
		drawBox("Redimensionnement d'image", []string{
			"Spécifiez les nouvelles dimensions de l'image",
			"(entrez 0 pour une dimension pour conserver le ratio)",
		}, 70)
		fmt.Println()

		bounds := img.Bounds()
		width := bounds.Max.X - bounds.Min.X
		height := bounds.Max.Y - bounds.Min.Y

		infoMessage(fmt.Sprintf("Dimensions actuelles: %d × %d pixels", width, height))
		
		fmt.Println("💡 Suggestions de redimensionnement:")
		fmt.Printf("  • Réduire de moitié: %d × %d\n", width/2, height/2)
		fmt.Printf("  • Doubler la taille: %d × %d\n", width*2, height*2)
		fmt.Printf("  • Format HD (16:9): 1920 × 1080\n")
		fmt.Printf("  • Format carré: %d × %d\n", min(width, height), min(width, height))
		fmt.Printf("  • Largeur fixe 800px: 800 × %d\n", (800*height)/width)
		fmt.Printf("  • Hauteur fixe 600px: %d × 600\n", (600*width)/height)
		fmt.Println("  • Entrez 0 pour une dimension pour conserver le ratio")
		fmt.Println()

		widthStr := readUserInput("Nouvelle largeur (pixels)")
		heightStr := readUserInput("Nouvelle hauteur (pixels)")

		newWidth, err1 := strconv.Atoi(widthStr)
		newHeight, err2 := strconv.Atoi(heightStr)

		if err1 != nil || err2 != nil || (newWidth < 0) || (newHeight < 0) {
			errorMessage("Dimensions invalides")
			time.Sleep(1 * time.Second)
			return nil, fmt.Errorf("dimensions invalides")
		}

		clearScreen()
		infoMessage("Redimensionnement en cours...")

		for i := 0; i <= 100; i += 2 {
			drawProgressBar(float64(i)/100.0, 40)
			time.Sleep(15 * time.Millisecond)
			if i < 100 {
				fmt.Print("\033[1A\r") 
			}
		}

		resizedImg := resizeImage(img, newWidth, newHeight)

		successMessage(fmt.Sprintf("Image redimensionnée avec succès: %d × %d pixels",
			resizedImg.Bounds().Max.X-resizedImg.Bounds().Min.X,
			resizedImg.Bounds().Max.Y-resizedImg.Bounds().Min.Y))

		time.Sleep(1 * time.Second)
		return resizedImg, nil
	}

	outputPath := readUserInput("Chemin du fichier de sortie (avec extension)")

	switch choice {
	case "1": 
		if !strings.HasSuffix(strings.ToLower(outputPath), ".png") {
			outputPath += ".png"
		}

		clearScreen()
		infoMessage("Conversion en PNG...")

		for i := 0; i <= 100; i += 5 {
			drawProgressBar(float64(i)/100.0, 40)
			time.Sleep(50 * time.Millisecond)
			if i < 100 {
				fmt.Print("\033[1A\r") 
			}
		}

		file, err := os.Create(outputPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		err = png.Encode(file, img)
		if err != nil {
			return nil, err
		}

	case "2": 
		if !strings.HasSuffix(strings.ToLower(outputPath), ".jpg") &&
			!strings.HasSuffix(strings.ToLower(outputPath), ".jpeg") {
			outputPath += ".jpg"
		}

		clearScreen()
		infoMessage("Conversion en JPEG (qualité standard)...")

		for i := 0; i <= 100; i += 5 {
			drawProgressBar(float64(i)/100.0, 40)
			time.Sleep(30 * time.Millisecond)
			if i < 100 {
				fmt.Print("\033[1A\r")
			}
		}

		file, err := os.Create(outputPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 75})
		if err != nil {
			return nil, err
		}

	case "3": 
		if !strings.HasSuffix(strings.ToLower(outputPath), ".jpg") &&
			!strings.HasSuffix(strings.ToLower(outputPath), ".jpeg") {
			outputPath += ".jpg"
		}

		clearScreen()
		infoMessage("Conversion en JPEG (haute qualité)...")

		for i := 0; i <= 100; i += 5 {
			drawProgressBar(float64(i)/100.0, 40)
			time.Sleep(40 * time.Millisecond)
			if i < 100 {
				fmt.Print("\033[1A\r")
			}
		}

		file, err := os.Create(outputPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 95})
		if err != nil {
			return nil, err
		}

	case "4":
		if !strings.HasSuffix(strings.ToLower(outputPath), ".jpg") &&
			!strings.HasSuffix(strings.ToLower(outputPath), ".jpeg") {
			outputPath += ".jpg"
		}

		fmt.Println("💡 Guide des qualités JPEG:")
		fmt.Println("  • 10-30 = Très faible (petite taille, qualité médiocre)")
		fmt.Println("  • 30-50 = Faible (pour aperçus)")
		fmt.Println("  • 50-70 = Moyenne (usage web)")
		fmt.Println("  • 70-90 = Bonne (recommandé pour la plupart des usages)")
		fmt.Println("  • 90-100 = Excellente (gros fichiers, qualité maximale)")
		fmt.Println()
		
		qualityStr := readUserInput("Qualité (1-100)")
		quality, err := strconv.Atoi(qualityStr)
		if err != nil || quality < 1 || quality > 100 {
			errorMessage("Valeur invalide, utilisation de la qualité par défaut (75)")
			quality = 75
		}

		clearScreen()
		infoMessage(fmt.Sprintf("Conversion en JPEG (qualité %d)...", quality))

		for i := 0; i <= 100; i += 5 {
			drawProgressBar(float64(i)/100.0, 40)
			time.Sleep(35 * time.Millisecond)
			if i < 100 {
				fmt.Print("\033[1A\r")
			}
		}

		file, err := os.Create(outputPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		err = jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
		if err != nil {
			return nil, err
		}

	case "5":
		if !strings.HasSuffix(strings.ToLower(outputPath), ".gif") {
			outputPath += ".gif"
		}

		clearScreen()
		infoMessage("Conversion en GIF...")

		for i := 0; i <= 100; i += 5 {
			drawProgressBar(float64(i)/100.0, 40)
			time.Sleep(45 * time.Millisecond)
			if i < 100 {
				fmt.Print("\033[1A\r")
			}
		}

		file, err := os.Create(outputPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		var palette color.Palette
		if q, ok := img.(image.PalettedImage); ok {
			palette = q.ColorModel().(color.Palette)
		} else {
			palette = color.Palette{
				color.RGBA{0, 0, 0, 255},       // Noir
				color.RGBA{255, 255, 255, 255}, // Blanc
				color.RGBA{255, 0, 0, 255},     // Rouge
				color.RGBA{0, 255, 0, 255},     // Vert
				color.RGBA{0, 0, 255, 255},     // Bleu
				color.RGBA{255, 255, 0, 255},   // Jaune
				color.RGBA{255, 0, 255, 255},   // Magenta
				color.RGBA{0, 255, 255, 255},   // Cyan
			}
		}

		bounds := img.Bounds()
		palettedImg := image.NewPaletted(bounds, palette)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				palettedImg.Set(x, y, img.At(x, y))
			}
		}

		err = gif.Encode(file, palettedImg, &gif.Options{
			NumColors: len(palette),
			Quantizer: nil,
			Drawer:    nil,
		})

		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("option de conversion invalide")
	}

	successMessage(fmt.Sprintf("Image convertie avec succès: %s", outputPath))
	time.Sleep(1 * time.Second)
	return nil, nil
}
