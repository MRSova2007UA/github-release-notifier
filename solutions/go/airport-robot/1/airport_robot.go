package airportrobot

import "fmt"

// Write your code here.
// This exercise does not have tests for each individual task.
// Try to solve all the tasks first before running the tests.
type Greeter interface {
	LanguageName() string
	Greet(string) string
}


type Italian struct {}
func (i Italian) LanguageName() string {
    return "Italian"
}
func (i Italian) Greet(name string) string {
    return "Ciao " + name + "!"
}
type Portuguese struct {}
func (p Portuguese) LanguageName() string {
    return "Portuguese"
}
func (p Portuguese) Greet(name string) string {
    return "Olá " + name + "!"
}

func SayHello(name string, g Greeter) string {
	return fmt.Sprintf("I can speak %s: %s", g.LanguageName(), g.Greet(name))
}


