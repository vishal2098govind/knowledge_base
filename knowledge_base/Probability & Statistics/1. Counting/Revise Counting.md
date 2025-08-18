#gate-da #counting #revise
## Rule of Sum & Product
- If Task T can be divided into T1, T2, ... Tn
	- No.of ways of doing T1 = W1
	- No.of ways of doing T2 = W2
	- ... No.of ways of doing Tn = Wn
	- Then, no.of ways of doing 
		- any one of them = $W1+W2+...Wn$ #rule-of-sum 
		- all of them sequentially = $W1*W2*W3*\dots *Wn$ #rule-of-product 

## Factorial
- $n!$ is the no.of ways of 
	- arranging n objects in a row
	- OR distributing n objects among n people/containers

## Divisors
- $n=p_{1}^{e_{1}}*p_{2}^{e_{2}}*\dots*p_{k}^{e_{k}}$ - prime factorisation of number n
- No.of divisors
	- $(e_{1}+1)(e_{2}+1)(e_{3}+1)\dots(e_{k}+1)$ #no-of-divisors 
- Sum of divisors
	- $(p_{1}^{0}+p_{1}^{2}+\dots+p_{1}^{e_{1}})(p_{2}^{0}+p_{2}^{2}+\dots+p_{2}^{e_{2}})\dots(p_{k}^{0}+p_{k}^{2}+\dots+p_{k}^{e_{k}}) = \prod_{i=1}^{k}(\frac{p_{i}^{e_{i}+1}-1}{p_{i}-1})$ #sum-of-divisors

## P & C
- $^nP_{r}=\frac{n!}{(n-r)!}$
- $^nC_{r}=\frac{^nP_{r}}{r!}$
- $^nC_{r}=^nC_{n-r}$
- $^nCr=^{n-1}Cr+^{n-1}C_{r-1}$
	- choose 1st and find no.of ways of selecting r-1 among remaining n-1 ($^{n-1}C_{}{r-1}$)
	- **OR** don't select 1st and find no.of ways of selecting r among remaining n ($^{n-1}C_{r}$) #rule-of-sum
- #boxes-sticks 
## Binomial & Multinomial Theorem
- Binomial Theorem
	- $(X+Y)^n=\sum_{k=0}^n(^n_{k}X^kY^{n-k})$
		- no.of terms = 
			- no.of ways in which n1+n2=n
			- 2 partitions, 1 stick => 1 box #boxes-sticks 
			- n more boxes for each of n => n+1 boxes in total
			- thus, $^{n+1}C_{1}$ ways => $^{n+1}C_{1}$ terms
		- sum of terms = $2^n$
- Multinomial Theorem
	- $(^n_{a_{,1}a_{2},a_{3},\dots,a_{k}})$ is read as $n$ choose $a_{1},a_{2},a_{3},\dots,a_{k}$ 
		- is equal to $^nC_{a_{1}}*^{n-a_{1}}C_{a_{2}}*^{(n-a_{1}-a_{2})}C_{a_{3}}\dots*^{(n-a_{1}-a_{2}-\dots a_{k-1})}C_{a_{k}}$
		- it's the multinomial coefficient of a term in multinomial expansion
	- $$(X_{1}+X_{2}+X_{3}+\dots+X_{m})^{n}=\sum_{a_{1}+a_{2}+\dots a_{k}=n}(^n_{a_{1},a_{2},a_{3}\dots a_{k}})\prod_{i=1}^m(X_{i}^{a_{i}})$$
	- $(X_{1}+X_{2}+X_{3}+X_{4})^n$ = 
		- $(^4C_{4}*^0C_{0}*^0C_{0}*^0C_{0})*(X1^4*X_{2}^0*X_{3}^0*X_{4}^0)$
		- + $(^4C_{3}*^1C_{1}*^0C_{0}*^0C_{0})*(X_{1}^3*X_{2}^1*X_{3}^0*X_{4}^0)$
		- + ...
	- No.of terms: 
		- no.of ways in which $n_{1}+n_{2}+n_{3}+\dots n_{m}=n$
		- $m$ partitions => $m-1$ sticks => $m-1$ boxes #boxes-sticks 
		- $n$ more boxes for each of $n$ => $n+m-1$ boxes in total
		- thus, $^{n+m-1}C_{m-1}$ terms or ways
	- Sum of terms: $m^n$