#go #concurrency #article

Before talking about goroutines, channels, or worker pools, we must ask a basic question:

> **Can this problem be solved serially in a simple, obvious way?**

Take the counting problem.

`Count from 1 to 1,000,000`

A single loop does this perfectly fine.

- Easy to read
- Easy to reason about
- Zero coordination cost
- No overhead

If your program only needs to count **10,000 numbers**, introducing concurrency would be like hiring **8 people to count 10 pages**—you’d spend more time coordinating than counting.

> This way of thinking applies not just to abstract numbers, but to how humans naturally organize work at scale.

---
#### When Concurrency Becomes a Question

Now change the problem slightly:

`Count from 1 to 1,000,000,000`

Suddenly, the serial solution has a cost:
- The CPU is busy for a long time
- Other work might be blocked
- The program becomes _unavailable_ for extended periods

This is the key transition point.

> **We don’t introduce concurrency because the task is big.  
> We introduce concurrency because the serial execution makes something unavailable.**

I realized that the problem of counting currency notes maps extremely well to how concurrency in Go — and even the Go scheduler — is designed.

If u think about the currency notes counting problem, if we had just a few number of bundles, we need not think much and having just one person and a counting table should be good enough. 

But if we are supplied with a huge number of currency note bundles — say, suitcases full of them — counting alone no longer makes sense, even if the person is perfectly capable.

So, here, a key thing to understand would be unless the scale of the problem is large enough, we need not jump into thinking along the lines of concurrency

## Counting Hall Story
### Deciding on the Cast
- The hall has a limited number of counting tables (say 8).
- The hall receives many currency note bundles of different denominations.
- There are many people who know how to count notes — **workers**.
- **Condition**: a worker can count only when they acquire both a table and a bundle.
- Sometimes, we deliberately assign **more workers than tables**.  
	- This ensures that tables never sit idle — if someone finishes early or pauses, another worker can immediately take the chair.
- Even if a worker has been waiting near one table, they may notice another table becoming free and rush to occupy it.
- In addition, we introduce two special roles:
	- **Producers**, who distribute bundles to workers
	- **Accumulators**, who collect and aggregate results
 
### Deciding on the implementation
- After the cast of the Counting Hall is Ready, I moved ahead thinking of the implementation.
- Thinking of it's implementation, think of it like there are two pipes
	- job pipe
	- result pipe
- This separation ensures unidirectional flow and clear ownership of responsibilities.
- Every worker is connected to 
	- **receiving end** of the job pipe
		- Every worker, when acquires one of the table, will be able to **receive** a currency bundle from this job pipe
	- also the **sending end** of the result pipe
		- Every worker, on acquiring the table and receiving a bundle, starts calculating notes, and after finishing, prepares a chit containing its results and **sends** that result chit in the result pipe, for it to reach the accumulator
- The producer is connected to sending end of the job pipe
- The accumulator is connected to the receiving end of the result pipe

## Co-relating with Go
### The corresponding cast:
If I have to co-relate this with go-implementation, the cast would be
- the machine has limited cores/CPUs (say 8)
- we have a bunch of note bundles of different denominations
- each **worker** go-routine we launch (in **Running** state) is capable of counting notes of a give bundle, when given CPU
- **condition**: each go-routine can only count when they are given a CPU and also a bundle
- The **extra** worker go-routines would be in the **Runnable/Ready** states waiting near one of the available CPUs, currently not counting but ready to count.
- Conceptually, behavior of behavior of acquiring the chair of the first table they find empty, maps closely to Go’s work-stealing scheduler, where idle processors pull runnable go-routines from other queues to keep CPUs busy.
- The **producers** and **accumulators** go-routines will be special go-routines who do not count but are responsible in facilitating the supply to the workers and gathering overall results

### The corresponding implementation:
- The pipes can be implemented using channels

Now, if we look at the below Golang implementation, things would look way more clear.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Bundle represents a physical bundle of notes (work unit).
type Bundle struct {
	ID    int
	Notes []int // each int is a denomination: 10, 20, 50, 100, 200, 500...
}

