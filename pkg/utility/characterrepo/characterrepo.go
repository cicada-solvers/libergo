package runelib

import (
	"sort"
	"strings"
)

// CharacterRepo is a repository for ASCII and ANSI characters.
type CharacterRepo struct {
	asciiAnsiItems map[string]map[int]string
}

// NewCharacterRepo creates a new CharacterRepo.
func NewCharacterRepo() *CharacterRepo {
	repo := &CharacterRepo{
		asciiAnsiItems: make(map[string]map[int]string),
	}

	repo.asciiAnsiItems["ASCII"] = map[int]string{
		0: "<NUL>", 1: "<SOH>", 2: "<STX>", 3: "<ETX>", 4: "<EOT>", 5: "<ENQ>", 6: "<ACK>", 7: "<BEL>",
		8: "<BS>", 9: "<HT>", 10: "\n", 11: "<VT>", 12: "<FF>", 13: "\r", 14: "<SO>", 15: "<SI>",
		16: "<DLE>", 17: "<DC1>", 18: "<DC2>", 19: "<DC3>", 20: "<DC4>", 21: "<NAK>", 22: "<SYN>", 23: "<ETB>",
		24: "<CAN>", 25: "<EM>", 26: "<SUB>", 27: "<ESC>", 28: "<FS>", 29: "<GS>", 30: "<RS>", 31: "<US>",
		32: " ", 33: "!", 34: "\"", 35: "#", 36: "$", 37: "%", 38: "&", 39: "'",
		40: "(", 41: ")", 42: "*", 43: "+", 44: ",", 45: "-", 46: ".", 47: "/",
		48: "0", 49: "1", 50: "2", 51: "3", 52: "4", 53: "5", 54: "6", 55: "7",
		56: "8", 57: "9", 58: ":", 59: ";", 60: "<", 61: "=", 62: ">", 63: "?",
		64: "@", 65: "A", 66: "B", 67: "C", 68: "D", 69: "E", 70: "F", 71: "G",
		72: "H", 73: "I", 74: "J", 75: "K", 76: "L", 77: "M", 78: "N", 79: "O",
		80: "P", 81: "Q", 82: "R", 83: "S", 84: "T", 85: "U", 86: "V", 87: "W",
		88: "X", 89: "Y", 90: "Z", 91: "[", 92: "\\", 93: "]", 94: "^", 95: "_",
		96: "`", 97: "a", 98: "b", 99: "c", 100: "d", 101: "e", 102: "f", 103: "g",
		104: "h", 105: "i", 106: "j", 107: "k", 108: "l", 109: "m", 110: "n", 111: "o",
		112: "p", 113: "q", 114: "r", 115: "s", 116: "t", 117: "u", 118: "v", 119: "w",
		120: "x", 121: "y", 122: "z", 123: "{", 124: "|", 125: "}", 126: "~", 127: "<DEL>",
	}

	repo.asciiAnsiItems["ANSI"] = map[int]string{
		0: "<NUL>", 1: "<SOH>", 2: "<STX>", 3: "<ETX>", 4: "<EOT>", 5: "<ENQ>", 6: "<ACK>", 7: "<BEL>",
		8: "<BS>", 9: "<HT>", 10: "\n", 11: "<VT>", 12: "<FF>", 13: "\r", 14: "<SO>", 15: "<SI>",
		16: "<DLE>", 17: "<DC1>", 18: "<DC2>", 19: "<DC3>", 20: "<DC4>", 21: "<NAK>", 22: "<SYN>", 23: "<ETB>",
		24: "<CAN>", 25: "<EM>", 26: "<SUB>", 27: "<ESC>", 28: "<FS>", 29: "<GS>", 30: "<RS>", 31: "<US>",
		32: " ", 33: "!", 34: "\"", 35: "#", 36: "$", 37: "%", 38: "&", 39: "'",
		40: "(", 41: ")", 42: "*", 43: "+", 44: ",", 45: "-", 46: ".", 47: "/",
		48: "0", 49: "1", 50: "2", 51: "3", 52: "4", 53: "5", 54: "6", 55: "7",
		56: "8", 57: "9", 58: ":", 59: ";", 60: "<", 61: "=", 62: ">", 63: "?",
		64: "@", 65: "A", 66: "B", 67: "C", 68: "D", 69: "E", 70: "F", 71: "G",
		72: "H", 73: "I", 74: "J", 75: "K", 76: "L", 77: "M", 78: "N", 79: "O",
		80: "P", 81: "Q", 82: "R", 83: "S", 84: "T", 85: "U", 86: "V", 87: "W",
		88: "X", 89: "Y", 90: "Z", 91: "[", 92: "\\", 93: "]", 94: "^", 95: "_",
		96: "`", 97: "a", 98: "b", 99: "c", 100: "d", 101: "e", 102: "f", 103: "g",
		104: "h", 105: "i", 106: "j", 107: "k", 108: "l", 109: "m", 110: "n", 111: "o",
		112: "p", 113: "q", 114: "r", 115: "s", 116: "t", 117: "u", 118: "v", 119: "w",
		120: "x", 121: "y", 122: "z", 123: "{", 124: "|", 125: "}", 126: "~", 127: "<DEL>",
		128: "�", 129: "", 130: "‚", 131: "ƒ", 132: "„", 133: "…", 134: "†", 135: "‡",
		136: "ˆ", 137: "‰", 138: "Š", 139: "‹", 140: "Œ", 141: "", 142: "Ž", 143: "",
		144: "", 145: "‘", 146: "’", 147: "“", 148: "”", 149: "•", 150: "–", 151: "—",
		152: "˜", 153: "™", 154: "š", 155: "›", 156: "œ", 157: "", 158: "ž", 159: "Ÿ",
		160: "", 161: "¡", 162: "¢", 163: "£", 164: "¤", 165: "¥", 166: "¦", 167: "§",
		168: "¨", 169: "©", 170: "ª", 171: "«", 172: "¬", 173: "", 174: "®", 175: "¯",
		176: "°", 177: "±", 178: "²", 179: "³", 180: "´", 181: "µ", 182: "¶", 183: "·",
		184: "¸", 185: "¹", 186: "º", 187: "»", 188: "¼", 189: "½", 190: "¾", 191: "¿",
		192: "À", 193: "Á", 194: "Â", 195: "Ã", 196: "Ä", 197: "Å", 198: "Æ", 199: "Ç",
		200: "È", 201: "É", 202: "Ê", 203: "Ë", 204: "Ì", 205: "Í", 206: "Î", 207: "Ï",
		208: "Ð", 209: "Ñ", 210: "Ò", 211: "Ó", 212: "Ô", 213: "Õ", 214: "Ö", 215: "×",
		216: "Ø", 217: "Ù", 218: "Ú", 219: "Û", 220: "Ü", 221: "Ý", 222: "Þ", 223: "ß",
		224: "à", 225: "á", 226: "â", 227: "ã", 228: "ä", 229: "å", 230: "æ", 231: "ç",
		232: "è", 233: "é", 234: "ê", 235: "ë", 236: "ì", 237: "í", 238: "î", 239: "ï",
		240: "ð", 241: "ñ", 242: "ò", 243: "ó", 244: "ô", 245: "õ", 246: "ö", 247: "÷",
		248: "ø", 249: "ù", 250: "ú", 251: "û", 252: "ü", 253: "ý", 254: "þ", 255: "ÿ",
	}

	return repo
}

