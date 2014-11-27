package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	T_MAX               = 10000
	Q_MAX               = 1000
	N_MAX               = 10000
	X_MIN               = 0.0
	X_MAX               = 1000000.0
	Y_MIN               = 0.0
	Y_MAX               = 1000000.0
	ID_MIN              = 0
	ID_MAX              = 100000
	QUESTION_TOPICS_MAX = 10
)

var (
	T = T_MAX / 1
	Q = Q_MAX / 1
	N = N_MAX / 1
)

func RandomPoint() []float64 {
	return []float64{rand.Float64(), rand.Float64()}
}

func RandomTopicIds(n int, m int) []string {
	ids := make([]string, 0)
	for i := 0; i < n; i++ {
		ids = append(ids, strconv.Itoa(rand.Intn(m)))
	}
	return ids
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	rand.Seed(time.Now().Unix())
	fmt.Println(T, Q, N)
	for i := 0; i < T; i++ {
		p := RandomPoint()
		fmt.Printf("%d %.2f %.2f\n", i, p[0], p[1])
	}

	for i := 0; i < Q; i++ {
		numTopics := rand.Intn(QUESTION_TOPICS_MAX)
		topicIds := RandomTopicIds(numTopics, T)
		if numTopics == 0 {
			fmt.Printf("%d %d\n", i, numTopics)
		} else {
			fmt.Printf("%d %d %s\n", i, numTopics, strings.Join(topicIds, " "))
		}

	}

	choices := []string{"t", "q"}
	for i := 0; i < N; i++ {
		p := RandomPoint()
		fmt.Printf("%s %d %.2f %.2f\n", choices[random(0, 2)], rand.Intn(100), p[0], p[1])
	}

}
