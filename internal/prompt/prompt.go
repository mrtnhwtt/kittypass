package prompt

import (
    "fmt"
    "os"
    "syscall"

    "golang.org/x/term"
)

// PasswordPrompt asks for a string value using the label.
// The entered value will not be displayed on the screen
// while typing.
func PasswordPrompt(label string) string {
    var s string
    fmt.Fprint(os.Stderr, label+" ")
    for {
        b, _ := term.ReadPassword(int(syscall.Stdin))
        s = string(b)
        if s != "" {
            break
        }
    }
    fmt.Println()
    return s
}