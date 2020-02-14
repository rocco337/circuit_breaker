package main

import "fmt"
import "time"

import "errors"

const failureTreshold = 5
const resetTimeout = 5 * time.Second

//CircuitBreaker ...
type CircuitBreaker struct {
	IsOpen          bool
	ErrorsCount     int
	LastFailureTime time.Time
}

//Call ...
func (breaker *CircuitBreaker) Call(call func() error) error {
	now := time.Now()

	if breaker.IsOpen && now.Sub(breaker.LastFailureTime).Seconds() >= resetTimeout.Seconds() {
		breaker.close()
	}

	if !breaker.IsOpen {
		err := call()
		if err != nil {
			breaker.ErrorsCount++
			breaker.LastFailureTime = now
			if breaker.ErrorsCount >= failureTreshold {
				breaker.open()
			}

			return err
		} else {
			return nil
		}

	} else {
		fmt.Println("Breaker is open, cannot process calls")
		return errors.New("Breaker is open, cannot process calls!")
	}

}

func (breaker *CircuitBreaker) open() {
	breaker.IsOpen = true
}

func (breaker *CircuitBreaker) close() {
	breaker.IsOpen = false
	breaker.ErrorsCount = 0
	breaker.LastFailureTime = time.Time{}
}
