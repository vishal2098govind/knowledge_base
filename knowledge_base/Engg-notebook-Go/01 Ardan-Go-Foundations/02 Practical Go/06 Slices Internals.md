#engineers-notebook #golang #slices 

```go
package main

import "fmt"

func main() {
	s1 := make([]int, 10)
	s2 := s1[3:7]
	fmt.Println(s1)
	fmt.Println(s2)
	s2 = append(s2, 10)
	fmt.Println(&s1[3], &s2[0])
	fmt.Println(s1)
	fmt.Println(s2)
	s2 = append(s2, 20)
	fmt.Println(&s1[3], &s2[0])
	fmt.Println(s1)
	fmt.Println(s2)
	s2 = append(s2, 30)
	fmt.Println(&s1[3], &s2[0])
	fmt.Println(s1)
	fmt.Println(s2)
	s2 = append(s2, 40)
	fmt.Println(&s1[3], &s2[0])
	fmt.Println(s1)
	fmt.Println(s2)
}
```

```sh
➜  slices git:(vishal-kb) ✗ go run .
[0 0 0 0 0 0 0 0 0 0]
[0 0 0 0]
0xc00001a068 0xc00001a068
[0 0 0 0 0 0 0 10 0 0]
[0 0 0 0 10]
0xc00001a068 0xc00001a068
[0 0 0 0 0 0 0 10 20 0]
[0 0 0 0 10 20]
0xc00001a068 0xc00001a068
[0 0 0 0 0 0 0 10 20 30]
[0 0 0 0 10 20 30]
0xc00001a068 0xc000026150 -> new underlying array created here for s2
[0 0 0 0 0 0 0 10 20 30]
[0 0 0 0 10 20 30 40]
```
```
```

```
                s1 := make([]int, 10)
              array, len=10, cap=10
               │                    
               └──────────────┬──────────────────────────┐
                              │                          │
                              ↓                          ↓ 
            Underlying array [0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
               indices        0  1  2  3  4  5  6  7  8  9
                                       ↑        ↑
                                       │        │
                  ┌────────────────────┴────────┘
                  │
               array, len=4 (7-3), cap=7 (10-3)
                 s2 := s1[3:7]


			s2 = s2.append(10)
			  
			  s1        
              array, len=10, cap=10
               │                    
               └──────────────┬───────────────────────────┐
                              │                           │
                              ↓                           ↓ 
			Underlying array [0, 0, 0, 0, 0, 0, 0, 10, 0, 0]
			   indices        0  1  2  3  4  5  6   7  8  9
									   ↑            ↑
                                       │            │
                  ┌────────────────────┴────────────┘
                  │
			   array, len=4 (7-3), cap=7 (10-3)
			   	s2 := s1[3:8]

			s2 = s2.append(20)
			  
			  s1        
              array, len=10, cap=10
               │                    
               └──────────────┬────────────────────────────┐
                              │                            │
                              ↓                            ↓ 
			Underlying array [0, 0, 0, 0, 0, 0, 0, 10, 20, 0]
			   indices        0  1  2  3  4  5  6   7   8  9
									   ↑                ↑
                                       │                │
                  ┌────────────────────┴────────────────┘
                  │
			   array, len=4 (7-3), cap=7 (10-3)
			   	s2 := s1[3:9]
		

			s2 = s2.append(30)
			  
			  s1        
              array, len=10, cap=10
               │                    
               └──────────────┬─────────────────────────────┐
                              │                             │
                              ↓                             ↓ 
			Underlying array [0, 0, 0, 0, 0, 0, 0, 10, 20, 30]
			   indices        0  1  2  3  4  5  6   7   8   9
									   ↑                    ↑
                                       │                    │
                  ┌────────────────────┴────────────────────┘
                  │
			   array, len=7 (7-3), cap=7 (10-3)
			   	s2 := s1[3:10]

			s2 = s2.append(40)
			  
			  s1        
              array, len=10, cap=10
               │                    
               └──────────────┬─────────────────────────────┐
                              │                             │
                              ↓                             ↓ 
			Underlying array [0, 0, 0, 0, 0, 0, 0, 10, 20, 30]
			   indices        0  1  2  3  4  5  6   7   8   9
			   ------------------------------------ (creates new underlying array now)
			Underlying array 		  [0, 0, 0, 0, 10, 20, 30, 40]
			   indices       		   0  1  2  3  4  5  6  7   8
									   ↑                        ↑
                                       │                        │
                  ┌────────────────────┴────────────────────────┘
                  │
			   array, len=8
			   	s2
```


```go
package main

import "fmt"

func main() {
	cart := []string{"apple", "orange", "banana"}
	fmt.Println("len:", len(cart), "cap:", cap(cart))
	fmt.Println("cart[1]:", cart[1])

	// indices
	for i := range cart {
		fmt.Println(i)
	}

	// index + value
	for i, c := range cart {
		fmt.Println(i, c)
	}

	// values
	for _, c := range cart {
		fmt.Println(c)
	}

	cart = append(cart, "milk")
	fmt.Println(cart, "cap:", cap(cart))

	// cart: {apple, orange, banana, milk}
	// fruit: {apple, orange, banana}
	// append(fruit, lemon): {apple, orange, banana, lemon}

	// slicing operator, half-open
	fruit := cart[:2] // f:{apple, orange}, c:{apple, orange, banana, milk}
	fmt.Println("fruit:", fruit, "&fruit:", &fruit, "cap(fruit)", cap(fruit), "\n&cart:", &cart, "cap(cart)", cap(cart))

	fruit = append(fruit, "lemon") // f:{apple, orange, lemon}, c:{apple, orange, lemon, milk}
	fmt.Println("fruit:", fruit, "cart", cart, "\n&fruit", &fruit, "&cart", &cart)
	fmt.Println("cap(fruit)", cap(fruit), "len(fruit)", len(fruit))
	fmt.Println("cap(cart)", cap(cart), "len(cart)", len(cart))

	fruit = append(fruit, "bread") // f:{apple, orange, lemon, bread}, c:{apple, orange, lemon, bread}
	fmt.Println("cap(fruit)", cap(fruit), "len(fruit)", len(fruit))
	fmt.Println("cap(cart)", cap(cart), "len(cart)", len(cart))
	fmt.Println("fruit:", fruit, "cart", cart, "\n&fruit", &fruit, "&cart", &cart)

	fruit = append(fruit, "butter") // f:{apple, orange, lemon, bread, butter}, c:{apple, orange, lemon, bread}
	fmt.Println("fruit:", fruit, "cart", cart, "\n&fruit", &fruit, "&cart", &cart)
	// fmt.Println("cart:", cart)

	// items := make([]string, 4)
	// copy(items, cart)
	// items = append(items, "bread")
	// fmt.Println("cart:", cart)
	// fmt.Println("fruit:", fruit)
	// fmt.Println("items:", items)
}

```