package go_database_bench

import (
	"os"
	"testing"

	"github.com/zippoxer/bow"
)

type Bow struct {
	db  *bow.DB
	dir string
}

func NewBow() *Bow {
	dir := tempfile("bow-")
	db, err := bow.Open(dir)
	if err != nil {
		panic(err)
	}
	return &Bow{db: db, dir: dir}
}

func (b *Bow) Put(p *Person) error {
	return b.db.Bucket("people").Put(p)
}

func (b *Bow) Get(id []byte) (*Person, error) {
	var p Person
	err := b.db.Bucket("people").Get(id, &p)
	return &p, err
}

func (b *Bow) Iter(c chan<- Person) error {
	defer close(c)
	iter := b.db.Bucket("people").Iter()
	var p Person
	for iter.Next(&p) {
		c <- p
	}
	return iter.Err()
}

func (b *Bow) Count() (int, error) {
	n := 0
	iter := b.db.Bucket("people").Iter()
	var p Person
	for iter.Next(&p) {
		n++
	}
	return n, iter.Err()
}

func (b *Bow) Drop() error {
	defer os.RemoveAll(b.dir)
	return b.db.Close()
}

func BenchmarkBowPut(b *testing.B) {
	b.StopTimer()
	bow := NewBow()
	defer bow.Drop()
	b.StartTimer()
	benchPut(b, bow)
}

func BenchmarkBowGet(b *testing.B) {
	b.StopTimer()
	bow := NewBow()
	defer bow.Drop()
	b.StartTimer()
	benchGet(b, bow)
}

func BenchmarkBowIter(b *testing.B) {
	b.StopTimer()
	bow := NewBow()
	defer bow.Drop()
	b.StartTimer()
	benchIter(b, bow)
}
