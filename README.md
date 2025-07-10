# GoImage - Convertisseur et éditeur d'images en TUI

## Présentation

GoImage est un outil en ligne de commande (TUI) pour convertir, éditer et dessiner sur des images, écrit en Go **sans aucune dépendance externe**, uniquement avec les bibliothèques standard. Il permet d'appliquer des effets, de dessiner des formes, de redimensionner des images et de convertir entre formats.

---

## Architecture du projet

```
goimage/
│
├── cmd/
│   └── goimage/
│       ├── main.go         # Point d'entrée, effets, logique de traitement d'images
│       └── tui.go          # Interface utilisateur TUI (menus, couleurs, layout)
│
├── test/
│   ├── test_image.png      # Image de test
│   └── fond_blanc.png      # Image de test
│
├── go.mod                  # Configuration du module Go (sans dépendances)
└── README.md
```

---

## Fonctionnalités

- **Interface utilisateur** : TUI colorée et interactive inspirée de Bubble Tea
- **Manipulation d'images** :
  - Chargement d'images (PNG, JPEG, GIF)
  - Sauvegarde dans différents formats
  - Affichage des métadonnées d'image
- **Effets** :
  - Négatif
  - Niveaux de gris
  - Sépia
- **Dessin** :
  - Carré
  - Cercle
- **Conversion** :
  - Entre PNG, JPEG (plusieurs niveaux de qualité) et GIF
  - Redimensionnement avec préservation du ratio

---

## Utilisation

### Installation

Pour compiler le projet en un exécutable unique :

```bash
go build -o goimage cmd/goimage/main.go cmd/goimage/tui.go
```

Puis exécutez le programme avec :

```bash
./goimage
```

### Menu principal

- **Charger une image** : Ouvrir une image depuis le système de fichiers
- **Appliquer un effet** : Appliquer un effet (négatif, niveaux de gris, sépia)
- **Dessiner une forme** : Ajouter une forme à l'image (carré, cercle)
- **Convertir l'image** : Modifier le format ou redimensionner
- **Sauvegarder l'image** : Enregistrer l'image modifiée
- **Quitter** : Fermer l'application

---

## Implémentation technique

### Interface Effect

Tous les effets implémentent l'interface `Effect` :

```go
type Effect interface {
    Apply(img image.Image) image.Image
    Name() string
    Description() string
}
```

### Exemple d'effet (Négatif)

```go
type NegativeEffect struct{}

func (n *NegativeEffect) Name() string { return "Négatif" }
func (n *NegativeEffect) Description() string { return "Inverse toutes les couleurs de l'image" }
func (n *NegativeEffect) Apply(img image.Image) image.Image {
    bounds := img.Bounds()
    result := image.NewRGBA(bounds)
    
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r, g, b, a := img.At(x, y).RGBA()
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
```

### Interface TUI

L'interface utilisateur est implémentée avec des séquences d'échappement ANSI pour les couleurs et les bordures, sans dépendance externe.

```go
// Exemple d'affichage d'un cadre coloré
func drawBox(title string, content []string, width int) {
    // Ligne supérieure
    fmt.Print(ColorCyan)
    fmt.Print("╭")
    // ... suite du code ...
    fmt.Println("╯" + ColorReset)
}
```

---

## Limitations et bonnes pratiques

- Formats supportés limités à JPEG, PNG et GIF
- Le redimensionnement utilise l'algorithme du plus proche voisin (rapide mais moins précis)
- Utilisez les barres de progression pour suivre l'état des opérations longues
- Préférez des images de taille raisonnable (<10 Mo) pour de meilleures performances

---

## Caractéristiques techniques

- Écrit en Go natif (100% bibliothèque standard)
- Aucune dépendance externe
- Interface TUI colorée et intuitive
- Algorithmes d'effets implémentés manuellement
- Traitement pixel par pixel pour une compatibilité maximale
- Compatible avec les formats d'image standard

---

## Pour aller plus loin

- Ajouter d'autres effets (contraste, luminosité, flou...)
- Améliorer l'algorithme de redimensionnement (bilinéaire, bicubique)
- Ajouter des outils de sélection de zones
- Implémenter des filtres de convolution
- Ajouter des outils de dessin plus avancés (ligne, triangle, texte)

---

## Auteur

Projet GoImage, architecture modulaire et évolutive, inspirée par les bonnes pratiques Go.