// GetANSICharFromBin returns the ANSI character for the given binary value.
func (repo *CharacterRepo) GetANSICharFromBin(bin string, includeControlCharacters bool) string {
	for _, char := range repo.asciiAnsiItems["ANSI"] {
		if strings.Contains(char, bin) {
			if !includeControlCharacters && strings.HasPrefix(char, "<") && strings.HasSuffix(char, ">") {
				return ""
			}
			return char
		}
	}
	return ""
}

// GetANSICharFromDec returns the ANSI character for the given decimal value.
func (repo *CharacterRepo) GetANSICharFromDec(dec int, includeControlCharacters bool) string {
	char, exists := repo.asciiAnsiItems["ANSI"][dec]
	if exists {
		if !includeControlCharacters && strings.HasPrefix(char, "<") && strings.HasSuffix(char, ">") {
			return ""
		}
		return char
	}
	return ""
}

// GetASCIICharFromBin returns the ASCII character for the given binary value.
func (repo *CharacterRepo) GetASCIICharFromBin(bin string, includeControlCharacters bool) string {
	for _, char := range repo.asciiAnsiItems["ASCII"] {
		if strings.Contains(char, bin) {
			if !includeControlCharacters && strings.HasPrefix(char, "<") && strings.HasSuffix(char, ">") {
				return ""
			}
			return char
		}
	}
	return ""
}

