package main

import (
	"bufio"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Config flags
type cfg struct {
	headerPath  string
	dataPath    string
	alphabetStr string
	bitOrder    string // "lsb" or "msb"
	maxOut      int
	timeoutSec  int
	printHex    bool
}

func main() {
	var c cfg
	flag.StringVar(&c.headerPath, "header", "", "Path to header byte array (comma-separated integers 0-255)")
	flag.StringVar(&c.dataPath, "data", "", "Path to compressed data file")
	flag.StringVar(&c.alphabetStr, "alphabet", "", "Alphabet as a string (symbols order defines indices)")
	flag.StringVar(&c.bitOrder, "bitorder", "lsb", "Bit order for raw Huffman decode attempts: lsb or msb")
	flag.IntVar(&c.maxOut, "maxout", 1<<20, "Max decompressed output bytes per attempt")
	flag.IntVar(&c.timeoutSec, "timeout", 20, "Max seconds per attempt")
	flag.BoolVar(&c.printHex, "hex", false, "Also print decompressed output as hex")
	flag.Parse()

	if c.headerPath == "" || c.dataPath == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -header header.txt -data data.bin [options]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(2)
	}

	headerBytes, err := loadCommaIntBytes(c.headerPath)
	if err != nil {
		fail("failed to read header: %v", err)
	}
	data, err := os.ReadFile(c.dataPath)
	if err != nil {
		fail("failed to read data: %v", err)
	}

	alphabet := loadAlphabet(c.alphabetStr)

	fmt.Println("Loaded header bytes:")
	fmt.Printf("  count=%d\n", len(headerBytes))
	fmt.Printf("  preview=%s\n", previewBytes(headerBytes, 64))
	fmt.Println()

	// Try parsing header in multiple ways
	parses := parseHeaderCandidates(headerBytes, alphabet)

	// Print parsed header/tree candidates
	fmt.Println("Interpreted header candidates:")
	for i, p := range parses {
		fmt.Printf("- [%d] %s\n", i+1, p.Label)
		fmt.Print(p.Pretty)
	}
	if len(parses) == 0 {
		fmt.Println("- (no recognized interpretations)")
	}
	fmt.Println()

	// Attempts list
	attempts := []attempt{
		attemptGzip(data, c),
		attemptZlib(data, c),
		attemptDeflateRaw(data, c),
		attemptInflateRFC1951(data, c), // new: raw DEFLATE via custom inflate
	}
	// Raw Huffman attempts using parsed header candidates
	for _, p := range parses {
		attempts = append(attempts, attemptHuffman(data, c, p))
	}

	anySuccess := false
	for _, a := range attempts {
		fmt.Printf("Attempt: %s\n", a.Name)
		out, det, err := runWithLimitAndTimeout(a.Run, c.maxOut, time.Duration(c.timeoutSec)*time.Second)
		if err != nil {
			fmt.Printf("  Result: failure: %v\n\n", err)
			continue
		}
		anySuccess = true
		fmt.Println("  Result: success")
		if det != "" {
			fmt.Printf("  Details: %s\n", det)
		}
		printOutput(out, c.printHex)
		fmt.Println()
	}
	if !anySuccess {
		os.Exit(1)
	}
}

// Utilities

func fail(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(1)
}

func previewBytes(b []byte, n int) string {
	if len(b) <= n {
		return intsCSV(b)
	}
	return intsCSV(b[:n]) + ", ..."
}

func intsCSV(b []byte) string {
	parts := make([]string, len(b))
	for i, v := range b {
		parts[i] = strconv.Itoa(int(v))
	}
	return strings.Join(parts, ",")
}

