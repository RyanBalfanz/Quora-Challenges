test:
	go run main.go < sample_input.txt

random_problem:
	go run probgen/main.go > generated_input.txt

test_random_problem: random_problem
	go run main.go < generated_input.txt