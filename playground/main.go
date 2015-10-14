package main

func main() {
	a := NewApp()
	err := a.Run()
	if err != nil {
		println("Run Error:", err.Error())
	}
	println("Exits.")
}
