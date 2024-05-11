package main 

import ( 
	"fmt"
	"time"
	"strconv"
) 

func display(str string) { 
 
	fmt.Println(str) 
	
} 

func main() { 
	time.Sleep(1 * time.Millisecond)
	// Calling Goroutine 
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Millisecond)
		go display("Welcome" + strconv.Itoa(i))
	}
	for i := 0; i < 10; i++ {
		display(strconv.Itoa(i))
	}
} 
