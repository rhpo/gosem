# A Go SE2 (L3/ISIL) Semaphore Simulation

## Installation
```bash
go get -u https://github.com/rhpo/gosem
```

## Code Example (Classic Pont Semaphore)
```go
package main

import (
	"fmt"
	. "github.com/rhpo/gosem"
)

var (
	pont      = NomSemaphore()
	SemNbNord = NomSemaphore()
	SemNbSud  = NomSemaphore()
)

func main() {
	I(pont, 1)
	I(SemNbNord, 1)
	I(SemNbSud, 1)

	nbNord := 0
	nbSud := 0

	// Processus Nord
	Repeat(5, func() {
		Process(func() {

			P(SemNbNord)

			nbNord++
			if nbNord == 1 {
				P(pont)
			}

			V(SemNbNord)

			fmt.Println("Une voiture nord a passé")

			nbNord--
			if nbNord == 0 {
				V(pont)
			}

		})
	})

	// Processus Sud
	Repeat(5, func() {
		Process(func() {

			P(SemNbSud)

			nbSud++
			if nbSud == 1 {
				P(pont)
			}

			V(SemNbSud)

			fmt.Println("Une voiture sud a passé")

			nbSud--
			if nbSud == 0 {
				V(pont)
			}

		})
	})

	Wait()

}
```
