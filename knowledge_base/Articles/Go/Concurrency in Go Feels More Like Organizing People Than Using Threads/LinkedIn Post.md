> For some reason, concurrency in Go has never felt like “threads” to me.
> It has always felt more like **organizing people**.

When I think about go-routines, I don’t picture threads and locks.  
I picture:
- people who know how to do work,
- limited tables they can sit at and work,
- some people actively working,
- others waiting nearby,
- and people moving to free tables when they notice one.

That mental model made a lot of things _click_ for me:
- running vs runnable go-routines
- why we start more go-routines than CPU cores
- how Go’s work-stealing scheduler keeps CPUs busy
- **availability over speed** perspective on concurrency

I wrote a detailed piece explaining this intuition using a currency-counting hall analogy and mapped it directly to Go code.
\[Substack link]
If concurrency has ever felt abstract or mechanical to you, this way of thinking might help.


