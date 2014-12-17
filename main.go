package main

import (
	"bufio"
	"fmt"
	"net/mail"
	"os"
)

func main() {
	mf, err := os.Open("/corpus/enron_mail_20110402/maildir/allen-p/_sent_mail/1.")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	msg, err := mail.ReadMessage(mf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(msg.Header)
	scanner := bufio.NewScanner(msg.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
