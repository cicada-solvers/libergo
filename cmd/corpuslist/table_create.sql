CREATE TABLE IF NOT EXISTS dictionary_words (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    dict_word TEXT COLLATE utf8mb4_general_ci NOT NULL,
    dict_runeglish TEXT COLLATE utf8mb4_general_ci NOT NULL,
    dict_rune TEXT COLLATE utf8mb4_general_ci NOT NULL,
    dict_rune_no_doublet TEXT COLLATE utf8mb4_general_ci NOT NULL,
    gem_sum BIGINT NOT NULL,
    gem_sum_prime TINYINT(1) NOT NULL,
    gem_product TEXT COLLATE utf8mb4_general_ci NOT NULL,
    gem_product_prime TINYINT(1) NOT NULL,
    dict_word_length INT NOT NULL,
    dict_runeglish_length INT NOT NULL,
    dict_rune_length INT NOT NULL,
    dict_rune_no_doublet_length INT NOT NULL,
    rune_pattern TEXT COLLATE utf8mb4_general_ci NOT NULL,
    rune_pattern_no_doublet TEXT COLLATE utf8mb4_general_ci NOT NULL,
    language VARCHAR(255) COLLATE utf8mb4_general_ci NOT NULL
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE INDEX idx_dict_word ON dictionary_words (dict_word(255));
CREATE INDEX idx_gem_sum ON dictionary_words (gem_sum);
CREATE INDEX idx_dict_word_length ON dictionary_words (dict_word_length);
CREATE INDEX idx_dict_runeglish_length ON dictionary_words (dict_runeglish_length);
CREATE INDEX idx_dict_rune_length ON dictionary_words (dict_rune_length);
CREATE INDEX idx_dict_rune_no_doublet_length ON dictionary_words (dict_rune_no_doublet_length);
CREATE INDEX idx_rune_pattern ON dictionary_words (rune_pattern(255));
CREATE INDEX idx_rune_pattern_no_doublet ON dictionary_words (rune_pattern_no_doublet(255));

CREATE TABLE IF NOT EXISTS dict_sentences (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    file_name VARCHAR(1024) COLLATE utf8mb4_general_ci NOT NULL,
    dict_sentence TEXT COLLATE utf8mb4_general_ci NOT NULL
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE INDEX idx_file_name ON dict_sentences (file_name(255));








