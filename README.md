# GoImage - Éditeur d'Images TUI 🎨

## 🚀 Présentation

GoImage est un éditeur d'images TUI (Terminal User Interface) écrit en Go **100% natif** sans dépendances externes.

**Fonctionnalités principales :**
- 🔍 Navigation de fichiers interactive
- ✨ 5 effets d'image (négatif, gris, sépia, luminosité, contraste)
- 🔶 Dessin de formes (carré, cercle)
- 🔄 Conversion multi-formats (PNG, JPEG, GIF)
- 💡 Système d'aide contextuel ('h')
- 📊 Barres de progression animées

---

## 🏗️ Architecture du Projet

```
goimage/
│
├── cmd/goimage/
│   ├── main.go         # Logique métier, effets, workflows
│   ├── tui.go          # Interface TUI (couleurs, menus, progression)
│   └── fileutils.go    # Navigation de fichiers interactive
│
├── pkg/effects/
│   ├── interface.go    # Interface Effect commune
│   ├── negative.go     # Effet négatif
│   ├── grayscale.go    # Conversion niveaux de gris
│   ├── sepia.go        # Effet sépia vintage
│   ├── brightness.go   # Ajustement luminosité
│   ├── contrast.go     # Ajustement contraste
│   └── shapes.go       # Formes géométriques
│
├── test/
│   ├── test_image.png  # Image de test
│   └── fond_blanc.png  # Image de test
│
├── go.mod              # Configuration module Go (sans dépendances)
└── README.md
```

---

## 🚀 Démarrage Rapide

### Installation et Test

```bash
# Compilation
go build -o goimage ./cmd/goimage/

# Lancement
./goimage

# Test avec image fournie
# 1. Choisir option 1 (Charger une image)
# 2. Navigation interactive → sélectionner test/test_image.png
# 3. Appliquer effet sépia (option 2 → 3)
# 4. Sauvegarder (option 5 → test_sepia.png)
```

### Workflow Principal

1. **📥 Charger** : Option 1 → Navigation interactive → `test/test_image.png`
2. **✨ Appliquer effet** : Option 2 → Choisir un effet
3. **🔶 Dessiner forme** : Option 3 → Carré/Cercle (optionnel)
4. **💾 Sauvegarder** : Option 5 → Nom du fichier

### Raccourcis Clavier

- **1-6** : Sélection options
- **h** : Aide contextuelle
- **q** : Quitter

---

## 🎯 Fonctionnalités

### Effets d'Image
- **Négatif** : Inversion couleurs
- **Niveaux de gris** : Conversion N&B
- **Sépia** : Effet vintage
- **Luminosité** : Paramétrable (0.5-3.0)
- **Contraste** : Paramétrable (0.5-3.0)

### Formes
- **Carré** : Position X,Y + taille
- **Cercle** : Centre X,Y + rayon
- **Couleurs RGB** : Format `255,0,0` (rouge)

### Conversion
- **PNG** : Qualité max, transparence
- **JPEG** : Qualité 75/95/personnalisée
- **GIF** : Palette optimisée
- **Redimensionnement** : Préservation ratio

---

## 🛠️ Implémentation

### Interface Effect
```go
type Effect interface {
    Apply(img image.Image) image.Image
    Name() string
    Description() string
}
```

### Ajouter un Effet
```go
// 1. Créer pkg/effects/nouvel_effet.go
type NouvelEffect struct { /* paramètres */ }

func (e *NouvelEffect) Apply(img image.Image) image.Image {
    // Algorithme de traitement
}

// 2. Ajouter dans main.go applyEffectEnhanced()
case "6":
    effect = &effects.NouvelEffect{}
```

### Navigation de Fichiers
```go
type FileInfo struct {
    Name    string
    IsDir   bool
    Size    int64
    ModTime string
}
```

---

## 📋 Caractéristiques Techniques

- **100% Go natif** - Zéro dépendance externe
- **Formats supportés** : PNG, JPEG, GIF
- **Compatibilité** : Windows, macOS, Linux
- **Terminal** : Unicode et couleurs ANSI
- **Performance** : Algorithmes optimisés pixel par pixel

---

## 💡 Exemples d'Usage

### Test Rapide
```bash
./goimage
# 1 → Navigation → test/test_image.png
# 2 → 3 (Sépia)
# 5 → test_sepia.png
```

### Workflow Complet
```bash
# 1. Charger image HD
# 2. Luminosité (factor: 1.2)
# 3. Contraste (factor: 1.5)
# 4. Cercle rouge au centre (255,0,0)
# 5. Redimensionner HD (1920×1080)
# 6. Sauvegarder JPEG qualité 95
```
