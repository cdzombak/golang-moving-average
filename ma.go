package movingaverage

import (
	"math"

	"github.com/montanaflynn/stats"
)

// MovingStats holds the most recently added N values (N = Options.Window)
// and provides a bridge from those data to the github.com/montanaflynn/stats
// Float64Data type for statistical calculations.
//
// Avg() and Median() are the only statistical methods provided by this package;
// they wrap the relevant functions from github.com/montanaflynn/stats but return
// 0.0 if an error occurs.
//
// (Avg() is the (non-geometric) mean of the values, and Median() is the median.)
type MovingStats interface {
	// Add adds the given values to the moving stats instance.
	Add(values ...float64)

	// Window returns the number of values kept in the moving stats instance.
	Window() int

	// SlotsFilled returns whether all slots in the moving stats instance have been filled.
	SlotsFilled() bool

	// Values returns the values in the moving stats instance, as stats.Float64Data.
	Values() stats.Float64Data

	// Count returns the number of values in the moving stats instance.
	Count() int

	// Avg returns the average of the values in the moving stats instance.
	// If no values have been added or any other error occurs, 0.0 is returned.
	Avg() float64

	// Median returns the median of the values in the moving stats instance.
	// If no values have been added or any other error occurs, 0.0 is returned.
	Median() float64

	// Min returns the minimum of the values in the moving stats instance.
	// If no values have been added or any other error occurs, 0.0 is returned.
	Min() float64

	// Max returns the maximum of the values in the moving stats instance.
	// If no values have been added or any other error occurs, 0.0 is returned.
	Max() float64

	// UnsafeDoStat runs the given function on the values in the moving stats instance.
	// If the function returns an error, that error is returned.
	// Functions passed to UnsafeDoStat must not modify the values slice or call Add(). This will result in undefined behavior.
	UnsafeDoStat(func(stats.Float64Data) (float64, error)) (float64, error)

	// UnsafeDo runs the given function on the values in the moving stats instance.
	// If the function returns an error, that error is returned.
	// Functions passed to UnsafeDo must not modify the values slice or call Add(). This will result in undefined behavior.
	UnsafeDo(func(stats.Float64Data) error) error
}

// Options configures a new movingStats instance.
type Options struct {
	// Whether to ignore NaN values when adding values to the moving stats instance.
	IgnoreNanValues bool

	// Whether to ignore Inf values when adding values to the moving stats instance.
	IgnoreInfValues bool

	// The number of values to keep in the moving stats instance.
	Window int
}

// New returns a new MovingStats instance with the given options.
func New(opts Options) MovingStats {
	return &movingStats{
		values:          make([]float64, opts.Window),
		valPos:          0,
		slotsFilled:     false,
		window:          opts.Window,
		ignoreInfValues: opts.IgnoreInfValues,
		ignoreNanValues: opts.IgnoreNanValues,
	}
}

type movingStats struct {
	window          int
	values          []float64
	valPos          int
	slotsFilled     bool
	ignoreNanValues bool
	ignoreInfValues bool
}

func (ma *movingStats) filledValues() stats.Float64Data {
	var c = ma.window - 1

	// Are all slots filled? If not, ignore unused
	if !ma.slotsFilled {
		c = ma.valPos - 1
		if c < 0 {
			// Empty register
			return nil
		}
	}
	return ma.values[0 : c+1]
}

func (ma *movingStats) Add(values ...float64) {
	for _, val := range values {
		// ignore NaN?
		if ma.ignoreNanValues && math.IsNaN(val) {
			continue
		}

		// ignore Inf?
		if ma.ignoreInfValues && math.IsInf(val, 0) {
			continue
		}

		// Put into values array
		ma.values[ma.valPos] = val

		// Increment value position
		ma.valPos = (ma.valPos + 1) % ma.window

		// Did we just go back to 0, effectively meaning we filled all registers?
		if !ma.slotsFilled && ma.valPos == 0 {
			ma.slotsFilled = true
		}
	}
}

func (ma *movingStats) Window() int {
	return ma.window
}

func (ma *movingStats) SlotsFilled() bool {
	return ma.slotsFilled
}

func (ma *movingStats) Values() stats.Float64Data {
	internal := ma.filledValues()
	retv := make(stats.Float64Data, len(internal))
	_ = copy(retv, internal)
	return retv

}

func (ma *movingStats) Count() int {
	return len(ma.filledValues())
}

func (ma *movingStats) Avg() float64 {
	retv, err := ma.filledValues().Mean()
	if err != nil {
		return 0.0
	}
	return retv
}

func (ma *movingStats) Median() float64 {
	retv, err := ma.filledValues().Median()
	if err != nil {
		return 0.0
	}
	return retv
}

func (ma *movingStats) Min() float64 {
	retv, err := ma.filledValues().Min()
	if err != nil {
		return 0.0
	}
	return retv
}

func (ma *movingStats) Max() float64 {
	retv, err := ma.filledValues().Max()
	if err != nil {
		return 0.0
	}
	return retv
}

func (ma *movingStats) UnsafeDoStat(f func(stats.Float64Data) (float64, error)) (float64, error) {
	return f(ma.filledValues())
}

func (ma *movingStats) UnsafeDo(f func(stats.Float64Data) error) error {
	return f(ma.filledValues())
}
