package main

import (
	"errors"
	"testing"
	"time"
)

func TestCircuitBreaker(t *testing.T) {

	t.Run("Successfull calls - should pass", func(t *testing.T) {
		cicuitBreaker := new(CircuitBreaker)
		for i := 0; i < 10; i++ {
			cicuitBreaker.Call(func() error {
				return nil
			})

			if cicuitBreaker.IsOpen {
				t.Error("Circuit should be closed")
			}

			if cicuitBreaker.ErrorsCount > 0 {
				t.Error("Should not be any errors")
			}
		}
	})

	t.Run("Unsuccessfull calls - after 5 calls, circuit should open", func(t *testing.T) {
		cicuitBreaker := new(CircuitBreaker)

		for i := 0; i < 5; i++ {
			cicuitBreaker.Call(func() error {
				return errors.New("Error")
			})
		}

		if !cicuitBreaker.IsOpen {
			t.Error("Circuit should be open")
		}

		if cicuitBreaker.ErrorsCount != 5 {
			t.Error("Should have 5 errors")
		}
	})

	t.Run("After 5 unsuccessfull calls, it should open, than after timout it should close", func(t *testing.T) {
		cicuitBreaker := new(CircuitBreaker)

		for i := 0; i <= 6; i++ {
			if i <= 5 {
				cicuitBreaker.Call(func() error {
					return errors.New("Error")
				})
			}

			if i == 5 {
				if !cicuitBreaker.IsOpen {
					t.Error("Circuit should be open")
				}

				if cicuitBreaker.ErrorsCount != 5 {
					t.Error("Should have 5 errors")
				}
			}

			if i == 6 {
				time.Sleep(5 * time.Second)
				cicuitBreaker.Call(func() error {
					return nil
				})

				if cicuitBreaker.IsOpen {
					t.Error("Circuit should be closed")
				}

				if cicuitBreaker.ErrorsCount != 0 {
					t.Error("Should have 0 errors")
				}
			}

		}

	})

}
