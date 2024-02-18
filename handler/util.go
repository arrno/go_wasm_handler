package wasmhandler

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/jaevor/go-nanoid"
)

type FaultCode int

const (
	Clear FaultCode = iota
	UserErr
	ServerErr
)

func WorkingDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path := strings.Split(dir, "/")
	if len(path) == 1 {
		return path[0], nil
	}
	return strings.Join(path[:len(path)-1], "/"), nil
}

func NewNanoid() (string, error) {
	if canonicID, err := nanoid.Standard(21); err != nil {
		return "", err
	} else {
		return canonicID(), nil
	}
}

func RemoveDir(uid string) error {
	return os.RemoveAll(uid)
}

func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func CopyDir(src string, dst string) error {
	var err error
	var fds []fs.DirEntry
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = os.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if !fd.IsDir() {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

// <----------- proc ----------->

type Proc struct {
	uid        string // unique ID
	program    string
	workingDir string
}

func NewProc(program string) (*Proc, error) {
	p := new(Proc)
	var err error
	p.uid, err = NewNanoid()
	if err != nil {
		return p, err
	}
	p.program = program
	p.workingDir, err = WorkingDir()
	if err != nil {
		return p, err
	}
	return p, nil

}

func (p *Proc) Inject() error {
	blob, err := os.ReadFile(fmt.Sprintf("%s/%s/%s/template.txt", p.workingDir, "temp", p.uid))
	if err != nil {
		return err
	}
	sblob := string(blob)
	begin := strings.Split(sblob, "// <-- begin run -->")[0]
	end := strings.Split(sblob, "// <-- end run -->")[1]
	injected := begin + strings.Replace(p.program, "func main() {", "func run() {", 1) + end
	return os.WriteFile(fmt.Sprintf("%s/%s/%s/hosted.go", p.workingDir, "temp", p.uid), []byte(injected), 0644)
}

func (p *Proc) Compile() error {
	cmd := exec.Command("go", "fmt", ".")
	cmd.Dir = fmt.Sprintf("%s/%s/%s/", p.workingDir, "temp", p.uid)
	if _, err := cmd.Output(); err != nil {
		return err
	}
	cmd = exec.Command("sh", "compile.sh")
	cmd.Dir = fmt.Sprintf("%s/%s/%s/", p.workingDir, "temp", p.uid)
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}

func (p *Proc) Upload() error {
	blob, err := os.ReadFile(fmt.Sprintf("%s/%s/%s/main.wasm", p.workingDir, "temp", p.uid))
	if err != nil {
		return err
	}
	return bucketWriter.WriteToBucket(p.uid, blob)
}

// DO ALL
func (p *Proc) DoProcess() (string, FaultCode, error) {

	if err := CopyDir(p.workingDir+"/wasm", p.workingDir+"/temp/"+p.uid); err != nil {
		return "", ServerErr, err
	}

	defer func() {
		go func() {
			if err := RemoveDir(p.workingDir + "/temp/" + p.uid); err != nil {
				fmt.Println("Failed to delete " + p.uid + ". Err: " + err.Error())
			}
		}()
	}()

	if err := p.Inject(); err != nil {
		return "", ServerErr, err
	}

	if err := p.Compile(); err != nil {
		return "", UserErr, err
	}

	if err := p.Upload(); err != nil {
		return "", ServerErr, err
	}

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s/%s/main.wasm", BUCKET_ROOT, WASM_FOLDER, p.uid), Clear, nil
}
