package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(10)
		var testKey Key = "aaa"
		testValue := 100

		c.Set(testKey, testValue)
		cacheValue, ok := c.Get(testKey)
		require.True(t, ok)
		require.Equal(t, 100, cacheValue)

		c.Clear()
		_, ok = c.Get(testKey)
		require.False(t, ok)
	})

	t.Run("set overflow", func(t *testing.T) {
		c := NewCache(3)
		for i := 0; i < 4; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}

		_, ok := c.Get("0")
		require.False(t, ok)
	})

	t.Run("overflow not used elements", func(t *testing.T) {
		c := NewCache(3)
		for i := 0; i < 3; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}

		c.Get("0")
		c.Set("3", 3)

		_, ok := c.Get("0")
		require.True(t, ok)

		_, ok = c.Get("1")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
			if rand.Intn(1_000_000) < 500_000 {
				c.Clear()
			}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
