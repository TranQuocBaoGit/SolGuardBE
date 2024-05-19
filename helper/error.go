package helper

import (
	"fmt"
	"log"
)


func CheckError(err error){
    if err != nil{
        log.Fatal(err)
    }
}

func MakeError(err error, typ string) error {
	return fmt.Errorf("error type \"%s\": %v", typ, err)
}