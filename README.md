# Sudoku-GO
 Ce projet est un jeu de sudoku en Golang. L'ordinateur génère des grilles de Sudoku et permet à l'utilisateur de les résoudre à l'aide d'une interface graphique.

 - Le **Sudoku** est un jeu où le but est de remplir un grille avec des chiffres allant de 1 à 9 en respectant certaines conditions.

---

## Règles du jeu :

1. **Grille** : Le jeu se joue sur une grille de 9x9 cases, divisée en 9 sous grilles de 3x3 cases appelées "régions".

2. **Remplissage des cases** : L'objectif principal est de remplir chaque case vide avec un chiffre de 1 à 9.

3. **Règles de remplissage** : 
    - Chaque chiffre de 1 à 9 doit apparaître **exactement une fois** dans chaque **ligne** de la grille.
    - Chaque chiffre de 1 à 9 doit apparaître **exactement une fois** dans chaque **colonne** de la grille.
    - Chaque chiffre de 1 à 9 doit apparaître **exactement une fois** dans chaque **région** de la grille (c'est-à-dire, chaque sous grille de 3x3).

4. **Pré-remplissage** : Certaines cases de la grille sont pré-remplies avec des chiffres. Ces chiffres sont des **indices**. Ces cases ne peuvent pas être modifiées par l'utilisateur.

5. **Solution unique** : ⚠️ **Attention** : La grille est conçue pour que, si correctement résolue, il n'y ait qu'une seule réponse possible par case.
---
## Objectif

L'objectif est de remplir toute la grille en respectant les règles énoncées ci-dessus.

> ⚠️ **Pas de panique !!** le jeu vous informe si jamais vous vous êtes trompé.
---
## Conseils : 
- **Commencer par les cases les plus évidentes** : Cherchez les cases où une seule réponse est possible.
- **Procédez par élimination** : Si une case ne peut pas être remplie avec certains chiffres en raison des règles, éliminez ces chiffres comme options pour cette case.
- **Vérifier souvent** : Assurez-vous que vos réponses respectent toujours les règles du Sudoku à chaque étape.
---
## Prérequis

Avant de pouvoir utilisez ce projet, veillez à bien installer tout les prérequis.

### Prérequis pour éxécuter l'application

1. **Installer Go** : 
    - Vous devez installer [GO (Golang)](https://golang.org/dl/). Assurez-vous que la version installée est **1.23.5** ou une version plus récente.

2. **Dépendances du projet** :
    - Ce projet utilise une bibliothèque spécifique pour la gestion de l'interface graphique.

    En clonant le projet, vous serez amené à installer divers modules et dépendances.

## Installation
### 1. Cloner le dépot, installer les dépendances & compiler (si nécessaire)

```bash
git clone https://github.com/rayourama/Sudoku-GO.git
cd Sudoku-GO/details
go mod tidy
go build -o ../sudoku.exe main.go

```
> ⚠️ **Remarque importante** : Assurez-vous d'avoir installé Go avant de tenter de compiler le projet. Si vous avez des erreurs, veuillez également vérifier que Go soit correctement configuré dans les variables d'environnement.

## Lancement de l'application :
- Il ne vous reste plus qu'à lancer le fichier exécutable qui a été généré

## Crédits 
- [Ebiten](https://github.com/hajimehoshi/ebiten) - Pour la création de l'interface graphique.
- [Go Image](https://pkg.go.dev/golang.org/x/image) - Pour la gestion des images et des polices de caractères.
