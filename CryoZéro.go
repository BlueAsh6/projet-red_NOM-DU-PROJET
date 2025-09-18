package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Character struct {
	Nom              string
	Classe           string 
	Niveau           int
	PointsVieMax     int
	PointsVieActuels int
	Inventaire       map[string]int
	Coins            int
	XP               int
	Quetes           []string
	Succes           []string
}

type Land struct {
	Nom         string
	Description string
}

type Arme struct {
    Nom         string
    Prix        int
    Degats      int
    Description string
}

// compteur de mort 
var deathCount int
func loadDeaths() int {
    data, err := os.ReadFile("save.txt")
    if err != nil {
        return 0 // si pas de fichier → 0 mort
    }
    val, err := strconv.Atoi(strings.TrimSpace(string(data)))
    if err != nil {
        return 0
    }
    return val
}

func saveDeaths(count int) {
    os.WriteFile("save.txt", []byte(strconv.Itoa(count)), 0644)
}



// Création des personnages
func initCharacter(nom, classe string, niveau, pvMax, pvActuels, coins int, inventaire map[string]int) Character {
	return Character{
		Nom:              nom,
		Classe:           classe,
		Niveau:           niveau,
		PointsVieMax:     pvMax,
		PointsVieActuels: pvActuels,
		Inventaire:       inventaire,
		Coins:            coins,
		XP:               0,
		Quetes:           []string{},
		Succes:           []string{},
	}
}

// Création des lieux
func initLand(nom, description string) Land {
	return Land{
		Nom:         nom,
		Description: description,
	}
}

// Animation textuel
func printSlow(text string, delay time.Duration) {
	for _, r := range text {
		fmt.Printf("%c", r)
		time.Sleep(delay)
	}
	fmt.Println()
}

//Pour initialiser les armes de l'armurerie
func creerArmes() []Arme {
    return []Arme{
        {"Épée Segfault", 50, 15, "Une épée bugguée qui inflige 15 dégâts."},
        {"Bouclier Firewall", 75, 0, "Réduit les dégâts subis de 20%."},
        {"Arc Latence 300ms", 120, 20, "Arc futuriste infligeant 20 dégâts."},
        {"Trident Proxy", 200, 30, "Trident divin des réseaux, 30 dégâts."},
        {"AirMax 244p", 300, 0, "Augmente vos esquives, 10% chance d’éviter un coup."},
    }
}

// Choix 0-6 achat ou retour - argent = ou pas
func torLeMarchand(c *Character, reader *bufio.Reader) {
    armes := creerArmes()

    for {
        printSlow("Bienvenue chez ⚒️ TOR LE MARCHAND ⚒️", 40*time.Millisecond)
        fmt.Printf("Vous avez actuellement 💰 %d coins.\n", c.Coins)
        fmt.Println("Armes disponibles :")

        for i, arme := range armes {
            fmt.Printf("%d. %s — Prix : %d coins | Dégâts : %d | %s\n",
                i+1, arme.Nom, arme.Prix, arme.Degats, arme.Description)
        }
        fmt.Printf("%d. Quitter la boutique\n", len(armes)+1)
        fmt.Print("👉 Votre choix : ")

        choixStr, _ := reader.ReadString('\n')
        choixStr = strings.TrimSpace(choixStr)
        choix, err := strconv.Atoi(choixStr)
        if err != nil || choix < 1 || choix > len(armes)+1 {
            fmt.Println("❌ Choix invalide.")
            continue
        }

        if choix == len(armes)+1 {
            printSlow("Merci d’être passé chez TOR ! ❄️", 40*time.Millisecond)
            return
        }

        arme := armes[choix-1]
        if c.Coins >= arme.Prix {
            c.Coins -= arme.Prix
            c.Inventaire[arme.Nom]++
            printSlow(fmt.Sprintf("✅ Vous avez acheté %s !", arme.Nom), 40*time.Millisecond)
        } else {
            printSlow("❌ Pas assez de coins.", 40*time.Millisecond)
        }
    }
}

// renvoie le nombre d'XP requis pour le niveau en cours
func xpNeededForLevel(niveau int) int {
	return niveau * 100
}

// renvoie combien d'XP il manque pour le prochain niveau (sur une copie)
func xpToNextLevel(c Character) int {
	return xpNeededForLevel(c.Niveau) - c.XP
}

