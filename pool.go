package captchapool

import (
	"context"
	"sync"
	"time"

	"github.com/gammazero/deque"
)

type Options struct {
	// Number of captcha tokens to generate at a time
	Count int
	// Refresh allows regenerating captcha pool at a given interval
	Refresh bool
	// Interval duration for generating new captchas to the pool
	RefreshInterval time.Duration
	// Duration of captcha validity, i.e. time to live
	TTL time.Duration
}

type Captcha struct {
	created time.Time
	ttl     time.Duration
	token   string
}

type CaptchaPool struct {
	mu              sync.Mutex
	cond            *sync.Cond
	ctx             context.Context
	pool            *deque.Deque[Captcha]
	solve           func() (string, error)
	count           int
	refresh         bool
	refreshInterval time.Duration
	ttl             time.Duration
	maxRetry        int
}

// Creates a new captcha pool with given options
// Solve is a captcha solver that returns a captcha token string and an error
func New(solve func() (string, error), options *Options) *CaptchaPool {
	deque := new(deque.Deque[Captcha])
	deque.SetBaseCap(options.Count)
	
	newPool := &CaptchaPool{
		ctx:             context.Background(),
		pool:            deque,
		solve:           solve,
		count:           options.Count,
		refresh:         options.Refresh,
		refreshInterval: options.RefreshInterval,
		ttl:             options.TTL,
		maxRetry:        3,
	}
	newPool.cond = sync.NewCond(&newPool.mu)
	return newPool
}

// Start solving captchas with specified configuration
// Solved captchas are added to the pool
func (c *CaptchaPool) Start() {
	go func() {
		for range c.count {
			go c.solveCaptcha()
		}
		if !c.refresh {
			return
		}
		ticker := time.NewTicker(c.refreshInterval)
		defer ticker.Stop()

		for {
			select {
			case <-c.ctx.Done():
				return
			case <-ticker.C:
				for range c.count {
					go c.solveCaptcha()
				}
			}
		}
	}()
}

// Stop captcha pool refresh execution
func (c *CaptchaPool) Stop() {
	if !c.refresh {
		return
	}
	c.ctx.Done()
}

// Get a solved captcha from the pool
func (c *CaptchaPool) GetToken() string {
	return c.pop().token
}

// Runs captcha solver and pushes token to pool if successful
func (c *CaptchaPool) solveCaptcha() {
	for range c.maxRetry {
		token, err := c.solve()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		c.push(Captcha{
			created: time.Now(),
			ttl:     c.ttl,
			token:   token,
		})
		return
	}
}

// Internal function to safely push a new captcha to pool
func (c *CaptchaPool) push(captcha Captcha) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.pool.Len() == c.count {
		c.pool.PopFront()
	}
	c.pool.PushBack(captcha)
	c.cond.Signal()
}

// Internal function to safely pop captcha from pool
// Blocks until pool non-empty and unexpired captcha can be retrieved
func (c *CaptchaPool) pop() Captcha {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Remove expired captchas from pool
	for c.pool.Len() > 0 {
		cap := c.pool.Front()
		if time.Now().After(cap.created.Add(cap.ttl)) {
			c.pool.PopFront()
		} else {
			break
		}
	}
	if c.pool.Len() == 0 {
		c.cond.Wait()
	}
	return c.pool.PopFront()
}
