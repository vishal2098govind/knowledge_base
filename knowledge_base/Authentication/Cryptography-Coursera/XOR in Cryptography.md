#cryptography 

## Theorem: 
Let Y be a random variable with arbitrary distribution over {0, 1}^n
Let X be an independent random variable with uniform distribution over {0, 1}^n
Then, Z := Y (XOR) X is will be a random variable with uniform distribution over {0, 1}^n

i.e. if we take any arbitrarily or maliciously distributed random variable Y
and XOR with an **independent** uniformly distributed random variable X, 
we end up with a uniformly distributed random variable Z

This property of XOR makes it very useful for crypto