package go_database_bench

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/dgraph-io/badger"
)

type Badger struct {
	db  *badger.DB
	dir string
}

func NewBadger() *Badger {
	dir := tempfile("badger-")
	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = dir
	db, err := badger.Open(opts)
	if err != nil {
		panic(err)
	}
	return &Badger{db: db, dir: dir}
}

func (b *Badger) Put(p *Person) error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set(p.Id, data)
	})
}

func (b *Badger) Get(id []byte) (*Person, error) {
	var p Person
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(id)
		if err != nil {
			return err
		}
		data, err := item.Value()
		if err != nil {
			return err
		}
		return json.Unmarshal(data, &p)
	})
	return &p, err
}

func (b *Badger) Iter(c chan<- Person) error {
	defer close(c)
	return b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		var p Person
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			v, err := item.Value()
			if err != nil {
				return err
			}
			err = json.Unmarshal(v, &p)
			if err != nil {
				return err
			}
			c <- p
		}
		return nil
	})
}

func (b *Badger) Count() (int, error) {
	n := 0
	c := make(chan Person)
	done := make(chan error)
	go func() {
		done <- b.Iter(c)
	}()
	for _ = range c {
		n++
	}
	return n, <-done
}

func (b *Badger) Drop() error {
	defer os.RemoveAll(b.dir)
	return b.db.Close()
}

func BenchmarkBadgerPut(b *testing.B) {
	b.StopTimer()
	badger := NewBadger()
	defer badger.Drop()
	b.StartTimer()
	benchPut(b, badger)
}

func BenchmarkBadgerGet(b *testing.B) {
	b.StopTimer()
	badger := NewBadger()
	defer badger.Drop()
	b.StartTimer()
	benchGet(b, badger)
}

func BenchmarkBadgerIter(b *testing.B) {
	b.StopTimer()
	badger := NewBadger()
	defer badger.Drop()
	b.StartTimer()
	benchIter(b, badger)
}
