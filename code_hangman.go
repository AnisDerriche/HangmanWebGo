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

var lutil []string
var statusjeu GameState
var save map[string]interface{}

func Index(w http.ResponseWriter, r *http.Request) {
	save = map[string]interface{}{
		"MaskedWord":        strings.Join(statusjeu.MaskedWord, " "),
		"RemainingAttempts": statusjeu.RemainingAttempts,
		"Message":           "",
		"GameLoose":         false,
		"GameWin":           false,
		"Word":              statusjeu.Word,
		"ShowImage":         statusjeu.AfficheImage,
		"LetterUse":         strings.Join(lutil, " , "),
	}
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
	http.HandleFunc("/charger-sauvegarde/", charge_save)
	fmt.Println("http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func charge_save(w http.ResponseWriter, r *http.Request) {
	data := save
	t := template.Must(template.ParseFiles("templates/jeux.html"))
	t.Execute(w, data)
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

func simplelettre(l string) bool {
	for _, i := range l {
		if (i < 'a' || i > 'z') && (i < 'A' || i > 'Z') {
			return false
		}
	}
	return true
}
