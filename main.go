package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/mail"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/wkharold/corpusserver/catalog"
)

const MailDir = "/corpus/enron_mail_20110402/maildir"

func main() {
	mc := make(chan catalog.MbxMsg)
	done := make(chan struct{})

	go catalog.Cataloger(mc, done)

	mc <- catalog.MbxMsg{Mailbox: "lay-k", Folder: "inbox", Msgfile: "1."}
	done <- struct{}{}
}

func mainer() {
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

		return IndexMailbox(rel)
	})
}

func IndexMailbox(mbx string) error {
	mbxpath := path.Join(MailDir, mbx)
	filepath.Walk(mbxpath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(mbxpath, path)
		if err != nil {
			return err
		}

		if strings.Contains(rel, string(filepath.Separator)) {
			return filepath.SkipDir
		}

		if rel != "." {
			fmt.Printf("%s::%s\n", mbxpath, rel)
		}
		return nil
	})

	return nil
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
