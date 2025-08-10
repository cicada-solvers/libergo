package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type Crawler struct {
	startURL   *url.URL
	maxDepth   int
	sameDomain bool

	httpClient *http.Client

	sem chan struct{} // concurrency limiter

	mu         sync.Mutex
	visited    map[string]struct{} // normalized URL -> set membership
	discovered map[string]struct{} // all found URLs (including external, if not filtered)

	wg sync.WaitGroup
}

func NewCrawler(start string, concurrency int, timeout time.Duration, maxDepth int, sameDomain bool) (*Crawler, error) {
	if concurrency <= 0 {
		concurrency = 8
	}
	if maxDepth < 0 {
		maxDepth = 0
	}
	u, err := url.Parse(start)
	if err != nil {
		return nil, fmt.Errorf("invalid start url: %w", err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, errors.New("start url must include scheme and host, e.g., https://example.com")
	}
	return &Crawler{
		startURL:   u,
		maxDepth:   maxDepth,
		sameDomain: sameDomain,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		sem:        make(chan struct{}, concurrency),
		visited:    make(map[string]struct{}),
		discovered: make(map[string]struct{}),
	}, nil
}

func (c *Crawler) Start(ctx context.Context) {
	// seed with the start URL
	c.enqueue(ctx, c.startURL, 0)

	c.wg.Wait()
}

func (c *Crawler) enqueue(ctx context.Context, u *url.URL, depth int) {
	nu := normalizeURL(u)
	if nu == "" {
		return
	}

	// same-domain restriction for traversal
	if c.sameDomain && !sameHost(c.startURL, u) {
		// Still record discovery if found, but don't traverse
		c.mu.Lock()
		c.discovered[nu] = struct{}{}
		c.mu.Unlock()
		return
	}

	c.mu.Lock()
	if _, ok := c.visited[nu]; ok {
		c.mu.Unlock()
		return
	}
	c.visited[nu] = struct{}{}
	c.discovered[nu] = struct{}{}
	c.mu.Unlock()

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.crawl(ctx, u, depth)
	}()
}

func (c *Crawler) crawl(ctx context.Context, u *url.URL, depth int) {
	if depth > c.maxDepth {
		return
	}

	select {
	case c.sem <- struct{}{}:
		// acquired
	case <-ctx.Done():
		return
	}
	defer func() { <-c.sem }()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Go-Recursive-Link-Crawler/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Only parse HTML content
	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(strings.ToLower(ct), "text/html") {
		return
	}

	// Parse and extract links
	links, err := extractLinks(resp.Body, u)
	if err != nil {
		return
	}

	// Record all discovered links, and enqueue eligible ones
	for _, link := range links {
		nu := normalizeURL(link)
		if nu == "" {
			continue
		}

		c.mu.Lock()
		c.discovered[nu] = struct{}{}
		c.mu.Unlock()

		// Traverse if within constraints
		if depth < c.maxDepth {
			if !c.sameDomain || sameHost(c.startURL, link) {
				c.enqueue(ctx, link, depth+1)
			}
		}
	}
}

func extractLinks(r io.Reader, base *url.URL) ([]*url.URL, error) {
	root, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	var out []*url.URL
	var walker func(*html.Node)
	walker = func(n *html.Node) {
		if n.Type == html.ElementNode && strings.EqualFold(n.Data, "a") {
			for _, a := range n.Attr {
				if strings.EqualFold(a.Key, "href") {
					raw := strings.TrimSpace(a.Val)
					if raw == "" || skipScheme(raw) {
						continue
					}
					ref, err := url.Parse(raw)
					if err != nil {
						continue
					}
					abs := base.ResolveReference(ref)
					out = append(out, abs)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walker(c)
		}
	}
	walker(root)
	return out, nil
}

func skipScheme(raw string) bool {
	l := strings.ToLower(raw)
	switch {
	case strings.HasPrefix(l, "javascript:"),
		strings.HasPrefix(l, "mailto:"),
		strings.HasPrefix(l, "tel:"),
		strings.HasPrefix(l, "#"):
		return true
	default:
		return false
	}
}

func sameHost(a, b *url.URL) bool {
	// Treat empty port vs default port as equal; compare hostnames case-insensitively
	ha := strings.ToLower(hostWithoutDefaultPort(a))
	hb := strings.ToLower(hostWithoutDefaultPort(b))
	return ha == hb
}

func hostWithoutDefaultPort(u *url.URL) string {
	h := u.Host
	if strings.HasSuffix(h, ":80") && u.Scheme == "http" {
		return strings.TrimSuffix(h, ":80")
	}
	if strings.HasSuffix(h, ":443") && u.Scheme == "https" {
		return strings.TrimSuffix(h, ":443")
	}
	return h
}

func normalizeURL(u *url.URL) string {
	if u.Scheme == "" || u.Host == "" {
		return ""
	}
	nu := *u
	nu.Fragment = "" // drop fragment
	// Normalize path: remove default index pages is opinionated; we'll keep full path but remove trailing slash except root
	if nu.Path == "" {
		nu.Path = "/"
	}
	// Remove default ports
	nu.Host = hostWithoutDefaultPort(&nu)
	// Remove query params that are just tracking? Not safe generically; keep as-is.
	// Lower-case scheme and host
	nu.Scheme = strings.ToLower(nu.Scheme)
	nu.Host = strings.ToLower(nu.Host)
	// Normalize trailing slash: keep if root, remove if not root and present
	if nu.Path != "/" {
		nu.Path = strings.TrimRight(nu.Path, "/")
		if nu.Path == "" {
			nu.Path = "/"
		}
	}
	return nu.String()
}

func writeToFile(path string, urls map[string]struct{}) error {
	// Ensure directory exists
	if dir := filepath.Dir(path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Sort URLs for stable output
	list := make([]string, 0, len(urls))
	for u := range urls {
		list = append(list, u)
	}
	sort.Strings(list)

	for _, u := range list {
		if _, err := io.WriteString(f, u+"\n"); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var (
		start       string
		outPath     string
		concurrency int
		timeoutSec  int
		maxDepth    int
		sameDomain  bool
	)
	flag.StringVar(&start, "url", "", "Start URL to crawl (e.g., https://example.com)")
	flag.StringVar(&outPath, "out", "urls.txt", "Output file to write discovered URLs")
	flag.IntVar(&concurrency, "concurrency", 10, "Maximum concurrent HTTP requests")
	flag.IntVar(&timeoutSec, "timeout", 15, "HTTP client timeout in seconds")
	flag.IntVar(&maxDepth, "max-depth", 2, "Maximum crawl depth (0 = only the start page)")
	flag.BoolVar(&sameDomain, "same-domain", true, "Only crawl within the same domain as the start URL")
	flag.Parse()

	if start == "" {
		fmt.Fprintln(os.Stderr, "error: -url is required")
		os.Exit(1)
	}

	c, err := NewCrawler(start, concurrency, time.Duration(timeoutSec)*time.Second, maxDepth, sameDomain)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startTime := time.Now()
	c.Start(ctx)
	elapsed := time.Since(startTime)

	if err := writeToFile(outPath, c.discovered); err != nil {
		fmt.Fprintln(os.Stderr, "error writing output:", err)
		os.Exit(1)
	}

	fmt.Printf("Crawled %d unique URLs in %s. Results written to %s\n", len(c.discovered), elapsed.Truncate(time.Millisecond), outPath)
}
