package dao

import (
	"fmt"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	InitDb()
	m.Run()
}

func TestXxx(t *testing.T) {
	fmt.Println( time.Now())

}