func loadCommaIntBytes(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []byte
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1024), 1024*1024*32)
	for sc.Scan() {
		line := sc.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// allow whitespace and newlines; split on commas
		for _, tok := range strings.Split(line, ",") {
			tok = strings.TrimSpace(tok)
			if tok == "" {
				continue
			}
			v, err := strconv.Atoi(tok)
			if err != nil || v < 0 || v > 255 {
				return nil, fmt.Errorf("invalid byte value %q", tok)
			}
			out = append(out, byte(v))
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func loadAlphabet(inline string) []string {
	return strings.Split(inline, ",")
}

func printOutput(out []byte, alsoHex bool) {
	// Print as UTF-8 best effort
	fmt.Printf("  Output (len=%d):\n", len(out))
	fmt.Println("  -------- BEGIN TEXT --------")
	fmt.Println(string(out))
	fmt.Println("  --------- END TEXT ---------")
	if alsoHex {
		fmt.Println("  Hex:")
		fmt.Println(strings.ToUpper(hex.EncodeToString(out)))
	}
}

// Attempt framework

type attempt struct {
	Name string
	Run  func(limit int, deadline time.Time) (out []byte, details string, err error)
}

func runWithLimitAndTimeout(run func(limit int, deadline time.Time) ([]byte, string, error), limit int, timeout time.Duration) ([]byte, string, error) {
	deadline := time.Now().Add(timeout)
	return run(limit, deadline)
}

// Gzip/zlib/deflate attempts

func attemptGzip(data []byte, c cfg) attempt {
	return attempt{
		Name: "gzip (RFC1952)",
		Run: func(limit int, deadline time.Time) ([]byte, string, error) {
			r, err := gzip.NewReader(bytes.NewReader(data))
			if err != nil {
				return nil, "", err
			}
			defer r.Close()
			limited := &limitedDeadlineReader{r: r, limit: int64(limit), deadline: deadline}
			out, err := io.ReadAll(limited)
			if err != nil {
				return nil, "", err
			}
			return out, fmt.Sprintf("gzip header: name=%q modtime=%v", r.Name, r.ModTime), nil
		},
	}
}

func attemptZlib(data []byte, c cfg) attempt {
	return attempt{
		Name: "zlib (RFC1950)",
		Run: func(limit int, deadline time.Time) ([]byte, string, error) {
			r, err := zlib.NewReader(bytes.NewReader(data))
			if err != nil {
				return nil, "", err
			}
			defer r.Close()
			limited := &limitedDeadlineReader{r: r, limit: int64(limit), deadline: deadline}
			out, err := io.ReadAll(limited)
			if err != nil {
				return nil, "", err
			}
			return out, "", nil
		},
	}
}

func attemptDeflateRaw(data []byte, c cfg) attempt {
	return attempt{
		Name: "deflate raw (RFC1951, no zlib/gzip wrapper)",
		Run: func(limit int, deadline time.Time) ([]byte, string, error) {
			r := flate.NewReader(bytes.NewReader(data))
			defer r.Close()
			limited := &limitedDeadlineReader{r: r, limit: int64(limit), deadline: deadline}
			out, err := io.ReadAll(limited)
			if err != nil {
				return nil, "", err
			}
			return out, "", nil
		},
	}
}

// New: RFC1951 inflate (raw deflate blocks using our bit reader/Huffman)
func attemptInflateRFC1951(data []byte, c cfg) attempt {
	return attempt{
		Name: "inflate RFC1951 (raw DEFLATE blocks via custom Huffman)",
		Run: func(limit int, deadline time.Time) ([]byte, string, error) {
			if limit <= 0 {
				return nil, "", errors.New("limit reached")
			}
			br := newBitReader(data, true) // RFC1951: LSB-first within bytes
			out := make([]byte, 0, min(limit, 1<<16))

			var fixedLitLen, fixedDist *canonicalDecoder

			for {
				if time.Now().After(deadline) {
					return nil, "", errors.New("timeout")
				}
				// BFINAL, BTYPE
				bfinal, ok := br.readBit()
				if !ok {
					if len(out) == 0 {
						return nil, "", errors.New("truncated before first block")
					}
					break
				}
				b0, ok := br.readBit()
				if !ok {
					return nil, "", errors.New("truncated in BTYPE")
				}
				b1, ok := br.readBit()
				if !ok {
					return nil, "", errors.New("truncated in BTYPE")
				}
				btype := int(b0 | (b1 << 1))

				switch btype {
				case 0: // stored
					// align to next byte
					for br.bit != 8 {
						if _, ok := br.readBit(); !ok {
							return nil, "", errors.New("truncated aligning stored block")
						}
					}
					if br.pos+4 > len(br.b) {
						return nil, "", errors.New("truncated stored header")
					}
					lenLE := uint16(br.b[br.pos]) | uint16(br.b[br.pos+1])<<8
					nlenLE := uint16(br.b[br.pos+2]) | uint16(br.b[br.pos+3])<<8
					br.pos += 4
					br.bit = 8
					if lenLE^0xFFFF != nlenLE {
						return nil, "", errors.New("stored LEN/NLEN mismatch")
					}
					L := int(lenLE)
					if br.pos+L > len(br.b) {
						return nil, "", errors.New("truncated stored payload")
					}
					copyLen := L
					if remain := limit - len(out); copyLen > remain {
						copyLen = remain
					}
					out = append(out, br.b[br.pos:br.pos+copyLen]...)
					br.pos += L
					if len(out) >= limit {
						return out, "output limit reached", nil
					}
					if bfinal == 1 {
						return out, "", nil
					}

				case 1, 2: // fixed or dynamic Huffman
					var llDec, ddDec *canonicalDecoder
					if btype == 1 {
						// fixed
						if fixedLitLen == nil || fixedDist == nil {
							llLens := make([]int, 288)
							for i := 0; i <= 143; i++ {
								llLens[i] = 8
							}
							for i := 144; i <= 255; i++ {
								llLens[i] = 9
							}
							for i := 256; i <= 279; i++ {
								llLens[i] = 7
							}
							for i := 280; i <= 287; i++ {
								llLens[i] = 8
							}
							ddLens := make([]int, 32)
							for i := range ddLens {
								ddLens[i] = 5
							}
							var err error
							fixedLitLen, err = newCanonicalDecoder(llLens)
							if err != nil {
								return nil, "", err
							}
							fixedDist, err = newCanonicalDecoder(ddLens)
							if err != nil {
								return nil, "", err
							}
						}
						llDec, ddDec = fixedLitLen, fixedDist
					} else {
						var err error
						llDec, ddDec, err = readDynamicHuffman(br)
						if err != nil {
							return nil, "", err
						}
					}

					// decode symbols
					for {
						if time.Now().After(deadline) {
							return nil, "", errors.New("timeout")
						}
						s, ok := llDec.decode(br)
						if !ok {
							return nil, "", errors.New("truncated in literals/lengths")
						}
						if s < 256 {
							out = append(out, byte(s))
							if len(out) >= limit {
								return out, "output limit reached", nil
							}
							continue
						}
						if s == 256 { // end-of-block
							break
						}
						// length-distance copy
						length, ok := readLength(br, s)
						if !ok {
							return nil, "", errors.New("invalid length symbol")
						}
						ds, ok := ddDec.decode(br)
						if !ok {
							return nil, "", errors.New("truncated in distance code")
						}
						dist, ok := readDistance(br, ds)
						if !ok || dist <= 0 {
							return nil, "", errors.New("invalid distance symbol")
						}
						if dist > len(out) {
							return nil, "", errors.New("distance exceeds output size")
						}
						start := len(out) - dist
						for i := 0; i < length; i++ {
							out = append(out, out[start+i%dist])
							if len(out) >= limit {
								return out, "output limit reached", nil
							}
						}
					}
					if bfinal == 1 {
						return out, "", nil
					}

				default:
					return nil, "", errors.New("unsupported BTYPE")
				}
			}
			return nil, "", nil
		},
	}
}

// Dynamic Huffman reader (RFC1951)
func readDynamicHuffman(br *bitReader) (*canonicalDecoder, *canonicalDecoder, error) {
	readBits := func(n int) (uint32, bool) {
		var v uint32
		for i := 0; i < n; i++ {
			b, ok := br.readBit()
			if !ok {
				return 0, false
			}
			v |= uint32(b) << uint(i) // LSB-first
		}
		return v, true
	}
	hlitv, ok := readBits(5)
	if !ok {
		return nil, nil, errors.New("truncated HLIT")
	}
	hdistv, ok := readBits(5)
	if !ok {
		return nil, nil, errors.New("truncated HDIST")
	}
	hclenv, ok := readBits(4)
	if !ok {
		return nil, nil, errors.New("truncated HCLEN")
	}
	HLIT := int(hlitv) + 257
	HDIST := int(hdistv) + 1
	HCLEN := int(hclenv) + 4

	order := []int{16, 17, 18, 0, 8, 7, 9, 6, 10, 5, 11, 4, 12, 3, 13, 2, 14, 1, 15}
	clLens := make([]int, 19)
	for i := 0; i < HCLEN; i++ {
		v, ok := readBits(3)
		if !ok {
			return nil, nil, errors.New("truncated code length code")
		}
		clLens[order[i]] = int(v)
	}
	clDec, err := newCanonicalDecoder(clLens)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid code length Huffman: %w", err)
	}

	// RLE decode lit/len and dist code lengths
	readRun := func(dst []int) error {
		for i := 0; i < len(dst); {
			s, ok := clDec.decode(br)
			if !ok {
				return errors.New("truncated in code lengths")
			}
			switch s {
			case 0:
				dst[i] = 0
				i++
			case 16:
				if i == 0 {
					return errors.New("repeat(16) with no previous length")
				}
				prev := dst[i-1]
				extra, ok := readBits(2)
				if !ok {
					return errors.New("truncated repeat(16)")
				}
				reps := int(extra) + 3
				if i+reps > len(dst) {
					return errors.New("repeat(16) overruns")
				}
				for r := 0; r < reps; r++ {
					dst[i] = prev
					i++
				}
			case 17:
				extra, ok := readBits(3)
				if !ok {
					return errors.New("truncated repeat(17)")
				}
				reps := int(extra) + 3
				if i+reps > len(dst) {
					return errors.New("repeat(17) overruns")
				}
				for r := 0; r < reps; r++ {
					dst[i] = 0
					i++
				}
			case 18:
				extra, ok := readBits(7)
				if !ok {
					return errors.New("truncated repeat(18)")
				}
				reps := int(extra) + 11
				if i+reps > len(dst) {
					return errors.New("repeat(18) overruns")
				}
				for r := 0; r < reps; r++ {
					dst[i] = 0
					i++
				}
			default:
				if s < 1 || s > 15 {
					return errors.New("invalid code length symbol")
				}
				dst[i] = s
				i++
			}
		}
		return nil
	}

	litlenLens := make([]int, 288)
	if HLIT > len(litlenLens) {
		return nil, nil, errors.New("HLIT too large")
	}
	if err := readRun(litlenLens[:HLIT]); err != nil {
		return nil, nil, err
	}

	distLens := make([]int, 32)
	if HDIST > len(distLens) {
		return nil, nil, errors.New("HDIST too large")
	}
	if err := readRun(distLens[:HDIST]); err != nil {
		return nil, nil, err
	}

	// If all distances are zero, invalid block
	allZero := true
	for i := 0; i < HDIST; i++ {
		if distLens[i] != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		return nil, nil, errors.New("no distance codes in dynamic block")
	}

	llDec, err := newCanonicalDecoder(litlenLens)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid lit/len Huffman: %w", err)
	}
	ddDec, err := newCanonicalDecoder(distLens)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid dist Huffman: %w", err)
	}
	return llDec, ddDec, nil
}

