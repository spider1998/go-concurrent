//虽然系统中引入队列非常有用，但他通常是优化程序时希望采用的最后一种技术之一
//队列几乎不会加快程序的总运行时间，他只会让程序的行为有所不同
//--它并不是减少一个stage（或进程）的运行时间，所以队列的真正作用是将各个
//--stage分离，使其运行时间互不影响，从而解偶stage，级联改变整个系统的运行时行为
package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)


//时间：2041ns （无缓冲）
func BenchmarkUnBufferedWrite(b *testing.B)  {
	performWrite(b,tmpFileOrFatal())
}

//时间：597ns (有缓冲)
func BenchmarkBufferedWrite(b *testing.B)  {
	bufferredFile := bufio.NewWriter(tmpFileOrFatal())
	performWrite(b,bufio.NewWriter(bufferredFile))
}

func tmpFileOrFatal() *os.File {
	file,err := ioutil.TempFile("","tmp")
	if err != nil{
		log.Fatalf("error: %v",err)
	}
	return file
}

func performWrite(b *testing.B,write io.Writer)  {

	repeat := func(done <- chan interface{},values ...interface{})<-chan interface {}{
		valsStream := make(chan interface{})
		go func() {
			for {
				for _,v := range values{
					select {
					case <- done:
						return
					case valsStream <- v:

					}
				}
			}
		}()
		return valsStream
	}

	take := func(done <- chan interface{}, valueStream <-chan interface{},num int) <-chan interface{}{
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i:=0;i<num;i++{
				select {
				case <-done:
					return
				case takeStream <- <- valueStream:
				}
			}
		}()
		return takeStream
	}

	done := make(chan interface{})
	defer close(done)

	b.ResetTimer()
	for bt := range take(done,repeat(done,byte(0)),b.N){
		write.Write([]byte{bt.(byte)})
	}
}