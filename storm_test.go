package go_database_bench

import (
	"os"
	"testing"

	"github.com/asdine/storm"
)

type Storm struct {
	db       *storm.DB
	filename string
}

func NewStorm() *Storm {
	filename := tempfile("storm-")
	db, err := storm.Open(filename)
	if err != nil {
		panic(err)
	}
	return &Storm{db: db, filename: filename}
}
func (s *Storm) Put(p *Person) error {
	return s.db.Set("Person", p.Id, p)
}

func (s *Storm) Get(id []byte) (*Person, error) {
	var p Person
	err := s.db.Get("Person", id, &p)
	return &p, err
}

func (s *Storm) Iter(c chan<- Person) error {
	defer close(c)
	return s.db.Select().Each(new(Person), func(v interface{}) error {
		c <- *(v.(*Person))
		return nil
	})
}

func (s *Storm) Count() (int, error) {
	n := 0
	err := s.db.Select().Each(new(Person), func(v interface{}) error {
		n++
		return nil
	})
	return n, err
}

func (s *Storm) Drop() error {
	defer os.Remove(s.filename)
	return s.db.Close()
}

func BenchmarkStormPut(b *testing.B) {
	b.StopTimer()
	s := NewStorm()
	defer s.Drop()
	b.StartTimer()
	benchPut(b, s)
}

func BenchmarkStormGet(b *testing.B) {
	b.StopTimer()
	s := NewStorm()
	defer s.Drop()
	b.StartTimer()
	benchGet(b, s)
}

func BenchmarkStormIter(b *testing.B) {
	b.StopTimer()
	s := NewStorm()
	defer s.Drop()
	b.StartTimer()
	benchIter(b, s)
}
