package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func calcHash(dirPath string, hashMap []FileLen) ([]FileLen, error) {
	err := filepath.Walk(dirPath, func(filePath string, f os.FileInfo, err error) error {
		if f == nil {
			return nil
		}

		if f.IsDir() {
			return nil
		}

		if strings.HasSuffix(filePath, ".DS_Store") {
			return nil
		}

		hashMap = append(hashMap, FileLen{filePath, f.Size()})
		return nil
	})

	return hashMap, err
}

func hashFile(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		return nil, err
	}

	return md5hash.Sum(nil), nil
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func snapshotMetaFile(hashMap map[string]int64, metaFilePath string) {
	b, err := json.Marshal(hashMap)
	if err != nil {
		fmt.Sprintln("json.Marshal err", err)
		return
	}
	ioutil.WriteFile(metaFilePath, b, 0666)
}

type FileLen struct {
	Name string
	Len  int64
}

type FileLenSlice []FileLen

func (t FileLenSlice) Len() int { // 重写 Len() 方法
	return len(t)
}
func (t FileLenSlice) Swap(i, j int) { // 重写 Swap() 方法
	t[i], t[j] = t[j], t[i]
}
func (t FileLenSlice) Less(i, j int) bool { // 重写 Less() 方法， 从小到大
	return t[i].Len < t[j].Len
}

func checkDup(hashMap []FileLen) {
	var f1 *FileLen
	md5Map := make(map[string]string)
	pathMap := make(map[string]string)
	for i, _ := range hashMap {
		if f1 == nil || f1.Len != hashMap[i].Len {
			f1 = &hashMap[i]
			continue
		}

		//如果f1的md5没算过计算一下
		if _, ok := pathMap[f1.Name]; !ok {
			b, _ := hashFile(f1.Name)
			md5Str := fmt.Sprintf("%x", b)
			pathMap[f1.Name] = md5Str
			md5Map[md5Str] = f1.Name
		}

		//如果hashMap[i]的md5没算过计算一下
		md5Str, ok := pathMap[hashMap[i].Name]
		if !ok {
			b, _ := hashFile(f1.Name)
			md5Str = fmt.Sprintf("%x", b)
			pathMap[hashMap[i].Name] = md5Str
		}

		if _, ok = md5Map[md5Str]; !ok {
			md5Map[md5Str] = hashMap[i].Name
			continue
		}

		//重复
		fmt.Printf("%s,%s,%d\n", md5Map[md5Str], hashMap[i].Name, hashMap[i].Len)
		os.Remove(hashMap[i].Name)
	}
}

func main() {
	dirPath := "/Volumes/D2/MOVIE/AV/"
	hashMap := make([]FileLen, 0)
	hashMap, err := calcHash(dirPath, hashMap)
	if err != nil {
		fmt.Sprintln("calcHash err", err)
		return
	}

	sort.Sort(FileLenSlice(hashMap))
	checkDup(hashMap)
	//fmt.Printf("%v\n", hashMap)
	return
}
