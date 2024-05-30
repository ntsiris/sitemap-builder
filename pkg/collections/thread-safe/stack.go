package syncstack

import "sync"

type Item interface{}

type Stack struct {
	items []Item
	mutex sync.Mutex
}

func (stack *Stack) Push(item Item) {
	stack.mutex.Lock()
	defer stack.mutex.Unlock()

	stack.items = append(stack.items, item)
}

func (stack *Stack) Pop() Item {
	stack.mutex.Lock()
	defer stack.mutex.Unlock()

	if len(stack.items) == 0 {
		return nil
	}

	popItem := stack.items[len(stack.items)-1]
	stack.items = stack.items[:len(stack.items)-1]

	return popItem
}

func (stack *Stack) isEmpty() bool {
	stack.mutex.Lock()
	defer stack.mutex.Unlock()

	return len(stack.items) == 0
}

func (stack *Stack) Peek() Item {
	stack.mutex.Lock()
	defer stack.mutex.Unlock()

	if len(stack.items) == 0 {
		return nil
	}

	return stack.items[len(stack.items)-1]
}

func (stack *Stack) Len() int {
	stack.mutex.Lock()
	defer stack.mutex.Unlock()

	return len(stack.items)
}
