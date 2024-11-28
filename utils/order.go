package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateOrderNumber() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d%d", 
		time.Now().Unix()%1000000,
		rand.Intn(1000))
}