// Length and distance tables (RFC1951)

func readLength(br *bitReader, sym int) (int, bool) {
	if sym < 257 || sym > 285 {
		return 0, false
	}
	if sym == 285 {
		return 258, true
	}
	var base, extra int
	switch {
	case sym <= 264: // 257..264
		base = 3 + (sym - 257)
		extra = 0
	case sym <= 268: // 265..268
		base = 11 + 2*(sym-265)
		extra = 1
	case sym <= 272: // 269..272
		base = 19 + 4*(sym-269)
		extra = 2
	case sym <= 276: // 273..276
		base = 35 + 8*(sym-273)
		extra = 3
	case sym <= 280: // 277..280
		base = 67 + 16*(sym-277)
		extra = 4
	case sym <= 284: // 281..284
		base = 131 + 32*(sym-281)
		extra = 5
	default:
		return 0, false
	}
	val := 0
	for i := 0; i < extra; i++ {
		b, ok := br.readBit()
		if !ok {
			return 0, false
		}
		val |= int(b) << i
	}
	return base + val, true
}

func readDistance(br *bitReader, sym int) (int, bool) {
	if sym < 0 || sym > 29 {
		return 0, false
	}
	if sym <= 3 {
		return sym + 1, true
	}
	var base, extra int
	switch sym {
	case 4, 5:
		base, extra = 5+(sym-4)*2, 1
	case 6, 7:
		base, extra = 9+(sym-6)*4, 2
	case 8, 9:
		base, extra = 17+(sym-8)*8, 3
	case 10, 11:
		base, extra = 33+(sym-10)*16, 4
	case 12, 13:
		base, extra = 65+(sym-12)*32, 5
	case 14, 15:
		base, extra = 129+(sym-14)*64, 6
	case 16, 17:
		base, extra = 257+(sym-16)*128, 7
	case 18, 19:
		base, extra = 513+(sym-18)*256, 8
	case 20, 21:
		base, extra = 1025+(sym-20)*512, 9
	case 22, 23:
		base, extra = 2049+(sym-22)*1024, 10
	case 24, 25:
		base, extra = 4097+(sym-24)*2048, 11
	case 26, 27:
		base, extra = 8193+(sym-26)*4096, 12
	case 28, 29:
		base, extra = 16385+(sym-28)*8192, 13
	default:
		return 0, false
	}
	val := 0
	for i := 0; i < extra; i++ {
		b, ok := br.readBit()
		if !ok {
			return 0, false
		}
		val |= int(b) << i
	}
	return base + val, true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Header parsing and Huffman attempt

type headerParse struct {
	Label  string
	Pretty string
	Kind   parseKind
	// For canonical: code lengths per symbol index (aligned to alphabet)
	CodeLengths []int
	// For serialized tree: preorder nodes with leaves bearing symbol indices
	Tree *node
	// Optional end-of-block symbol index (-1 if none)
	HasEOF bool
	EOFIdx int
}

type parseKind int

const (
	parseUnknown           parseKind = iota
	parseCanonicalCodelens           // assume header is code-lengths for entire alphabet
	parseSerializedTree              // assume header is [flags and symbol indices] preorder
)

// We do not know the exact format. Try two common ones:
// 1) Canonical code lengths: header length should be >= len(alphabet). If so, take first N as code lengths (0..32 accepted).
// 2) Serialized tree: stream of nodes: [tag, ...]; tag 1=leaf then [symbolIndex (1 or 2 bytes?)]; tag 0=internal. We'll try 1-byte symbol index if it fits alphabet.
func parseHeaderCandidates(header []byte, alphabet []string) []headerParse {
	var out []headerParse

	// Candidate 1: canonical code-length table
	if len(header) >= len(alphabet) {
		ok := true
		codeLens := make([]int, len(alphabet))
		for i := 0; i < len(alphabet); i++ {
			cl := int(header[i])
			if cl > 32 {
				ok = false
				break
			}
			codeLens[i] = cl
		}
		if ok {
			label := "canonical Huffman code-lengths over supplied alphabet"
			pretty := prettyPrintCodeLengths(codeLens, alphabet)
			out = append(out, headerParse{
				Label:       label,
				Pretty:      pretty,
				Kind:        parseCanonicalCodelens,
				CodeLengths: codeLens,
				HasEOF:      false,
				EOFIdx:      -1,
			})
		}
	}

	// Candidate 2: simple serialized binary tree:
	// Format assumption:
	//   Node: [tag]
	//     tag=1 -> leaf: [symbolIndex:1 byte]
	//     tag=0 -> internal: then Node(left), Node(right)
	// Stop when one full tree is consumed.
	if len(header) >= 2 {
		if tr, n := tryParseSerializedTree(header, len(alphabet)); tr != nil && n == len(header) {
			label := "serialized binary Huffman tree (preorder: tag 0=internal, 1=leaf+symbolIndex)"
			var b strings.Builder
			fmt.Fprintf(&b, "  Tree size bytes: %d\n", n)
			fmt.Fprintf(&b, "  Tree structure:\n")
			printTree(&b, tr, alphabet, "", true)
			out = append(out, headerParse{
				Label:  label,
				Pretty: b.String(),
				Kind:   parseSerializedTree,
				Tree:   tr,
				HasEOF: false,
				EOFIdx: -1,
			})
		}
	}

	return out
}

type node struct {
	leaf bool
	sym  int
	l    *node
	r    *node
}

func tryParseSerializedTree(b []byte, alphLen int) (*node, int) {
	var pos int
	var parse func() (*node, bool)
	parse = func() (*node, bool) {
		if pos >= len(b) {
			return nil, false
		}
		tag := b[pos]
		pos++
		if tag == 1 {
			if pos >= len(b) {
				return nil, false
			}
			si := int(b[pos])
			pos++
			if si < 0 || si >= alphLen {
				return nil, false
			}
			return &node{leaf: true, sym: si}, true
		} else if tag == 0 {
			left, ok := parse()
			if !ok {
				return nil, false
			}
			right, ok := parse()
			if !ok {
				return nil, false
			}
			return &node{leaf: false, l: left, r: right}, true
		}
		return nil, false
	}
	root, ok := parse()
	if !ok {
		return nil, 0
	}
	// Must consume exactly the buffer to consider it a strong candidate
	if pos != len(b) {
		return nil, pos
	}
	return root, pos
}

func prettyPrintCodeLengths(codeLens []int, alphabet []string) string {
	type entry struct {
		idx  int
		r    string
		cl   int
		code uint32
	}
	// Build canonical codes for pretty-print
	codes, _ := buildCanonical(codeLens)
	var rows []entry
	for i := range codeLens {
		rows = append(rows, entry{idx: i, r: alphabet[i], cl: codeLens[i], code: codes[i]})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].cl != rows[j].cl {
			return rows[i].cl < rows[j].cl
		}
		return rows[i].idx < rows[j].idx
	})
	var b strings.Builder
	fmt.Fprintf(&b, "  Symbols: %d\n", len(alphabet))
	for _, e := range rows {
		if e.cl == 0 {
			continue
		}
		fmt.Fprintf(&b, "    idx=%d sym=%q len=%d code=%s\n", e.idx, printableRune(e.r), e.cl, binCode(e.code, e.cl))
	}
	return b.String()
}

