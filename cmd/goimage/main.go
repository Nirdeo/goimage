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
)

// Interface Effect que tous les effets doivent implémenter
type Effect interface {
	Apply(img image.Image) image.Image
	Name() string
	Description() string
}

// Effet Négatif
type NegativeEffect struct{}

func (n *NegativeEffect) Name() string        { return "Négatif" }
func (n *NegativeEffect) Description() string { return "Inverse toutes les couleurs de l'image" }
func (n *NegativeEffect) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			// Inverser les couleurs en soustrayant de 65535 (max pour uint16)
			result.Set(x, y, color.RGBA{
				uint8(255 - uint8(r>>8)),
				uint8(255 - uint8(g>>8)),
				uint8(255 - uint8(b>>8)),
				uint8(a >> 8),
			})
		}
	}
	return result
}

// Effet Niveaux de gris
type GrayscaleEffect struct{}

func (g *GrayscaleEffect) Name() string        { return "Niveaux de gris" }
func (g *GrayscaleEffect) Description() string { return "Convertit l'image en noir et blanc" }
func (g *GrayscaleEffect) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// Formule standard pour convertir RGB en niveaux de gris
			grayValue := uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256)
			result.Set(x, y, color.Gray{Y: grayValue})
		}
	}
	return result
}

// Effet Sépia
type SepiaEffect struct{}

func (s *SepiaEffect) Name() string        { return "Sépia" }
func (s *SepiaEffect) Description() string { return "Applique un effet sépia vintage à l'image" }
func (s *SepiaEffect) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			// Convertir à 8 bits par canal
			r8 := float64(r >> 8)
			g8 := float64(g >> 8)
			b8 := float64(b >> 8)

			// Appliquer la transformation sépia
			newR := uint8(min(0.393*r8+0.769*g8+0.189*b8, 255))
			newG := uint8(min(0.349*r8+0.686*g8+0.168*b8, 255))
			newB := uint8(min(0.272*r8+0.534*g8+0.131*b8, 255))

			result.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
		}
	}
	return result
}

// Fonction utilitaire min
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Effet Carré
type SquareEffect struct {
	X, Y, Size int
	Color      color.Color
}

func (s *SquareEffect) Name() string        { return "Carré" }
func (s *SquareEffect) Description() string { return "Dessine un carré rempli" }
func (s *SquareEffect) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	// Copie de l'image originale
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			result.Set(x, y, img.At(x, y))
		}
	}

	// Dessin du carré
	for y := s.Y; y < s.Y+s.Size && y < bounds.Max.Y; y++ {
		for x := s.X; x < s.X+s.Size && x < bounds.Max.X; x++ {
			if x >= bounds.Min.X && y >= bounds.Min.Y {
				result.Set(x, y, s.Color)
			}
		}
	}
	return result
}

// Effet Cercle
type CircleEffect struct {
	CenterX, CenterY, Radius int
	Color                    color.Color
}

func (c *CircleEffect) Name() string        { return "Cercle" }
func (c *CircleEffect) Description() string { return "Dessine un cercle rempli" }
func (c *CircleEffect) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	// Copie de l'image originale
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			result.Set(x, y, img.At(x, y))
		}
	}

	// Dessin du cercle
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

func main() {
	StartTUI()
}

