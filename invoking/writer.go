package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	log "github.com/cihub/seelog"
	"github.com/vrecan/life"
)

// Generator generates things to write
type Generator struct {
	*life.Life
	Out chan<- int
}

// Writer writes things to files
type Writer struct {
	*life.Life
	In    <-chan int
	FName string
}

// NewGenerator creates a new generator
func NewGenerator(out chan<- int) Generator {
	g := Generator{
		Life: life.NewLife(),
		Out:  out,
	}
	g.Life.SetRun(g.generateNums)
	return g
}

// Close satisfies the io.Closer interface
func (g Generator) Close() error {
	close(g.Out)
	return nil
}

// GenerateNums generates random number until told to stop
func (g Generator) generateNums() {
	mT := time.NewTicker(1 * time.Millisecond)
	defer mT.Stop()

	for {
		select {
		case <-g.Life.Done:
			return
		case <-mT.C:
			r := rand.Int()
			g.Out <- r
		}
	}
}

// NewWriter creates a writer for a file
func NewWriter(in <-chan int, name string) Writer {
	w := Writer{
		Life:  life.NewLife(),
		In:    in,
		FName: name,
	}
	w.SetRun(w.writeNums)
	return w
}

// Close satisfies the io.Closer interface
func (w Writer) Close() error {
	return nil
}

// writeNums writes numbers to a file until told to stop
func (w Writer) writeNums() {
	for {
		select {
		case num := <-w.In:
			f, err := os.Open(w.FName)
			if os.IsNotExist(err) {
				f, err = os.Create(w.FName)
				if err != nil {
					log.Errorf("Error creating %s: %v", w.FName, err)
				}
			} else if err != nil {
				log.Errorf("Error opening: %s: %v", w.FName, err)
				continue
			}

			numStr := strconv.Itoa(num)
			f.Write([]byte(fmt.Sprintf("Number is now: %s\n", numStr)))
			err = f.Close()
			if err != nil {
				log.Error("Error Closing file: ", err)
			}
		case <-w.Life.Done:
			return
		}
	}
}
