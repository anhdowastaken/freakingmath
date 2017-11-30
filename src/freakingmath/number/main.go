package number

import (
    "sync"
    "math"
    "math/rand"
    "time"
    // "fmt"
)

type Number struct {
    min int
    max int
}

var instance *Number
var once sync.Once

func New() *Number {
    once.Do(func() {
        instance = &Number{}
    })

    return instance
}

func (n *Number) Set_level(level int) {
    if level > 0 {
        if level == 1 {
            n.min = 0
        } else {
            n.min = int(math.Pow10(level - 1))
        }
        n.max = int(math.Pow10(level))
    } else {
        return
    }
}

func (n Number) Random(max int, min int) (int) {
    if max == min {
        return max
    }

    var tmp int
    if min > max {
        min = tmp
        min = max
        max = tmp
    }

    rand.Seed(time.Now().UTC().UnixNano())

    return (rand.Intn((max + 1) - min) + min)
}

func (n Number) Random_default() (int) {
    return n.Random(n.max, n.min)
}

func (n Number) Random_around(i int, around int) (int) {
    min := i - around
    max := i + around

    return n.Random(max, min)
}