// CountResult is what a "person" (goroutine) reports back.
type CountResult struct {
	BundleID int
	Total    int
	Notes    int
	WorkerID int
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Imagine we received many bundles on the desk.
	bundles := makeBundles(12, 80) // 12 bundles, ~80 notes each (mixed)

	// We intentionally start more workers than tables (cores). - Work Stealing of the Go-Scheduler takes care of utilization
	workers := 10

	jobs := make(chan Bundle)
	results := make(chan CountResult)

	// Fan-out: start worker goroutines (people).
	for i := 1; i <= workers; i++ {
		go worker(i, jobs, results)
	}

	// Send bundles to be counted (like handing bundles to people).
	go func() {
		for _, b := range bundles {
			jobs <- b
		}
		close(jobs) // no more bundles
	}()

	// Fan-in: collect exactly N results (one per bundle).
	var grandTotal int
	var totalNotes int
	for i := 0; i < len(bundles); i++ {
		r := <-results
		fmt.Printf("[worker-%d] counted bundle-%d: notes=%d total=%d\n",
			r.WorkerID, r.BundleID, r.Notes, r.Total)

		grandTotal += r.Total
		totalNotes += r.Notes
	}

	fmt.Println("=========================================")
	fmt.Printf("FINAL: bundles=%d notes=%d grandTotal=%d\n", len(bundles), totalNotes, grandTotal)
}

// worker is a "person" who repeatedly:
// 1) takes a bundle,
// 2) counts it,
// 3) reports the result.
func worker(workerID int, jobs <-chan Bundle, results chan<- CountResult) {
	for b := range jobs {
		sum := 0
		for i, denom := range b.Notes {
			sum += denom
			// Simulate variable human speed / interruptions
			// say, if each human worker takes a break after counting some notes, other available workers can take acquire the table and count
			if i%3 == 2 {
				time.Sleep(time.Duration(rand.Intn(40)) * time.Millisecond)
			}
		}

		results <- CountResult{
			BundleID: b.ID,
			Total:    sum,
			Notes:    len(b.Notes),
			WorkerID: workerID,
		}
	}
}

// makeBundles creates N bundles with random denominations.
// Each bundle has up to maxNotes notes (randomized count per bundle).
func makeBundles(n int, maxNotes int) []Bundle {
	denoms := []int{10, 20, 50, 100, 200, 500}
	out := make([]Bundle, 0, n)

	for i := 1; i <= n; i++ {
		nNotes := 10 + rand.Intn(maxNotes-9) // at least 10 notes
		notes := make([]int, 0, nNotes)
		for j := 0; j < nNotes; j++ {
			notes = append(notes, denoms[rand.Intn(len(denoms))])
		}
		out = append(out, Bundle{ID: i, Notes: notes})
	}
	return out
}
```

### Co-relation
- Each **worker go-routine** counts a bundle of notes **received** through the **jobs channel (pipe)** from the **producer go-routine** and **sends** the result through the **results channel (pipe)**.
    This behavior is enforced by the signature of the `worker` function:
    - `jobs` has the type `<-chan Bundle`, a **read-only** channel, which **explicitly restricts** the worker to only _receiving_ from the `jobs` channel and not sending into it.
    - `results` has the type `chan<- CountResult`, a **send-only** channel, which **explicitly restricts** the worker to only _sending_ results into the `results` channel and not receiving from it.
    >This makes the worker’s role unambiguous: it consumes work from one pipe and produces results into another.
- the **main** go-routine is acting as the **accumulator** here as it is holding the **receiving** end of the **results channel** into which each of the **worker** go-routines pass in their results
- each core/CPU of our machine is acting as the counting table
- the `[]Bundle{}` (bundle array) is acting as the bunch of bundles we need to count
- `CountResult` is acting as the result **chit** that each worker prepares at the end of their counting at the table and **send** it in the **result channel (pipe)**, which is received by the **accumulator**

## Conclusion

In the end, I don’t want us to lose sight of the **availability over speed** perspective on concurrency.

Throughout the counting hall story, the goal was never to make a single person count faster.

The goal was to make sure that **_counting never stops_**.

We added more workers not to increase individual speed, but to ensure that:
- when one worker pauses, another can immediately take a table,
- when one table becomes free, some worker is always ready to occupy it,
- and no table ever sits idle while work is waiting.

That is exactly what Go’s concurrency model optimizes for.

Go-routines are cheap to create not so that they all run at once, but so that there is **always** work ready when CPU time becomes available.

**Runnable** go-routines exist so that the system can immediately make progress when capacity frees up.

**Work stealing** exists so that idle CPUs don’t wait while work is stuck elsewhere.

All of this is about **keeping the system responsive and productive**, even as the amount of work scales.

That brings us back to the core idea:

> **Concurrency is not about making one person faster.**
> 
> **It’s about keeping the hall productive even when work scales.**
> 
> And that, more than raw performance, is what good concurrency design is about.

At its core, good concurrency prioritizes:

- **availability** — the system is always ready to do useful work,
- **clarity** — each component has a well-defined role,
- **structure** — work flows in one direction with clear coordination,
- and **human-like organization** — because that’s how large systems naturally scale.