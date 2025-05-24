package bst

import (
	"testing"
)

func compareInt(a, b int) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}
func TestAll(t *testing.T) {
	t.Run("BSTMap", func(t *testing.T) {
		t.Run("TestPut_EmptyMap", TestPut_EmptyMap)
		t.Run("TestPut_SingleInsert", TestPut_SingleInsert)
		t.Run("TestPut_SingleInsert", TestPut_SingleInsert)
		t.Run("TestPut_SingleInsert", TestPut_SingleInsert)
		t.Run("TestPut_MultipleInserts", TestPut_MultipleInserts)
		t.Run("TestPut_UpdateExistingKey", TestPut_UpdateExistingKey)
		t.Run("TestPut_UpdateExistingKey", TestPut_UpdateExistingKey)
		t.Run("TestPut_NilCompareFunction", TestPut_NilCompareFunction)

		t.Run("Test1", TestBST_CheckNodesDepthAndValues1)
		t.Run("Test2", TestBST_CheckNodesDepthAndValues2)
		t.Run("Test3", TestBST_CheckNodesAfterDeletion)

	})
}
func TestPut_EmptyMap(t *testing.T) {
	m := NewBSTMap[int, string](compareInt).(*Map[int, string])

	if m.Size() != 0 {
		t.Errorf("New map should be empty, got size %d", m.Size())
	}
}

func TestPut_SingleInsert(t *testing.T) {
	m := NewBSTMap[int, string](compareInt).(*Map[int, string])

	err := m.Put(1, "one")
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	if m.Size() != 1 {
		t.Errorf("Expected size 1, got %d", m.Size())
	}

	val, err := m.Get(1)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if val != "one" {
		t.Errorf("Expected value 'one', got '%s'", val)
	}
}

func TestPut_MultipleInserts(t *testing.T) {
	m := NewBSTMap[int, string](compareInt).(*Map[int, string])

	testCases := []struct {
		key   int
		value string
	}{
		{5, "five"},
		{3, "three"},
		{7, "seven"},
	}

	for _, tc := range testCases {
		err := m.Put(tc.key, tc.value)
		if err != nil {
			t.Fatalf("Put(%d) failed: %v", tc.key, err)
		}
	}

	if m.Size() != len(testCases) {
		t.Errorf("Expected size %d, got %d", len(testCases), m.Size())
	}

	for _, tc := range testCases {
		val, err := m.Get(tc.key)
		if err != nil {
			t.Errorf("Get(%d) failed: %v", tc.key, err)
		}
		if val != tc.value {
			t.Errorf("Expected value '%s' for key %d, got '%s'", tc.value, tc.key, val)
		}
	}

	m.AsciiPrint()
}

func TestPut_UpdateExistingKey(t *testing.T) {
	m := NewBSTMap[int, string](compareInt).(*Map[int, string])

	err := m.Put(1, "one")
	if err != nil {
		t.Fatalf("First Put failed: %v", err)
	}

	err = m.Put(1, "updated_one")
	if err != nil {
		t.Fatalf("Second Put failed: %v", err)
	}

	if m.Size() != 1 {
		t.Errorf("Size should remain 1 after update, got %d", m.Size())
	}

	val, err := m.Get(1)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if val != "updated_one" {
		t.Errorf("Expected updated value 'updated_one', got '%s'", val)
	}
}

func TestPut_NilCompareFunction(t *testing.T) {
	m := &Map[int, string]{} // No compare function provided

	err := m.Put(1, "one")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedErr := "comparison function not provided"
	if err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%v'", expectedErr, err)
	}
}

