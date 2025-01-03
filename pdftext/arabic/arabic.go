package arabic

// copied from https://github.com/abdullahdiaa/garabic and modified

import (
	"bytes"
	"fmt"
	"github.com/raceresult/gopdf/pdf"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

// letterGroup represents the letter and bounding letters
type letterGroup struct {
	backLetter  rune
	letter      rune
	frontLetter rune
}

// letterShape represents all shapes of arabic letters in a word
// https://web.stanford.edu/dept/lc/arabic/alphabet/incontextletters.html
type letterShape struct {
	Independent, Initial, Medial, Final rune
}

// Map of different shapes of arabic alphabet
var arabicAlphabetShapes = map[rune]letterShape{
	// Letter (ﺃ)
	'\u0623': {Independent: '\uFE83', Initial: '\u0623', Medial: '\uFE84', Final: '\uFE84'},
	// Letter (ﺍ)
	'\u0627': {Independent: '\uFE8D', Initial: '\u0627', Medial: '\uFE8E', Final: '\uFE8E'},
	// Letter (ﺁ)
	'\u0622': {Independent: '\uFE81', Initial: '\u0622', Medial: '\uFE82', Final: '\uFE82'},
	// Letter (ﺀ)
	'\u0621': {Independent: '\uFE80', Initial: '\u0621', Medial: '\u0621', Final: '\u0621'},
	// Letter (ﺅ)
	'\u0624': {Independent: '\uFE85', Initial: '\u0624', Medial: '\uFE86', Final: '\uFE86'},
	// Letter (ﺇ)
	'\u0625': {Independent: '\uFE87', Initial: '\u0625', Medial: '\uFE88', Final: '\uFE88'},
	// Letter (ﺉ)
	'\u0626': {Independent: '\uFE89', Initial: '\uFE8B', Medial: '\uFE8C', Final: '\uFE8A'},
	// Letter (ﺏ)
	'\u0628': {Independent: '\uFE8F', Initial: '\uFE91', Medial: '\uFE92', Final: '\uFE90'},
	// Letter (ﺕ)
	'\u062A': {Independent: '\uFE95', Initial: '\uFE97', Medial: '\uFE98', Final: '\uFE96'},
	// Letter (ﺓ)
	'\u0629': {Independent: '\uFE93', Initial: '\u0629', Medial: '\u0629', Final: '\uFE94'},
	// Letter (ﺙ)
	'\u062B': {Independent: '\uFE99', Initial: '\uFE9B', Medial: '\uFE9C', Final: '\uFE9A'},
	// Letter (ﺝ)
	'\u062C': {Independent: '\uFE9D', Initial: '\uFE9F', Medial: '\uFEA0', Final: '\uFE9E'},
	// Letter (ﺡ)
	'\u062D': {Independent: '\uFEA1', Initial: '\uFEA3', Medial: '\uFEA4', Final: '\uFEA2'},
	// Letter (ﺥ)
	'\u062E': {Independent: '\uFEA5', Initial: '\uFEA7', Medial: '\uFEA8', Final: '\uFEA6'},
	// Letter (ﺩ)
	'\u062F': {Independent: '\uFEA9', Initial: '\u062F', Medial: '\uFEAA', Final: '\uFEAA'},
	// Letter (ﺫ)
	'\u0630': {Independent: '\uFEAB', Initial: '\u0630', Medial: '\uFEAC', Final: '\uFEAC'},
	// Letter (ﺭ)
	'\u0631': {Independent: '\uFEAD', Initial: '\u0631', Medial: '\uFEAE', Final: '\uFEAE'},
	// Letter (ﺯ)
	'\u0632': {Independent: '\uFEAF', Initial: '\u0632', Medial: '\uFEB0', Final: '\uFEB0'},
	// Letter (ﺱ)
	'\u0633': {Independent: '\uFEB1', Initial: '\uFEB3', Medial: '\uFEB4', Final: '\uFEB2'},
	// Letter (ﺵ)
	'\u0634': {Independent: '\uFEB5', Initial: '\uFEB7', Medial: '\uFEB8', Final: '\uFEB6'},
	// Letter (ﺹ)
	'\u0635': {Independent: '\uFEB9', Initial: '\uFEBB', Medial: '\uFEBC', Final: '\uFEBA'},
	// Letter (ﺽ)
	'\u0636': {Independent: '\uFEBD', Initial: '\uFEBF', Medial: '\uFEC0', Final: '\uFEBE'},
	// Letter (ﻁ)
	'\u0637': {Independent: '\uFEC1', Initial: '\uFEC3', Medial: '\uFEC4', Final: '\uFEC2'},
	// Letter (ﻅ)
	'\u0638': {Independent: '\uFEC5', Initial: '\uFEC7', Medial: '\uFEC8', Final: '\uFEC6'},
	// Letter (ﻉ)
	'\u0639': {Independent: '\uFEC9', Initial: '\uFECB', Medial: '\uFECC', Final: '\uFECA'},
	// Letter (ﻍ)
	'\u063A': {Independent: '\uFECD', Initial: '\uFECF', Medial: '\uFED0', Final: '\uFECE'},
	// Letter (ﻑ)
	'\u0641': {Independent: '\uFED1', Initial: '\uFED3', Medial: '\uFED4', Final: '\uFED2'},
	// Letter (ﻕ)
	'\u0642': {Independent: '\uFED5', Initial: '\uFED7', Medial: '\uFED8', Final: '\uFED6'},
	// Letter (ﻙ)
	'\u0643': {Independent: '\uFED9', Initial: '\uFEDB', Medial: '\uFEDC', Final: '\uFEDA'},
	// Letter (ﻝ)
	'\u0644': {Independent: '\uFEDD', Initial: '\uFEDF', Medial: '\uFEE0', Final: '\uFEDE'},
	// Letter (ﻡ)
	'\u0645': {Independent: '\uFEE1', Initial: '\uFEE3', Medial: '\uFEE4', Final: '\uFEE2'},
	// Letter (ﻥ)
	'\u0646': {Independent: '\uFEE5', Initial: '\uFEE7', Medial: '\uFEE8', Final: '\uFEE6'},
	// Letter (ﻩ)
	'\u0647': {Independent: '\uFEE9', Initial: '\uFEEB', Medial: '\uFEEC', Final: '\uFEEA'},
	// Letter (ﻭ)
	'\u0648': {Independent: '\uFEED', Initial: '\u0648', Medial: '\uFEEE', Final: '\uFEEE'},
	// Letter (ﻱ)
	'\u064A': {Independent: '\uFEF1', Initial: '\uFEF3', Medial: '\uFEF4', Final: '\uFEF2'},
	// Letter (ﻯ)
	'\u0649': {Independent: '\uFEEF', Initial: '\u0649', Medial: '\uFEF0', Final: '\uFEF0'},
	// Letter (ـ)
	'\u0640': {Independent: '\u0640', Initial: '\u0640', Medial: '\u0640', Final: '\u0640'},
	// Letter (ﻻ)
	'\uFEFB': {Independent: '\uFEFB', Initial: '\uFEFB', Medial: '\uFEFC', Final: '\uFEFC'},
	// Letter (ﻷ)
	'\uFEF7': {Independent: '\uFEF7', Initial: '\uFEF7', Medial: '\uFEF8', Final: '\uFEF8'},
	// Letter (ﻹ)
	'\uFEF9': {Independent: '\uFEF9', Initial: '\uFEF9', Medial: '\uFEFA', Final: '\uFEFA'},
	// Letter (ﻵ)
	'\uFEF5': {Independent: '\uFEF5', Initial: '\uFEF5', Medial: '\uFEF6', Final: '\uFEF6'},
}

// Normalizable Arabic letters
var normalizable = &unicode.RangeTable{
	R16: []unicode.Range16{
		/*
			Arabic Harakat (Harakat تَشْكِيل)
		*/
		//Tatweel => ـ
		{0x0640, 0x0640, 1},
		//TanwinFatḥah
		{0x064B, 0x64B, 1},
		//TanwinDammah
		{0x064C, 0x64C, 1},
		//TanwinKasrah
		{0x064D, 0x64D, 1},
		//Fatḥah
		{0x064E, 0x64E, 1},
		//Dammah
		{0x064F, 0x64F, 1},
		//Kasrah
		{0x0650, 0x650, 1},
		//Shaddah
		{0x0651, 0x651, 1},
		//Sukun
		{0x0652, 0x652, 1},
		//DaggerAlif =>
		{0x0670, 0x0670, 1},
	},
}

// Normalizable letters [alef/Yae/Hae]
const (
	//Alef  => ا
	Alef = '\u0627'
	//AlefMad =>  آ
	AlefMad = '\u0622'
	//AlefHamzaAbove => أ
	AlefHamzaAbove = '\u0623'
	//AlefHamzaBelow إ
	AlefHamzaBelow = '\u0625'
	//Yae =>  ي
	Yae = '\u064A'
	//DotlessYae =>  ى
	DotlessYae = '\u0649'
	//TehMarbuta => ة
	TehMarbuta = '\u0629'
	//Hae => ه
	Hae = '\u0647'
	//AlefWaslah ٱ / Waslah is considered part of harakat/تَشْكِيل ?
	AlefWaslah = '\u0671'
)

// RemoveHarakat will remove harakat from arabic text
func RemoveHarakat(input string) string {
	input = normalizeTransform(input)
	runes := bytes.Runes([]byte(input))
	for i := 0; i < len(runes); i++ {
		//fmt.Println(string(runes[i]))
		switch runes[i] {
		//Remove Waslah from AlefWaslah / Waslah is considered part of harakat/تَشْكِيل ?
		case AlefWaslah:
			runes[i] = Alef
		}
	}
	return string(runes)
}

// Normalize will prepare an arabic text to search and index
func Normalize(input string) string {
	input = normalizeTransform(input)
	runes := bytes.Runes([]byte(input))
	for i := 0; i < len(runes); i++ {
		//fmt.Println(string(runes[i]))
		switch runes[i] {
		//Normalizable letters
		case AlefMad, AlefHamzaAbove, AlefHamzaBelow, AlefWaslah:
			runes[i] = Alef
		case DotlessYae:
			runes[i] = Yae
		case TehMarbuta:
			runes[i] = Hae
		}
	}
	//@TODO: optimize runes by converting it to bytes, arabic letters use only 2 bytes
	return string(runes)
}

// Use text/transform algorithm for faster normalization
func normalizeTransform(input string) string {
	//Use text/transform algorithm for faster normalization
	tm := transform.Chain(runes.Remove(runes.In(normalizable)))
	input, _, _ = transform.String(tm, input)
	return input
}

// deleteRune will delete a rune from the slice while keeping the order of runes
func deleteRune(runes []rune, i int) []rune {
	if i >= len(runes) {
		return runes
	}
	runes = append(runes[:i], runes[i+1:]...)
	return runes
}

// Tashkeel will add matching diacritics to arabic text
func Tashkeel(input string) string {
	JarrWords := []string{"من", "الي", "عن", "على", "مذ", "خلا", "عدا", "حاشا"}
	words := strings.Fields(input)
	for i, word := range words {
		// يُجَرُّ الاسم إذا سُبِق بأحد حروف جرٍّ، مثل كلمة الشركة في جملة: توجّهْتُ إلى الشركةِ
		fmt.Println(Normalize(word))
		if contains(JarrWords, Normalize(word)) {
			words[i+1] += string('\u0650')
		}
	}
	return strings.Join(words, " ")
}

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// Shape will reconstruct arabic text to be connected correctly
func Shape(input string, font pdf.FontHandler) string {
	var foundArabic bool
	for _, letter := range input {
		if IsArabicLetter(letter) {
			foundArabic = true
			break
		}
	}
	if !foundArabic {
		return input
	}

	var langSections []string
	var continousLangAr string
	var continousLangLt string

	for _, letter := range input {
		if IsArabicLetter(letter) {
			if len(continousLangLt) > 0 {
				langSections = append(langSections, strings.TrimSpace(continousLangLt))
			}
			continousLangLt = ""
			continousLangAr += string(letter)
		} else {
			if len(continousLangAr) > 0 {
				langSections = append(langSections, strings.TrimSpace(continousLangAr))
			}
			continousLangAr = ""
			continousLangLt += string(letter)
		}
	}
	if len(continousLangLt) > 0 {
		langSections = append(langSections, strings.TrimSpace(continousLangLt))
	}
	if len(continousLangAr) > 0 {
		langSections = append(langSections, strings.TrimSpace(continousLangAr))
	}

	var shapedSentence []string
	for _, section := range langSections {
		if IsArabic(section) {
			for _, word := range strings.Fields(section) {
				shapedSentence = append(shapedSentence, shapeWord(word, font))
			}
		} else {
			shapedSentence = append(shapedSentence, section)
		}
	}
	//Reverse words
	for i := len(shapedSentence)/2 - 1; i >= 0; i-- {
		opp := len(shapedSentence) - 1 - i
		shapedSentence[i], shapedSentence[opp] = shapedSentence[opp], shapedSentence[i]
	}
	return strings.Join(shapedSentence, " ")
}

// shapeWord will reconstruct an arabic word to be connected correctly
func shapeWord(input string, font pdf.FontHandler) string {
	if !IsArabic(input) {
		return input
	}

	var shapedInput bytes.Buffer

	//Convert input into runes
	inputRunes := []rune(RemoveHarakat(input))
	for i := range inputRunes {
		//Get Bounding back and front letters
		var backLetter, frontLetter rune
		if i-1 >= 0 {
			backLetter = inputRunes[i-1]
		}
		if i != len(inputRunes)-1 {
			frontLetter = inputRunes[i+1]
		}
		//Fix the letter based on bounding letters
		if _, ok := arabicAlphabetShapes[inputRunes[i]]; ok {
			adjustedLetter := adjustLetter(letterGroup{backLetter, inputRunes[i], frontLetter})
			if font.HasGylph([]rune{adjustedLetter})[0] {
				shapedInput.WriteRune(adjustedLetter)
			} else {
				shapedInput.WriteRune(inputRunes[i])
			}
		} else {
			shapedInput.WriteRune(inputRunes[i])
		}
	}

	//In case no Tashkeel deteted, same size of runes
	if len([]rune(shapedInput.String())) == len([]rune(input)) {
		return reverse(shapedInput.String())
	}

	var shapedInputTashkeel bytes.Buffer
	inputTashkeelRunes := []rune(input)

	letterIndex := 0
	//Restore Tashkeel
	for _, r := range inputTashkeelRunes {
		if _, ok := arabicAlphabetShapes[r]; ok {
			rr := []rune(shapedInput.String())
			if letterIndex >= len(rr) {
				shapedInputTashkeel.WriteRune(r)
				continue
			}
			shapedInputTashkeel.WriteRune(rr[letterIndex])
			letterIndex++
		} else {
			shapedInputTashkeel.WriteRune(r)
		}
	}

	return reverse(shapedInputTashkeel.String())
}

// reverse the arabic string for RTL support in rendering
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// adjustLetter will adjust the arabic letter depending on its position
func adjustLetter(g letterGroup) rune {
	switch {
	//Inbetween 2 letters
	case g.backLetter > 0 && g.frontLetter > 0:
		if isAlwaysInitial(g.backLetter) {
			return arabicAlphabetShapes[g.letter].Initial
		}
		return arabicAlphabetShapes[g.letter].Medial

		//Not preceded by any letter
	case g.backLetter == 0 && g.frontLetter > 0:
		return arabicAlphabetShapes[g.letter].Initial

		//Not followed by any letter
	case g.backLetter > 0 && g.frontLetter == 0:
		if isAlwaysInitial(g.backLetter) {
			return arabicAlphabetShapes[g.letter].Independent
		}
		return arabicAlphabetShapes[g.letter].Final

	default:
		return arabicAlphabetShapes[g.letter].Independent
	}
}

// Check if the letter is always .Initial
func isAlwaysInitial(letter rune) bool {
	alwaysInitial := [13]rune{'\u0627', '\u0623', '\u0622', '\u0625', '\u0649', '\u0621', '\u0624', '\u0629', '\u062f', '\u0630', '\u0631', '\u0632', '\u0648'}
	for _, item := range alwaysInitial {
		if item == letter {
			return true
		}
	}
	return false
}

// IsArabicLetter checks if the letter is arabic
func IsArabicLetter(ch rune) bool {
	return (ch >= 0x600 && ch <= 0x6FF)
}

// IsArabic checks if the input string contains arabic unicode only
func IsArabic(input string) bool {
	for _, v := range input {
		if !unicode.IsSpace(v) && !IsArabicLetter(v) {
			return false
		}
	}
	return true
}
