package mailbox

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/wkharold/corpusserver/catalog"
)

func Catalog(rootdir, mbx string, wg *sync.WaitGroup) {
	defer wg.Done()

	wg.Add(1)

	mbxpath := path.Join(rootdir, mbx)
	filepath.Walk(mbxpath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		rp, err := filepath.Rel(mbxpath, path)
		if err != nil {
			return err
		}

		if strings.Contains(rp, string(filepath.Separator)) {
			return filepath.SkipDir
		}

		if rp == "." {
			return nil
		}

		if err := catalogFolder(rootdir, mbx, rp); err != nil {
			return err
		}

		return nil
	})
}

func catalogFolder(rootdir, mbx, folder string) error {
	folderpath := path.Join(path.Join(rootdir, mbx), folder)
	filepath.Walk(folderpath, func(msgpath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(folderpath, msgpath)
		if err != nil {
			return err
		}

		if strings.Contains(rel, string(filepath.Separator)) {
			return nil
		}

		catalog.Msgs <- catalog.MbxMsg{Mailbox: path.Join(rootdir, mbx), Folder: folder, Msgfile: rel}

		return nil
	})

	return nil
}
