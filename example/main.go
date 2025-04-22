package main

import (
	"fmt"
	"os"
	"time"

	api2captcha "github.com/2captcha/2captcha-go"
	"github.com/Ahmed-Labs/captcha-pool"
	"github.com/joho/godotenv"
)

var pool *captchapool.CaptchaPool

func solve() (string, error) {
	startTime := time.Now()
	client := api2captcha.NewClient(os.Getenv("2CAPTCHA_KEY"))

	cap := api2captcha.ReCaptcha{
		SiteKey:   "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
		Url:       "https://www.google.com/recaptcha/api2/demo",
		Invisible: true,
		Action:    "verify",
	}

	req := cap.ToRequest()
	token, _, err := client.Solve(req)
	fmt.Printf("captcha request took %.4f seconds\n", time.Since(startTime).Seconds())
	return token, err
}

func runTask(solved chan string) {
	time.Sleep(time.Second * 60)
	token := pool.GetToken()
	solved <- token
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("could not load env")
	}

	numTasks := 10
	pool = captchapool.New(solve,
		&captchapool.Options{
			Count:           numTasks,
			Refresh:         true,
			RefreshInterval: time.Second * 30,
			TTL:             time.Second * 60,
		},
	)

	solved := make(chan string)
	numSolved := 0

	startTime := time.Now()
	pool.Start()

	defer pool.Stop()
	for range numTasks {
		go runTask(solved)
	}

	for {
		token := <-solved
		numSolved++
		fmt.Printf("solved token (%d): %s...\n", numSolved, token[:30])
		if numSolved == numTasks {
			fmt.Printf("finished solving %d captchas in %.4f seconds\n", numSolved, time.Since(startTime).Seconds())
			return
		}
	}
}
