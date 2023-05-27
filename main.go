package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type player struct {
	name  string
	score int
}

func (p *player) increment_score(count int) {
	p.score += count
}

type quiz struct {
	problem  string
	solution string
}

func read_csv(path string) [][]string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file_content, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return file_content
}

func get_input(prompt string, r *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	input, err := r.ReadString('\n')
	return strings.TrimSpace(input), err
}

func main() {
	// Init:
	quiz_timer := time.Now().Add(30 * time.Second)
	reader := bufio.NewReader(os.Stdin)
	name, _ := get_input("Please Enter Your name: ", reader)
	player := player{name: name, score: 0}
	var quiz_set []quiz
	for _, v := range read_csv("./quiz.csv") {
		quiz_set = append(quiz_set, quiz{problem: v[0], solution: v[1]})
	}
	//Game loop:
	fmt.Printf("%v, you have 30s to complete the following quiz\n", name)
	for _, v := range quiz_set {
		go func() {
			quiz_timer.Sub(time.Now()).Round(time.Second)
		}()
		reader := bufio.NewReader(os.Stdin)
		input, _ := get_input(v.problem+" = ", reader) // blocks till next iteration
		// fmt.Println(input, v.solution)

		if v.solution == input {
			player.increment_score(1)
			quiz_timer = quiz_timer.Add(5 * time.Second)
		}

		if time.Now().After(quiz_timer) {
			log.Fatal("Times up bois")
		}
	}
	fmt.Println(player.score, "/", len(quiz_set))
}