// Affiche l'écran de mort et termine le programme
func vousEtesMort() {
    deathCount++
    saveDeaths(deathCount)

    fmt.Println("══════════════════════════════════════════════")
    printSlow("☠️  Vous êtes mort... Fin de l’aventure.", 40*time.Millisecond)
    fmt.Printf("💀 Nombre de morts cumulées : %d\n", deathCount)
    fmt.Println("══════════════════════════════════════════════")
    os.Exit(0)
}

// Ajoute de l'XP au joueur, gère le passage de niveau et restaure PV/augmente PV max - 
// boucle au cas où on monte plusieurs niveaux d'un coup
func gagnerXP(joueur *Character, xp int) {
	joueur.XP += xp

	for {
		xpPourNiveau := xpNeededForLevel(joueur.Niveau)
		if joueur.XP < xpPourNiveau {
			break
		}
		joueur.XP -= xpPourNiveau
		joueur.Niveau++
		joueur.PointsVieMax += 20                
		joueur.PointsVieActuels = joueur.PointsVieMax 
		printSlow(fmt.Sprintf("⬆️ Niveau %d atteint ! PV max augmentés à %d 🎉", joueur.Niveau, joueur.PointsVieMax), 40*time.Millisecond)
	}

	
	restant := xpNeededForLevel(joueur.Niveau) - joueur.XP
	printSlow(fmt.Sprintf("🔸 Il vous reste %d XP pour atteindre le niveau %d.", restant, joueur.Niveau+1), 30*time.Millisecond)
}


// Affiche les stats d’un personnage 
func afficherPersonnage(c Character) {
	fmt.Println("══════════════════════════════════════════════")
	printSlow("🧍 Personnage sélectionné :", 40*time.Millisecond)
	printSlow(fmt.Sprintf(" Nom           : %s", c.Nom), 40*time.Millisecond)
	printSlow(fmt.Sprintf(" Classe        : %s", c.Classe), 40*time.Millisecond)
	printSlow(fmt.Sprintf(" Niveau        : %d", c.Niveau), 40*time.Millisecond)
	printSlow(fmt.Sprintf(" Points de vie : %d / %d", c.PointsVieActuels, c.PointsVieMax), 40*time.Millisecond)
	printSlow(fmt.Sprintf(" Coins         : %d", c.Coins), 40*time.Millisecond)
	printSlow(fmt.Sprintf(" XP            : %d", c.XP), 40*time.Millisecond)

	printSlow(" Inventaire    :", 40*time.Millisecond)
	for item, qty := range c.Inventaire {
		printSlow(fmt.Sprintf("  - %s : %d", item, qty), 40*time.Millisecond)
	}
	printSlow(" Quêtes        :", 40*time.Millisecond)
	for i, q := range c.Quetes {
		printSlow(fmt.Sprintf("  %d. %s", i+1, q), 40*time.Millisecond)
	}
	printSlow(" Succès        :", 20*time.Millisecond)
	for i, s := range c.Succes {
		printSlow(fmt.Sprintf("  %d. %s", i+1, s), 20*time.Millisecond)
	}
	fmt.Println("══════════════════════════════════════════════")
}

// Affiche le menu des lieux
func afficherLieux(lieux []Land) {
	fmt.Println("══════════════════════════════════════════════")
	fmt.Println("À quel endroit voulez-vous aller ?")
	fmt.Println()
	for i, l := range lieux {
		fmt.Printf(" %d. %s\n", i+1, l.Nom)
		fmt.Printf("    %s\n\n", l.Description)
	}
	fmt.Println("══════════════════════════════════════════════")
}

