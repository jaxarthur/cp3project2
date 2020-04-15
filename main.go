package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
)

//Structs

//Person stores information about a person.
type Person struct {
	name           string
	ratings        []Rating
	friends        []Friend
	recomendations []Recomendation
}

//Friend stores information about a person with simular intrests.
type Friend struct {
	likeness int
	data     *Person
}

//Rating stores information about how a book was rated.
type Rating struct {
	value int
	book  *Book
}

//Book stores information about a book.
type Book struct {
	name   string
	author string
}

//Recomendation stores information about recomendations for a person
type Recomendation struct {
	strength int //Calculated with ranking * 10 + likeness sumed for everyone
	book     *Book
}

// Functions

func main() {
	CompileData()
}

//CompileData A function that returns both book and person data based on provided text files.
func CompileData() ([]Book, []Person) {
	books := compileBooks()
	people := compilePeople(books)
	people = linkFriends(people)
	people = linkRecomendations(books, people)

	return books, people
}

func compileBooks() []Book {
	books := make([]Book, 0)
	file, _ := os.Open("booklist.txt")
	reader := csv.NewReader(file)
	lines, error := reader.ReadAll()
	if error != nil {
		log.Fatal(error)
	}

	for _, line := range lines {
		newBook := Book{line[1], line[0]}
		books = append(books, newBook)
	}

	return books
}

func compilePeople(books []Book) []Person {
	people := make([]Person, 0)
	file, _ := os.Open("ratings.txt")
	reader := csv.NewReader(file)
	lines, error := reader.ReadAll()
	if error != nil {
		log.Fatal(error)
	}

	for i := 0; i < len(lines); i += 2 {
		name := lines[i][0]
		ratingValues := strings.Split(lines[i+1][0], " ")
		ratingValues = ratingValues[:len(ratingValues)-1]
		ratings := make([]Rating, 0)

		for bookInd, ratingValueStr := range ratingValues {
			ratingValue, _ := strconv.Atoi(ratingValueStr)
			newRating := Rating{ratingValue, &books[bookInd]}
			ratings = append(ratings, newRating)
		}
		newPerson := Person{name, ratings, make([]Friend, 0), make([]Recomendation, 0)}
		people = append(people, newPerson)
	}
	return people
}

func linkFriends(people []Person) []Person {
	for pi, person := range people {
		for oi, otherPerson := range people {
			if !(person.name == otherPerson.name) {
				likeness := 0
				for i, p := range person.ratings {
					o := otherPerson.ratings[i]
					likeness += p.value * o.value
				}
				newFriend := Friend{likeness, &people[oi]}
				people[pi].friends = append(people[pi].friends, newFriend)
			}
		}
	}
	return people
}

func linkRecomendations(books []Book, people []Person) []Person {
	for pi, person := range people {
		for _, rating := range person.ratings {
			strength := 0
			for _, friend := range person.friends {
				friendRatings := friend.data.ratings
				for _, friendRating := range friendRatings {
					if rating.book.name == friendRating.book.name {
						if rating.value == 0 && !(friendRating.value == 0) {
							strength += friendRating.value*10 + friend.likeness
						}
					}
				}
			}
			newRecomendation := Recomendation{strength, rating.book}
			people[pi].recomendations = append(people[pi].recomendations, newRecomendation)
		}
	}
	return people
}
