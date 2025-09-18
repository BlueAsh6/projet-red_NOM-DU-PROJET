# Projet RED — Rapport Final

## Introduction

Le **Projet RED** est un mini-jeu en ligne de commande (CLI) développé dans le cadre de l’Ymmersion.  
Il met en pratique toutes les compétences abordées : **programmation Go, structures, logique de jeu, économie, inventaire, PNJ et combats tour par tour**.  

Nous avons choisi l’univers **CryoZéro™**, une dystopie cyberpunk glacée où le joueur doit survivre, progresser et affronter des ennemis de plus en plus puissants.  

👥 Créé par : **RIO Killian & ANCELIN Baptiste**

---

## Direction Artistique (DA)

- **Univers** : Monde futuriste gelé suite au crash mondial d’une IA. Mélange cyberpunk et heroic-fantasy.  
- **Graphisme** : Style minimaliste via l’affichage terminal, enrichi par des textes stylisés et de l’ASCII art.  
- **Thématique** : Conflit entre *White Hat* et *Black Hat*, erreurs 404, bugs système et artefacts informatiques.  
- **Expérience** : Progression par étapes : création du personnage, quêtes, combats dans l’arène, boss final et succès.  

---

## Présentation des personnages

### Joueur
- Personnage personnalisable : **nom** et **classe** (*Hacker* ou *Analyste SOC*).  
- Chaque classe a des points de vie et un inventaire spécifiques au départ.  
- Le joueur peut gagner de l’XP, monter de niveau, augmenter ses PV max et débloquer des succès.  

### Ennemis
- *Script Kiddie*, *Botnet Zombie*, *Admin Corrompu*.  
- **Boss final** : *Le Divin Sylvain*, ennemi redoutable lié à la quête principale.  

### PNJ
- 🏛 **Mairie** : attribution de quêtes.  
- 🛒 **Épicerie** : vente de nourriture contre des gains d’XP.  
- ⚒️ **Armurerie** : distribution d’armes et équipements spéciaux.  
- 🧑‍💻 **Tor le Marchand** : vente et craft d’armes avancées.  
- 🏠 **Base du survivant** : lieu de repos et de soin.  

---

## But du jeu

L’objectif du joueur est de :  
- Créer et personnaliser son personnage.  
- Gérer son inventaire, ses ressources et ses armes.  
- Remplir des quêtes obtenues à la mairie.  
- Combattre dans l’**Arène du Néon** contre des ennemis de plus en plus puissants.  
- Vaincre le **boss final** (*Le Divin Sylvain*).  
- Récupérer les **artefacts rares** pour débloquer les succès ultimes.  
- Accomplir la réparation symbolique de l’**Erreur 404**.  

---

## Fonctionnalités implémentées

### 1. Création du personnage & base
- ✅ Structure `Character` avec : nom, classe, niveau, PV max/actuels, inventaire, monnaie, XP, quêtes, succès.  
- ✅ Création interactive (`createCharacterInteractive`) avec choix du nom et de la classe.  
- ✅ Affichage complet des informations du personnage.  
- ✅ Système d’XP et de montée de niveau (gain de PV max et restauration des PV).  

### 2. Lieux et PNJ
- ✅ Structure `Land` (nom + description).  
- ✅ **Mairie** → attribution de quêtes.  
- ✅ **Épicerie** → objets consommables contre XP.  
- ✅ **Armurerie** → acquisition d’armes et objets.  
- ✅ **Base du survivant** → restauration des PV.  
- ✅ **Tor le Marchand** → boutique avancée et crafting d’armes.  

### 3. Économie et progression
- ✅ Monnaie interne (`Coins`).  
- ✅ Gain d’XP lors des combats et de l’achat de nourriture.  
- ✅ Succès spéciaux : *Trophée Admin Ultime*, *Vainqueur de l’Arène*, etc.  
- ✅ Loot unique pour éviter les doublons (artefacts rares).  

### 4. Combat et arène
- ✅ Structure `Ennemi` avec PV, attaque, récompenses et statut de boss.  
- ✅ Combats **tour par tour** : attaquer, utiliser un objet, ou fuir.  
- ✅ Pénalités en cas de fuite ou défaite (perte de coins, PV réduits).  
- ✅ Boss final *Le Divin Sylvain* avec succès et récompenses majeures.  
- ✅ Quête spéciale liée à la victoire sur l’arène complète.  

### 5. Gestion des objets
- ✅ Inventaire basé sur un `map[string]int`.  
- ✅ Potions de soins, armes spéciales (Épée Segfault, Arc Latence, Trident Proxy, etc.).  
- ✅ Utilisation d’objets directement en combat (`utiliserObjet`).  
- ✅ Système de loot progressif (artefacts rares).  

### 6. Interface et narration
- ✅ Intro scénarisée (année 3099, planète gelée, crash de l’IA).  
- ✅ Menu interactif clair pour naviguer entre les lieux.  
- ✅ Écran de mort (`vousEtesMort`).  
- ✅ Effet de texte lent pour renforcer l’immersion.  

---

## Conclusion

Le projet **CryoZéro™** concrétise l’ensemble des objectifs du **Projet RED** :  
- ✅ Création et gestion complète d’un personnage.  
- ✅ Monde structuré avec PNJ, quêtes et économie interne.  
- ✅ Combats stratégiques tour par tour avec progression et boss final.  
- ✅ Direction artistique cohérente et narration immersive.  

🎯 Le joueur progresse dans un univers cyber-gelé, bat les ennemis de l’arène, obtient les artefacts rares et débloque le **succès ultime** : la réparation de l’erreur 404.  

Un projet à la fois **technique, narratif et ludique**, qui illustre notre maîtrise de la programmation et notre créativité.  

---
