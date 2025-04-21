package captchapool

import (
	"context"
	"time"
)

// CaptchaPool should be a queue essentially
// A -> feeding that queue with captcha tokens
// B -> Reading from queue and returning

type CaptchaPoolOptions struct {
	// Number of captcha tokens to generate
	Count int
	// Perpetual refreshes captcha pool at given refresh rate
	Perpetual bool
	// Refresh is an interval duration for refreshing the pool with new captchas
	RefreshRate time.Duration
	// Duration of captcha validity, i.e. time to live
	TTL time.Duration
}

type Captcha struct {
	created time.Time
	ttl time.Duration
	token string
}

type CaptchaPool struct {
	ctx context.Context
	pool []Captcha
	solve func() string
	count int
	perpetual bool
	refreshRate time.Duration
	ttl time.Duration
}

// Creates a new captcha pool with given options
// Solve is a blocking captcha solver that returns a captcha token string
func New(solve func() string, options *CaptchaPoolOptions) *CaptchaPool {
	return &CaptchaPool{}
}

// Start solving captchas with specified configuration
// Solved captchas are added to captcha pool
func (c *CaptchaPool) Start() {
	// If not perpetual, get n tokens once and add them to pool
	// go routine for captcha pool to run in the background
}

// Stop captcha pool execution
// Use context to cancel pool execution
func (c *CaptchaPool) Stop() {
	
}

// Internal function for pushing solved captcha to pool
func (c *CaptchaPool) push(captcha Captcha) {

}

// Internal function for popping solved captcha from pool
// Pop from queue but only if captcha valid
// 		->  Blocks until a valid unexpired captcha is retrieved from the
func (c *CaptchaPool) pop() Captcha {
	return Captcha{}
}

// Get a solved captcha from the pool
func (c *CaptchaPool) Get() string {
	return c.pop().token
}
