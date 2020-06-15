// Copyright 2010 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package llrb

import (
	"math"
	"math/rand"
	"testing"
)

func TestCases(t *testing.T) {
	tree := New()
	tree.ReplaceOrInsert(Int(1))
	tree.ReplaceOrInsert(Int(1))
	if tree.Len() != 1 {
		t.Errorf("expecting len 1")
	}
	if !tree.Has(Int(1)) {
		t.Errorf("expecting to find key=1")
	}

	tree.Delete(Int(1))
	if tree.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	if tree.Has(Int(1)) {
		t.Errorf("not expecting to find key=1")
	}

	tree.Delete(Int(1))
	if tree.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	if tree.Has(Int(1)) {
		t.Errorf("not expecting to find key=1")
	}
}

func TestReverseInsertOrder(t *testing.T) {
	tree := New()
	n := 100
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(Int(n - i))
	}
	i := 0
	tree.AscendGreaterOrEqual(Int(0), func(item Item) bool {
		i++
		if item.(Int) != Int(i) {
			t.Errorf("bad order: got %d, expect %d", item.(Int), i)
		}
		return true
	})
}

func TestRange(t *testing.T) {
	tree := New()
	order := []String{
		"ab", "aba", "abc", "a", "aa", "aaa", "b", "a-", "a!",
	}
	for _, i := range order {
		tree.ReplaceOrInsert(i)
	}
	k := 0
	tree.AscendRange(String("ab"), String("ac"), func(item Item) bool {
		if k > 3 {
			t.Fatalf("returned more items than expected")
		}
		i1 := order[k]
		i2 := item.(String)
		if i1 != i2 {
			t.Errorf("expecting %s, got %s", i1, i2)
		}
		k++
		return true
	})
}

func TestRandomInsertOrder(t *testing.T) {
	tree := New()
	n := 1000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(Int(perm[i]))
	}
	j := 0
	tree.AscendGreaterOrEqual(Int(0), func(item Item) bool {
		if item.(Int) != Int(j) {
			t.Fatalf("bad order")
		}
		j++
		return true
	})
}

func TestRandomReplace(t *testing.T) {
	tree := New()
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(Int(perm[i]))
	}
	perm = rand.Perm(n)
	for i := 0; i < n; i++ {
		if replaced := tree.ReplaceOrInsert(Int(perm[i])); replaced == nil || replaced.(Int) != Int(perm[i]) {
			t.Errorf("error replacing")
		}
	}
}

func TestRandomInsertSequentialDelete(t *testing.T) {
	tree := New()
	n := 1000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(Int(perm[i]))
	}
	for i := 0; i < n; i++ {
		tree.Delete(Int(i))
	}
}

func TestRandomInsertDeleteNonExistent(t *testing.T) {
	tree := New()
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(Int(perm[i]))
	}
	if tree.Delete(Int(200)) != nil {
		t.Errorf("deleted non-existent item")
	}
	if tree.Delete(Int(-2)) != nil {
		t.Errorf("deleted non-existent item")
	}
	for i := 0; i < n; i++ {
		if u := tree.Delete(Int(i)); u == nil || u.(Int) != Int(i) {
			t.Errorf("delete failed")
		}
	}
	if tree.Delete(Int(200)) != nil {
		t.Errorf("deleted non-existent item")
	}
	if tree.Delete(Int(-2)) != nil {
		t.Errorf("deleted non-existent item")
	}
}

func TestRandomInsertPartialDeleteOrder(t *testing.T) {
	tree := New()
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(Int(perm[i]))
	}
	for i := 1; i < n-1; i++ {
		tree.Delete(Int(i))
	}
	j := 0
	tree.AscendGreaterOrEqual(Int(0), func(item Item) bool {
		switch j {
		case 0:
			if item.(Int) != Int(0) {
				t.Errorf("expecting 0")
			}
		case 1:
			if item.(Int) != Int(n-1) {
				t.Errorf("expecting %d", n-1)
			}
		}
		j++
		return true
	})
}

func TestRandomInsertStats(t *testing.T) {
	tree := New()
	n := 100000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(Int(perm[i]))
	}
	avg, _ := tree.HeightStats()
	expAvg := math.Log2(float64(n)) - 1.5
	if math.Abs(avg-expAvg) >= 2.0 {
		t.Errorf("too much deviation from expected average height")
	}
}

func BenchmarkInsert(b *testing.B) {
	tree := New()
	for i := 0; i < b.N; i++ {
		tree.ReplaceOrInsert(Int(b.N - i))
	}
}

func BenchmarkDelete(b *testing.B) {
	b.StopTimer()
	tree := New()
	for i := 0; i < b.N; i++ {
		tree.ReplaceOrInsert(Int(b.N - i))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.Delete(Int(i))
	}
}

func BenchmarkDeleteMin(b *testing.B) {
	b.StopTimer()
	tree := New()
	for i := 0; i < b.N; i++ {
		tree.ReplaceOrInsert(Int(b.N - i))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.DeleteMin()
	}
}

func TestInsertNoReplace(t *testing.T) {
	tree := New()
	n := 1000
	for q := 0; q < 2; q++ {
		perm := rand.Perm(n)
		for i := 0; i < n; i++ {
			tree.InsertNoReplace(Int(perm[i]))
		}
	}
	j := 0
	tree.AscendGreaterOrEqual(Int(0), func(item Item) bool {
		if item.(Int) != Int(j/2) {
			t.Fatalf("bad order")
		}
		j++
		return true
	})
}

