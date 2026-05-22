package pubnub

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func trustTestServerCert(existing *tls.Config, cert *x509.Certificate) *tls.Config {
	var cfg *tls.Config
	if existing != nil {
		cfg = existing.Clone()
	} else {
		cfg = &tls.Config{MinVersion: tls.VersionTLS12}
	}
	if cfg.RootCAs == nil {
		cfg.RootCAs = x509.NewCertPool()
	}
	cfg.RootCAs.AddCert(cert)
	cfg.MinVersion = tls.VersionTLS12
	return cfg
}

func TestNewHTTP2Client_plainHTTPUsesHTTP11(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := NewHTTP2Client(5, 10)
	resp, err := client.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	if got := resp.Proto; got != "HTTP/1.1" {
		t.Fatalf("plaintext http: Proto = %q, want HTTP/1.1", got)
	}
}

func TestNewHTTP1Client_httpsUsesHTTP11(t *testing.T) {
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	srv.StartTLS()
	defer srv.Close()

	client := NewHTTP1Client(5, 10, 10)
	tr := client.Transport.(*http.Transport)
	tr.TLSClientConfig = trustTestServerCert(nil, srv.Certificate())

	resp, err := client.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	if got := resp.Proto; got != "HTTP/1.1" {
		t.Fatalf("Proto = %q, want HTTP/1.1", got)
	}
}

func TestNewHTTP2Client_httpsUsesHTTP2WhenServerSupportsHTTP2(t *testing.T) {
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	srv.EnableHTTP2 = true
	srv.StartTLS()
	defer srv.Close()

	client := NewHTTP2Client(5, 10)
	tr := client.Transport.(*http.Transport)
	tr.TLSClientConfig = trustTestServerCert(tr.TLSClientConfig, srv.Certificate())

	resp, err := client.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	if got := resp.Proto; got != "HTTP/2.0" {
		t.Fatalf("Proto = %q, want HTTP/2.0", got)
	}
}

func TestNewHTTP2Client_httpsUsesHTTP11WhenServerDoesNotAdvertiseHTTP2(t *testing.T) {
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	srv.EnableHTTP2 = false
	srv.StartTLS()
	defer srv.Close()

	client := NewHTTP2Client(5, 10)
	tr := client.Transport.(*http.Transport)
	tr.TLSClientConfig = trustTestServerCert(tr.TLSClientConfig, srv.Certificate())

	resp, err := client.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	if got := resp.Proto; got != "HTTP/1.1" {
		t.Fatalf("fallback Proto = %q, want HTTP/1.1", got)
	}
}
