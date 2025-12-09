#go #go-interface

## Draft-1: Pure and concrete implementation - No interfaces anywhere
### Primitive Layer
```go
// data being pulled/stored
type Data struct {
	Line string
}

// pull from
type Xenia struct {
	Host    string
	Timeout time.Duration
}

// Pull knows how to pull data out of Xenia
func (x *Xenia) Pull(d *Data) error {
	switch rand.Intn(10) {
	case 1, 9:
		return io.EOF

	case 5:
		return errors.New("error reading data from Xenia")

	default:
		d.Line = "Data"
		fmt.Println("In:", d.Line)
		return nil
	}
}

// store at
type Pillar struct {
	Host    string
	Timeout time.Duration
}

// Store knows how to store data into Pillar
func (p *Pillar) Store(d *Data) error {
	fmt.Println("Out: ", d.Line)
	return nil
}

// defined concrete types using which we can pull from source data-store and store at the destination data-store
```
### Low Level Layer
```go
// pull knows how to pull bulk data from Xenia
func pull(x *Xenia, data []Data) (int, error) {
	for i := range data {
		if err := x.Pull(&data[i]); err != nil {
			return i, err
		}
	}

	return len(data), nil
}

// store knows how to store data into Pillar
func store(p *Pillar, data []Data) (int, error) {
	for i := range data {
		if err := p.Store(&data[i]); err != nil {
			return i, err
		}
	}

	return len(data), nil
}
```
### High Level Layer
```go
// define a system combining the concrete pull and store types
type System struct {
	Xenia
	Pillar
}

func Copy(sys *System, batch int) error {
	data := make([]Data, batch)

	for {

		i, err := pull(&sys.Xenia, data)
		if i > 0 {
			if _, err := store(&sys.Pillar, data[:i]); err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}

	}
}
```

### API Client
```go
package main

func main() {
	s := System{
		Xenia:  Xenia{Host: "localhost:8000", Timeout: time.Second},
		Pillar: Pillar{Host: "", Timeout: time.Second},
	}

	if err := Copy(&s, 3); err != nil {
		fmt.Println(err)
	}
}
```

## Draft-2: Discover interfaces
Extract behaviors to interfaces
### Primitive Layer
```go
// Puller
type Puller interface {
	Pull(d *Data) error
}

type Storer interface {
	Store(d *Data) error
}
```
### Low-Level API Layer
```go
// pull knows how to pull bulk data from Xenia
func pull(x Puller, data []Data) (int, error) {
	for i := range data {
		if err := x.Pull(&data[i]); err != nil {
			return i, err
		}
	}

	return len(data), nil
}

// store knows how to store data into Pillar
func store(s Storer, data []Data) (int, error) {
	for i := range data {
		if err := s.Store(&data[i]); err != nil {
			return i, err
		}
	}

	return len(data), nil
}
```

## Draft-3: Readability Code Review
Get rid of unnecessary types to improve readability and simplicity
### High Level Layer
```go
// Copy knows how to pull and store data from any Puller to Storer
func Copy(p Puller, s Storer, batch int) error {
	data := make([]Data, batch)

	for {

		i, err := pull(p, data)
		if i > 0 {
			if _, err := store(s, data[:i]); err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}

	}
}
```
### API Client
```go
func main() {
	p := Xenia{Host: "localhost:8000", Timeout: time.Second}
	s := Pillar{Host: "", Timeout: time.Second}

	if err := Copy(&p, &s, 3); err != nil {
		fmt.Println(err)
	}

	// can easily extend to any other future combinations of concrete pullers and storers
	// p = Postgres{}
	// s = Pillar{}
	// if err := Copy(&p, &s, 3); err != nil {
	// 	fmt.Println(err)
	// }
}

```