func TestBST_CheckNodesDepthAndValues1(t *testing.T) {
	m := NewBSTMap[int, string](compareInt).(*Map[int, string])

	// Вставка элементов как в первом C++ тесте
	m.Put(5, "a")
	m.Put(2, "b")
	m.Put(15, "c")
	m.Put(3, "d")
	m.Put(14, "e")
	m.Put(1, "l")

	// Ожидаемые данные: глубина, ключ, значение
	expected := []struct {
		depth int
		key   int
		value string
	}{
		{2, 1, "l"},
		{1, 2, "b"},
		{2, 3, "d"},
		{0, 5, "a"},
		{2, 14, "e"},
		{1, 15, "c"},
	}

	for _, exp := range expected {
		// Проверяем значение
		val, err := m.Get(exp.key)
		if err != nil {
			t.Errorf("Key %d not found", exp.key)
			continue
		}
		if val != exp.value {
			t.Errorf("For key %d: expected value '%s', got '%s'", exp.key, exp.value, val)
		}

		// Проверяем глубину
		depth, err := m.GetDepth(exp.key)
		if err != nil {
			t.Errorf("Failed to get depth for key %d: %v", exp.key, err)
			continue
		}
		if depth != exp.depth {
			t.Errorf("For key %d: expected depth %d, got %d", exp.key, exp.depth, depth)
		}
	}
	m.AsciiPrint()

}

func TestBST_CheckNodesDepthAndValues2(t *testing.T) {
	m := NewBSTMap[int, int](compareInt).(*Map[int, int])

	// Вставка элементов как во втором C++ тесте
	m.Put(1, 5)
	m.Put(2, 12)
	m.Put(15, 1)
	m.Put(3, 67)
	m.Put(4, 45)

	// Ожидаемые данные: глубина, ключ, значение
	expected := []struct {
		depth int
		key   int
		value int
	}{
		{0, 1, 5},
		{1, 2, 12},
		{2, 15, 1},
		{3, 3, 67},
		{4, 4, 45},
	}

	for _, exp := range expected {
		val, err := m.Get(exp.key)
		if err != nil {
			t.Errorf("Key %d not found", exp.key)
			continue
		}
		if val != exp.value {
			t.Errorf("For key %d: expected value %d, got %d", exp.key, exp.value, val)
		}

		depth, err := m.GetDepth(exp.key)
		if err != nil {
			t.Errorf("Failed to get depth for key %d: %v", exp.key, err)
			continue
		}
		if depth != exp.depth {
			t.Errorf("For key %d: expected depth %d, got %d", exp.key, exp.depth, depth)
		}
	}
	m.AsciiPrint()

}

func TestBST_CheckNodesAfterDeletion(t *testing.T) {
	m := NewBSTMap[int, string](compareInt).(*Map[int, string])

	// Вставка элементов как в третьем C++ тесте
	m.Put(6, "a")
	m.Put(8, "c")
	m.Put(15, "x")
	m.Put(11, "j")
	m.Put(19, "i")
	m.Put(12, "l")
	m.Put(17, "b")
	m.Put(18, "e")

	// Удаление узла
	err := m.Delete(15)
	if err != nil {
		t.Fatalf("Failed to delete key 15: %v", err)
	}

	// Ожидаемые данные после удаления: глубина, ключ, значение
	expected := []struct {
		depth int
		key   int
		value string
	}{
		{0, 6, "a"},
		{1, 8, "c"},
		{3, 11, "j"},
		{2, 12, "l"},
		{4, 17, "b"},
		{5, 18, "e"},
		{3, 19, "i"},
	}

	for _, exp := range expected {
		val, err := m.Get(exp.key)
		if err != nil {
			t.Errorf("Key %d not found after deletion", exp.key)
			continue
		}
		if val != exp.value {
			t.Errorf("For key %d: expected value '%s', got '%s'", exp.key, exp.value, val)
		}

		depth, err := m.GetDepth(exp.key)
		if err != nil {
			t.Errorf("Failed to get depth for key %d: %v", exp.key, err)
			continue
		}
		if depth != exp.depth {
			t.Errorf("For key %d: expected depth %d, got %d", exp.key, exp.depth, depth)
		}
	}

	// Проверяем, что удаленный элемент отсутствует
	_, err = m.Get(15)
	if err == nil {
		t.Error("Deleted key 15 still exists in the tree")
	}
	m.AsciiPrint()
}
