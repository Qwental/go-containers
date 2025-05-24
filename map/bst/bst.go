package bst

import (
	"errors"
	"fmt"
	"github.com/Qwental/go-containers/map"
	"reflect"
	"strings"
)

type CompareFunc[K comparable] func(k1, k2 K) int

type Node[K comparable, V any] struct {
	Key    K
	Value  V
	Parent *Node[K, V]
	Left   *Node[K, V]
	Right  *Node[K, V]
}

type Map[K comparable, V any] struct {
	root    *Node[K, V]
	size    int
	compare CompareFunc[K]
}

func NewBSTMap[K comparable, V any](compare CompareFunc[K]) _map.Map[K, V] {
	return &Map[K, V]{compare: compare}
}

func (m *Map[K, V]) Put(key K, value V) error {
	if m.compare == nil {
		return errors.New("comparison function not provided")
	}
	if m.root == nil {
		m.root = &Node[K, V]{Key: key, Value: value}
		m.size++
		return nil
	}
	_, err := m.Get(key)
	if err == nil {
		m.size-- // Чтоб заново два раза size не увеличить
	}
	m.put(key, value)
	m.size++
	return nil
}

func (m *Map[K, V]) put(key K, value V) {
	current := m.root
	var parent *Node[K, V]

	for current != nil {
		parent = current
		cmp := m.compare(key, current.Key)
		if cmp < 0 {
			current = current.Left
		} else if cmp > 0 {
			current = current.Right
		} else {
			current.Value = value
			return
		}
	}

	newNode := &Node[K, V]{Key: key, Value: value, Parent: parent}
	if parent != nil {
		if m.compare(key, parent.Key) < 0 {
			parent.Left = newNode
		} else {
			parent.Right = newNode
		}
	} else {
		panic("Dereferenced key not found in BSTMap")
	}
}

func (m *Map[K, V]) Get(key K) (V, error) {
	if m.compare == nil {
		return *new(V), errors.New("comparison function not provided")
	}
	n := m.root
	for n != nil {
		cmp := m.compare(key, n.Key)
		if cmp < 0 {
			n = n.Left
		} else if cmp > 0 {
			n = n.Right
		} else {
			return n.Value, nil
		}
	}
	return *new(V), errors.New("key not found")
}

func (m *Map[K, V]) Delete(key K) error {
	if m.compare == nil {
		return errors.New("comparison function not provided")
	}

	// Поиск узла для удаления
	//var node *Node[K, V]
	//current := m.root
	//for current != nil {
	//	cmp := m.compare(key, current.Key)
	//	if cmp < 0 {
	//		current = current.Left
	//	} else if cmp > 0 {
	//		current = current.Right
	//	} else {
	//		node = current
	//		break
	//	}
	//}
	node, err := m.getNode(key)
	if err != nil {
		return err
	}

	if node == nil {
		return errors.New("key not found")
	}

	m.size--

	// Случай 1: Нет потомков
	if node.Left == nil && node.Right == nil {
		if node.Parent == nil {
			m.root = nil
		} else {
			if node.Parent.Left == node {
				node.Parent.Left = nil
			} else {
				node.Parent.Right = nil
			}
		}
		return nil
	}

	// Случай 2: Только один потомок
	if node.Left == nil || node.Right == nil {
		var child *Node[K, V]
		if node.Left != nil {
			child = node.Left
		} else {
			child = node.Right
		}

		child.Parent = node.Parent

		if node.Parent == nil {
			m.root = child
		} else {
			if node.Parent.Left == node {
				node.Parent.Left = child
			} else {
				node.Parent.Right = child
			}
		}
		return nil
	}

	// Случай 3: Два потомка (используем predecessor)
	predecessor := node.Left
	for predecessor.Right != nil {
		predecessor = predecessor.Right
	}

	// Если predecessor не является прямым левым потомком
	if predecessor.Parent != node {
		predecessor.Parent.Right = predecessor.Left
		if predecessor.Left != nil {
			predecessor.Left.Parent = predecessor.Parent
		}

		predecessor.Left = node.Left
		if predecessor.Left != nil {
			predecessor.Left.Parent = predecessor
		}
	}

	// Переносим правое поддерево
	predecessor.Right = node.Right
	if predecessor.Right != nil {
		predecessor.Right.Parent = predecessor
	}

	// Обновляем родителя
	predecessor.Parent = node.Parent
	if node.Parent == nil {
		m.root = predecessor
	} else {
		if node.Parent.Left == node {
			node.Parent.Left = predecessor
		} else {
			node.Parent.Right = predecessor
		}
	}

	return nil
}

