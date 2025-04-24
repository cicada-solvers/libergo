-- Table for the header file
CREATE TABLE LiberWordHeader (
                                 LiberWordGuid VARCHAR(36) NOT NULL,
                                 LiberWord VARCHAR(255) NOT NULL,
                                 LiberWordLength INT NOT NULL,
                                 LiberWordPosition INT NOT NULL,
                                 LiberWordSection VARCHAR(255) NOT NULL,
                                 PRIMARY KEY (LiberWordGuid)
);

-- Table for the detail file
CREATE TABLE LiberWordDetail (
                                 DictionaryWord VARCHAR(255) NOT NULL,
                                 DictionaryWordDistancePattern TEXT NOT NULL,
                                 WordDistancePatternGuid VARCHAR(36) NOT NULL,
                                 WordListOrigin VARCHAR(255) NOT NULL,
                                 LiberWordGuid VARCHAR(36) NOT NULL,
                                 TranslatedLatin VARCHAR(255) NOT NULL,
                                 PRIMARY KEY (WordDistancePatternGuid),
                                 FOREIGN KEY (LiberWordGuid) REFERENCES LiberWordHeader(LiberWordGuid)
);