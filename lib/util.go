package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// func copyFile(dst, src string, perm os.FileMode) error {
// 	in, err := os.Open(src)
// 	if err != nil {
// 		return err
// 	}
// 	defer in.Close()
// 	tmp, err := TempFile(filepath.Dir(dst), "")
// 	if err != nil {
// 		return err
// 	}
// 	_, err = io.Copy(tmp, in)
// 	if err != nil {
// 		tmp.Close()
// 		os.Remove(tmp.Name())
// 		return err
// 	}
// 	if err = tmp.Close(); err != nil {
// 		os.Remove(tmp.Name())
// 		return err
// 	}
// 	if err = os.Chmod(tmp.Name(), perm); err != nil {
// 		os.Remove(tmp.Name())
// 		return err
// 	}
// 	return os.Rename(tmp.Name(), dst)
// }

func WriteFile(path string, data string) error {
	dir, _ := filepath.Split(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("unable to make dir %q: %w", dir, err)
	}

	err = ioutil.WriteFile(path, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("unable to write file %q: %w", path, err)
	}

	return nil
}

func RunCmd(cmd *exec.Cmd) ([]byte, error) {
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("error executing %q, %v", cmd.Path, cmd.Args)
		log.Printf(string(output))
		return nil, fmt.Errorf("error executing %q: %w", cmd.Path, err)
	}
	return output, nil
}

func Cp(src, dst string) error {
	cpCmd := exec.Command("cp", "-rf", src, dst)
	_, err := RunCmd(cpCmd)
	return err
}

func TerraformCmd(tfpath, wd string, args ...string) ([]byte, error) {
	cmd := exec.Command(tfpath, args...)
	cmd.Dir = wd

	return RunCmd(cmd)
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func ProviderShortName(n string) string {
	return strings.TrimPrefix(n, "terraform-provider-")
}

func resourceShortName(name, providerName string) string {
	psn := ProviderShortName(providerName)
	return strings.TrimPrefix(name, psn+"_")
}