func printableRune(r string) string {
	if r == "\n" {
		return "\\n"
	}
	if r == "\r" {
		return "\\r"
	}
	if r == "\t" {
		return "\\t"
	}

	rValue, _ := strconv.Atoi(r)

	if rValue < 32 || rValue == 0x7f {
		return fmt.Sprintf("\\x%02X", r)
	}
	return string(r)
}

func printTree(b *strings.Builder, n *node, alphabet []string, pad string, last bool) {
	if n == nil {
		return
	}
	conn := "├─"
	nextPad := pad + "│ "
	if last {
		conn = "└─"
		nextPad = pad + "  "
	}
	if n.leaf {
		fmt.Fprintf(b, "%s%s leaf idx=%d sym=%q\n", pad, conn, n.sym, printableRune(alphabet[n.sym]))
		return
	}
	fmt.Fprintf(b, "%s%s internal\n", pad, conn)
	printTree(b, n.l, alphabet, nextPad, false)
	printTree(b, n.r, alphabet, nextPad, true)
}

func binCode(code uint32, length int) string {
	var b strings.Builder
	for i := length - 1; i >= 0; i-- {
		if (code>>uint(i))&1 == 1 {
			b.WriteByte('1')
		} else {
			b.WriteByte('0')
		}
	}
	return b.String()
}

