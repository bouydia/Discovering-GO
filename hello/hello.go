package main

import( "fmt"
  "example.com/greetings"
)

func main() {
    message:=greetings.Hello("Asta")
    fmt.Println(message)
}