// Création interactive d’un personnage / Nom / Classe / Inventaire
func createCharacterInteractive(reader *bufio.Reader) Character {
	fmt.Println()
	fmt.Println("✨ Création de personnage ✨")
	fmt.Print("Entrez le nom de votre personnage : ")
	nom, _ := reader.ReadString('\n')
	nom = strings.TrimSpace(nom)
	if nom == "" {
		nom = "Héros Anonyme"
	}

	var classe string
	for {
		fmt.Print("Choisissez une classe (Hacker / Analyste SOC) : ")
		classeInp, _ := reader.ReadString('\n')
		classeInp = strings.ToLower(strings.TrimSpace(classeInp))
		if classeInp == "hacker" {
			classe = "Hacker"
			break
		} else if classeInp == "analyste soc" {
			classe = "Analyste SOC"
			break
		} else {
			fmt.Println("Classe invalide. Veuillez entrer 'Hacker' ou 'Analyste SOC'.")
		}
	}

	var pvMax, pvAct int
	inventaire := map[string]int{
		"Potion de soins":                         2,
		"Sauvegarde Repo du “Demi-Dieu M.Berger”": 1,
	}

	if classe == "Analyste SOC" {
		pvMax = 100
		pvAct = 50
		inventaire["Épée Segfault"] = 1
	} else {
		pvMax = 100
		pvAct = 40
		inventaire["Arc Latence"] = 1
	}

	coins := 20
	newChar := initCharacter(nom, classe, 1, pvMax, pvAct, coins, inventaire)
	fmt.Println("Personnage créé !")
	afficherPersonnage(newChar)
	return newChar
}

// Fonction de calcul arène dégats
// Calcul des effets
func calculerDegatsJoueur(c *Character) int {
    base := 5
    if c.Inventaire["Épée Segfault"] > 10 {
        base += 15
    }
    if c.Inventaire["Arc Latence 300ms"] > 10 {
        base += 20
    }
    if c.Inventaire["Trident Proxy"] > 15 {
        base += 30
    }
    return base
}


// Initialisation de la mairie
func mairie(c *Character, reader *bufio.Reader) {
	for {
		printSlow("Bienvenue à la 🏛  MAIRIE GLACIALE", 40*time.Millisecond)
		fmt.Println("1. Récupérer une quête")
		fmt.Println("2. Repartir")
		fmt.Print("Votre choix : ")
		var choix int
		fmt.Scan(&choix)
		if choix == 1 {
			c.Quetes = append(c.Quetes, "Nouvelle quête mystérieuse")
			printSlow("📜 Vous avez reçu une nouvelle quête !", 40*time.Millisecond)
			printSlow("👉 Éliminer au moins une fois chaque Black Hat de l’arène du Néon : 0/4 👈", 60*time.Millisecond)
		} else if choix == 2 {
			return
		} else {
			fmt.Println("Choix invalide")
		}
	}
}

// Initialisation de l'epicerie
func epicerie(c *Character, reader *bufio.Reader) {
	for {
		printSlow("Bienvenue à l’🛒  ÉPICERIE HIVERNALE", 40*time.Millisecond)
		fmt.Println("1. 🥪 Sandwich Packet Perdu — 5 coins, +10 XP")
		fmt.Println("2. 🍗 Poulet Frit Overclocké — 15 coins, +50 XP")
		fmt.Println("3. 🧃 Kombucha Firmware — 3 coins, +6 XP")
		fmt.Println("4. 🍔 Happy Meal .zip — 50 coins, +200 XP")
		fmt.Println("5. 🎖️ Trophée Admin Ultime — 10000 coins → Succès spécial")
		fmt.Println("6. Repartir")
		fmt.Print("Votre choix : ")
		var choix int
		fmt.Scan(&choix)
		switch choix {
		case 1:
			if c.Coins >= 5 {
				c.Coins -= 5
				gagnerXP(c, 10)
				printSlow("🥪 Vous mangez le Sandwich (+10 XP)", 40*time.Millisecond)
			} else {
				printSlow("❌ Pas assez de coins !", 40*time.Millisecond)
			}
		case 2:
			if c.Coins >= 15 {
				c.Coins -= 15
				gagnerXP(c, 50)
				printSlow("🍗 Vous mangez le Poulet (+50 XP)", 40*time.Millisecond)
			} else {
				printSlow("❌ Pas assez de coins !", 40*time.Millisecond)
			}
		case 3:
			if c.Coins >= 3 {
				c.Coins -= 3
				gagnerXP(c, 6)
				printSlow("🧃 Vous buvez le Kombucha (+6 XP)", 40*time.Millisecond)
			} else {
				printSlow("❌ Pas assez de coins !", 40*time.Millisecond)
			}
		case 4:
			if c.Coins >= 50 {
				c.Coins -= 50
				gagnerXP(c, 200)
				printSlow("🍔 Vous mangez le Happy Meal (+200 XP)", 40*time.Millisecond)
			} else {
				printSlow("❌ Pas assez de coins !", 40*time.Millisecond)
			}
		case 5:
			if c.Coins >= 10000 {
				c.Coins -= 10000
				c.Succes = append(c.Succes, "Trophée Admin Ultime")
				printSlow("🎖️ Succès débloqué : Trophée Admin Ultime !", 40*time.Millisecond)
			} else {
				printSlow("❌ Pas assez de coins !", 40*time.Millisecond)
			}
		case 6:
			return
		default:
			fmt.Println("Choix invalide")
		}
	}
}