// Huffman attempt

func attemptHuffman(data []byte, c cfg, hp headerParse) attempt {
	name := "raw Huffman using header: " + hp.Label + " [" + c.bitOrder + "-first bits]"
	return attempt{
		Name: name,
		Run: func(limit int, deadline time.Time) ([]byte, string, error) {
			if limit <= 0 {
				return nil, "", errors.New("limit reached")
			}
			br := newBitReader(data, c.bitOrder == "lsb")
			var out []byte
			switch hp.Kind {
			case parseCanonicalCodelens:
				dec, err := newCanonicalDecoder(hp.CodeLengths)
				if err != nil {
					return nil, "", err
				}
				for len(out) < limit && time.Now().Before(deadline) {
					s, ok := dec.decode(br)
					if !ok {
						break
					}
					if hp.HasEOF && s == hp.EOFIdx {
						break
					}
					out = append(out, byte(s)) // map symbol index to byte: use low 8 bits
				}
				if len(out) == 0 {
					return nil, "", errors.New("no symbols decoded")
				}
				return out, "", nil
			case parseSerializedTree:
				for len(out) < limit && time.Now().Before(deadline) {
					s, ok := decodeByTree(br, hp.Tree)
					if !ok {
						break
					}
					out = append(out, byte(s))
				}
				if len(out) == 0 {
					return nil, "", errors.New("no symbols decoded")
				}
				return out, "", nil
			default:
				return nil, "", errors.New("unsupported header kind")
			}
		},
	}
}

