package lgstructs

// RuneglishTextDocumentCharacter represents a character in a text document
type RuneglishTextDocumentCharacter struct {
	ID             string
	TextDocumentId string
	Character      string
	Count          int64
}