// Fonction de l'armurerie 
func armurerie(c *Character, reader *bufio.Reader) {
	for {
		printSlow("Bienvenue à l’⚒️  ARMURERIE", 40*time.Millisecond)
		fmt.Println("1. ⚔️ 	Épée Segfault")
		fmt.Println("2. 🛡️  Bouclier Firewall")
		fmt.Println("3. 🏹 	Arc Latence 300ms")
		fmt.Println("4. 🔱 	Trident Proxy")
		fmt.Println("5. 👟 	AirMax 244p")
		fmt.Println("6. Repartir")
		fmt.Print("Votre choix : ")
		var choix int
		fmt.Scan(&choix)
		switch choix {
		case 1:
			c.Inventaire["Épée Segfault"]++
			printSlow("⚔️ Vous avez obtenu une Épée Segfault !", 40*time.Millisecond)
		case 2:
			c.Inventaire["Bouclier Firewall"]++
			printSlow("🛡️ Vous avez obtenu un Bouclier Firewall !", 40*time.Millisecond)
		case 3:
			c.Inventaire["Arc Latence"]++
			printSlow("🏹 Vous avez obtenu un Arc Latence !", 40*time.Millisecond)
		case 4:
			c.Inventaire["Trident Proxy"]++
			printSlow("🔱 Vous avez obtenu un Trident Proxy !", 40*time.Millisecond)
		case 5:
			c.Inventaire["AirMax 244p"]++
			printSlow("👟 Vous avez obtenu des AirMax 244p !", 40*time.Millisecond)
		case 6:
			return
		default:
			fmt.Println("Choix invalide")
		}
	}
}

// Initialisation de la base du survivant
func baseSurvivant(c *Character, reader *bufio.Reader) {
    for {
        printSlow("Bienvenue à votre 🏠 BASE DU SURVIVANT", 40*time.Millisecond)
        fmt.Println("1. Se reposer (restaure PV)")
        fmt.Println("2. Quitter la base")
        fmt.Print("Votre choix : ")

        var choix int
        fmt.Scan(&choix)

        switch choix {
        case 1:
            c.PointsVieActuels = c.PointsVieMax
            fmt.Println("💤 Vous vous reposez... Vos PV sont restaurés !")

        case 2:
            return
        }
    }
}

// ARÈNE : types & fonctions 

type Ennemi struct {
	Nom   string
	PV    int
	Att   int
	Coins int
	Boss  bool
}

// Stats des ennemis
func creerEnnemis() []Ennemi {
	return []Ennemi{
		{"Script Kiddie", 50, 15, 10, false},
		{"Botnet Zombie", 60, 25, 20, false},
		{"Admin Corrompu", 100, 55, 40, false},
		{"Le Divin Sylvain (Boss)", 400, 45, 500, true},
	}
}

