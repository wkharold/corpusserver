package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/wkharold/corpusserver/mailbox"
)

const MailDir = "/corpus/enron_mail_20110402/maildir"

func main() {
	var wg sync.WaitGroup

	filepath.Walk(MailDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(MailDir, path)
		if err != nil {
			return err
		}

		if strings.Contains(rel, string(filepath.Separator)) {
			return filepath.SkipDir
		}

		go mailbox.Catalog(MailDir, rel, &wg)

		return nil
	})

	wg.Wait()
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
