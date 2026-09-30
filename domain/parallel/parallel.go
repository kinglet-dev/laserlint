// Package parallel spreads independent work over all processors. Callers
// write each result to its own slot, so results don't depend on how the
// work is split and runs stay repeatable.
package parallel

import (
	"runtime"
	"sync"
)

// For calls f(i) once for every i from 0 to n-1, on all processors, and
// returns when all calls are done. Indexes are dealt out in turn, since
// neighbouring indexes often cost about the same (dense areas cluster).
func For(n int, f func(i int)) {
	workers := min(runtime.GOMAXPROCS(0), n)
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := w; i < n; i += workers {
				f(i)
			}
		}()
	}
	wg.Wait()
}
