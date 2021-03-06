package saturate

import (
	"errors"
	"sync"
	"testing"
)

var errAFails = errors.New("i'm a.  I fail")

func TestSaturation(t *testing.T) {
	l := sync.Mutex{}
	task := map[string]int{}
	worker := map[string]int{}
	successes := map[string]int{}
	names := []string{"a", "b", "c", "d"}

	sat := New(names, func(name string) Worker {
		return WorkerFunc(func(i interface{}) error {
			l.Lock()
			defer l.Unlock()
			// Mark that we saw it
			worker[name]++

			// "a" always fails
			if name == "a" {
				return errAFails
			}
			// Mark that we did it
			x := i.(string)
			task[x]++
			successes[name]++

			return nil
		})
	}, nil)

	ch := make(chan WorkInput)
	go func() {
		defer close(ch)
		for i := 0; i < 100; i++ {
			for _, n := range names {
				ch <- WorkInput{n, names}
			}
		}
	}()

	err := sat.Saturate(ch)
	if err != nil {
		t.Fatalf("Error processing put: %v", err)
	}

	if successes["a"] != 0 {
		t.Fatalf("a never succeeded at anything: %v", successes)
	}

	total := 0
	for _, v := range task {
		total += v
	}

	if total != len(task)*100 {
		t.Fatalf("Expected total of %v, got %v: %v", len(task)*100, total, task)
	}

	t.Logf("%v/%v/%v/%v", task, worker, successes, total)
}

func TestSaturationFails(t *testing.T) {
	l := sync.Mutex{}
	m := map[string]int{}
	names := []string{"a", "b", "c", "d"}

	for _, n := range names {
		m[n] = 0
	}

	sat := New(names, func(name string) Worker {
		return WorkerFunc(func(i interface{}) error {
			// "a" always fails
			if name == "a" {
				return errAFails
			}

			l.Lock()
			defer l.Unlock()
			m[name]++

			return nil
		})
	}, nil)

	ch := make(chan WorkInput)
	go func() {
		for i := 0; i < 100; i++ {
			for j, n := range names {
				ch <- WorkInput{j, []string{n}}
			}
		}
		close(ch)
	}()

	err := sat.Saturate(ch)
	if err == nil {
		t.Fatalf("Expected failure.  Got none.")
	}

	t.Logf("Got expected failure: %v", err)
}
