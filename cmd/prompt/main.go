package main

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
)

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    client  := &http.Client{}

    fmt.Println("chronocache prompt — type your query, ctrl+c to exit\n")

    for {
        fmt.Print("> ")
        if !scanner.Scan() { break }

        query := strings.TrimSpace(scanner.Text())
        if query == "" { continue }

        body, _ := json.Marshal(map[string]any{
            "query": query,
        })

        resp, err := client.Post(
            "http://localhost:8080/v1/messages",
            "application/json",
            bytes.NewReader(body),
        )
        if err != nil {
            fmt.Println("error:", err)
            continue
        }

        respBody, _ := io.ReadAll(resp.Body)
        resp.Body.Close()

        fmt.Println(string(respBody))
        fmt.Println()
    }
}