package main

func main() {
	err := RunApp("motivational_bot_data.jsonl")
	if err != nil {
		panic(err)
	}
}
