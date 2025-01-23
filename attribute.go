package artwork

import "sort"

type Attribute struct {
	Name
}

type AttributeWeightMap map[Attribute]float64

// @TODO: Rework the WeightMap functions, either make them more generic, or make them based on an attribute type that is fairly generic/useful.

// Sum returns the sum of the weight map.
func (a AttributeWeightMap) Sum() (sum float64) {
	for _, w := range a {
		sum += w
	}
	return
}

/* Computes an array of intervals which represent the normalized, weighted distribution of an AttributeMap. The sum of the receiver, "a" should equal 1.0.
Here's a visual:

[]AttributeWeightMap{aV:0.2, aW:0.15, aX:0.25, aY:0.1, aZ:0.3}

Would become something like,

[]AttributeWeightIntervalMap{aV:0.2, aW:0.35, aX:0.60, aY:0.7, aZ:1.0}

Order is irrelevant.

This function returns the distribution as an *AttributeWeightInterval slice
*/
func (a AttributeWeightMap) Intervals() AttributeWeightIntervals {
	// Make sure we have more than one weight.
	if len(a) <= 1 {
		// @TODO: Maybe this should return nil?
		//return nil
	}
	// Make the CDF slice
	prevSum := 0.0
	i, ints := 0, make(AttributeWeightIntervals, len(a))
	for attr, w := range a {
		ints[i] = &AttributeWeightInterval{
			Weight:    w + prevSum,
			Attribute: attr,
		}
		prevSum += w
		i++
	}
	// Sort the CDF slice
	sort.Slice(ints, func(i, j int) bool {
		return ints[i].Weight < ints[j].Weight
	},
	)
	return ints
}

type AttributeWeightInterval struct {
	Attribute
	Weight float64
}

type AttributeWeightIntervals []*AttributeWeightInterval

// Attribute returns the first attribute for which f is less than its computed distribution threshold. It assumes itself to be sorted.
func (a AttributeWeightIntervals) Attribute(f float64) Attribute {
	for _, awi := range a {
		// Since slice is sorted, this is all we should need.
		if f < awi.Weight {
			return awi.Attribute
		}
	}
	return 0
}
