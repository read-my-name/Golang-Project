package conc

import (
	"fmt"
	"time"
)

func run_conc(){
	start := time.Now()
	
	defer func(){
		//defer the function execution until the end of the other function
		fmt.Println("Hello World")
	}()
	defer func(){
		//defer the function execution until the end of the other function
		fmt.Println(time.Since(start))
	}()


	go count("Sheep")
	go count("fish")
	//sleep 2 seconds before exit the main
	// if not, it could have a chance to exit the main
	// before executing the two go routine
	time.Sleep(time.Second*2)
}

func count(thing string){
	for i:=1; i<5; i++{
		fmt.Println(i,thing)
		time.Sleep(time.Millisecond*500)
	}

}

