package main

func main() {
	err := RunApp("A 16th-century woman with long brown hair standing in front of a green vista with cloudy skies. She's looking at the viewer with a faint smile on her lips.")
	if err != nil {
		panic(err)
	}
}
