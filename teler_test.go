package teler

import (
	"os"
	"testing"

	"net/http"
	"net/http/httptest"
	"path/filepath"

	"github.com/kitabisa/teler-waf/request"
	"github.com/kitabisa/teler-waf/threat"
	"github.com/stretchr/testify/assert"
)

// Prepraring handler for all cases
var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
})

func init() {
	cache := os.Getenv("CACHE")

	switch cache {
	case "1", "true", "TRUE":
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	cacheDir := filepath.Join(homeDir, ".cache", "teler-waf")
	os.RemoveAll(cacheDir)
}

func TestNewDefaultOptions(t *testing.T) {
	// Initialize teler
	telerMiddleware := New(Options{})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the user agent to "X"
	req.Header.Set("User-Agent", "X")

	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewWithNoStderr(t *testing.T) {
	// Initialize teler
	telerMiddleware := New(Options{NoStderr: true})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewWithNoUpdateCheck(t *testing.T) {
	// Initialize teler
	telerMiddleware := New(Options{NoStderr: true, NoUpdateCheck: true})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewWithLogFile(t *testing.T) {
	// Initialize teler
	telerMiddleware := New(Options{NoStderr: true, LogFile: "/dev/null"})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewWithWhitelist(t *testing.T) {
	// Initialize teler
	telerMiddleware := New(Options{
		Whitelists: []string{"Go-http-client"},
		NoStderr:   true,
	})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewCustom(t *testing.T) {
	// Initialize teler
	telerMiddleware := New(Options{
		Excludes: []threat.Threat{
			threat.CommonWebAttack,
			threat.CVE,
			threat.BadIPAddress,
			threat.BadReferrer,
			threat.BadCrawler,
			threat.DirectoryBruteforce,
		},
		Customs: []Rule{
			{
				Name:      "Log4j Attack",
				Condition: "or",
				Rules: []Condition{
					{
						Method:  request.GET,
						Element: request.URI,
						Pattern: `\$\{.*:\/\/.*\/?\w+?\}`,
					},
				},
			},
			{
				Name:      "And condition",
				Condition: "and",
				Rules: []Condition{
					{
						Element: request.Headers,
						Pattern: `Go-http-client`,
					},
					{
						Element: request.URI,
						Pattern: `.`,
					},
				},
			},
			{
				Name:      "Headers element",
				Condition: "and",
				Rules: []Condition{
					{
						Element: request.Headers,
						Pattern: `.`,
					},
				},
			},
			{
				Name:      "Body element",
				Condition: "and",
				Rules: []Condition{
					{
						Element: request.Headers,
						Pattern: `.`,
					},
				},
			},
		},
		NoStderr: true,
	})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewCommonWebAttackOnly(t *testing.T) {
	// Initialize teler
	telerMiddleware := New(Options{
		Excludes: []threat.Threat{
			threat.CVE,
			threat.BadIPAddress,
			threat.BadReferrer,
			threat.BadCrawler,
			threat.DirectoryBruteforce,
		},
		NoStderr: true,
	})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewCVEOnly(t *testing.T) {
	// Initialize teler
	telerMiddleware := New(Options{
		Excludes: []threat.Threat{
			threat.CommonWebAttack,
			threat.BadIPAddress,
			threat.BadReferrer,
			threat.BadCrawler,
			threat.DirectoryBruteforce,
		},
		NoStderr: true,
	})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the request path to "/vcac/" (CVE-2022-22972)
	req.URL.Path = "/vcac/"

	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewBadIPAddressOnly(t *testing.T) {
	// Initialize teler
	telerMiddleware := New(Options{
		Excludes: []threat.Threat{
			threat.CommonWebAttack,
			threat.CVE,
			threat.BadReferrer,
			threat.BadCrawler,
			threat.DirectoryBruteforce,
		},
		NoStderr: true,
	})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the custom header for X-Real-Ip
	req.Header.Set("X-Real-Ip", "1.14.77.81")

	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewBadReferrerOnly(t *testing.T) {
	// Initialize teler
	telerMiddleware := New(Options{
		Excludes: []threat.Threat{
			threat.CommonWebAttack,
			threat.CVE,
			threat.BadIPAddress,
			threat.BadCrawler,
			threat.DirectoryBruteforce,
		},
		NoStderr: true,
	})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the HTTP referrer of the request
	req.Header.Set("Referer", "http://34.gs/")

	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewBadCrawlerOnly(t *testing.T) {
	// Initialize teler
	telerMiddleware := New(Options{
		Excludes: []threat.Threat{
			threat.CommonWebAttack,
			threat.CVE,
			threat.BadIPAddress,
			threat.BadReferrer,
			threat.DirectoryBruteforce,
		},
		NoStderr: true,
	})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewDirectoryBruteforceOnly(t *testing.T) {
	// Initialize teler
	telerMiddleware := New(Options{
		Excludes: []threat.Threat{
			threat.CommonWebAttack,
			threat.CVE,
			threat.BadIPAddress,
			threat.BadReferrer,
			threat.BadCrawler,
		},
		NoStderr: true,
	})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the request path
	req.URL.Path = "/.git"

	_, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewInvalidWhitelist(t *testing.T) {
	defer func() {
		// Check that the teler function panics
		if r := recover(); r != nil {
			assert.Panics(t, func() { panic(r) })
		}
	}()

	// Initialize teler
	telerMiddleware := New(Options{
		Whitelists: []string{`foo(?!bar)`},
		NoStderr:   true,
	})

	telerMiddleware.Handler(handler)
}

func TestNewInvalidCustomRuleName(t *testing.T) {
	defer func() {
		// Check that the teler function panics
		if r := recover(); r != nil {
			assert.Panics(t, func() { panic(r) })
		}
	}()

	// Initialize teler
	telerMiddleware := New(Options{
		Customs: []Rule{
			{
				Name:      "",
				Condition: "or",
				Rules: []Condition{
					{
						Method:  request.GET,
						Element: request.URI,
						Pattern: `.`,
					},
				},
			},
		},
		NoStderr: true,
	})

	telerMiddleware.Handler(handler)
}

func TestNewInvalidCustomRuleCondition(t *testing.T) {
	defer func() {
		// Check that the teler function panics
		if r := recover(); r != nil {
			assert.Panics(t, func() { panic(r) })
		}
	}()

	// Initialize teler
	telerMiddleware := New(Options{
		Customs: []Rule{
			{
				Name:      "foo",
				Condition: "bar",
				Rules: []Condition{
					{
						Method:  request.GET,
						Element: request.URI,
						Pattern: `.`,
					},
				},
			},
		},
		NoStderr: true,
	})

	telerMiddleware.Handler(handler)
}

func TestNewBlankCustomRulePattern(t *testing.T) {
	defer func() {
		// Check that the teler function panics
		if r := recover(); r != nil {
			assert.Panics(t, func() { panic(r) })
		}
	}()

	// Initialize teler
	telerMiddleware := New(Options{
		Customs: []Rule{
			{
				Name:      "foo",
				Condition: "or",
				Rules: []Condition{
					{
						Method:  request.GET,
						Element: request.URI,
						Pattern: "",
					},
				},
			},
		},
		NoStderr: true,
	})

	telerMiddleware.Handler(handler)
}

func TestNewInvalidCustomRulePattern(t *testing.T) {
	defer func() {
		// Check that the teler function panics
		if r := recover(); r != nil {
			assert.Panics(t, func() { panic(r) })
		}
	}()

	// Initialize teler
	telerMiddleware := New(Options{
		Customs: []Rule{
			{
				Name:      "foo",
				Condition: "or",
				Rules: []Condition{
					{
						Method:  request.GET,
						Element: request.URI,
						Pattern: `foo(?!bar)`,
					},
				},
			},
		},
		NoStderr: true,
	})

	telerMiddleware.Handler(handler)
}

func BenchmarkTelerDefaultOptions(b *testing.B) {
	// Initialize teler
	telerMiddleware := New(Options{NoStderr: true})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		b.Fatal(err)
	}

	// Set the custom User-Agent so that the operation does
	// not stop at the BadCrawler check
	req.Header.Set("User-Agent", "awikwok")

	// Run the benchmark
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Send the request to the test server and discard the response
		_, err := client.Do(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTelerCommonWebAttackOnly(b *testing.B) {
	// Initialize teler
	telerMiddleware := New(Options{
		Excludes: []threat.Threat{
			threat.CVE,
			threat.BadIPAddress,
			threat.BadReferrer,
			threat.BadCrawler,
			threat.DirectoryBruteforce,
		},
		NoStderr: true,
	})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		b.Fatal(err)
	}

	// Run the benchmark
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Send the request to the test server and discard the response
		_, err := client.Do(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTelerCVEOnly(b *testing.B) {
	// Initialize teler
	telerMiddleware := New(Options{
		Excludes: []threat.Threat{
			threat.CommonWebAttack,
			threat.BadIPAddress,
			threat.BadReferrer,
			threat.BadCrawler,
			threat.DirectoryBruteforce,
		},
		NoStderr: true,
	})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		b.Fatal(err)
	}

	// Run the benchmark
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Send the request to the test server and discard the response
		_, err := client.Do(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTelerBadIPAddressOnly(b *testing.B) {
	// Initialize teler
	telerMiddleware := New(Options{
		Excludes: []threat.Threat{
			threat.CommonWebAttack,
			threat.CVE,
			threat.BadReferrer,
			threat.BadCrawler,
			threat.DirectoryBruteforce,
		},
		NoStderr: true,
	})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		b.Fatal(err)
	}

	// Run the benchmark
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Send the request to the test server and discard the response
		_, err := client.Do(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTelerBadReferrerOnly(b *testing.B) {
	// Initialize teler
	telerMiddleware := New(Options{
		Excludes: []threat.Threat{
			threat.CommonWebAttack,
			threat.CVE,
			threat.BadIPAddress,
			threat.BadCrawler,
			threat.DirectoryBruteforce,
		},
		NoStderr: true,
	})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		b.Fatal(err)
	}

	// Run the benchmark
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Send the request to the test server and discard the response
		_, err := client.Do(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTelerBadCrawlerOnly(b *testing.B) {
	// Initialize teler
	telerMiddleware := New(Options{
		Excludes: []threat.Threat{
			threat.CommonWebAttack,
			threat.CVE,
			threat.BadIPAddress,
			threat.BadReferrer,
			threat.DirectoryBruteforce,
		},
		NoStderr: true,
	})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		b.Fatal(err)
	}

	// Set the user agent to "X"
	req.Header.Set("User-Agent", "X")

	// Run the benchmark
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Send the request to the test server and discard the response
		_, err := client.Do(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTelerDirectoryBruteforceOnly(b *testing.B) {
	// Initialize teler
	telerMiddleware := New(Options{
		Excludes: []threat.Threat{
			threat.CommonWebAttack,
			threat.CVE,
			threat.BadIPAddress,
			threat.BadReferrer,
			threat.BadCrawler,
		},
		NoStderr: true,
	})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		b.Fatal(err)
	}

	// Set the user agent to "X"
	req.Header.Set("User-Agent", "X")

	// Run the benchmark
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Send the request to the test server and discard the response
		_, err := client.Do(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTelerCustomRule(b *testing.B) {
	// Initialize teler
	telerMiddleware := New(Options{
		Excludes: []threat.Threat{
			threat.CommonWebAttack,
			threat.CVE,
			threat.BadIPAddress,
			threat.BadReferrer,
			threat.BadCrawler,
			threat.DirectoryBruteforce,
		},
		Customs: []Rule{
			{
				Name:      "Log4j Attack",
				Condition: "or",
				Rules: []Condition{
					{
						Method:  request.GET,
						Element: request.URI,
						Pattern: `\$\{.*:\/\/.*\/?\w+?\}`,
					},
				},
			},
		},
		NoStderr: true,
	})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		b.Fatal(err)
	}

	// Run the benchmark
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Send the request to the test server and discard the response
		_, err := client.Do(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTelerWithoutCommonWebAttack(b *testing.B) {
	// Initialize teler
	telerMiddleware := New(Options{
		Excludes: []threat.Threat{
			threat.CommonWebAttack,
		},
		NoStderr: true,
	})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		b.Fatal(err)
	}

	// Set the user agent to "X"
	req.Header.Set("User-Agent", "X")

	// Run the benchmark
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Send the request to the test server and discard the response
		_, err := client.Do(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTelerWithoutCVE(b *testing.B) {
	// Initialize teler
	telerMiddleware := New(Options{
		Excludes: []threat.Threat{
			threat.CVE,
		},
		NoStderr: true,
	})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		b.Fatal(err)
	}

	// Set the user agent to "X"
	req.Header.Set("User-Agent", "X")

	// Run the benchmark
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Send the request to the test server and discard the response
		_, err := client.Do(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTelerWithoutBadIPAddress(b *testing.B) {
	// Initialize teler
	telerMiddleware := New(Options{
		Excludes: []threat.Threat{
			threat.BadIPAddress,
		},
		NoStderr: true,
	})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		b.Fatal(err)
	}

	// Set the user agent to "X"
	req.Header.Set("User-Agent", "X")

	// Run the benchmark
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Send the request to the test server and discard the response
		_, err := client.Do(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTelerWithoutBadReferrer(b *testing.B) {
	// Initialize teler
	telerMiddleware := New(Options{
		Excludes: []threat.Threat{
			threat.BadReferrer,
		},
		NoStderr: true,
	})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		b.Fatal(err)
	}

	// Set the user agent to "X"
	req.Header.Set("User-Agent", "X")

	// Run the benchmark
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Send the request to the test server and discard the response
		_, err := client.Do(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTelerWithoutBadCrawler(b *testing.B) {
	// Initialize teler
	telerMiddleware := New(Options{
		Excludes: []threat.Threat{
			threat.BadCrawler,
		},
		NoStderr: true,
	})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		b.Fatal(err)
	}

	// Run the benchmark
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Send the request to the test server and discard the response
		_, err := client.Do(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTelerWithoutDirectoryBruteforce(b *testing.B) {
	// Initialize teler
	telerMiddleware := New(Options{
		Excludes: []threat.Threat{
			threat.DirectoryBruteforce,
		},
		NoStderr: true,
	})
	wrappedHandler := telerMiddleware.Handler(handler)

	// Create a test server with the wrapped handler
	ts := httptest.NewServer(wrappedHandler)
	defer ts.Close()

	// Create a client to send requests to the test server
	client := &http.Client{}

	// Create a request to send to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		b.Fatal(err)
	}

	// Set the user agent to "X"
	req.Header.Set("User-Agent", "X")

	// Run the benchmark
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Send the request to the test server and discard the response
		_, err := client.Do(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func ExampleNew_default() {
	var myHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello world"))
	})

	telerMiddleware := New()

	app := telerMiddleware.Handler(myHandler)
	go func() {
		_ = http.ListenAndServe("127.0.0.1:3000", app)
	}()
}

func ExampleNew_setHandler() {
	var forbidden = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "We're sorry, but your request has been denied for security reasons.", http.StatusForbidden)
	})

	var myHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello world"))
	})

	telerMiddleware := New()
	telerMiddleware.SetHandler(forbidden)

	app := telerMiddleware.Handler(myHandler)
	go func() {
		_ = http.ListenAndServe("127.0.0.1:3000", app)
	}()
}

func ExampleNew_custom() {
	var myHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello world"))
	})

	telerMiddleware := New(Options{
		Excludes: []threat.Threat{
			threat.BadReferrer,
			threat.BadCrawler,
		},
		Whitelists: []string{
			`(curl|Go-http-client|okhttp)/*`,
			`^/wp-login\.php`,
			`(?i)Referer: https?:\/\/www\.facebook\.com`,
			`192\.168\.0\.1`,
		},
		Customs: []Rule{
			{
				Name:      "Log4j Attack",
				Condition: "or",
				Rules: []Condition{
					{
						Method: request.GET,
						// if Method is not set or invalid, defaulting to request.GET.
						Element: request.URI,
						// you can use request.Any: it useful when you want to
						// match against multiple elements of the request at once,
						// rather than just a single element.
						Pattern: `\$\{.*:\/\/.*\/?\w+?\}`,
					},
				},
			},
		},
		LogFile: "/tmp/teler.log",
	})

	app := telerMiddleware.Handler(myHandler)
	go func() {
		_ = http.ListenAndServe("127.0.0.1:3000", app)
	}()
}
