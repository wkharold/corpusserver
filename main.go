package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/mail"
	"os"
	"path/filepath"
)

func main() {
	filepath.Walk("/corpus/enron_mail_20110402/maildir", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			rdr, err := os.Open(path)
			if err == nil {
				msg, err := Format(rdr)
				if err == nil {
					fmt.Println(msg)
				}
			}
		}
		return nil
	})
}

func Format(msgrdr io.Reader) (string, error) {
	msg, err := mail.ReadMessage(msgrdr)
	if err != nil {
		return "", err
	}

	mb := new(bytes.Buffer)

	scanner := bufio.NewScanner(msg.Body)
	for scanner.Scan() {
		mb.WriteString(scanner.Text())
	}

	mm := make(map[string]interface{}, 2)
	mm["body"] = mb.String()

	hm := make(map[string]string)
	for k, v := range msg.Header {
		if len(v[0]) > 0 {
			hm[k] = v[0]
		}
	}

	mm["headers"] = hm

	b, err := json.MarshalIndent(mm, "", "   ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}
