package movingaverage

import (
	"sync"

	"github.com/montanaflynn/stats"
)

type concurrentMovingStats struct {
	ma  MovingStats
	mux sync.RWMutex
}

// NewConcurrent returns a new concurrency-safe MovingStats instance
// with the given options.
func NewConcurrent(opts Options) MovingStats {
	return &concurrentMovingStats{
		ma: New(opts),
	}
}

func (c *concurrentMovingStats) Add(values ...float64) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.ma.Add(values...)
}

func (c *concurrentMovingStats) Window() int {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.ma.Window()
}

func (c *concurrentMovingStats) SlotsFilled() bool {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.ma.SlotsFilled()
}

func (c *concurrentMovingStats) Values() stats.Float64Data {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.ma.Values()
}

func (c *concurrentMovingStats) Count() int {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.ma.Count()
}

func (c *concurrentMovingStats) Avg() float64 {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.ma.Avg()
}

func (c *concurrentMovingStats) Median() float64 {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.ma.Median()
}

func (c *concurrentMovingStats) Min() float64 {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.ma.Min()
}

func (c *concurrentMovingStats) Max() float64 {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.ma.Max()
}

func (c *concurrentMovingStats) UnsafeDoStat(f func(stats.Float64Data) (float64, error)) (float64, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.ma.UnsafeDoStat(f)
}

func (c *concurrentMovingStats) UnsafeDo(f func(stats.Float64Data) error) error {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.ma.UnsafeDo(f)
}
