package number

import (
    "sync"
    "math"
    "math/rand"
    "time"
    "fmt"
)

type Number struct {
    min int
    max int
}

const AROUND_NUMBER = 3

var instance *Number
var once sync.Once

func New() *Number {
    once.Do(func() {
        instance = &Number{}
    })

    return instance
}

func (o *Number) Set_level(level int) {
    if level > 0 {
        if level == 1 {
            o.min = 0
        } else {
            o.min = int(math.Pow10(level - 1))
        }
        o.max = int(math.Pow10(level))
    } else {
        return
    }
}

func (o Number) To_string(n int) (string) {
    return fmt.Sprintf("%d", n)
}

func (o Number) Random() (int) {
    rand.Seed(time.Now().UTC().UnixNano())

    return (rand.Intn(o.max - o.min) + o.min)
}

func (o Number) Random_around(n int) (int) {
    min := n - AROUND_NUMBER
    max := n + AROUND_NUMBER

    rand.Seed(time.Now().UTC().UnixNano())

    return (rand.Intn(max - min) + min)
}
