package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Les informations d'un dossier ou fichier
type FileInfo struct {
	Name    string
	IsDir   bool
	Size    int64
	ModTime string
}

// liste contenu répertoire avec filtrage images
func listDirectory(dirPath string) ([]FileInfo, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var files []FileInfo
	
	// Ajout répertoire parent si pas à la racine
	if dirPath != "." && dirPath != "/" {
		files = append(files, FileInfo{
			Name:  "..",
			IsDir: true,
		})
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Pour les fichiers cachés
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		fileInfo := FileInfo{
			Name:    entry.Name(),
			IsDir:   entry.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime().Format("2006-01-02 15:04"),
		}

		// Vérif Image
		if !entry.IsDir() && !isImageFile(entry.Name()) {
			continue
		}

		files = append(files, fileInfo)
	}

	// Trie alphabétique
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir && !files[j].IsDir {
			return true
		}
		if !files[i].IsDir && files[j].IsDir {
			return false
		}
		return files[i].Name < files[j].Name
	})

	return files, nil
}

// extension supportée
func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif"
}

// Affichage UX
func displayFileList(files []FileInfo, currentDir string) {
	clearScreen()
	
	drawBox("Navigation de fichiers", []string{
		"Répertoire actuel: " + currentDir,
		"Fichiers et dossiers disponibles (images uniquement):",
		"",
		"Navigation: tapez le nom ou le numéro",
		".. = répertoire parent | [Enter] = répertoire actuel",
	}, 70)
	
	fmt.Println()
	
	if len(files) == 0 {
		errorMessage("Aucun fichier image trouvé dans ce répertoire")
		return
	}

	// UX avec numérotation des dossiers
	for i, file := range files {
		var icon, sizeStr string
		
		if file.IsDir {
			icon = ColorBlue + "📁" + ColorReset
			sizeStr = "<DIR>"
		} else {
			icon = ColorGreen + "🖼️" + ColorReset
			sizeStr = formatFileSize(file.Size)
		}
		
		fmt.Printf("%s [%d] %-30s %s %10s %s\n",
			icon,
			i+1,
			file.Name,
			ColorCyan+file.ModTime+ColorReset,
			sizeStr,
			ColorReset,
		)
	}
	
	fmt.Println()
}

// format taille fichier
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// selectionner fichier
func navigateToFile() (string, error) {
	currentDir := "."
	
	for {
		// chemin absolu
		absPath, _ := filepath.Abs(currentDir)
		
		files, err := listDirectory(currentDir)
		if err != nil {
			return "", fmt.Errorf("impossible de lire le répertoire %s: %v", currentDir, err)
		}

		displayFileList(files, absPath)
		
		input := readUserInput("Sélection (nom/numéro/Enter pour '.') ou 'q' pour annuler")
		
		if input == "q" || input == "Q" {
			return "", fmt.Errorf("navigation annulée")
		}
		
		if input == "" {
			input = "."
		}
		
		// vérif number
		if num := parseNumber(input); num > 0 && num <= len(files) {
			selectedFile := files[num-1]
			
			if selectedFile.IsDir {
				if selectedFile.Name == ".." {
					currentDir = filepath.Dir(currentDir)
				} else {
					currentDir = filepath.Join(currentDir, selectedFile.Name)
				}
				continue
			} else {
				return filepath.Join(currentDir, selectedFile.Name), nil
			}
		}
		
		found := false
		for _, file := range files {
			if strings.EqualFold(file.Name, input) {
				if file.IsDir {
					if file.Name == ".." {
						currentDir = filepath.Dir(currentDir)
					} else {
						currentDir = filepath.Join(currentDir, file.Name)
					}
					found = true
					break
				} else {
					return filepath.Join(currentDir, file.Name), nil
				}
			}
		}
		
		if !found {
			errorMessage("Fichier ou dossier non trouvé. Utilisez le nom exact ou le numéro.")
			readUserInput("Appuyez sur Entrée pour continuer")
		}
	}
}

// parse une chaine --> nombre
func parseNumber(s string) int {
	num := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0
		}
		num = num*10 + int(r-'0')
	}
	return num
} 