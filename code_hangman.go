package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strings"
)

type GameState struct {
	Word              string
	MaskedWord        []string
	RemainingAttempts int
	Difficulte        string
	Ascii             string
	AfficheImage      string
}

var difficulté string
var lutil []string
var statusjeu GameState

func Index(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/Index.html"))
	t.Execute(w, nil)
}

func Jeux(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// Récupère la difficulté choisie
	statusjeu.Difficulte = r.FormValue("diffi")
	if statusjeu.Difficulte == "" {
		statusjeu.Difficulte = "moyen" // Défaut si aucune difficulté n'est choisie
	}
	// Initialise le jeu
	initializeGame()
	if statusjeu.RemainingAttempts == 10 {
		statusjeu.AfficheImage = "pootis.jpg"
	} else {
		statusjeu.AfficheImage = "pootis" + string(57-statusjeu.RemainingAttempts) + ".jpg"
	}
	// Préparation des données pour le template
	data := map[string]interface{}{
		"MaskedWord":        strings.Join(statusjeu.MaskedWord, " "),
		"RemainingAttempts": statusjeu.RemainingAttempts,
		"Message":           "",
		"GameLoose":         false,
		"GameWin":           false,
		"Word":              statusjeu.Word,
		"ShowImage":         statusjeu.AfficheImage,
		"LetterUse":         strings.Join(lutil, " , "),
	}
	// Rend la page du jeu
	t := template.Must(template.ParseFiles("templates/jeux.html"))
	t.Execute(w, data)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", Index)
	http.HandleFunc("/jeux", Jeux)
	http.HandleFunc("/lettre/", Hang)
	http.HandleFunc("/test/", test)
	fmt.Println("http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func Hang(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// Récupère la lettre soumise
	letter := ToUpper(r.PostFormValue("input"))

	// Vérifie si la lettre est valide
	if !simplelettre(letter) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// Gère la logique du jeu
	message := ""
	if InTab(lutil, letter) {
		message = "Lettre déjà utilisée."
	} else {
		lutil = append(lutil, letter)
		trouver := false
		for i := 0; i < len(statusjeu.Word); i++ {
			if ToUpper(string(statusjeu.Word[i])) == letter {
				statusjeu.MaskedWord[i] = letter
				trouver = true
			}
		}
		if !trouver {
			statusjeu.RemainingAttempts--
			statusjeu.AfficheImage = "pootis" + string(57-statusjeu.RemainingAttempts) + ".jpg"
			message = "Lettre incorrecte."
		} else {
			message = "Bien joué !"
		}
	}

	// Vérifie si le jeu est terminé
	if MotFini(statusjeu.MaskedWord) {
		statusjeu.AfficheImage = "happy-hoovi.png"
		message = "Bravo, vous avez gagné !"
	} else if statusjeu.RemainingAttempts <= 0 {
		message = "Vous avez perdu !"
	}
	// Recharge la page avec les mises à jour
	data := map[string]interface{}{
		"MaskedWord":        strings.Join(statusjeu.MaskedWord, " "),
		"RemainingAttempts": statusjeu.RemainingAttempts,
		"Message":           message,
		"GameLoose":         statusjeu.RemainingAttempts <= 0,
		"GameWin":           MotFini(statusjeu.MaskedWord),
		"Word":              statusjeu.Word,
		"ShowImage":         statusjeu.AfficheImage,
		"LetterUse":         strings.Join(lutil, " , "),
	}
	t := template.Must(template.ParseFiles("templates/jeux.html"))
	t.Execute(w, data)
	strhtml := fmt.Sprintf(" %s ", letter)
	tmpl, _ := template.New("t").Parse(strhtml)
	tmpl.Execute(w, nil)

}

func initializeGame() {
	lutil = []string{}
	statusjeu.Word = choimot(statusjeu.Difficulte + ".txt")
	statusjeu.MaskedWord = motcache(statusjeu.Word)
	if statusjeu.Difficulte == "facile" {
		nbrand := (len(statusjeu.Word) / 2) - 1
		for nbrand != 0 {
			n := rand.IntN(len(statusjeu.Word) - 1)
			if statusjeu.MaskedWord[n] == "_" {
				statusjeu.MaskedWord[n] = ToUpper(string(statusjeu.Word[n]))
				nbrand--
			}
		}
		statusjeu.RemainingAttempts = 10
	} else if statusjeu.Difficulte == "moyen" {
		statusjeu.MaskedWord[0] = ToUpper(string(statusjeu.Word[0]))
		statusjeu.RemainingAttempts = 8
	} else {
		statusjeu.RemainingAttempts = 5
	}

}

func test(w http.ResponseWriter, r *http.Request) {
	hangman()
}

func hangman() {

	var err error
	lutil := []string{}
	position := gettxt("hangman")
	lettreascii := []string{}
	a := 0
	ascii := ' '
	if len(os.Args) > 1 {
		if os.Args[1] == "save" {
			statusjeu, err = chargeJeu()
			if err != nil {
				fmt.Print("Erreur avec la sauvegarde")
				return
			}
			a++
			if statusjeu.Ascii == "maj" {
				lettreascii = gettxt("maj")
				ascii = 'M'
			} else if statusjeu.Ascii == "min" {
				lettreascii = gettxt("min")
				ascii = 'm'
			} else {
				ascii = 'n'
			}
		} else {
			fmt.Print("Trop d'arguments ou arguments invalide !")
			return
		}
	}
	difficulté = "facile"
	if difficulté == "facile" {
		statusjeu.Word = choimot("facile.txt")
		statusjeu.MaskedWord = motcache(statusjeu.Word)
		nbrand := (len(statusjeu.Word) / 2) - 1
		for nbrand != 0 {
			n := rand.IntN(len(statusjeu.Word) - 1)
			if statusjeu.MaskedWord[n] == "_" {
				statusjeu.MaskedWord[n] = ToUpper(string(statusjeu.Word[n]))
				nbrand--
			}
		}
		statusjeu.RemainingAttempts = 10
	} else if difficulté == "moyen" {
		statusjeu.Word = choimot("moyen.txt")
		statusjeu.MaskedWord = motcache(statusjeu.Word)
		statusjeu.MaskedWord[0] = ToUpper(string(statusjeu.Word[0]))
		statusjeu.RemainingAttempts = 8
	} else {
		statusjeu.Word = choimot("difficile.txt")
		statusjeu.MaskedWord = motcache(statusjeu.Word)
		statusjeu.RemainingAttempts = 5
	}
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
					statusjeu.Ascii = "maj"
				} else if ascii == 'm' {
					a = 3
					lettreascii = gettxt("min")
					statusjeu.Ascii = "min"
				} else {
					fmt.Print("Caractère invalide\n")
				}
			}
		} else {
			fmt.Print("Caractère invalide\n")
		}
	}

	fmt.Printf("\n \n Tu as %d essais pour trouver le bon mot\n", statusjeu.RemainingAttempts)
	fmt.Print("\t \tBonne chance\n \n")
	affichemot(statusjeu.MaskedWord, lettreascii, ascii)
	for statusjeu.RemainingAttempts != 0 {
		if MotFini(statusjeu.MaskedWord) {
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
			sauvegarde(statusjeu)
			return
		}
		if len(l) > 1 {
			if l == statusjeu.Word {
				break
			} else {
				if statusjeu.RemainingAttempts == 1 {
					statusjeu.RemainingAttempts--
				} else {
					statusjeu.RemainingAttempts -= 2
				}
				if statusjeu.RemainingAttempts != 0 {
					fmt.Printf("Il te reste %d essais pour trouver le bon mot\n", statusjeu.RemainingAttempts)
					fmt.Print(string(position[10-statusjeu.RemainingAttempts-1]) + "\n")
					affichemot(statusjeu.MaskedWord, lettreascii, ascii)
					continue
				}
			}
		}
		if InTab(lutil, l) {
			fmt.Print("Lettre déjà utiliser réessayer\n")
			continue
		}
		lutil = append(lutil, l)
		for i := 0; i < len(statusjeu.Word); i++ {
			if string(statusjeu.Word[i]) == l && string(statusjeu.MaskedWord[i]) == "_" {
				statusjeu.MaskedWord[i] = ToUpper(l)
				danslemot++
			}
		}
		if !InTab(statusjeu.MaskedWord, "_") {
			continue
		} else if danslemot == 0 {
			statusjeu.RemainingAttempts--
		}
		if statusjeu.RemainingAttempts == 0 {
			continue
		} else {
			if statusjeu.RemainingAttempts != 10 {
				fmt.Print(string(position[10-statusjeu.RemainingAttempts-1]) + "\n")
			}
			fmt.Printf("Il te reste %d essais pour trouver le bon mot\n", statusjeu.RemainingAttempts)
			affichemot(statusjeu.MaskedWord, lettreascii, ascii)
		}
	}
	if statusjeu.RemainingAttempts == 0 {
		fmt.Print(string(position[9]))
		fmt.Print("Perdu !! Le mot était\n")
		for i := 0; i < len(statusjeu.Word); i++ {
			statusjeu.MaskedWord[i] = ToUpper(string(statusjeu.Word[i]))
		}
		affichemot(statusjeu.MaskedWord, lettreascii, ascii)
	} else {
		fmt.Print("bravo vous avez trouvez le mot cacher\n")
		for i := 0; i < len(statusjeu.Word); i++ {
			statusjeu.MaskedWord[i] = ToUpper(string(statusjeu.Word[i]))
		}
		affichemot(statusjeu.MaskedWord, lettreascii, ascii)
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
	if ascii == 'n' {
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

func sauvegarde(state GameState) error {
	fichier, err := os.Create("save.txt")
	if err != nil {
		return err
	}
	defer fichier.Close()

	_, err = fichier.WriteString(fmt.Sprintf("%s\n%s\n%d\n%s\n", state.Word, convertmotenstr(state.MaskedWord), state.RemainingAttempts, state.Ascii))
	if err != nil {
		return err
	}
	fmt.Println("Partie sauvegardée dans save.txt.")
	return nil
}

func chargeJeu() (GameState, error) {
	fichier, err := os.Open("save.txt")
	if err != nil {
		return GameState{}, err
	}
	defer fichier.Close()

	var state GameState
	scanner := bufio.NewScanner(fichier)

	if scanner.Scan() {
		state.Word = scanner.Text()
	}
	if scanner.Scan() {
		state.MaskedWord = convertmotentab(scanner.Text())
	}
	if scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d", &state.RemainingAttempts)
	}
	if scanner.Scan() {
		state.Ascii = scanner.Text()
	}
	return state, nil
}

func convertmotentab(mot string) []string {
	motf := []string{}
	for _, i := range mot {
		motf = append(motf, string(i))
	}
	return motf
}

func convertmotenstr(mot []string) string {
	motf := ""
	for _, i := range mot {
		motf += i
	}
	return motf
}