func utiliserObjet(joueur *Character, choix string, ennemi *Ennemi) {
	switch choix {
	case "Potion de soins":
		if joueur.Inventaire["Potion de soins"] > 0 {
			heal := 30
			joueur.PointsVieActuels += heal
			if joueur.PointsVieActuels > joueur.PointsVieMax {
				joueur.PointsVieActuels = joueur.PointsVieMax
			}
			fmt.Printf("💊 Vous buvez une potion et récupérez %d PV !\n", heal)
			joueur.Inventaire["Potion de soins"]--
		} else {
			fmt.Println("❌ Vous n'avez plus de potion.")
		}

	case "Épée Segfault":
		if joueur.Inventaire["Épée Segfault"] > 0 {
			dmg := 15 + rand.Intn(6) // 15 dégâts de base + petit bonus aléatoire
			ennemi.PV -= dmg
			fmt.Printf("⚔️ Vous frappez avec l'Épée Segfault et infligez %d dégâts !\n", dmg)
		} else {
			fmt.Println("❌ Vous n'avez pas cette arme.")
		}

	case "Arc Latence 300ms":
		if joueur.Inventaire["Arc Latence 300ms"] > 0 {
			dmg := 20 + rand.Intn(10)
			ennemi.PV -= dmg
			fmt.Printf("🏹 Vous tirez avec l'Arc Latence 300ms et infligez %d dégâts !\n", dmg)
		} else {
			fmt.Println("❌ Vous n'avez pas cette arme.")
		}

	case "Trident Proxy":
		if joueur.Inventaire["Trident Proxy"] > 0 {
			dmg := 30 + rand.Intn(15)
			ennemi.PV -= dmg
			fmt.Printf("🔱 Vous transpercez avec le Trident Proxy et infligez %d dégâts !\n", dmg)
		} else {
			fmt.Println("❌ Vous n'avez pas cette arme.")
		}

	default:
		fmt.Println("❌ Objet inconnu ou inutilisable en combat.")
	}
}


func combat(joueur *Character, ennemi Ennemi, reader *bufio.Reader) bool {
	fmt.Printf("\n⚔️ Combat contre %s ! (%d PV)\n", ennemi.Nom, ennemi.PV)
	pvEnnemi := ennemi.PV
	ennemiOriginalAtt := ennemi.Att

	for joueur.PointsVieActuels > 0 && pvEnnemi > 0 {
		fmt.Printf("\n[%s] PV: %d/%d | [%s] PV: %d\n",
			joueur.Nom, joueur.PointsVieActuels, joueur.PointsVieMax,
			ennemi.Nom, pvEnnemi,
		)

		fmt.Println("1. Attaquer")
		fmt.Println("2. Utiliser un objet")
		fmt.Println("3. Fuir")
		fmt.Print("👉 Choix : ")
		choix, _ := reader.ReadString('\n')
		choix = strings.TrimSpace(choix)

		if choix == "1" {
			degatsBase := calculerDegatsJoueur(joueur)
			dmg := rand.Intn(5) + degatsBase
			fmt.Printf("💥 Vous attaquez avec vos armes et infligez %d dégâts !\n", dmg)
		pvEnnemi -= dmg

		} else if choix == "2" {
			fmt.Println("Inventaire :")
			hasUsable := false
			for item, qty := range joueur.Inventaire {
				if qty > 0 {
					fmt.Printf(" - %s (%d)\n", item, qty)
					hasUsable = true
				}
			}
			if !hasUsable {
				fmt.Println(" Aucun item utilisable.")
			} else {
				fmt.Print("Quel objet utiliser ? : ")
				obj, _ := reader.ReadString('\n')
				obj = strings.TrimSpace(obj)
				utiliserObjet(joueur, obj, &ennemi)
			}
		} else if choix == "3" {
			fmt.Println("🏃 Vous fuyez le combat !")
			return false
		} else {
			fmt.Println("❌ Choix invalide.")
			continue
		}

		if pvEnnemi <= 0 {
			fmt.Printf("🎉 Vous avez vaincu %s !\n", ennemi.Nom)
			joueur.Coins += ennemi.Coins
			gagnerXP(joueur, 150)
			if ennemi.Boss {
				fmt.Println("🏆 Succès débloqué : Vainqueur de l’Arène !")
				joueur.Succes = append(joueur.Succes, "Vainqueur de l'Arène")
				joueur.Coins += 1500
				gagnerXP(joueur, 150)

				printSlow("📜 Quête accomplie : Vous avez éliminé les 4 Black Hat de l’Arène !", 40*time.Millisecond)

				recompenses := []string{
					"Clef SSH Maudite",
					"Fichier suspect.docx",
					"Sauvegarde Repo du “Demi-Dieu M.Berger”",
				}

				candidates := []string{}
				for _, r := range recompenses {
					if joueur.Inventaire[r] == 0 {
						candidates = append(candidates, r)
					}
				}

				var drop string
				if len(candidates) > 0 {
					drop = candidates[rand.Intn(len(candidates))]
					joueur.Inventaire[drop]++
					printSlow(fmt.Sprintf("🎁 Vous avez obtenu : %s", drop), 40*time.Millisecond)
				} else {
					printSlow("✨ Vous possédez déjà tous les artefacts rares !", 40*time.Millisecond)
					printSlow("À la place, vous recevez 5000 coins bonus 💰", 40*time.Millisecond)
					joueur.Coins += 50000
				}

				joueur.Succes = append(joueur.Succes, "Termine la quête des Black Hat")
				printSlow("===========================================================================", 40*time.Millisecond)
				printSlow("🌟 Succès débloqué : Fin de la Quête 'Black Hat 4/4' !", 40*time.Millisecond)
			}
			ennemi.Att = ennemiOriginalAtt
			return true
		}

		if ennemi.Att > 0 {
			dmg := rand.Intn(ennemi.Att) + 3
			joueur.PointsVieActuels -= dmg
			fmt.Printf("💀 %s vous inflige %d dégâts !\n", ennemi.Nom, dmg)
		} else {
			fmt.Printf("❄️ %s est momentanément immobilisé.\n", ennemi.Nom)
			ennemi.Att = ennemiOriginalAtt
		}

		if joueur.PointsVieActuels <= 0 {
			printSlow("☠️ Vous êtes mort...", 80*time.Millisecond)
			vousEtesMort()
		}
	}


	return false
}


