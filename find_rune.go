package main

import (
    "path/filepath"
    "fmt"
    "os"
    "io/ioutil"
    "runtime"
    "sync"
    "strings"
)

var workers = runtime.NumCPU()

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    if len(os.Args)!=3 || os.Args[1]=="-h" || os.Args[1]=="--help" {
        fmt.Printf("usage: %s <file1> <chinese>\n",filepath.Base(os.Args[0]))
        return
    }
    filename1:=os.Args[1]
    filename2:=os.Args[2]
    info,err:=os.Stat(filename1)
    if err!=nil {
        fmt.Printf("failed to open the file:%s ,%s\n",filename1,err)
        return
    }
    data,err:=ioutil.ReadFile(filename2)
    if err!=nil {
        fmt.Printf("failed to open the file:%s ,%s\n",filename2,err)
        return
    }
    chunkSize:=info.Size() / int64(workers)
    done:=sync.WaitGroup{}
    done.Add(workers)
    results:=make([][]rune,workers)
    for idx:=0;idx<workers;idx++ {
        offset:=int64(idx) * chunkSize
        if idx+1==workers {
            chunkSize*=2
        }
        go processLines(&results[idx],string(data),filename1,offset,chunkSize,&done)
    }
    done.Wait()
    all_data:=map[rune]bool{}
    for idx:=0;idx<workers;idx++ {
        result_unit:=results[idx]
        for _,character:=range result_unit {
            if all_data[character] {
                continue
            }
            all_data[character]=true
            fmt.Print(string(character))
        }
    }
}

func processLines(results *[]rune,compare_data string,filename string,offset,chunkSize int64,done *sync.WaitGroup){
    defer done.Done()
    file,err:=os.Open(filename)
    if err!=nil {
        fmt.Printf("failed to open the file:%s ,%s\n",filename,err)
    }
    defer file.Close()
    file.Seek(offset,0)
    data:=make([]byte,1024)
    len:=chunkSize/1024
    left:=chunkSize%1024
    var read_len int
    for idx:=int64(0);idx<len;idx++ {
        read_len,err=file.Read(data)
        characters:=[]rune(string(data[:read_len]))
        for _,character:=range characters {
            if strings.ContainsRune(compare_data,character) {
                continue
            }
            *results=append(*results,character)
        }
        if err!=nil {
            break
        }
    }
    if left>0 {
        read_len,err=file.Read(data)
        characters:=[]rune(string(data[:read_len]))
        for _,character:=range characters {
            if strings.ContainsRune(compare_data,character) {
                continue
            }
            *results=append(*results,character)
        }
    }
}