// GetASCIICharFromDec returns the ASCII character for the given decimal value.
func (repo *CharacterRepo) GetASCIICharFromDec(dec int, includeControlCharacters bool) string {
	char, exists := repo.asciiAnsiItems["ASCII"][dec]
	if exists {
		if !includeControlCharacters && strings.HasPrefix(char, "<") && strings.HasSuffix(char, ">") {
			return ""
		}
		return char
	}
	return ""
}

// GetGematriaRunes returns a list of runes used in the rune alphabet.
func (repo *CharacterRepo) GetGematriaRunes() []string {
	var retval []string

	// Create a slice of key-value pairs
	type kv struct {
		Key   string
		Value int
	}
	var sortedRunes []kv
	for k, v := range runeToValueMap {
		sortedRunes = append(sortedRunes, kv{k, v})
	}

	// Sort the slice by the integer values
	sort.Slice(sortedRunes, func(i, j int) bool {
		return sortedRunes[i].Value < sortedRunes[j].Value
	})

	// Print the sorted runes
	for _, kv := range sortedRunes {
		retval = append(retval, kv.Key)
	}

	return retval
}

var runeToCharMap = map[string]string{
	"ᛝ": "ING",
	"ᛟ": "OE",
	"ᛇ": "EO",
	"ᛡ": "IO",
	"ᛠ": "EA",
	"ᚫ": "AE",
	"ᚦ": "TH",
	"ᚠ": "F",
	"ᚢ": "U",
	"ᚩ": "O",
	"ᚱ": "R",
	"ᚳ": "C",
	"ᚷ": "G",
	"ᚹ": "W",
	"ᚻ": "H",
	"ᚾ": "N",
	"ᛁ": "I",
	"ᛄ": "J",
	"ᛈ": "P",
	"ᛉ": "X",
	"ᛋ": "S",
	"ᛏ": "T",
	"ᛒ": "B",
	"ᛖ": "E",
	"ᛗ": "M",
	"ᛚ": "L",
	"ᛞ": "D",
	"ᚪ": "A",
	"ᚣ": "Y",
	"•": " ",
	"⊹": ".",
}

func (repo *CharacterRepo) GetCharFromRune(value string) string {
	if char, exists := runeToCharMap[value]; exists {
		return char
	}
	return value
}

var charToRuneMap = map[string]string{
	"ING": "ᛝ",
	"NG":  "ᛝ",
	"OE":  "ᛟ",
	"EO":  "ᛇ",
	"IO":  "ᛡ",
	"IA":  "ᛡ",
	"EA":  "ᛠ",
	"AE":  "ᚫ",
	"TH":  "ᚦ",
	"F":   "ᚠ",
	"V":   "ᚢ",
	"U":   "ᚢ",
	"O":   "ᚩ",
	"R":   "ᚱ",
	"Q":   "ᚳ",
	"K":   "ᚳ",
	"C":   "ᚳ",
	"G":   "ᚷ",
	"W":   "ᚹ",
	"H":   "ᚻ",
	"N":   "ᚾ",
	"I":   "ᛁ",
	"J":   "ᛄ",
	"P":   "ᛈ",
	"X":   "ᛉ",
	"Z":   "ᛋ",
	"S":   "ᛋ",
	"T":   "ᛏ",
	"B":   "ᛒ",
	"E":   "ᛖ",
	"M":   "ᛗ",
	"L":   "ᛚ",
	"D":   "ᛞ",
	"A":   "ᚪ",
	"Y":   "ᚣ",
	" ":   "•",
	".":   "⊹",
}

func (repo *CharacterRepo) GetRuneFromChar(value string) string {
	if runeChar, exists := charToRuneMap[value]; exists {
		return runeChar
	}
	return value
}