func arene(joueur *Character, reader *bufio.Reader) {
	fmt.Println("Bienvenue dans l'🏟️ ARÈNE DU NÉON !")
	ennemis := creerEnnemis()
	for i, e := range ennemis {
		fmt.Printf("\n=== Niveau %d : %s ===\n", i+1, e.Nom)
		victoire := combat(joueur, e, reader)
		if !victoire {
			fmt.Println("Vous quittez l'arène...")
			
			// quand on perd on applique petite pénalité (perte moitié coins) et remet un peu de PV pour continuer le jeu
			loss := joueur.Coins / 2
			joueur.Coins -= loss
			if joueur.Coins < 0 {
				joueur.Coins = 0
			}
			minHP := jogadorSafeHP(joueur)
			joueur.PointsVieActuels = minHP
			fmt.Printf("Vous repartez avec %d coins, PV restaurés à %d.\n", joueur.Coins, joueur.PointsVieActuels)
			return
		}
		if i < len(ennemis)-1 {
			fmt.Print("Voulez-vous continuer ? (o/n) : ")
			choix, _ := reader.ReadString('\n')
			choix = strings.TrimSpace(strings.ToLower(choix))
			if choix != "o" && choix != "oui" {
				fmt.Println("Vous quittez l'arène avec vos gains !")
				return
			}
		}
	}
	// si on boucle tous les ennemis on revient
	fmt.Println("Vous avez nettoyé toutes les vagues de l'arène !")
}

// helper pour remettre un peu de PV au joueur après défaite
func jogadorSafeHP(joueur *Character) int {
	minHP := joueur.PointsVieMax / 4
	if minHP < 10 {
		minHP = 10
	}
	return minHP
}

