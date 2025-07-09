# GoImage - Convertisseur et éditeur d'images en TUI

## Présentation

GoImage est un outil en ligne de commande (TUI) pour convertir, éditer et dessiner sur des images, écrit en Go **sans dépendances externes** (hors standard et x/image). Il permet d'appliquer des effets, de dessiner des formes, de modifier des pixels, etc.

---

## Architecture du projet

```
goimage/
│
├── cmd/
│   └── goimage/
│       └── main.go         # Point d’entrée, interface TUI, logique utilisateur
│
├── pkg/
│   └── effects/
│       ├── interface.go    # Interface Effect
│       ├── negative.go     # Effet négatif
│       ├── grayscale.go    # Effet niveaux de gris
│       ├── shapes.go       # Formes géométriques (carré, cercle, triangle, ligne)
│       └── ...             # Autres effets (séparés par type)
│
├── test/
│   └── test_image.png      # Images de test
│
├── go.mod
├── go.sum
└── README.md
```

---

## Rôle des fichiers principaux

### `cmd/goimage/main.go`

- **Rôle** : Point d’entrée du programme, gère l’interface utilisateur (TUI), les menus, la navigation, la saisie utilisateur, et appelle les fonctions d’effets/dessin.
- **Contenu** :
  - Menus (conversion, édition, dessin…)
  - Fonctions pour demander les paramètres à l’utilisateur
  - Appels à la librairie d’effets (ex : `effects.Apply(img, params)`)

### `pkg/effects/interface.go`

- **Rôle** : Définit l’interface `Effect` que tous les effets doivent implémenter.
- **Contenu** :

```go
package effects
import "image"
type Effect interface {
    Apply(img image.Image) image.Image
    Name() string
    Description() string
}
```

### `pkg/effects/negative.go`, `grayscale.go`, ...

- **Rôle** : Chaque fichier contient un effet (ou une famille d’effets) spécifique.
- **Exemple** :

```go
package effects
import (
    "image"
    "image/color"
)
type NegativeEffect struct{}
func (n *NegativeEffect) Name() string { return "Négatif" }
func (n *NegativeEffect) Description() string { return "Inverse toutes les couleurs de l'image" }
func (n *NegativeEffect) Apply(img image.Image) image.Image { ... }
```

### `pkg/effects/shapes.go`

- **Rôle** : Contient les effets pour dessiner des formes géométriques (carré, cercle, triangle, ligne).
- **Exemple** :

```go
package effects
import (
    "image"
    "image/color"
)
type SquareEffect struct { X, Y, Size int; Color color.Color }
func (s *SquareEffect) Apply(img image.Image) image.Image { ... }
```

---

## Fonctionnement de la librairie d’effets

- **Interface** : Tous les effets implémentent `Effect`.
- **Application** :
  - L’utilisateur choisit un effet dans le TUI.
  - Le programme crée une instance de l’effet avec les bons paramètres.
  - Il appelle `Apply(img)` pour obtenir l’image modifiée.
  - Le résultat est sauvegardé.
- **Ajout d’un effet** :
  - Créer une nouvelle struct qui implémente `Effect` dans un fichier séparé.
  - Ajouter l’effet dans le menu du TUI.

---

## Découpage recommandé pour gros fichiers

Quand un fichier devient trop gros, découpe-le ainsi :

- **interface.go** : interface `Effect`
- **negative.go** : effet négatif
- **grayscale.go** : effet niveaux de gris
- **shapes.go** : formes géométriques
- **brightness.go**, **contrast.go**, etc. : autres effets
- **utils.go** : fonctions utilitaires (min, max, abs…)

Tous les fichiers d’un même dossier et du même package sont compilés ensemble.

---

## Exemple d’utilisation dans le main.go

```go
import "github.com/nirdeo/goimage/pkg/effects"

// Création d’un effet carré
square := &effects.SquareEffect{X: 10, Y: 10, Size: 50, Color: color.RGBA{255,0,0,255}}
result := square.Apply(img)
```

---

## Bonnes pratiques

- Un fichier = un type d’effet ou une famille d’effets
- Utilise des interfaces pour la flexibilité
- Passe les paramètres via des structs
- Mets les fonctions utilitaires dans un fichier à part
- Commente chaque effet, chaque fonction

---

## Pour aller plus loin

- Ajoute d’autres effets dans des fichiers séparés
- Ajoute des outils de dessin interactif dans `pkg/drawing/`
- Utilise la même logique de découpage pour le TUI si besoin (ex : menus dans un fichier, gestion des entrées dans un autre…)

---

## Auteur

Projet GoImage, architecture modulaire et évolutive, inspirée par les bonnes pratiques Go.
