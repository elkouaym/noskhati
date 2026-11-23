# Noskhati

Avant de ranger une clé USB, vérifiez quels contenus du dossier de cours se retrouvent dans la sauvegarde, même si vous avez renommé ou déplacé des fichiers. Le contrôle lit les deux dossiers hors ligne. Il affiche un rapport et laisse les fichiers en place.

## Essayer

Installez Go 1.20 ou plus récent. Ouvrez un terminal dans le dossier du projet.

```sh
go run . --demo
go test -v ./...
go run . "chemin/vers/cours" "chemin/vers/cle"
```

La démo contient quatre fichiers source fictifs. Le plan est identique. Le cours apparaît sous archives/chapitre.txt. Les notes ont changé et l'exercice manque. Résultat attendu : **2/4 fichiers source ont un contenu retrouvé**.

## Lire le rapport

Le programme calcule une empreinte SHA-256 pour chaque fichier ordinaire et compare d'abord le chemin relatif, puis le contenu dans toute la sauvegarde. AILLEURS indique le premier chemin identique par ordre alphabétique. DIFFÉRENT indique que le chemin existe mais que le contenu source est introuvable dans la sauvegarde. Le bilan porte sur les contenus des fichiers source. Une même copie peut correspondre à plusieurs originaux identiques.

Utilisez deux dossiers distincts et stables pendant le contrôle. Les sous-dossiers sont parcourus. Les liens symboliques internes, dossiers vides, permissions et dates de modification sont exclus. Une erreur de lecture interrompt le contrôle. Les fichiers sont lus par blocs et seules les empreintes et les chemins restent en mémoire.

## Inspirations

[Czkawka](https://github.com/qarmin/czkawka) recherche des doublons notamment par empreinte. Sa [version 12.0.2](https://github.com/qarmin/czkawka/releases/tag/12.0.2) date du 9 septembre 2026. Son README distingue la licence MIT des outils Czkawka GTK et CLI, GPL-3.0-only pour Krokiet et Cedinia, et CC BY 4.0 pour les médias. [dupeGuru](https://github.com/arsenetar/dupeguru) est un outil graphique multiplateforme en Python et Qt sous GPL-3.0. Les deux projets prennent en charge Linux, Windows et macOS. Sources consultées le 23 novembre 2026.

Noskhati applique cette idée au contrôle directionnel d'une sauvegarde. L'implémentation utilise exclusivement la bibliothèque standard de Go.