// ================= BOUCLE PRINCIPALE DU JEU =================
func jeu(joueur *Character, lieux []Land, reader *bufio.Reader) {
	for {
		afficherLieux(lieux)
		fmt.Print("Entrez le numéro de votre choix + 'Entrée' : ")
		choixStr, _ := reader.ReadString('\n')
		choixStr = strings.TrimSpace(choixStr)
		choix, err := strconv.Atoi(choixStr)
		if err != nil || choix < 0 || choix > len(lieux) {
			fmt.Println("❌ Choix invalide.")
			continue
		}
		if choix == 0 {
			printSlow("Merci d'avoir joué à CryoZéro™ !", 50*time.Millisecond)
			break
		}

		printSlow(fmt.Sprintf("Vous êtes à : %s", lieux[choix-1].Nom), 40*time.Millisecond)

		switch choix {
		case 1:
			mairie(joueur, reader)
		case 2:
			baseSurvivant(joueur, reader)
		case 3:
			arene(joueur, reader)
		case 4:
			epicerie(joueur, reader)
		case 5:
    		torLeMarchand(joueur, reader)

		default:
			printSlow("❄️ Fonctionnalité pas encore disponible pour ce lieu.", 40*time.Millisecond)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)

	// ===== INTRO =====
	fmt.Println("══════════════════════════════════════════════════════════════")
	printSlow("                   🧊 Bienvenue dans CryoZéro™ 🧊             ", 40*time.Millisecond)
	fmt.Println("══════════════════════════════════════════════════════════════")
	fmt.Println()

	printSlow("🗓️  Année 3099.", 50*time.Millisecond)
	printSlow("L’intelligence artificielle mondiale a crashé pendant une mise à jour système...", 40*time.Millisecond)
	time.Sleep(400 * time.Millisecond)
	printSlow("Résultat : La planète entière est entrée en “mode congélateur”.", 40*time.Millisecond)
	printSlow("Les rares survivants vivent désormais dans d'immenses serveurs gelés,", 40*time.Millisecond)
	printSlow("essayant de retrouver leur humanité entre deux bugs critiques.", 40*time.Millisecond)
	time.Sleep(500 * time.Millisecond)
	fmt.Println()
	printSlow("🌐 Ici, les erreurs 404 sont mortelles, et la neige est un fichier corrompu.", 40*time.Millisecond)
	printSlow("Bienvenue à CryoZéro™.", 40*time.Millisecond)
	fmt.Println()
	printSlow("⚔️  J'espère que vous êtes prêt à choisir votre personnage et à plonger dans ce monde glitché ... ", 50*time.Millisecond)
	fmt.Println()

	// ===== Création des personnages =====
	perso1 := initCharacter("Bit McCrypte", "Analyste SOC", 1, 100, 40, 20000000, map[string]int{
		"Potion de soins":  2,
		"Épée Segfault":    1,
		"Clef SSH Maudite": 1,
	})
	perso2 := initCharacter("Clippy Pwnz", "Hacker", 1, 100, 50, 20, map[string]int{
		"Potion de soins":      2,
		"Bouclier Firewall":    1,
		"Fichier suspect.docx": 1,
	})

	personnages := []Character{perso1, perso2}

	// ===== Sélection ou création =====
	printSlow("Quel personnage voulez-vous incarner ?", 50*time.Millisecond)
	for i, p := range personnages {
		fmt.Printf(" %d. %s (%s)\n", i+1, p.Nom, p.Classe)
	}
	fmt.Printf(" %d. Créer un nouveau personnage\n", len(personnages)+1)

	fmt.Print("Entrez le numéro de votre choix + 'Entrée' : ")
	choixStr, _ := reader.ReadString('\n')
	choixStr = strings.TrimSpace(choixStr)
	choix, err := strconv.Atoi(choixStr)
	if err != nil || choix < 1 || choix > len(personnages)+1 {
		fmt.Println("❌ Choix invalide. Le programme s'arrête.")
		return
	}

	var joueur Character
	if choix == len(personnages)+1 {
		joueur = createCharacterInteractive(reader)
	} else {
		joueur = personnages[choix-1]
		afficherPersonnage(joueur)
	}

	deathCount = loadDeaths()
fmt.Printf("💀 Nombre total de morts enregistrées : %d\n", deathCount)


	// ===== Définitions des lieux =====
	Lieu1 := initLand("🏛  MAIRIE GLACIALE", "Centre névralgique des quêtes et des missions. Certains objets que vous possédez pourraient s’avérer importants pour débloquer de nouveaux secrets…")
	Lieu2 := initLand("🏠  BASE DU SURVIVANT", "Votre refuge dans ce monde gelé. Ici, vous pouvez vous soigner, organiser vos affaires et stocker vos objets en toute sécurité.")
	Lieu3 := initLand("🏟️  ARÈNE DU NÉON", "Entrez si vous l’osez : Chaque victoire vous rapproche… du Boss ultime : 'Le Divin Sylvain'. Ici, on combat pour sauver l'humanité du bug.")
	Lieu4 := initLand("🛒  ÉPICERIE HIVERNALE", "Un bazar rempli de victuailles, potions et objets étranges. Mangez pour gagner de l’expérience. Certains trésors inattendus pourraient apparaître…")
	Lieu5 := initLand("⚒️  TOR LE MARCHAND", "Ici, vous pouvez acheter, crafter et fusionner vos armes. La qualité n’a pas de bug… seulement un coût.")

	Lieux := []Land{Lieu1, Lieu2, Lieu3, Lieu4, Lieu5}

	// ===== LANCEMENT DU JEU =====
	jeu(&joueur, Lieux, reader)
}
