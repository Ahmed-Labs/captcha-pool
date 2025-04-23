package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Proxy struct {
	Host     string // IP or domain
	Port     string
	Username string
	Password string
}

type ProxyList struct {
	curr int
	proxies []Proxy
	mu sync.Mutex
}

func (p Proxy) String() string {
	return fmt.Sprintf("%s:%s@%s:%s", p.Username, p.Password, p.Host, p.Port)
}

func (p *ProxyList) GetProxy() Proxy {
	proxies.mu.Lock()
	defer proxies.mu.Unlock()
	proxy := p.proxies[p.curr]
	p.curr = (p.curr + 1)%len(p.proxies)
	return proxy
}

func readProxies(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			proxies = append(proxies, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return proxies, nil
}

func parseProxy(raw string) (Proxy, error) {
	parts := strings.Split(raw, ":")
	if len(parts) < 4 {
		return Proxy{}, fmt.Errorf("invalid proxy format: %s", raw)
	}

	return Proxy{
		Host:     parts[0],
		Port:     parts[1],
		Username: parts[2],
		Password: parts[3],
	}, nil
}

func LoadProxies(path string) (*ProxyList, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var proxies []Proxy
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		proxy, err := parseProxy(line)
		if err != nil {
			return nil, err // or log and continue
		}
		proxies = append(proxies, proxy)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &ProxyList{proxies: proxies}, nil
}
