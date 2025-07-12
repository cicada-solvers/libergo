package lgstructs

// RuneTextDocumentCharacter represents a character within a text document and its associated metadata.
// ID is the unique identifier for the RuneTextDocumentCharacter instance.
// TextDocumentId links this character to a specific text document.
// Character stores the single textual character being tracked.
// Count indicates the number of times the character appears in the document.
type RuneTextDocumentCharacter struct {
	ID             string
	TextDocumentId string
	Character      string
	Count          int64
}
