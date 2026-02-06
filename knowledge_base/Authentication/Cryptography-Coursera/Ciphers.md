#cryptography 

## Define Cipher
- A cipher is **defined over a triple**
	- Set of all possible keys (K)
	- Set of all possible messages (M)
	- Set of all possible cipher text (C)
- This triple defines the environment over which the ciphers are defined
- A cipher **is a pair** of **efficient** algorithms
	- Encryption (E): 
		- E: K x M -> C
	- Decryption (D)
		- D : K x C -> M
	- such that
		- for all m in M, and k in K
		- D(k, E(k, m)) = m
- E is often a **randomized** algorithm
- D is always **deterministic** algorithm

### One Time Pad Cipher
- One of the first ciphers, defined back in starting of 20th century
- M = C = K ={0, 1}^n
- |k| = |c| i.e. key is a random bit string as long as message
- E(k, m) = k XOR m
- D(k, c) = k XOR c
- D(k, E(k, m)) = D(k, k XOR m) = k XOR (k XOR m) 
	- = (k XOR k) XOR m
	- = 0 XOR m
	- = m
- Thus, this qualifies it to be a cipher
- In fact, we an identify the k, given m and c by simply finding m XOR c
- Although One-Time-Pad (OTP) is a very fast cipher, it is not practical, since the key size has to match the message size
- If we could find a way to send such a long key to both parties (Alice and Bob) securely, we would send the message itself using that way as well
- So, is OTP a good/secure cipher? What is a secure Cipher?
- What makes a cipher secure?

# Cipher Security
## Information Theoretic Security
(by Shannon in 1949)

Basic Idea: 
- If all we get to see is a cipher text (CT), 
- then we should learn absolutely nothing about the plain text (PT)
- i.e. a CT should reveal no **information** about the PT

A cipher (E, D) over (K, M, C) has **perfect secrecy** if
- for every two messages of same length m0 and m1 in M, with |m0| = |m1| (same length)
- and for every cipher text c in C
- and for **random uniform distribution** of key k in K
- Probability of E(k, m0) being c 
	- should be same as
	- Probability of E(k, m1) being c
- i.e. **P\[E(k, m0) = c] = P\[E(k, m1) = c]**

### What does it mean by the two probabilities to be same?
- If an attacker intercepts a particular cipher text c, 
- then in reality, the probability that this cipher text is the encryption of the message m0 
- is exactly same as the probability that this cipher text is the encryption of message m1

- So, if all the attacker has is "c" that he could intercept, then we have no idea if the cipher text came from either m0 or m1

### OTP has perfect secrecy
### Does cipher with perfect secrecy mean secure cipher to use?
- No
- cipher with perfect secrecy only says it that given cipher text, we learn nothing about the plain text it was encrypted from, i.e. no **CT only attack** possible on ciphers with perfect secrecy
- Although other attacks are possible
- The problem with OTP is, long keys
- Are there other ciphers that have perfect secrecy and have possibly much much shorter keys?

### The Bad news with perfect secrecy
- Shannon proved another theorem that:
- perfect secrecy => |K| >= |M|
- i.e. if a cipher has perfect secrecy, then the number of keys in cipher **must** be at least equal to the number of messages that the cipher can handle
- i.e. key length must be at least the length of message
- thus, the notion of perfect-secrecy is hard to use in practice, and doesn't really tell if practical ciphers are actually going to be secure


## Stream Cipher: making perfect secrecy practical
idea: replace "random" key by "pseudo random" key

**PRG : Pseudo-Random-Generator**
PRG is a function G: {0, 1}^s-> {0, 1}^n, n >> s
![[Pasted image 20260206210537.png]]
 - Now, since size of key is less than size of message, a stream cipher is no more a perfectly secure cipher
 - security of stream cipher will depend on specific PRG
 - **PRG must be unpredictable**
 - E(k, m) = m XOR G(k)
 - D(k, c) = c XOR G(k)

### Attacks on One-Time-Pad
- Two-Time pad is insecure
	- never use same stream cipher key more than once
		- c1 = m1 XOR G(k)
		- c2 = m2 XOR G(k)
		- if attacker has c1 and c2
			- attacker can compute c1 XOR c2
			- i.e. (m1 XOR G(k)) XOR (m2 XOR G(k))
			- = m1 XOR m2 XOR 0
			- = m1 XOR m2
		- then, the attacker could recover m1 and m2 both since english and ASCII have enough redundancy
- No Integrity 