var runeSet = map[string]struct{}{
	"ᛝ": {}, "ᛟ": {}, "ᛇ": {}, "ᛡ": {}, "ᛠ": {}, "ᚫ": {}, "ᚦ": {}, "ᚠ": {},
	"ᚢ": {}, "ᚩ": {}, "ᚱ": {}, "ᚳ": {}, "ᚷ": {}, "ᚹ": {}, "ᚻ": {}, "ᚾ": {},
	"ᛁ": {}, "ᛄ": {}, "ᛈ": {}, "ᛉ": {}, "ᛋ": {}, "ᛏ": {}, "ᛒ": {}, "ᛖ": {},
	"ᛗ": {}, "ᛚ": {}, "ᛞ": {}, "ᚪ": {}, "ᚣ": {},
}

var dunkusSet = map[string]struct{}{
	"•": {}, "⊹": {},
}

func (repo *CharacterRepo) IsDinkus(value string) bool {
	if _, exists := dunkusSet[value]; exists {
		return true
	}
	return false
}

var seperatorSet = map[string]struct{}{
	" ": {},
	".": {},
	",": {},
	"!": {},
	"?": {},
	":": {},
	";": {},
	"(": {},
	")": {},
}

func (repo *CharacterRepo) IsSeperator(value string) bool {
	if _, exists := seperatorSet[value]; exists {
		return true
	}
	return false
}

func (repo *CharacterRepo) IsRune(value string, includeDunkus bool) bool {
	if includeDunkus {
		if _, exists := dunkusSet[value]; exists {
			return true
		}
	}
	_, exists := runeSet[value]
	return exists
}

var runeToValueMap = map[string]int{
	"ᛝ": 79, "ᛟ": 83, "ᛇ": 41, "ᛡ": 107, "ᛠ": 109, "ᚫ": 101, "ᚦ": 5, "ᚠ": 2,
	"ᚢ": 3, "ᚩ": 7, "ᚱ": 11, "ᚳ": 13, "ᚷ": 17, "ᚹ": 19, "ᚻ": 23, "ᚾ": 29,
	"ᛁ": 31, "ᛄ": 37, "ᛈ": 43, "ᛉ": 47, "ᛋ": 53, "ᛏ": 59, "ᛒ": 61, "ᛖ": 67,
	"ᛗ": 71, "ᛚ": 73, "ᛞ": 89, "ᚪ": 97, "ᚣ": 103,
}

func (repo *CharacterRepo) GetValueFromRune(rune string) int {
	if value, exists := runeToValueMap[rune]; exists {
		return value
	}
	return 0
}

var valueToRuneMap = map[int]string{
	2: "ᚠ", 3: "ᚢ", 5: "ᚦ", 7: "ᚩ", 11: "ᚱ", 13: "ᚳ", 17: "ᚷ", 19: "ᚹ",
	23: "ᚻ", 29: "ᚾ", 31: "ᛁ", 37: "ᛄ", 41: "ᛇ", 43: "ᛈ", 47: "ᛉ", 53: "ᛋ",
	59: "ᛏ", 61: "ᛒ", 67: "ᛖ", 71: "ᛗ", 73: "ᛚ", 79: "ᛝ", 83: "ᛟ", 89: "ᛞ",
	97: "ᚪ", 101: "ᚫ", 103: "ᚣ", 107: "ᛡ", 109: "ᛠ",
}

func (repo *CharacterRepo) GetRuneFromValue(value int) string {
	if runeChar, exists := valueToRuneMap[value]; exists {
		return runeChar
	}
	return ""
}

// GetDoubletCount returns the count of doublets (two consecutive identical characters) in a string.
func (repo *CharacterRepo) GetDoubletCount(input string, alphabet []string) int {
	count := 0

	for _, character := range alphabet {
		doubletChar := character + character

		// Now we are going to count the times the doublet appears in the input string
		count += strings.Count(input, doubletChar)
	}

	return count
}

// GetRunglishAlphabet returns the runglish alphabet.
func (repo *CharacterRepo) GetRunglishAlphabet() []string {
	return []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "L", "M", "N", "O", "P", "R", "S", "T", "U", "W", "X", "Y"}
}
