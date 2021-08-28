package bench

import (
	"time"
)

type Task struct {
	url      string
	num      int
	isError  bool
	duration time.Duration
}


