package store

type Node struct {
	val  string
	prev *Node
	next *Node
}

func NewNode(value string) *Node {
	return &Node{
		val: value,
	}
}

type LinkedList struct {
	head *Node
	tail *Node
	size int
}

func NewLinkedList(value string) *LinkedList {
	newNode := NewNode(value)

	newLinkedList := &LinkedList{
		head: newNode,
		tail: newNode,
		size: 1,
	}

	return newLinkedList
}

//think about the implementation of structs using sync.Pool

func (ll *LinkedList) llen() int {
	return ll.size
}

func (ll *LinkedList) rpush(value string) {
	node := NewNode(value)

	if ll.size == 0 {
		ll.head = node
		ll.tail = node
	} else {
		ll.tail.next = node
		node.prev = ll.tail
		ll.tail = node
	}
	ll.size++
}

func (ll *LinkedList) lpush(value string) {
	node := NewNode(value)

	if ll.size == 0 {
		ll.head = node
		ll.tail = node
	} else {
		ll.head.prev = node
		node.next = ll.head
		ll.head = node
	}
	ll.size++
}

func (ll *LinkedList) rpop(count int) []string {
	ret := make([]string, min(ll.size, count))
	for i := range len(ret) {
		ret[i] = ll.tail.val
		ll.tail = ll.tail.prev
		if ll.tail != nil {
			ll.tail.next = nil
		} else {
			ll.head = nil
		}
		ll.size--
	}
	return ret
}

func (ll *LinkedList) lpop(count int) []string {
	ret := make([]string, min(ll.size, count))
	for i := range len(ret) {
		ret[i] = ll.head.val
		ll.head = ll.head.next
		if ll.head != nil {
			ll.head.prev = nil
		} else {
			ll.tail = nil
		}
		ll.size--
	}
	return ret
}

func (ll *LinkedList) lrange(start int, stop int) []string {
	// case 1 start +ve stop +ve -> if start < size -> iterate till min(stop, size)
	// case 2 start -ve stop +ve -> find start, start < 0, start becomes 0, if start
	// case 3 start +ve stop -ve -> normalize the stop,
	// case 4 start -ve stop -ve ->
	// combine these to -> if start is -ve, noarmalize it, then change it to 0 if it is still negative, then check if its less than size, then iterate.
	if start < 0 {
		start = ll.size + start
	}
	start = max(start, 0)

	if stop < 0 {
		stop = ll.size + stop
	}

	stop = min(stop, ll.size-1)

	if start >= ll.size || stop < 0 {
		return []string{}
	}

	retArr := make([]string, stop-start+1)
	temp := ll.head
	for i, ind := 0, 0; i < ll.size; i++ {
		if i > stop {
			break
		}
		if i >= start {
			retArr[ind] = temp.val
			ind++
		}
		temp = temp.next
	}

	return retArr
}

// exposed functions