// Bit reader

type bitReader struct {
	b     []byte
	pos   int // byte index
	bit   uint8
	lsb   bool
	curr  byte
	ended bool
}

func newBitReader(b []byte, lsbFirst bool) *bitReader {
	return &bitReader{b: b, pos: 0, bit: 8, lsb: lsbFirst}
}

func (r *bitReader) readBit() (uint8, bool) {
	if r.ended {
		return 0, false
	}
	if r.bit >= 8 {
		if r.pos >= len(r.b) {
			r.ended = true
			return 0, false
		}
		r.curr = r.b[r.pos]
		r.pos++
		r.bit = 0
	}
	var v uint8
	if r.lsb {
		v = (r.curr >> r.bit) & 1
	} else {
		v = (r.curr >> (7 - r.bit)) & 1
	}
	r.bit++
	return v, true
}

// Canonical Huffman builder/decoder

func buildCanonical(codeLens []int) ([]uint32, []int) {
	// RFC 1951 style canonical assignment
	maxLen := 0
	for _, l := range codeLens {
		if l > maxLen {
			maxLen = l
		}
	}
	blCount := make([]int, maxLen+1)
	for _, l := range codeLens {
		if l > 0 {
			blCount[l]++
		}
	}
	nextCode := make([]uint32, maxLen+1)
	code := uint32(0)
	for bits := 1; bits <= maxLen; bits++ {
		code = (code + uint32(blCount[bits-1])) << 1
		nextCode[bits] = code
	}
	codes := make([]uint32, len(codeLens))
	for i, l := range codeLens {
		if l == 0 {
			continue
		}
		codes[i] = nextCode[l]
		nextCode[l]++
	}
	return codes, blCount
}

