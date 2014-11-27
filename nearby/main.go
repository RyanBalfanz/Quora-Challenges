package main

import (
	"bufio"
	"fmt"
	// "github.com/davecheney/profile"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type ProblemDescription struct {
	T int // Number of topics
	Q int // Number of questions
	N int // Number of queries
}

type Topic struct {
	Id int
	X  float64
	Y  float64
}

type Question struct {
	Id       int
	qn       int
	TopicIds []int
}

// type QuestionInterface interface {
// 	GetTopics() []Topic
// }

type Query struct {
	ResultType      string  // "t" for topic or "q" for question)
	RequiredResults int     // the number of results required
	X               float64 // the x-location to be used as the query
	Y               float64 // the y-location to be used as the query
}

type Topics map[int]*Topic
type Questions map[int]*Question
type Queries []*Query

func NewTopic(s string) *Topic {
	tokens := strings.Split(s, " ")
	id, err := strconv.Atoi(tokens[0])
	if err != nil {
		log.Fatal(err)
	}
	x, err := strconv.ParseFloat(tokens[1], 64)
	if err != nil {
		log.Fatal(err)
	}
	y, err := strconv.ParseFloat(tokens[2], 64)
	if err != nil {
		log.Fatal(err)
	}

	return &Topic{Id: id, X: x, Y: y}
}

func NewQuestion(s string) *Question {
	tokens := strings.Split(s, " ")
	id, err := strconv.Atoi(tokens[0])
	if err != nil {
		log.Fatal(err)
	}
	qn, err := strconv.Atoi(tokens[1])
	if err != nil {
		log.Fatal(err)
	}
	topicIds := make([]int, qn)
	for t, tStr := range tokens[2:] {
		tId, err := strconv.Atoi(tStr)
		if err != nil {
			log.Fatal(err)
		}
		topicIds[t] = tId
	}

	return &Question{id, qn, topicIds}
}

func NewQuery(s string) *Query {
	tokens := strings.Split(s, " ")
	resultType := tokens[0]
	requiredResults, err := strconv.Atoi(tokens[1])
	if err != nil {
		log.Fatal(err)
	}
	x, err := strconv.ParseFloat(tokens[2], 64)
	if err != nil {
		log.Fatal(err)
	}
	y, err := strconv.ParseFloat(tokens[3], 64)
	if err != nil {
		log.Fatal(err)
	}

	return &Query{resultType, requiredResults, x, y}
}

type Point struct {
	X, Y float64
}

func (p Point) Distance(q Point) float64 {
	xDist := p.X - q.X
	yDist := p.Y - q.Y
	radicand := (math.Pow(xDist, 2) + math.Pow(yDist, 2))
	return math.Sqrt(radicand)
}

type DistanceTo struct {
	Id       int
	Distance float64
}

type DistanceToTopicQuestion struct {
	TopicId    int
	QuestionId int
	Distance   float64
}

type ByDistance []DistanceTo
type ByDistanceThroughTopic []DistanceTo // Distances from a Point to a Question

func (a ByDistance) Len() int           { return len(a) }
func (a ByDistance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDistance) Less(i, j int) bool { return a[i].Distance < a[j].Distance }

func main() {
	// log.Println("Nearby")
	// defer profile.Start(profile.CPUProfile).Stop()

	problemDescription := new(ProblemDescription)
	topics := make(Topics)
	questions := make(Questions)
	queries := make(Queries, 0)

	scanner := bufio.NewScanner(os.Stdin)

	// f, _ := os.Open("larger_input.txt")
	// scanner := bufio.NewScanner(f)

	scanner.Split(bufio.ScanLines)
	count := 0
	for scanner.Scan() {
		count++
		line := scanner.Text()

		// Read the first line, containing the problem description
		if count == 1 {
			tokens := strings.Split(line, " ")
			var ts []int
			for _, token := range tokens {
				tokenInt, _ := strconv.Atoi(token)
				ts = append(ts, tokenInt)
			}
			problemDescription = &ProblemDescription{ts[0], ts[1], ts[2]}
			continue
		}

		// Read the topic lines
		topicsStartOffset := 1
		if count > topicsStartOffset && count <= problemDescription.T+1 {
			topic := NewTopic(line)
			topics[topic.Id] = topic
			continue
		}

		// Read the question lines
		questionsStartOffset := topicsStartOffset + problemDescription.T
		if count > questionsStartOffset && count <= questionsStartOffset+problemDescription.Q {
			question := NewQuestion(line)
			questions[question.Id] = question
			continue
		}

		// Read the query lines
		queriesStartOffset := questionsStartOffset + problemDescription.Q
		if count > queriesStartOffset && count <= queriesStartOffset+problemDescription.Q {
			query := NewQuery(line)
			queries = append(queries, query)
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	// For each query compute and print the result
	for _, query := range queries {
		distances := make([]DistanceTo, 0)

		if query.ResultType == "t" {
			// Computes the result of a topic-based query
			queryLocation := Point{X: query.X, Y: query.Y}
			for _, topic := range topics {
				topicLocation := Point{X: topic.X, Y: topic.Y}
				d := DistanceTo{
					Id:       topic.Id,
					Distance: queryLocation.Distance(topicLocation),
				}
				distances = append(distances, d)
			}

			// Print the result of this topic-based query to stdout
			// Print each Topic.Id by distance in descending order up to the number of requested reults
			sort.Sort(ByDistance(distances))
			for _, dist := range distances[:query.RequiredResults] {
				fmt.Print(dist.Id, " ")
			}
			fmt.Println()
		} else {
			// Computes the result of a question-based query
			queryLocation := Point{X: query.X, Y: query.Y}
			for _, topic := range topics {
				topicLocation := Point{X: topic.X, Y: topic.Y}
				d := DistanceTo{
					Id:       topic.Id,
					Distance: queryLocation.Distance(topicLocation),
				}
				distances = append(distances, d)
			}

			sort.Sort(ByDistance(distances))

			questionIds := make([]int, 0)
			questionIdsSet := make(map[int]bool)

			// Only consider the X-nearest topics, this may not yield a correct answer
			X := 1000
			opts := []int{X, len(distances)}
			sort.Ints(opts)
			size := opts[0]
			nearestTopicsDistances := distances[:size]
			for _, distance := range nearestTopicsDistances {
				for i := problemDescription.Q - 1; i >= 0; i-- {
					question := questions[i]
					// log.Println(i, question)
					for _, topicId := range question.TopicIds {
						if topicId == distance.Id {
							// log.Println(len(questionIds))
							if len(questionIds) == 0 {
								questionIds = append(questionIds, question.Id)
								questionIdsSet[question.Id] = true
								// log.Println("Added: ", questionIds)
							} else {
								if _, exists := questionIdsSet[question.Id]; !exists {
									questionIds = append(questionIds, question.Id)
									questionIdsSet[question.Id] = true
								}
							}
						}
					}
				}
			}

			// Print the result of this topic-based query to stdout
			// Print each Topic.Id by distance in descending order up to the number of requested reults
			for d, dist := range questionIds {
				fmt.Print(dist, " ")
				if d >= query.RequiredResults {
					break
				}
			}
			fmt.Println()
		}
	}
}
