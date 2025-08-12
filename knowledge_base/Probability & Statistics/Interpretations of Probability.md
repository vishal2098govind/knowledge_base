#gate-da #probability

### What is Probability?
Two ways of **looking at** or **interpreting** probability:

### Likeliness Interpretation of Probability
- (RR vs RCB) will be an experiment
	- Event $E$: RCB wins toss
	- Event $E'$: RR wins toss
- Probability of $E$ i.e. $P(E)$: here it is **likelihood interpretation** of probability
	- i.e. **the event $E$ has not yet happened**, but from past experiences, can we tell somehow the chances of $E$ to happen
- Generally we use this likelihood interpretation of probability itself 
	- i.e. what is the probability that today it'll rain
- There are lot of factors we base on to come up with probability of an $E$ (**conditional probability**)
- E.g.: $A$: RCB wins the toss,  $B$: RR wins the match
	- $P(B|A)$
	- Lot of events have **influence** on one another
	- Also there can be factors/events which **do not influence** occurrence of an event. 
		- E.g. $C$: DC is practicing today. $P(C)$ can be anything but $C$ has no influence on $B$ and $A$
	- Here, $P(A)$ is easy to calculate, given the coin is fair, i.e. $P(A)$ = $1/2$
		- but $P(B)$ is difficult to calculate as $B$ depends on lot of other factors influencing $B$
		- There are not a lot of factors influencing $A$, other than the fairness of the coin
- Thus, in **statistical world**, 
	- we do not solve problems like $P(B)$
		- We will solve problems like finding $P(B)$, which are kind of real-life problems (with many influencing factors), in machine learning, because, these problems are so difficult to solve as we cannot enumerate all possible problems that can influence $B$ to happen/occur
		- Finding $P(B)$ becomes a **prediction problem** which is very difficult because we do not know the factors influencing the occurrence of the event, there can be many like 
			- Toss, Pitch, Player Injuries, Team performance, .......
	- and rather solve problems like $P(A)$
		- Here, other than fairness of coin, there are no much events influencing $A$
		- Also here event $A$ is **truly random**, as 50% chance of landing on tails and 50% chance of landing on heads
		- Btw, no event in the world is truly random as there will be some factors influencing it for sure, but for our simplicity we say tossing a coin is random

### Combinatorial Interpretation of Probability
- For the same event, we can have both interpretations of probability
- Let here, 
	- **Experiment** here is toss is happening between RCB vs RR 
	- **Event** $A$: RCB wins the toss
- Combinatorial interpretation says that, 
	- let's say if the experiment is performed, many many no.of times
		- 1st time: A happens i.e. RCB wins the toss
		- 2nd time: A happens i.e. RCB wins the toss
		- 3rd time: A does not happen i.e. RR wins the toss
		- .... 99 times
		- 100th time: A happens i.e. RCB wins the toss
	- Combinatorial interpretation says that $$P(A) = \frac{no.of \space times \space A \space happens, n_{A}}{no.of\space experiment\space is\space performed, n_{tosses}}$$
	- Here, the condition is the experiment has to be performed very very large number of, i.e. $\infty$ number of times
	- The higher number of times we perform the experiment, the closer we get to the probability $P(A)$
	- This is also called the **frequency interpretation**
	- Here, there's no talk about likelihood, here we are just performing the experiment many many no.of times and take count of what we want. The idea here is, if we perform the experiment around say 1 million times, that ratio will be very close to the $P(A)$
```python
import numpy as np

n_tosses = 10
results = np.random.randint(0, 2, n_tosses) # can get 0 or 1 (50% probabilty ~ coin toss)
# [1, 1, 0, 0, 0, 1, 1, 0, 0, 1]
# 1: RCB wins the toss (A) 
# 0: RCB looses the toss or RR wins the toss

# according to frequence/combinatorial interpretation, 
# P(A) = (#1s in the tosses array) / (n_tosses)

n_ones = 0
for r in results:
	if r == 1:
		n_ones += 1

prob_A = n_ones/n_tosses
print("P(A):", pa)
```

- Here, the probability might come different every time, but the law of large numbers says, 
	- if we perform this experiment $\infty$ number of times
	- i.e. $$P(A) = \lim_{ n_{tosses} \to \infty } \frac{n_{A}}{n_{tosses}}$$