type canonicalDecoder struct {
	minLen int
	maxLen int
	// map: length -> code -> symbol
	tables map[int]map[uint32]int
	// For fast decode: build a prefix table up to prefixBits
	prefixBits int
	prefix     []int // -1 means need slow path
	prefixMask uint32
}

func newCanonicalDecoder(codeLens []int) (*canonicalDecoder, error) {
	codes, _ := buildCanonical(codeLens)
	minL, maxL := 0, 0
	for _, l := range codeLens {
		if l == 0 {
			continue
		}
		if minL == 0 || l < minL {
			minL = l
		}
		if l > maxL {
			maxL = l
		}
	}
	if maxL == 0 {
		return nil, errors.New("no codes")
	}
	// Build tables
	tables := make(map[int]map[uint32]int)
	for sym, l := range codeLens {
		if l == 0 {
			continue
		}
		if tables[l] == nil {
			tables[l] = make(map[uint32]int)
		}
		tables[l][codes[sym]] = sym
	}
	prefixBits := 10
	if maxL < prefixBits {
		prefixBits = maxL
	}
	size := 1 << prefixBits
	prefix := make([]int, size)
	for i := range prefix {
		prefix[i] = -1
	}
	// Fill prefix: for any code of length <= prefixBits, expand to all suffixes
	for sym, l := range codeLens {
		if l == 0 || l > prefixBits {
			continue
		}
		code := codes[sym]
		shift := uint(prefixBits - l)
		base := int(code << shift)
		reps := 1 << shift
		for i := 0; i < reps; i++ {
			prefix[base+i] = sym
		}
	}
	return &canonicalDecoder{
		minLen:     minL,
		maxLen:     maxL,
		tables:     tables,
		prefixBits: prefixBits,
		prefix:     prefix,
		prefixMask: uint32(size - 1),
	}, nil
}

