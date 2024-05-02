# `github.com/cdzombak/golang-moving-average`

Moving average/median/stats implementation for Go. This project provides a moving window of `n` values, for which you can calculate the mean and median, or use any of the statistical functions provided by the excellent & well-tested [github.com/montanaflynn/stats](montanaflynn/stats) library.

This project's repo and package names come from [RobinUS2/golang-moving-average](https://github.com/RobinUS2/golang-moving-average), from which it was originally forked. @cdzombak's fork:

- adds a moving-median function
- integrates with [montanaflynn/stats](https://github.com/montanaflynn/stats)
- improves the concurrency-safe wrapper to prevent accidental misuse

**Documentation:** [pkg.go.dev/github.com/cdzombak/golang-moving-average](https://pkg.go.dev/github.com/cdzombak/golang-moving-average)

## Installation

```shell
go get github.com/cdzombak/golang-moving-average
```

## Usage

Create a `MovingStats` instance via `movingaverage.New()`, then add values to it via its `Add()` method:

```go
package main

import "github.com/cdzombak/golang-moving-average"

func main() {
	ms := movingaverage.New(movingaverage.Options{Window: 4})
	ms.Add(10)
	ms.Add(2)
	ms.Add(4)
	ms.Add(6)
	ms.Add(8) // This effectively overwrites the first value (10 in this example)

	avg := ms.Avg() // 5.0
}
```

### Basic stats

The `MovingStats` interface provides four statistical calculations directly: `Avg()`, `Median()`, `Min()`, and `Max()`. These call through to the relevant functions from [montanaflynn/stats](https://github.com/montanaflynn/stats).

If an error occurs (i.e. no values have been added yet), they return `0.0` (the `float64` zero value).

> [!TIP]
> For `Avg()` and `Median()`, this (more Golang-idiomatic) API provides the same behavior as the `Avg()` function in [RobinUS2/golang-moving-average](https://github.com/RobinUS2/golang-moving-average), from which this project was forked.
> 
> However, that project returned errors for `Min()` and `Max()`. For API consistency, this project returns `0.0` in those cases.
>
> If you prefer the [montanaflynn/stats](https://github.com/montanaflynn/stats) APIs' behavior, you can use its functions instead of these convenience wrappers, via the methods described in "Extended stats," below.

### Extended stats

To use statistical functions from [montanaflynn/stats](https://github.com/montanaflynn/stats) or implement entirely custom ones, read the current values from the `MovingStats` instance.

`Values()` returns the values currently stored in the moving stats instance. You can pass this slice to any of the functions in [montanaflynn/stats](https://github.com/montanaflynn/stats) or call its methods on the slice directly:

```go
package main

import "github.com/cdzombak/golang-moving-average"

func main() {
	ms := movingaverage.New(movingaverage.Options{Window: 3})
	ms.Add(1)
	ms.Add(2)
	ms.Add(3)
	ms.Add(2)

	mode, err := ms.Values().Mode() // [2], nil
}
```

#### Performance considerations

`Values()` returns a copy of the values in the `MovingStats` instance. If there are a large number of values and/or you're calling it extremely frequently, this could be a bottleneck.

To avoid this, you can use the `UnsafeDoStat()` and `UnsafeDo()` methods. These methods allow running a function that receives the values slice directly, without copying it.

> [!IMPORTANT]
> Functions passed to `UnsafeDoStat` or `UnsafeDo` **must not modify the values slice or call `Add()`**. This will result in undefined behavior.

Example:

```go
package main

import (
	"github.com/cdzombak/golang-moving-average"
	"github.com/montanaflynn/stats"
)

func main() {
	ms := movingaverage.New(movingaverage.Options{Window: 3})
	ms.Add(1)
	ms.Add(2)
	ms.Add(3)
	ms.Add(4)

	mean, _ := ms.UnsafeDoStat(stats.Mean) // 3.0
}
```

### Concurrency

`MovingStats` instances created by `movingaverage.New()` are not safe for concurrent use by multiple goroutines.

To create a concurrency-safe `MovingStats` instance, use `movingaverage.NewConcurrent()`. This function accepts the same `Options` as `New()`.

> [!IMPORTANT]
> Functions passed to `UnsafeDoStat` or `UnsafeDo` **must not call `Add()`**. This will cause a deadlock.

### Other methods

Additional methods are available for inspecting the `MovingStats` interface:

```go
// Window returns the number of values kept in the moving stats instance.
Window() int

// SlotsFilled returns whether all slots in the moving stats instance have been filled.
SlotsFilled() bool

// Count returns the number of values in the moving stats instance.
Count() int
```

### Partially used windows

If you create a `MovingStats` instance and `Add` fewer values than its `Window` size, stats will be calculated only on the values you've added.

Meaning, for example, if you create an instance with `Window = 5` and add 2 values, only those 2 values will be included in calculations; the window is not "padded" with zeroes.

```go
package main

import "github.com/cdzombak/golang-moving-average"

func main() {
	ms := movingaverage.New(movingaverage.Options{Window: 5})
	ms.Add(10)
	ms.Add(20)

	avg := ms.Avg() // 15.0, not 6.0

	window := ms.Window() // 5
	filled := ms.SlotsFilled() // false
	count := ms.Count() // 2
}
```

## License

Apache 2.0; see [LICENSE](LICENSE) in this repo.

## Author

Originally based on & forked from [RobinUS2/golang-moving-average](https://github.com/RobinUS2/golang-moving-average).

Modifications (as described in the intro of this README) by Chris Dzombak ([dzombak.com](https://dzombak.com); [github @cdzombak](https://github.com/cdzombak)).
