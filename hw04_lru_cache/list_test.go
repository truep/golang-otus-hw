package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("first element to front", func(t *testing.T) {
		l := NewList()
		require.Equal(t, 0, l.Len())

		li := l.PushFront(11)
		require.Equal(t, l.Front(), li)
		require.Equal(t, l.Back(), li)

		require.Nil(t, li.Prev)
		require.Nil(t, li.Next)
	})

	t.Run("first element to back", func(t *testing.T) {
		l := NewList()
		require.Equal(t, 0, l.Len())

		li := l.PushBack(11)
		require.Equal(t, l.Front(), li)
		require.Equal(t, l.Back(), li)

		require.Nil(t, li.Prev)
		require.Nil(t, li.Next)
	})

	t.Run("more then one element in list", func(t *testing.T) {
		l := NewList()
		require.Equal(t, 0, l.Len())

		liFirst := l.PushFront(11)
		liSecond := l.PushFront(22)

		require.Nil(t, liFirst.Next)
		require.Nil(t, liSecond.Prev)

		require.Equal(t, liFirst.Prev, liSecond)
		require.Equal(t, liSecond.Next, liFirst)

		require.Equal(t, l.Front(), liSecond)
		require.Equal(t, l.Back(), liFirst)

		liThird := l.PushBack(33)
		require.Nil(t, liSecond.Prev)
		require.Nil(t, liThird.Next)

		require.Equal(t, liFirst.Prev, liSecond)
		require.Equal(t, liFirst.Next, liThird)

		require.Equal(t, liSecond.Next, liFirst)
		require.Equal(t, liThird.Prev, liFirst)

		require.Equal(t, l.Front(), liSecond)
		require.Equal(t, l.Back(), liThird)

		require.Equal(t, 11, liFirst.Value)
		require.Equal(t, 22, liSecond.Value)
		require.Equal(t, 33, liThird.Value)
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		last := l.PushBack(50) // [50]
		require.Equal(t, 1, l.Len())

		l.Remove(last)  // []
		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		l.Remove(l.Front()) // [30]
		require.Equal(t, 1, l.Len())

		l.Remove(l.Front()) // []
		l.PushFront(10)     // [10]
		l.PushBack(30)      // [10, 30]
		require.Equal(t, 2, l.Len())
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 30, l.Back().Value)

		l.Remove(l.Back()) // [10]
		require.Equal(t, 1, l.Len())
		require.Equal(t, 10, l.Back().Value)

		l.PushBack(30)
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)

		elems = nil
		l.MoveToFront(l.Front().Next) // [80, 70, 60, 40, 10, 30, 50]
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{80, 70, 60, 40, 10, 30, 50}, elems)

		elems = nil
		l.MoveToFront(l.Back().Prev) // [30, 80, 70, 60, 40, 10, 50]
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{30, 80, 70, 60, 40, 10, 50}, elems)
	})
}
