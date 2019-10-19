package rules

import "context"

func GenerateShipSubsets(ctx context.Context, max ShipCount, buckets int) <-chan ShipCount {
	resChan := make(chan ShipCount)
	go func() {
		defer close(resChan)
		buckets0 := buckets
		if max[0]+1 < buckets0 {
			buckets0 = max[0] + 1
		}
		buckets1 := buckets
		if max[1]+1 < buckets0 {
			buckets0 = max[1] + 1
		}
		buckets2 := buckets
		if max[2]+1 < buckets0 {
			buckets0 = max[2] + 1
		}
		var res ShipCount
		for i0 := 0; i0 <= buckets0; i0++ {
			res[0] = (i0 * max[0]) / buckets0
			for i1 := 0; i1 <= buckets1; i1++ {
				res[1] = (i1 * max[1]) / buckets1
				for i2 := 0; i2 <= buckets2; i2++ {
					res[2] = (i2 * max[2]) / buckets2
					select {
					case resChan <- res:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()
	return resChan
}
