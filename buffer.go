package xcache

type Buffer interface {
	IsFull() bool
	Evict() *Block
	Add(key Key, value Value) *Block
	Remove(block *Block)
	Access(block *Block) *Block
	Size() int
	Keys() []Key
	Values() []Value
	Has(block *Block) bool
}

type LRU struct {
	head *Block
	tail *Block
	size int
	capacity int
}

func (l LRU) Keys() []Key {
	keys := make([]Key,0)

	for curr := l.head.next; curr != nil; curr = curr.next {
		keys = append(keys, curr.key)
	}

	return keys
}

func (l LRU) Values() []Value {
	values := make([]Value,0)

	for curr := l.head.next; curr != nil; curr = curr.next {
		values = append(values, curr.val)
	}

	return values
}

func (l LRU) Has(block *Block) bool {
	for curr := l.head; curr != nil; curr = curr.next {
		if curr == block {
			return true
		}
	}

	return false
}

func (l *LRU) Access(block *Block) *Block {
	l.Remove(block)

	key := block.key
	value := block.val

	return l.Add(key, value)
}

func (l LRU) Size() int {
	return l.size
}

func (l *LRU) Add(key Key, value Value) *Block {
	newBlock := Block{
		prev: l.tail,
		key: key,
		val:value,
	}
	l.tail.next = &newBlock

	l.tail = &newBlock
	l.size++

	return &newBlock
}

func (l *LRU) Remove(block *Block) {
	if block == nil {
		return
	}

	prev := block.prev
	prev.next = block.next

	if prev.next != nil {
		prev.next.prev = prev
	}

	if l.tail == block {
		l.tail = block.prev
	}

	l.size--
}

func (l LRU) IsFull() bool {
	return l.size == l.capacity
}

func (l *LRU) Evict() *Block {
	block := l.head.next
	l.Remove(l.head.next)
	return block
}

type Block struct {
	prev *Block
	next *Block
	key Key
	val Value
}

func NewLRUBuffer(capacity int) Buffer {
	dummyBlock :=  Block{}

	return &LRU{
		head: &dummyBlock,
		tail: &dummyBlock,
		size:0,
		capacity:capacity,
	}
}