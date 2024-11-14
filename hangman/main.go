package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
)

func main() {
	mot := ""
	motcacher := []string{}
	nbessais := 0
	lutil := []string{}
	if len(os.Args) > 1 {
		fmt.Print("trop d'arguments")
		return
	}
	position := gettxt("hangman")
	d := 0
	difficulté := ' '
	for d == 0 {
		fmt.Print("Choisir une difficulté :")
		fmt.Scanf("%c\n", &difficulté)
		if difficulté == 'f' {
			d++
			mot = choimot("facile.txt")
			motcacher = motcache(mot)
			nbrand := (len(mot) / 2) - 1
			for nbrand != 0 {
				n := rand.IntN(len(mot) - 1)
				if motcacher[n] == "_" {
					motcacher[n] = ToUpper(string(mot[n]))
					nbrand--
				}
			}
			nbessais = 10
		} else if difficulté == 'm' {
			d++
			mot = choimot("moyen.txt")
			motcacher = motcache(mot)
			motcacher[0] = ToUpper(string(mot[0]))
			nbessais = 8
		} else if difficulté == 'd' {
			d++
			mot = choimot("difficile.txt")
			motcacher = motcache(mot)
			nbessais = 5
		} else {
			fmt.Print("Caractère invalide\n")
		}
	}
	lettreascii := []string{"0"}
	a := 0
	ascii := ' '
	for a == 0 {
		fmt.Print("Jouer en ascii ? : ")
		fmt.Scanf("%c\n", &ascii)
		if ascii == 'n' {
			a++
		} else if ascii == 'y' {
			a++
			for a == 1 {
				fmt.Print("Majuscule ou minuscule ? : ")
				fmt.Scanf("%c\n", &ascii)
				if ascii == 'M' {
					a++
					lettreascii = gettxt("maj")
				} else if ascii == 'm' {
					a = 3
					lettreascii = gettxt("min")
				} else {
					fmt.Print("Caractère invalide\n")
				}
			}
		} else {
			fmt.Print("Caractère invalide\n")
		}
	}
	fmt.Printf("\n \n Tu as %d essais pour trouver le bon mot\n", nbessais)
	fmt.Print("\t \tBonne chance\n \n")
	affichemot(motcacher, lettreascii, ascii)
	for nbessais != 0 {
		if MotFini(motcacher) {
			break
		}
		danslemot := 0
		l := ""
		fmt.Print("\n Choisir une lettre :")
		fmt.Scanf("%s\n", &l)
		if !simplelettre(l) {
			fmt.Print("Caractère invalide\n")
			continue
		}
		if l == "STOP" {
			return
		}
		if len(l) > 1 {
			if l == mot {
				break
			} else {
				if nbessais == 1 {
					nbessais--
				} else {
					nbessais -= 2
				}
				if nbessais != 0 {
					fmt.Printf("Il te reste %d essais pour trouver le bon mot\n", nbessais)
					fmt.Print(string(position[10-nbessais-1]) + "\n")
					affichemot(motcacher, lettreascii, ascii)
					continue
				}
			}
		}
		if InTab(lutil, l) {
			fmt.Print("Lettre déjà utiliser réessayer\n")
			continue
		}
		lutil = append(lutil, l)
		for i := 0; i < len(mot); i++ {
			if string(mot[i]) == l && string(motcacher[i]) == "_" {
				motcacher[i] = ToUpper(l)
				danslemot++
			}
		}
		if !InTab(motcacher, "_") {
			continue
		} else if danslemot == 0 {
			nbessais--
		}
		if nbessais == 0 {
			continue
		} else {
			if nbessais != 10 {
				fmt.Print(string(position[10-nbessais-1]) + "\n")
			}
			fmt.Printf("Il te reste %d essais pour trouver le bon mot\n", nbessais)
			affichemot(motcacher, lettreascii, ascii)
		}
	}
	if nbessais == 0 {
		fmt.Print(string(position[9]))
		fmt.Print("Perdu !! Le mot était\n")
		for i := 0; i < len(mot); i++ {
			motcacher[i] = ToUpper(string(mot[i]))
		}
		affichemot(motcacher, lettreascii, ascii)
	} else {
		fmt.Print("bravo vous avez trouvez le mot cacher\n")
		for i := 0; i < len(mot); i++ {
			motcacher[i] = ToUpper(string(mot[i]))
		}
		affichemot(motcacher, lettreascii, ascii)
	}
}

func choimot(fich string) string {
	fichier, err := os.Open(fich)
	if err != nil {
		fmt.Print(err)
	}
	fileScanner := bufio.NewScanner(fichier)
	fileScanner.Split(bufio.ScanLines)
	mots := []string{}
	for fileScanner.Scan() {
		mots = append(mots, fileScanner.Text())
	}
	fichier.Close()
	mot := mots[rand.IntN(len(mots)-1)]
	return mot
}

func motcache(mot string) []string {
	motcacher := []string{}
	for i := 0; i < len(mot); i++ {
		motcacher = append(motcacher, "_")
	}
	return motcacher
}

func ToUpper(s string) string {
	h := []rune(s)
	result := ""
	for i := 0; i < len(h); i++ {
		if (h[i] >= 'a') && (h[i] <= 'z') {
			h[i] = h[i] - 32
		}
		result += string(h[i])
	}
	return result
}

func MotFini(tab []string) bool {
	for _, i := range tab {
		if i == "_" {
			return false
		}
	}
	return true
}

func InTab(tab []string, lettre string) bool {
	for _, i := range tab {
		if i == lettre {
			return true
		}
	}
	return false
}

func affichemot(mot []string, tabascii []string, ascii rune) {
	if tabascii[0] == "0" {
		affichemotnormal(mot)
	} else {
		affichasciimot(tabascii, mot, ascii)
	}
}

func affichemotnormal(tab []string) {
	affiche := ""
	for _, i := range tab {
		affiche += i + " "
	}
	fmt.Print(affiche + "\n")
}

func affichasciimot(tab []string, mot []string, ascii rune) {
	lignes := make([]string, 8)
	charac := ""
	for _, lettre := range mot {
		if lettre == "_" {
			charac = tab[0]
		} else {
			if ascii == 'm' {
				charac = tab[lettre[0]-64]
			} else {
				charac = tab[lettre[0]-64]
			}
		}
		ligneascii := strings.Split(charac, "\n")
		for i, line := range ligneascii {
			if i < len(lignes) {
				lignes[i] += line
			}
		}
	}
	for _, ligne := range lignes {
		fmt.Printf(ligne + "\n")
	}
}

func simplelettre(l string) bool {
	for _, i := range l {
		if (i < 'a' || i > 'z') && (i < 'A' || i > 'Z') {
			return false
		}
	}
	return true
}

func gettxt(taille string) []string {
	tab := []string{}
	fichier, err := os.Open(taille + ".txt")
	if err != nil {
		fmt.Print(err)
	}
	fileScanner := bufio.NewScanner(fichier)
	fileScanner.Split(bufio.ScanLines)
	lettre := ""
	lscan := 0
	for fileScanner.Scan() {
		if lscan < 0 {
			lscan++
			continue
		}
		lettre += fileScanner.Text() + "\n"
		lscan++
		if lscan%7 == 0 {
			tab = append(tab, lettre)
			lettre = ""
			lscan -= 8
		}
	}
	fichier.Close()
	return tab
}
