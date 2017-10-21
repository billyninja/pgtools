package colors

import (
    "fmt"
)

func Red(st interface{}) string {
    return fmt.Sprintf("\x1b[31;1m%v\x1b[0m", st)
}

func Green(st interface{}) string {
    return fmt.Sprintf("\x1b[32;1m%v\x1b[0m", st)
}

func Blue(st interface{}) string {
    return fmt.Sprintf("\x1b[34;1m%v\x1b[0m", st)
}

func Yellow(st interface{}) string {
    return fmt.Sprintf("\x1b[33;1m%v\x1b[0m", st)
}

func Bold(st interface{}) string {
    return fmt.Sprintf("\x1b[1m%v\x1b[0m", st)
}
