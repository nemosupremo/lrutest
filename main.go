package lru

import (
	"container/list"
	"sync"
)

//keyvalue is an internal structure used by preparedLRU
type keyvalue struct {
	key   string
	value interface{}
}

func NewPreparedLRU() *preparedLRU {
	return &preparedLRU{stmts: make([]*keyvalue, 10)}
}

//preparedLRU is a structure that houses the logic for caching and maintaing the cache of prepared
//statements for all connections
type preparedLRU struct {
	stmts []*keyvalue
	count int
	mu    sync.Mutex
}

//get fetchs the request query from the cache. If query is not cached nil will be returned
func (c *preparedLRU) get(address, stmt string) (ifp interface{}) {
	c.mu.Lock()
	pos := -1
	for i := 0; i < c.count; i++ {
		if c.stmts[i].key == address+stmt {
			pos = i
			break
		}
	}
	if pos == -1 {
		ifp = nil
	} else {
		if pos > 0 {
			//Remove existing reference and place at position 0
			kv := c.stmts[pos]
			copy(c.stmts[pos:c.count-1], c.stmts[pos+1:c.count])
			copy(c.stmts[1:c.count], c.stmts[:c.count-1])
			c.stmts[0] = kv
		}
		ifp = c.stmts[0].value
	}
	c.mu.Unlock()
	return
}

//set puts the prepared statement into the cache. Set will update the existing reference
//if the query is already cached.
func (c *preparedLRU) set(address, stmt string, ifp interface{}) {
	c.mu.Lock()
	pos := 0
	for i := 0; i < c.count; i++ {
		if c.stmts[i].key == address+stmt {
			pos = i
			break
		}
	}
	if pos != 0 {
		copy(c.stmts[pos:c.count-1], c.stmts[pos+1:c.count])
		c.stmts[c.count] = nil
		c.count--
	} else if c.count == len(c.stmts) {
		c.stmts[len(c.stmts)-1] = nil
		c.count--
	}
	if c.count > 0 {
		copy(c.stmts[1:c.count+1], c.stmts[:c.count])
	}
	c.stmts[0] = &keyvalue{key: address + stmt, value: ifp}
	c.count++
	c.mu.Unlock()

}

//delete removes the query from the cache
func (c *preparedLRU) delete(address, stmt string) {
	c.mu.Lock()
	key := -1
	for i := 0; i < c.count; i++ {
		if c.stmts[i].key == address+stmt {
			key = i
			break
		}
	}
	if key != -1 {
		copy(c.stmts[key:c.count-1], c.stmts[key+1:c.count])
		if c.count == len(c.stmts) {
			c.stmts[c.count-1] = nil
		} else {
			c.stmts[c.count] = nil
		}
		c.count--
	}
	c.mu.Unlock()
}

//setMaxStmts controls the size of the cache. This is safe to modify at run time.
func (c *preparedLRU) setMaxStmts(max int) {
	c.mu.Lock()
	stmtsL := len(c.stmts)
	if max != stmtsL {
		newStmts := make([]*keyvalue, max)
		if max > stmtsL {
			copy(newStmts, c.stmts)
		} else if max < stmtsL {
			copy(newStmts, c.stmts[:max])
			c.count = max
		}
		c.stmts = newStmts
	}
	c.mu.Unlock()
}

func NewListLRU() *listLRU {
	return &listLRU{stmts: list.New()}
}

//preparedLRU is a structure that houses the logic for caching and maintaing the cache of prepared
//statements for all connections
type listLRU struct {
	stmts *list.List
	max   int
	mu    sync.Mutex
}

//get fetchs the request query from the cache. If query is not cached nil will be returned
func (c *listLRU) get(address, stmt string) (ifp interface{}) {
	c.mu.Lock()
	var elem *list.Element
	for e := c.stmts.Front(); e != nil; e = e.Next() {
		if e.Value.(*keyvalue).key == address+stmt {
			elem = e
			break
		}
	}
	if elem == nil {
		ifp = nil
	} else {
		if c.stmts.Front() != elem {
			c.stmts.MoveToFront(elem)
		}
		ifp = c.stmts.Front().Value
	}
	c.mu.Unlock()
	return
}

//set puts the prepared statement into the cache. Set will update the existing reference
//if the query is already cached.
func (c *listLRU) set(address, stmt string, ifp interface{}) {
	c.mu.Lock()
	var elem *list.Element
	for e := c.stmts.Front(); e != nil; e = e.Next() {
		if e.Value.(*keyvalue).key == address+stmt {
			elem = e
			break
		}
	}
	if elem != nil {
		c.stmts.Remove(elem)
	} else if c.stmts.Len() == c.max {
		c.stmts.Remove(c.stmts.Back())
	}
	c.stmts.PushFront(&keyvalue{key: address + stmt, value: ifp})
	c.mu.Unlock()

}

//delete removes the query from the cache
func (c *listLRU) delete(address, stmt string) {
	c.mu.Lock()
	for e := c.stmts.Front(); e != nil; e = e.Next() {
		if e.Value.(*keyvalue).key == address+stmt {
			c.stmts.Remove(e)
			break
		}
	}
	c.mu.Unlock()
}

//setMaxStmts controls the size of the cache. This is safe to modify at run time.
func (c *listLRU) setMaxStmts(max int) {
	c.mu.Lock()
	c.max = max
	for c.stmts.Len() > max {
		c.stmts.Remove(c.stmts.Back())
	}
	c.mu.Unlock()
}
