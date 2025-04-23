package main

import (
	"fmt"
	"os"
	"time"

	api2captcha "github.com/2captcha/2captcha-go"
	captchapool "github.com/Ahmed-Labs/captcha-pool"
	"github.com/joho/godotenv"
)

var pool *captchapool.CaptchaPool
var proxies *ProxyList

func buildSolver(params ...string) func() (string, error) {
	// Set up anything needed initially / called once
	client := api2captcha.NewClient(os.Getenv("2CAPTCHA_KEY"))

	// Solve
	return func() (string, error) {
		startTime := time.Now()
		defer func() {
			fmt.Printf("captcha request took %.4f seconds\n", time.Since(startTime).Seconds())
		}()

		cap := api2captcha.ReCaptcha{
			SiteKey:   "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
			Url:       "https://www.google.com/recaptcha/api2/demo",
			Invisible: true,
			Action:    "verify",
		}
		req := cap.ToRequest()
		proxy := proxies.GetProxy().String()
		req.SetProxy("HTTPS", proxy)

		token, _, err := client.Solve(req)
		return token, err
	}
}

func runTask(solved chan string) {
	time.Sleep(time.Second * 60)
	token := pool.GetToken()
	solved <- token
}

func main() {
	// Load env
	err := godotenv.Load()
	if err != nil {
		panic("could not load env")
	}

	// Load proxies
	proxies, err = LoadProxies("./example/proxies.txt")
	if err != nil {
		panic(err)
	}

	// Configure captcha pool
	numTasks := 5
	pool = captchapool.New(buildSolver(),
		&captchapool.Options{
			Count:           numTasks,
			Refresh:         true,
			RefreshInterval: time.Second * 40,
			TTL:             time.Second * 60,
		},
	)

	// Start the solver and use for tasks
	solved := make(chan string)
	startTime := time.Now()
	pool.Start()
	defer pool.Stop()

	for range numTasks {
		go runTask(solved)
	}

	for range numTasks {
		token := <-solved
		fmt.Printf("used captcha: %s...\n", token[:30])
	}

	fmt.Printf("finished solving %d captchas in %.4f seconds\n", numTasks, time.Since(startTime).Seconds())
}