// StartTUI lance l'interface utilisateur en ligne de commande
func StartTUI() {
	var img image.Image
	var err error
	var currentFilePath string

	menuItems := []string{
		"Charger une image",
		"Appliquer un effet",
		"Dessiner une forme",
		"Convertir l'image",
		"Sauvegarder l'image",
		"Quitter",
	}

	for {
		drawHeader()
		drawMenu("Menu Principal", menuItems, -1, 50)
		drawFooter()

		choice := prompt("Choisissez une option")

		switch choice {
		case "1":
			img, currentFilePath, err = loadImage()
			if err != nil {
				errorMessage(fmt.Sprintf("Erreur lors du chargement de l'image: %v", err))
			} else if img == nil {
				errorMessage("L'image n'a pas pu être chargée correctement.")
			} else {
				successMessage(fmt.Sprintf("Image chargée: %s (%dx%d)",
					currentFilePath, img.Bounds().Dx(), img.Bounds().Dy()))
			}
			time.Sleep(1 * time.Second)
		case "2":
			if img == nil {
				errorMessage("Veuillez d'abord charger une image.")
				time.Sleep(1 * time.Second)
				continue
			}
			img = applyEffect(img)
		case "3":
			if img == nil {
				errorMessage("Veuillez d'abord charger une image.")
				time.Sleep(1 * time.Second)
				continue
			}
			img = drawShape(img)
		case "4":
			if img == nil {
				errorMessage("Veuillez d'abord charger une image.")
				time.Sleep(1 * time.Second)
				continue
			}
			err := convertImage(img)
			if err != nil {
				errorMessage(fmt.Sprintf("Erreur lors de la conversion: %v", err))
				time.Sleep(1 * time.Second)
			}
		case "5":
			if img == nil {
				errorMessage("Veuillez d'abord charger une image.")
				time.Sleep(1 * time.Second)
				continue
			}
			err := saveImage(img)
			if err != nil {
				errorMessage(fmt.Sprintf("Erreur lors de la sauvegarde: %v", err))
				time.Sleep(1 * time.Second)
			}
		case "6", "q", "Q":
			clearScreen()
			successMessage("Merci d'avoir utilisé GoImage. Au revoir!")
			return
		default:
			errorMessage("Option invalide, veuillez réessayer.")
			time.Sleep(1 * time.Second)
		}
	}
}

