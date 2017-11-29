package main

import (
    "freakingmath/operator"
    "freakingmath/number"
    "fmt"
    "os"
    "os/exec"
    "syscall"
    "os/signal"
    "time"
)

const DEFAULT_LEVEL = 1
const DEFAULT_TIMEOUT = 5

func calculate_result(left int, right int, op operator.Arithmetic) (int) {
    switch op {
    case operator.PLUS:
        return (left + right)
    case operator.MINUS:
        return (left - right)
    case operator.TIMES:
        return (left * right)
    default:
        return 0
    }
}

func create_result() (int) {
    return 0
}

func main() {
    defer func() {
        // Display entered characters on the screen
        exec.Command("stty", "-F", "/dev/tty", "echo").Run()
    }()

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        sig := <-sigs
        if sig == syscall.SIGINT || sig == syscall.SIGTERM {
            // Display entered characters on the screen
            exec.Command("stty", "-F", "/dev/tty", "echo").Run()
            os.Exit(0)
        }
    }()

    // Disable input buffering
    exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
    // Do not display entered characters on the screen
    exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

    fmt.Println("FREAKING MATH")
    fmt.Println("Enter t/T (true) or f/F (false) to answer")

    o := operator.New()
    n := number.New()
    n.Set_level(2)

    score := 0

    for {
        left := n.Random()
        right := n.Random()
        op := o.Random()
        result := calculate_result(left, right, op)
        actual := n.Random_around(result)

        fmt.Printf("%d %s %d = %d? ",
            left,
            o.To_string(op),
            right,
            actual)

        expected := (result == actual)
        input := make(chan bool, 1)
        expired := make(chan bool, 1)

        go func(flag chan bool) {
            time.Sleep(time.Second * DEFAULT_TIMEOUT)
            flag <- true
        } (expired)

        go func() {
            var b []byte = make([]byte, 1)
            for {
                os.Stdin.Read(b)
                fmt.Println(string(b))
                switch string(b) {
                case "t":
                    fallthrough
                case "T":
                    input <- true
                    expired <- false
                    return
                case "f":
                    fallthrough
                case "F":
                    input <- false
                    expired <- false
                    return
                default:
                    fmt.Println("Your score:", score)
                    sigs <- syscall.SIGTERM
                    return
                }
            }
        }()

        is_expired := <-expired
        if is_expired == true {
            fmt.Println("\nTIMEOUT!")
            fmt.Println("Your score:", score)
            break
        }

        answer := <-input
        if answer != expected {
            fmt.Println("\nWRONG!")
            fmt.Println("Your score:", score)
            break
        } else {
            score = score + 1
        }
    }
}
