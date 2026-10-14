package cmd

import(
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"gmc/config"
	"gmc/db"
	"gmc/db/model"
)

var MAX_REDIRECTS int = 5

func checkURL(address string, client *http.Client, ctx context.Context) int {
	u, err := url.Parse(address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n\t%s\n", address, err)
		return 1
	}

	for i:=0; i <= MAX_REDIRECTS; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodHead, u.String(), nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n\t%s\n", address, err)
			return 1
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n\t%s\n", address, err)
			return 1
		} 
		defer resp.Body.Close()
		
		if (resp.StatusCode >= 200) && (resp.StatusCode <= 299) {
			return 0

		} else if (resp.StatusCode >= 300) && (resp.StatusCode <= 399) {
			if i == MAX_REDIRECTS {
				fmt.Fprintf(os.Stderr, "%s\n\t%d : %s\n", address, resp.StatusCode, http.StatusText(resp.StatusCode))
				return 1
			}
			
			redirect := resp.Header.Get("Location")
			new_url, err := url.Parse(redirect)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n\t%s\n", address, err)
				return 1
			}

			host := new_url.Hostname()
			if host == "" {
				new_url.Host = u.Hostname()
			}

			if new_url.IsAbs() == false {
				new_url.Scheme = u.Scheme
			}

			u = new_url

		} else if resp.StatusCode >= 400 {
			fmt.Fprintf(os.Stderr, "%s\n\t%d : %s\n", address, resp.StatusCode, http.StatusText(resp.StatusCode))
			return 1
		}
	}
	return 1
}

func HealthURL(exec string, cfg *config.Config) int {
	cache := map[string][]net.IPAddr {}
	transport := &http.Transport{
		IdleConnTimeout: 15 * time.Second,
		MaxIdleConns: 150,
		MaxIdleConnsPerHost: 10,
		TLSClientConfig : &tls.Config{
			InsecureSkipVerify: true,
		},
		DialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}

			ips, cached := cache[host]
			if cached == false {
				ips, err = net.DefaultResolver.LookupIPAddr(ctx, host)
				if err != nil {
					return nil, err
				}
				cache[host] = ips
			}

			target := net.JoinHostPort(ips[0].String(), port)

			return (&net.Dialer{}).DialContext(ctx, network, target)
		},
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 3 * time.Second,
		Transport: transport,
	}

	
	db, err := db.New(cfg.Database)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
		return 1
	}
	
	err = db.ListURLs(func(f *model.URL) error {
		checkURL(*f.URL, client, context.Background())
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", exec, err)
		return 1
	}

	return 0
}