// readUserInput lit une entrée utilisateur et renvoie la chaîne de caractères
func readUserInput(promptMessage string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(ColorYellow + "? " + promptMessage + ColorReset + " ")
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// loadImage demande à l'utilisateur un chemin d'image et tente de la charger
func loadImage() (image.Image, string, error) {
	clearScreen()
	drawBox("Charger une image", []string{
		"Entrez le chemin de l'image que vous souhaitez charger.",
		"Formats supportés: PNG, JPEG, GIF",
	}, 60)
	fmt.Println()

	filePath := readUserInput("Chemin de l'image à charger")

	// Affichage d'une barre de progression
	infoMessage("Chargement de l'image en cours...")

	// Vérification que le fichier existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, "", fmt.Errorf("le fichier n'existe pas: %s", filePath)
	}

	// Simulation de progression
	for i := 0; i <= 100; i += 10 {
		drawProgressBar(float64(i)/100.0, 40)
		time.Sleep(50 * time.Millisecond)
		if i < 100 {
			fmt.Print("\033[1A\r") // Remonte d'une ligne
		}
	}

	// Ouverture du fichier
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("impossible d'ouvrir le fichier: %v", err)
	}
	defer file.Close()

	// Tentative de décodage de l'image
	infoMessage("Décodage de l'image...")
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", fmt.Errorf("impossible de décoder l'image (format non supporté?): %v", err)
	}

	successMessage(fmt.Sprintf("Image décodée avec succès, format: %s", format))

	// Vérification que l'image n'est pas nil
	if img == nil {
		return nil, "", fmt.Errorf("l'image a été décodée comme nil")
	}

	return img, filePath, nil
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
	}, 60)

	readUserInput("Appuyez sur Entrée pour continuer")
	return
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
	}, 60)
	fmt.Println()

	drawMenu("Options disponibles", formatItems, -1, 50)

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
		}, 60)
		fmt.Println()

		// Obtenir les dimensions actuelles
		bounds := img.Bounds()
		width := bounds.Max.X - bounds.Min.X
		height := bounds.Max.Y - bounds.Min.Y

		infoMessage(fmt.Sprintf("Dimensions actuelles: %d × %d pixels", width, height))

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

		qualityStr := readUserInput("Qualité (1-100)")
		quality, err := strconv.Atoi(qualityStr)
		if err != nil || quality < 1 || quality > 100 {
			errorMessage("Valeur invalide, utilisation de la qualité par défaut (75)")
			quality = 75
		}

		clearScreen()
		infoMessage(fmt.Sprintf("Conversion en JPEG (qualité %d)...", quality))

		// Simulation de progression
		for i := 0; i <= 100; i += 5 {
			drawProgressBar(float64(i)/100.0, 40)
			time.Sleep(35 * time.Millisecond)
			if i < 100 {
				fmt.Print("\033[1A\r") // Remonte d'une ligne
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

	case "5": // GIF
		if !strings.HasSuffix(strings.ToLower(outputPath), ".gif") {
			outputPath += ".gif"
		}

		clearScreen()
		infoMessage("Conversion en GIF...")

		// Simulation de progression
		for i := 0; i <= 100; i += 5 {
			drawProgressBar(float64(i)/100.0, 40)
			time.Sleep(45 * time.Millisecond)
			if i < 100 {
				fmt.Print("\033[1A\r") // Remonte d'une ligne
			}
		}

		file, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer file.Close()

		// Création d'une palette de couleurs pour le GIF
		var palette color.Palette
		if q, ok := img.(image.PalettedImage); ok {
			palette = q.ColorModel().(color.Palette)
		} else {
			// Palette par défaut si l'image source n'est pas en couleurs indexées
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

		// Conversion de l'image en GIF (version basique, sans animation)
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

// applyEffect permet à l'utilisateur de choisir et d'appliquer un effet
func applyEffect(img image.Image) image.Image {
	clearScreen()
	effectItems := []string{
		"Négatif",
		"Niveaux de gris",
		"Sépia",
		"Retour",
	}

	drawBox("Appliquer un effet", []string{
		"Choisissez un effet à appliquer à l'image",
	}, 60)
	fmt.Println()

	drawMenu("Effets disponibles", effectItems, -1, 50)

	choice := readUserInput("Choisissez un effet")

	var effect Effect

	switch choice {
	case "1":
		effect = &NegativeEffect{}
	case "2":
		effect = &GrayscaleEffect{}
	case "3":
		effect = &SepiaEffect{}
	case "4", "0":
		return img
	default:
		errorMessage("Option invalide, retour au menu principal.")
		time.Sleep(1 * time.Second)
		return img
	}

	clearScreen()
	infoMessage(fmt.Sprintf("Application de l'effet %s...", effect.Name()))

	// Simulation de progression
	for i := 0; i <= 100; i += 2 {
		drawProgressBar(float64(i)/100.0, 40)
		time.Sleep(20 * time.Millisecond)
		if i < 100 {
			fmt.Print("\033[1A\r") // Remonte d'une ligne
		}
	}

	modifiedImg := effect.Apply(img)
	successMessage("Effet appliqué avec succès!")
	time.Sleep(1 * time.Second)
	return modifiedImg
}

// drawShape permet à l'utilisateur de dessiner une forme sur l'image
func drawShape(img image.Image) image.Image {
	clearScreen()
	shapeItems := []string{
		"Carré",
		"Cercle",
		"Retour",
	}

	drawBox("Dessiner une forme", []string{
		"Choisissez une forme à dessiner sur l'image",
	}, 60)
	fmt.Println()

	drawMenu("Formes disponibles", shapeItems, -1, 50)

	choice := readUserInput("Choisissez une forme")

	if choice == "3" || choice == "0" {
		return img
	}

	// Définition de la couleur
	colorInput := readUserInput("Couleur (format R,G,B) - Ex: 255,0,0 pour rouge")
	colorParts := strings.Split(colorInput, ",")
	if len(colorParts) != 3 {
		errorMessage("Format de couleur invalide, utilisation du rouge par défaut.")
		colorParts = []string{"255", "0", "0"}
	}

	r, _ := strconv.Atoi(colorParts[0])
	g, _ := strconv.Atoi(colorParts[1])
	b, _ := strconv.Atoi(colorParts[2])
	shapeColor := color.RGBA{uint8(r), uint8(g), uint8(b), 255}

	switch choice {
	case "1": // Carré
		clearScreen()
		drawBox("Paramètres du carré", []string{
			"Définissez la position et la taille du carré",
		}, 60)
		fmt.Println()

		xStr := readUserInput("Position X")
		yStr := readUserInput("Position Y")
		sizeStr := readUserInput("Taille")

		x, _ := strconv.Atoi(xStr)
		y, _ := strconv.Atoi(yStr)
		size, _ := strconv.Atoi(sizeStr)

		clearScreen()
		infoMessage("Dessin d'un carré en cours...")

		// Simulation de progression
		for i := 0; i <= 100; i += 5 {
			drawProgressBar(float64(i)/100.0, 40)
			time.Sleep(30 * time.Millisecond)
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
		time.Sleep(1 * time.Second)
		return modifiedImg

	case "2": // Cercle
		clearScreen()
		drawBox("Paramètres du cercle", []string{
			"Définissez la position et le rayon du cercle",
		}, 60)
		fmt.Println()

		xStr := readUserInput("Centre X")
		yStr := readUserInput("Centre Y")
		radiusStr := readUserInput("Rayon")

		x, _ := strconv.Atoi(xStr)
		y, _ := strconv.Atoi(yStr)
		radius, _ := strconv.Atoi(radiusStr)

		clearScreen()
		infoMessage("Dessin d'un cercle en cours...")

		// Simulation de progression
		for i := 0; i <= 100; i += 5 {
			drawProgressBar(float64(i)/100.0, 40)
			time.Sleep(30 * time.Millisecond)
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
		time.Sleep(1 * time.Second)
		return modifiedImg

	default:
		errorMessage("Option invalide, retour au menu principal.")
		time.Sleep(1 * time.Second)
		return img
	}
}

// saveImage permet à l'utilisateur de sauvegarder l'image modifiée
func saveImage(img image.Image) error {
	clearScreen()
	drawBox("Sauvegarder l'image", []string{
		"Entrez le chemin où sauvegarder l'image",
		"Extensions supportées: .png, .jpg, .jpeg",
	}, 60)
	fmt.Println()

	filePath := readUserInput("Chemin de sauvegarde")

	// Déterminer le format de sortie en fonction de l'extension
	ext := strings.ToLower(filepath.Ext(filePath))

	if ext == "" {
		errorMessage("Extension manquante, veuillez spécifier .png ou .jpg")
		time.Sleep(1 * time.Second)
		return fmt.Errorf("extension manquante")
	}

	clearScreen()
	infoMessage("Sauvegarde de l'image en cours...")

	// Simulation de progression
	for i := 0; i <= 100; i += 2 {
		drawProgressBar(float64(i)/100.0, 40)
		time.Sleep(15 * time.Millisecond)
		if i < 100 {
			fmt.Print("\033[1A\r") // Remonte d'une ligne
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	switch ext {
	case ".png":
		infoMessage("Encodage au format PNG...")
		err = png.Encode(file, img)
	case ".jpg", ".jpeg":
		infoMessage("Encodage au format JPEG...")
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	default:
		return fmt.Errorf("format non supporté: %s (utilisez .png, .jpg ou .jpeg)", ext)
	}

	if err != nil {
		errorMessage(fmt.Sprintf("Erreur lors de la sauvegarde: %v", err))
		time.Sleep(1 * time.Second)
		return err
	}

	successMessage(fmt.Sprintf("Image sauvegardée avec succès: %s", filePath))
	time.Sleep(1 * time.Second)
	return nil
}
