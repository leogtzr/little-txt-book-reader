package main

import (
	"strings"
	"unicode"
)

var bannedWords = []string{
	"Los",
	"El",
	"La",
	"A",
	"Al",
	"En",
	"Las",
	"Mi",
	"Que",
	"Se",
	"Su",
	"Una",
	"Uno",
	"Lo",
	"Y",
	"Esta",
	"De",
	"Es",
	"Sus",
	"Si",
	"Un",
	"Con",
	"No",
	"Por",
	"Yo",
	"Todo",
	"Me",
	"Alli",
	"Nada",
	"Algo",
	"O",
	"Te",
	"Ya",
	"Aun",
	"Aún",
	"Muy",
	"Mis",
	"Oye",
	"Para",
	"Ese",
	"Sin",
	"Pero",
	"Sí",
	"Esto",
	"Porque",
	"Él",
	"Ella",
	"Ellas",
	"Esa",
	"Hoy",
	"Fue",
	"He",
	"Hice",
	"Está",
	"Ésta",
	"Cuando",
	"Desde",
	"Dios",
	"Era",
	"Eran",
	"Qué",
	"Sobre",
	"Solo",
	"Sólo",
	"Solos",
	"Soy",
	"Todos",
	"Hace",
	"Debo",
	"Debe",
	"Como",
	"Eso",
	"Nos",
	"Tan",
	"Sé",
	"Hasta",
	"Hay",
	"Otro",
	"Nunca",
	"Nosotros",
	"Puede",
	"Puedo",
	"Le",
	"Toda",
	"Así",
	"Aquí",
	"Ahí",
	"Ahora",
	"Tu",
	"Tú",
	"Tus",
	"Del",
	"Más",
	"Unas",
	"Unos",
}

func extractWords(line string) []string {
	words := strings.Split(line, " ")
	return words
}

func sanitizeWord(line string) string {
	line = strings.ReplaceAll(line, ".", "")
	line = strings.ReplaceAll(line, ",", "")
	line = strings.ReplaceAll(line, "\"", "")
	line = strings.ReplaceAll(line, ")", "")
	line = strings.ReplaceAll(line, "(", "")
	return line
}

func isTitle(word string) bool {
	chars := []rune(word)
	return unicode.IsUpper(chars[0])
}

func contains(words []string, word string) bool {
	for _, w := range words {
		if w == word {
			return true
		}
	}
	return false
}

func shouldIgnoreWord(word string) bool {
	return contains(bannedWords, word)
}
