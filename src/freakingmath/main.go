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
    "strconv"
    "bufio"
)

const DEFAULT_LEVEL = 1
const DEFAULT_AROUND = 3
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

func random_result(expected int) (int) {
    n := number.New()
    choices := make([]int, 0)

    choices = append(choices, expected)
    choices = append(choices, n.Random_around(expected, DEFAULT_AROUND))
    choices = append(choices, n.Random(expected, 0))

    i := n.Random(2, 0)

    return choices[i]
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

    fmt.Println("FREAKING MATH")
    fmt.Println("-------------------------")
    fmt.Printf("Enter level (1, 2... etc): ")
    reader := bufio.NewReader(os.Stdin)
    text, _ := reader.ReadString('\n')
    level, err := strconv.Atoi(text[:len(text) - 1])

    for err != nil || level <= 0 {
        fmt.Printf("Enter level (1, 2... etc): ")
        text, _ = reader.ReadString('\n')
        level, err = strconv.Atoi(text[:len(text) - 1])
    }

    o := operator.New()
    n := number.New()
    n.Set_level(level)

    // Disable input buffering
    exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
    // Do not display entered characters on the screen
    exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

    fmt.Println("-------------------------")
    fmt.Println("Enter t/T (true) or f/F (false) to answer")
    fmt.Println("-------------------------")
    score := 0

    for {
        op := o.Random()
        left := n.Random_default()
        var right int
        if level > 1 && op == operator.TIMES && left > 10 {
            n.Set_level(1)
            right = n.Random_default()
            n.Set_level(level)
        } else {
            right = n.Random_default()
        }
        result := calculate_result(left, right, op)
        actual := random_result(result)

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
        }(expired)

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