func (m *Map[K, V]) findMin(n *Node[K, V]) *Node[K, V] {
	for n.Right != nil {
		n = n.Right
	}
	return n
}

func (m *Map[K, V]) Size() int {
	return m.size
}

func (m *Map[K, V]) Print() {
	if m.root == nil {
		fmt.Println("(empty tree)")
		return
	}

	type nodeLevel struct {
		node  *Node[K, V]
		level int
	}

	queue := []nodeLevel{{node: m.root, level: 0}}
	currentLevel := 0
	var levels [][]K

	for len(queue) > 0 {
		nl := queue[0]
		queue = queue[1:]

		if nl.level > currentLevel {
			currentLevel = nl.level
		}

		if len(levels) <= currentLevel {
			levels = append(levels, []K{})
		}

		if nl.node != nil {
			levels[currentLevel] = append(levels[currentLevel], nl.node.Key)
			queue = append(queue,
				nodeLevel{node: nl.node.Left, level: currentLevel + 1},
				nodeLevel{node: nl.node.Right, level: currentLevel + 1},
			)
		} else {
			levels[currentLevel] = append(levels[currentLevel], *new(K))
		}
	}

	fmt.Printf("Tree (size: %d):\n", m.size)
	for i, level := range levels {
		fmt.Printf("Level %d: ", i)
		for _, key := range level {
			var keyStr string
			if reflect.ValueOf(key).IsZero() {
				keyStr = "nil"
			} else {
				keyStr = fmt.Sprintf("%v", key)
			}
			fmt.Printf("%5s ", keyStr)
		}
		fmt.Println()
	}
}

func (m *Map[K, V]) AsciiPrint() {
	if m.root == nil {
		fmt.Println("(empty tree)")
		return
	}
	fmt.Printf("ASCII Print Tree with size: %d:\n", m.size)
	m.printNode(m.root, 0)
}

func (m *Map[K, V]) printNode(n *Node[K, V], depth int) {
	if n == nil {
		return
	}
	m.printNode(n.Right, depth+1)
	fmt.Print(strings.Repeat("    ", depth))
	fmt.Printf("%v\n", n.Key)
	m.printNode(n.Left, depth+1)
}

func (m *Map[K, V]) getNode(key K) (*Node[K, V], error) {
	if m.compare == nil {
		return nil, errors.New("comparison function not provided")
	}

	current := m.root
	for current != nil {
		cmp := m.compare(key, current.Key)
		if cmp < 0 {
			current = current.Left
		} else if cmp > 0 {
			current = current.Right
		} else {
			return current, nil
		}
	}
	return nil, errors.New("key not found")
}

func (m *Map[K, V]) GetRight(key K) (rightKey K, rightValue V, err error) {
	node, err := m.getNode(key)
	if err != nil {
		return rightKey, rightValue, err
	}
	if node.Right == nil {
		return rightKey, rightValue, errors.New("right child not found")
	}
	return node.Right.Key, node.Right.Value, nil
}

func (m *Map[K, V]) GetLeft(key K) (leftKey K, leftValue V, err error) {
	node, err := m.getNode(key)
	if err != nil {
		return leftKey, leftValue, err
	}
	if node.Left == nil {
		return leftKey, leftValue, errors.New("left child not found")
	}
	return node.Left.Key, node.Left.Value, nil
}

func (m *Map[K, V]) GetParent(key K) (parentKey K, parentValue V, err error) {
	node, err := m.getNode(key)
	if err != nil {
		return parentKey, parentValue, err
	}
	if node.Parent == nil {
		return parentKey, parentValue, errors.New("root has no parent")
	}
	return node.Parent.Key, node.Parent.Value, nil
}

func (m *Map[K, V]) GetDepth(key K) (depth int, err error) {
	node, err := m.getNode(key)
	if err != nil {
		return -1, err
	}

	for node.Parent != nil {
		depth++
		node = node.Parent
	}
	return depth, nil
}
