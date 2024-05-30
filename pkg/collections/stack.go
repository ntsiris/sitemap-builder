package stack

type Item interface{}

type Stack struct {
	items []Item
}

func (stack *Stack) Push(item Item) {
	stack.items = append(stack.items, item)
}

func (stack *Stack) Pop() Item {
	if len(stack.items) == 0 {
		return nil
	}

	popItem := stack.items[len(stack.items)-1]
	stack.items = stack.items[:len(stack.items)-1]

	return popItem
}

func (stack *Stack) isEmpty() bool {
	return len(stack.items) == 0
}

func (stack *Stack) Peek() Item {
	if len(stack.items) == 0 {
		return nil
	}

	return stack.items[len(stack.items)-1]
}

func (stack *Stack) Len() int {
	return len(stack.items)
}
