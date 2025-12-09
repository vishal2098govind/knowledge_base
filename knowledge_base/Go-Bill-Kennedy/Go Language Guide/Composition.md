#go #go-embed #composition-over-inheritance 

- Composition is one of the most important design aspects of Go
- Composition is something that we need to leverage all of the time
- Composition is really what allows us to be able to write code that can decouple itself from change but compose over time, so we can improve and extend all the things we do

## Quotes around design:

> **"A good API is one that is not just easy to use but also really hard to misuse."**

> **"You can always embed, but you cannot decompose big interfaces once they are out there. Keep interfaces small"**

> **"Don't design with interfaces, discover them"**

> **"Duplication is far cheaper than wrong abstraction"**

> **"Look for making things easier to understand than easier to do"**

## Grouping Types By States Vs Behavior

- This is an example of using type hierarchies with a OOP pattern.
- This is not something we want to do in Go. Go does not have the concept of sub-typing. 
- All types are their own and the concepts of base and derived types do not exist in Go. 
- This pattern does not provide a good design principle in a Go program.

```go
package main

import "fmt"

// Animal contains all the base fields for animals.
type Animal struct {
	Name     string
	IsMammal bool
}

// Speak provides generic behavior for all animals and
// how they speak.
func (a *Animal) Speak() {
	fmt.Printf(
		"UGH! My name is %s, it is %t I am a mammal\n",
		a.Name,
		a.IsMammal,
	)
}

// Dog contains everything an Animal is but specific
// attributes that only a Dog has.
type Dog struct {
	Animal
	PackFactor int
}

// Speak knows how to speak like a dog.
func (d *Dog) Speak() {
	fmt.Printf(
		"Woof! My name is %s, it is %t I am a mammal with a pack factor of %d.\n",
		d.Name,
		d.IsMammal,
		d.PackFactor,
	)
}

// Cat contains everything an Animal is but specific
// attributes that only a Cat has.
type Cat struct {
	Animal
	ClimbFactor int
}

// Speak knows how to speak like a cat.
func (c *Cat) Speak() {
	fmt.Printf(
		"Meow! My name is %s, it is %t I am a mammal with a climb factor of %d.\n",
		c.Name,
		c.IsMammal,
		c.ClimbFactor,
	)
}

func main() {

	// Create a list of Animals that know how to speak.
	animals := []Animal{

		// Create a Dog by initializing its Animal parts
		// and then its specific Dog attributes.
		// compiler error: cannot use Dog{…} (value of struct type Dog) as Animal value in array or slice literal
		Dog{
			Animal: Animal{
				Name:     "Fido",
				IsMammal: true,
			},
			PackFactor: 5,
		},

		// Create a Cat by initializing its Animal parts
		// and then its specific Cat attributes.
		// compiler error: cannot use Cat{…} (value of struct type Cat) as Animal value in array or slice literal
		Cat{
			Animal: Animal{
				Name:     "Milo",
				IsMammal: true,
			},
			ClimbFactor: 4,
		},
	}

	// Have the Animals speak.
	for _, animal := range animals {
		animal.Speak()
	}
}
```

### Smells:
- The Animal type is providing an abstraction layer of **reusable state**. Prefer working on **reusable behavior (interface)**
- The program never needs to create or solely use a value of type Animal.
- The implementation of the Speak method for the Animal type is a generalization.
- The Speak method for the Animal type is never going to be called.

### Fix:
- When working with concrete types, we care more about **what things are (state)**
- but when we are working with decoupling and working with multiple types, we care more about **what does these things do (behavior)**
- i.e. prefer grouping types based on behavior than grouping types based on state
- When embedding, **don't embed for state but for behavior**, which helps in decoupling
- Focus on behavior of things being dealt with, than the state of things being dealt with
	- Start thinking on **what animals do** rather than on **what animals are**
	- then go ahead and define interface(s) based on **what animals do**
- Interfaces describe **behavior of things** 
- and concrete types describe **what the things are**

- This is an example of using composition and interfaces. 
- This is something we want to do in Go. 
- We will group common types by their behavior and not by their state. 
- This pattern does provide a good design principle in a Go program.
```go
package main

import "fmt"

// Speaker provide a common behavior for all concrete types
// to follow if they want to be a part of this group. This
// is a contract for these concrete types to follow.
type Speaker interface {
	Speak()
}

// Dog contains everything a Dog needs.
type Dog struct {
	Name       string
	IsMammal   bool
	PackFactor int
}

// Speak knows how to speak like a dog.
// This makes a Dog now part of a group of concrete
// types that know how to speak.
func (d *Dog) Speak() {
	fmt.Printf(
		"Woof! My name is %s, it is %t I am a mammal with a pack factor of %d.\n",
		d.Name,
		d.IsMammal,
		d.PackFactor,
	)
}

// Cat contains everything a Cat needs.
type Cat struct {
	Name        string
	IsMammal    bool
	ClimbFactor int
}

// Speak knows how to speak like a cat.
// This makes a Cat now part of a group of concrete
// types that know how to speak.
func (c *Cat) Speak() {
	fmt.Printf(
		"Meow! My name is %s, it is %t I am a mammal with a climb factor of %d.\n",
		c.Name,
		c.IsMammal,
		c.ClimbFactor,
	)
}

func main() {

	// Create a list of Animals that know how to speak.
	speakers := []Speaker{

		// Create a Dog by initializing its Animal parts
		// and then its specific Dog attributes.
		&Dog{
			Name:       "Fido",
			IsMammal:   true,
			PackFactor: 5,
		},

		// Create a Cat by initializing its Animal parts
		// and then its specific Cat attributes.
		&Cat{
			Name:        "Milo",
			IsMammal:    true,
			ClimbFactor: 4,
		},
	}

	// Have the Animals speak.
	for _, spkr := range speakers {
		spkr.Speak()
	}
}
```
Here are some guidelines around declaring types:
* Declare types that represent something new or unique.
* Validate that a value of any type is created or used on its own.
* **Embed** types **to reuse existing behaviors** you need to satisfy.
* Question types that are an alias or abstraction for an existing type.
* Question types whose sole purpose is to share common state.