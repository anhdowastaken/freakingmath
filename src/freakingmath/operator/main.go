package operator

import (
    "sync"
    "math/rand"
    "time"
)

type Arithmetic int
const (
    ARITHMETIC_START = iota
    PLUS
    MINUS
    TIMES
    ARITHMETIC_END
)

type Operator struct {

}

var instance *Operator
var once sync.Once

func New() *Operator {
    once.Do(func() {
        instance = &Operator{}
    })

    return instance
}

func (o Operator) To_string(a Arithmetic) (string) {
    switch a {
    case PLUS:
        return "+"
    case MINUS:
        return "-"
    case TIMES:
        return "x"
    default:
        return ""
    }
}

func (o Operator) Random() (Arithmetic) {
    min := ARITHMETIC_START + 1
    max := ARITHMETIC_END

    rand.Seed(time.Now().Unix())

    return Arithmetic(rand.Intn(max - min) + min)
}
