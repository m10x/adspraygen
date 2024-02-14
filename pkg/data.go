package pkg

func getGermanSeasons() []string {
	return []string{"Frühling", "Sommer", "Herbst", "Winter"}
}

func getAmericanSeasons() []string {
	return []string{"Spring", "Summer", "Fall", "Winter"}
}

func getBritishSeasons() []string {
	return []string{"Spring", "Summer", "Autumn", "Winter"}
}

func getGermanMonths() []string {
	return []string{
		"Januar", "Februar", "März", "April", "Mai", "Juni",
		"Juli", "August", "September", "Oktober", "November", "Dezember",
	}
}

func getEnglishMonths() []string {
	return []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}
}

func getLeetBasic() map[rune]string {
	return map[rune]string{
		'E': "3",
		'O': "0",
		'I': "1",
		'A': "4",
	}
}

func getLeetBasicPlus() map[rune]string {
	return map[rune]string{
		'E': "3",
		'O': "0",
		'I': "1",
		'A': "@",
		'T': "7",
	}
}
