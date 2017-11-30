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
    "runtime"
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

    // We create a list of choices
    // 1. Correct choice
    // 2. Wrong choice but somewhere near the answer
    // 3. Random between 0 and the anser
    choices = append(choices, expected)
    choices = append(choices, n.Random_around(expected, DEFAULT_AROUND))
    choices = append(choices, n.Random(expected, 0))

    i := n.Random(2, 0)

    return choices[i]
}

func disable_display_character() {
    // stty commands on Mac OS and Linux a different
    if runtime.GOOS == "darwin" {
        exec.Command("stty", "-f", "/dev/tty", "cbreak", "min", "1").Run()
        exec.Command("stty", "-f", "/dev/tty", "-echo").Run()
    } else {
        exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
        exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
    }
}

func enable_display_character() {
    if runtime.GOOS == "darwin" {
        exec.Command("stty", "-f", "/dev/tty", "echo").Run()
    } else {
        exec.Command("stty", "-F", "/dev/tty", "echo").Run()
    }
}

func main() {
    // Restore setting of terminal finally
    defer func() {
        enable_display_character()
    }()

    // Handle SIGINT and SIGTERM signals to restore setting of terminal
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        sig := <-sigs
        if sig == syscall.SIGINT || sig == syscall.SIGTERM {
            enable_display_character()
            os.Exit(0)
        }
    }()

    fmt.Println("FREAKING MATH")
    fmt.Println("-------------------------")
    fmt.Printf("Enter level (1, 2... etc): ")
    reader := bufio.NewReader(os.Stdin)
    text, _ := reader.ReadString('\n')
    level, err := strconv.Atoi(text[:len(text) - 1])

    // We only accept an integer number which is greater than 0
    for err != nil || level <= 0 {
        fmt.Printf("Enter level (1, 2... etc): ")
        text, _ = reader.ReadString('\n')
        level, err = strconv.Atoi(text[:len(text) - 1])
    }

    o := operator.New()
    n := number.New()
    n.Set_level(level)

    disable_display_character()

    fmt.Println("-------------------------")
    fmt.Println("Enter t/T (true) or f/F (false) to answer")
    fmt.Println("-------------------------")
    score := 0

    timer := time.NewTimer(time.Duration(DEFAULT_TIMEOUT) * time.Second)
    continue_looping := true

    for continue_looping == true {
        op := o.Random()
        left := n.Random_default()
        // If one number is greater than 10 and arithmetic is TIMES, the other must be less than 10
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

        timer.Reset(time.Duration(DEFAULT_TIMEOUT) * time.Second)
        expected := (result == actual)
        input := make(chan bool, 1)

        // We only accept answer as a single character
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
                    return
                case "f":
                    fallthrough
                case "F":
                    input <- false
                    return
                default:
                    // Quit if input character is not t/T or f/F
                    fmt.Println("Your score:", score)
                    sigs <- syscall.SIGTERM
                    continue_looping = false
                    return
                }
            }
        }()

        select {
        case <-timer.C:
            fmt.Println("\nTIMEOUT!")
            fmt.Println("Your score:", score)
            sigs <- syscall.SIGTERM
            continue_looping = false

        case answer := <-input:
            if answer != expected {
                fmt.Println("\nWRONG!")
                fmt.Println("Your score:", score)
                sigs <- syscall.SIGTERM
                continue_looping = false
            } else {
                score = score + 1
            }
        }
    }
}
