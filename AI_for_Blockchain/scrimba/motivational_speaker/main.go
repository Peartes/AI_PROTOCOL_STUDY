package main

func main() {
	err := RunApp("motivational_bot_data.jsonl", "I am feeling down today. Can you give me some motivation?")
	if err != nil {
		panic(err)
	}
}
