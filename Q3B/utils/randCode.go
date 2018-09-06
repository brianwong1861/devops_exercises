package utils

import (
        "time"
        _ "fmt"
        "math/rand"
)
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func init(){
        rand.Seed(time.Now().UnixNano())
}

func RandStringBytes(n int) string {

        b := make([]byte, n)
        for i := range b {
                b[i] = letterBytes[rand.Intn(len(letterBytes))]
        }
        return string(b)
}



