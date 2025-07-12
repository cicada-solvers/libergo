package lgstructs

// TextDocumentCharacter represents a character within a text document, including its ID, document ID, value, and count.
type TextDocumentCharacter struct {
	ID             string
	TextDocumentId string
	Character      string
	Count          int64
}
