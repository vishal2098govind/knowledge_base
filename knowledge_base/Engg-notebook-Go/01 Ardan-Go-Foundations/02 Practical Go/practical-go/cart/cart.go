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