func TestNDescendants(t *testing.T) {
	for outerIdx, shuffle := range [][]Int{
		[]Int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20},
		[]Int{16, 2, 20, 10, 8, 14, 6, 4, 12, 18},
	} {
		t.Log("outerIdx:", outerIdx)
		tree := New()
		for i := 1; i <= 10; i++ {
			tree.InsertNoReplace(shuffle[i-1])
		}
		//t.Log(tree.printBFS())
		if tree.root.NDescendants != 10 ||
			tree.root.Left.NDescendants != 3 ||
			tree.root.Right.NDescendants != 6 ||
			tree.root.Left.Left.NDescendants != 1 ||
			tree.root.Left.Right.NDescendants != 1 ||
			tree.root.Right.Left.NDescendants != 3 ||
			tree.root.Right.Right.NDescendants != 2 ||
			tree.root.Right.Left.Left.NDescendants != 1 ||
			tree.root.Right.Left.Right.NDescendants != 1 ||
			tree.root.Right.Right.Left.NDescendants != 1 {
			t.Error(tree.printBFS())
		}

		for i := 1; i <= 10; i++ {
			if reality := tree.GetByRank(i); reality.(Int) != Int(2*i) {
				t.Error(Int(2*i), reality)
			}
		}

		for i := 1; i <= 10; i++ {
			item := Int(2 * i)
			r, foundItem := tree.GetRankOf(item)
			//t.Logf("item: %v, rank: %v, foundItem: %v", item, r, foundItem)
			if r != i || foundItem == nil {
				t.Error(i, r, foundItem)
			}
		}

		r, foundItem := tree.GetRankOf(Int(17))
		if foundItem != nil || r != 9 {
			t.Error(r, foundItem, foundItem == nil)
		}
		r, foundItem = tree.GetRankOf(Int(5))
		if foundItem != nil || r != 3 {
			t.Error(r, foundItem)
		}
	}
}

func TestLLRB_RankEmpty(t *testing.T) {
	tree := New()
	if tree.GetByRank(10) != nil {
		t.Error()
	}
	rank, item := tree.GetRankOf(Int(10))
	if rank != 0 || item != nil {
		t.Error()
	}
}

func BenchmarkLLRB_GetRankOf(b *testing.B) {
	b.StopTimer()
	tree := New()
	for i := 0; i < b.N; i++ {
		tree.InsertNoReplace(Int(i))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.GetRankOf(Int(i))
		tree.GetByRank(rand.Intn(b.N))
	}
}

func TestLLRB_Delete(t *testing.T) {
	tree := New()
	for i := 1; i <= 10; i++ {
		tree.InsertNoReplace(Int(2 * i))
	}
	tree.Delete(Int(18))
	if tree.root.NDescendants != 9 {
		t.Errorf("root.NDescendants: expect 9, reality: %v", tree.root.NDescendants)
	}
	if tree.root.Right.NDescendants != 5 {
		t.Errorf("root.Right.Right.NDescendants: expect 5, reality: %v",
			tree.root.Right.NDescendants)
	}
	if tree.root.Right.Right.NDescendants != 1 {
		t.Errorf("root.Right.Right.NDescendants: expect 1, reality: %v",
			tree.root.Right.Right.NDescendants)
	}
}

func TestLLRB_DeleteMin(t *testing.T) {
	tree := New()
	tree.InsertNoReplace(Int(2))
	tree.InsertNoReplace(Int(4))
	//t.Log(tree.printBFS())
	tree.DeleteMin()
	if tree.root.NDescendants != 1 {
		t.Errorf("expect %v, reality: %v", 1, tree.root.NDescendants)
		t.Log(tree.printBFS())
	}
}

func TestLLRB_DeleteMin2(t *testing.T) {
	tree := New()
	shuffle := []Int{6, 2, 10, 8, 4}
	for _, e := range shuffle {
		tree.InsertNoReplace(e)
	}

	tree.DeleteMin()
	if tree.root.NDescendants != 4 ||
		tree.root.Left.NDescendants != 1 {
		t.Error(tree.printBFS())
	}

	tree.DeleteMin()
	if tree.root.NDescendants != 3 ||
		tree.root.Left.NDescendants != 1 {
		t.Error(tree.printBFS())
	}
}

func TestLLRB_DeleteMax(t *testing.T) {
	tree := New()
	tree.InsertNoReplace(Int(2))
	tree.InsertNoReplace(Int(4))
	//t.Log(tree.printBFS())
	tree.DeleteMax()
	if tree.root.NDescendants != 1 {
		t.Errorf(tree.printBFS())
	}
}

func TestLLRB_DeleteMax2(t *testing.T) {
	tree := New()
	shuffle := []Int{6, 2, 10, 8, 4}
	for _, e := range shuffle {
		tree.InsertNoReplace(e)
	}
	tree.DeleteMax()
	if tree.root.NDescendants != 4 ||
		tree.root.Right.NDescendants != 1 {
		t.Error(tree.printBFS())
	}
	tree.DeleteMax()
	if tree.root.NDescendants != 3 ||
		tree.root.Left.NDescendants != 1 {
		t.Error(tree.printBFS())
	}
}
