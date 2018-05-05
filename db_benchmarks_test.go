package go_database_bench

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

type DB interface {
	Put(*Person) error
	Get(id []byte) (*Person, error)
	Iter(chan<- Person) error
	Count() (int, error)
	Drop() error
}

type Person struct {
	Id       []byte `bow:"key" storm:"id"`
	Name     string
	BirthDay time.Time
	Phone    string
	Siblings int
	Spouse   bool
	Money    float64
}

func randString(l int) string {
	buf := make([]byte, l)
	for i := 0; i < (l+1)/2; i++ {
		buf[i] = byte(rand.Intn(256))
	}
	return fmt.Sprintf("%x", buf)[:l]
}

func generate() *Person {
	return &Person{
		Id:       []byte(randString(32)),
		Name:     randString(16),
		BirthDay: time.Now(),
		Phone:    randString(10),
		Siblings: rand.Intn(5),
		Spouse:   rand.Intn(2) == 1,
		Money:    rand.Float64(),
	}
}

func benchPut(b *testing.B, db DB) {
	b.ReportAllocs()
	var n int
	for n = 0; n < b.N; n++ {
		err := db.Put(generate())
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	count, err := db.Count()
	if err != nil {
		b.Fatal(err)
	}
	if n != count {
		b.Fatalf("inserted %d, counted %d", n, count)
	}
}

func benchGet(b *testing.B, db DB) {
	b.StopTimer()
	ids := make([][]byte, 1000)
	for i := 0; i < 1000; i++ {
		p := generate()
		ids[i] = p.Id
		db.Put(p)
	}
	b.ReportAllocs()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		id := ids[rand.Intn(len(ids))]
		p, err := db.Get(id)
		if err != nil {
			b.Fatal(err)
		}
		if !bytes.Equal(id, p.Id) {
			b.Fatalf("expected id %x, got %x", id, p.Id)
		}
	}
}

func benchIter(b *testing.B, db DB) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		ids := make(map[string]bool)
		for i := 0; i < 100; i++ {
			p := generate()
			ids[string(p.Id)] = false
			db.Put(p)
		}
		c := make(chan Person)
		go func() {
			err := db.Iter(c)
			if err != nil {
				panic(err)
			}
		}()
		b.StartTimer()
		for p := range c {
			ids[string(p.Id)] = true
		}
		b.StopTimer()
		for id, ok := range ids {
			if !ok {
				b.Fatalf("id %x not found", id)
			}
		}
	}
}

// tempfile returns a temporary file path.
func tempfile(prefix string) string {
	f, err := ioutil.TempFile("", prefix)
	if err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	if err := os.Remove(f.Name()); err != nil {
		panic(err)
	}
	return f.Name()
}
