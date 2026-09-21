package main
import "fmt"
import "os"
import "runtime"

func main(){
	helloWorld()
	userName()
	printArgs()
	printGoVersion()
}

func helloWorld(){
	fmt.Println("Hello, World!")
}

func userName(){
	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME")
	}
	fmt.Println(user)	
}

func printArgs() {
    args := os.Args[1:]
    if len(args) == 0 {
        fmt.Println("Аргументы: (не переданы)")
        return
    }
    fmt.Println("Аргументы:")
    for i, arg := range args {
        fmt.Printf("  %d: %s\n", i+1, arg)
    }
}

func printGoVersion() {
    fmt.Println("Версия Go:", runtime.Version())
}