func (d *canonicalDecoder) decode(br *bitReader) (int, bool) {
	// Read up to prefixBits and try fast path
	var acc uint32
	for i := 0; i < d.prefixBits; i++ {
		b, ok := br.readBit()
		if !ok {
			return 0, false
		}
		acc = (acc << 1) | uint32(b)
	}
	idx := d.prefix[acc&d.prefixMask]
	if idx >= 0 {
		return idx, true
	}
	// Slow path: continue reading until some length hits
	code := acc
	length := d.prefixBits
	for {
		b, ok := br.readBit()
		if !ok {
			return 0, false
		}
		code = (code << 1) | uint32(b)
		length++
		if length > d.maxLen {
			return 0, false
		}
		if tbl := d.tables[length]; tbl != nil {
			if sym, ok := tbl[code]; ok {
				return sym, true
			}
		}
	}
}

// Tree-based decode

func decodeByTree(br *bitReader, root *node) (int, bool) {
	n := root
	for {
		if n == nil {
			return 0, false
		}
		if n.leaf {
			return n.sym, true
		}
		b, ok := br.readBit()
		if !ok {
			return 0, false
		}
		if b == 0 {
			n = n.l
		} else {
			n = n.r
		}
	}
}

// IO limit wrapper

type limitedDeadlineReader struct {
	r        io.Reader
	read     int64
	limit    int64
	deadline time.Time
}

func (l *limitedDeadlineReader) Read(p []byte) (int, error) {
	if l.limit >= 0 && l.read >= l.limit {
		return 0, io.ErrUnexpectedEOF
	}
	if !l.deadline.IsZero() && time.Now().After(l.deadline) {
		return 0, errors.New("timeout")
	}
	if int64(len(p)) > l.limit-l.read && l.limit >= 0 {
		p = p[:l.limit-l.read]
	}
	n, err := l.r.Read(p)
	l.read += int64(n)
	return n, err
}
