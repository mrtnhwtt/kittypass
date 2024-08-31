package kittypass

import "math/rand"

type PasswordGenerator struct {
	Length      int
	SpecialChar bool
	Numeral     bool
	Uppercase   bool
}

func (g *PasswordGenerator) GeneratePassword() string {
	lowercase := []rune("abcdefghijklmnopqrstuvwxyz")
	numbers := []rune("0123456789")
	uppercase := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	specialChars := []rune("!#$%&*+-?@^_~")

	var selectedChars []rune

	selectedChars = append(selectedChars, lowercase[rand.Intn(len(lowercase))])

	if g.Numeral {
		selectedChars = append(selectedChars, numbers[rand.Intn(len(numbers))])
	}
	if g.Uppercase {
		selectedChars = append(selectedChars, uppercase[rand.Intn(len(uppercase))])
	}
	if g.SpecialChar {
		selectedChars = append(selectedChars, specialChars[rand.Intn(len(specialChars))])
	}

	allChars := lowercase
	if g.Numeral {
		allChars = append(allChars, numbers...)
	}
	if g.Uppercase {
		allChars = append(allChars, uppercase...)
	}
	if g.SpecialChar {
		allChars = append(allChars, specialChars...)
	}

	remainingLength := g.Length - len(selectedChars)
	for i := 0; i < remainingLength; i++ {
		selectedChars = append(selectedChars, allChars[rand.Intn(len(allChars))])
	}

	rand.Shuffle(len(selectedChars), func(i, j int) {
		selectedChars[i], selectedChars[j] = selectedChars[j], selectedChars[i]
	})

	return string(selectedChars)
}
