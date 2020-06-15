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
		//t.Log("outerIdx:", outerIdx)
		_ = outerIdx
		tree := New()
		for i := 1; i <= 10; i++ {
			tree.InsertNoReplace(shuffle[i-1])
		}
		//t.Log(tree.stringBFS())
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
			t.Error(tree.stringBFS())
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
	deleted := tree.Delete(Int(18))
	if deleted == nil {
		t.Fatal()
	}
	deleted = tree.Delete(Int(17))
	if deleted != nil {
		t.Fatal()
	}
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
	//t.Log(tree.stringBFS())
	deleted = tree.Delete(Int(12))
	if deleted == nil {
		t.Fatal()
	}
	if tree.root.NDescendants != 8 || // item 8
		tree.root.Right.NDescendants != 4 || // item 16
		tree.root.Right.Right.NDescendants != 1 || // item 20
		tree.root.Right.Left.NDescendants != 2 { // item 14
		t.Error(tree.stringBFS())
	}
}

func TestLLRB_Delete2(t *testing.T) {
	tree := New()
	for _, e := range []Int{32, 12, 4, 24, 2} {
		tree.InsertNoReplace(e)
	}
	tree.Delete(Int(4))
	for _, e := range []Int{14, 8, 36, 20, 34} {
		tree.InsertNoReplace(e)
	}
	tree.Delete(Int(20))
	for _, e := range []Int{40, 16, 30, 28, 26} {
		tree.InsertNoReplace(e)
	}
	tree.Delete(Int(26))
	tree.InsertNoReplace(Int(20))
	for _, e := range []Int{10, 38, 22, 18, 6} {
		tree.InsertNoReplace(e)
	}
	tree.Delete(Int(38))
	for _, c := range []struct {
		item Int
		rank int
	}{
		{item: 2, rank: 1}, {item: 6, rank: 2}, {item: 8, rank: 3},
		{item: 10, rank: 4}, {item: 12, rank: 5}, {item: 14, rank: 6},
		{item: 16, rank: 7}, {item: 18, rank: 8}, {item: 20, rank: 9},
		{item: 22, rank: 10}, {item: 24, rank: 11}, {item: 28, rank: 12},
		{item: 30, rank: 13}, {item: 32, rank: 14}, {item: 34, rank: 15},
		{item: 36, rank: 16}, {item: 40, rank: 17},
	} {
		r, _ := tree.GetRankOf(c.item)
		if r != c.rank {
			t.Errorf("item: %v, expected: %v, reality: %v", c.item, c.rank, r)
		}
	}
}

func TestLLRB_DeleteMin(t *testing.T) {
	tree := New()
	tree.InsertNoReplace(Int(2))
	tree.InsertNoReplace(Int(4))
	//t.Log(tree.stringBFS())
	tree.DeleteMin()
	if tree.root.NDescendants != 1 {
		t.Errorf("expect %v, reality: %v", 1, tree.root.NDescendants)
		t.Log(tree.stringBFS())
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
		t.Error(tree.stringBFS())
	}

	tree.DeleteMin()
	if tree.root.NDescendants != 3 ||
		tree.root.Left.NDescendants != 1 {
		t.Error(tree.stringBFS())
	}
}

func TestLLRB_DeleteMax(t *testing.T) {
	tree := New()
	tree.InsertNoReplace(Int(2))
	tree.InsertNoReplace(Int(4))
	//t.Log(tree.stringBFS())
	tree.DeleteMax()
	if tree.root.NDescendants != 1 {
		t.Errorf(tree.stringBFS())
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
		t.Error(tree.stringBFS())
	}
	tree.DeleteMax()
	if tree.root.NDescendants != 3 ||
		tree.root.Left.NDescendants != 1 {
		t.Error(tree.stringBFS())
	}
}

func TestLLRB_ReplaceOrInsert(t *testing.T) {
	array := []Int{22, 11, 13, 6, 9, 12, 4, 7, 1, 23, 13, 1, 18, 10, 18, 8, 19,
		16, 15, 4, 19, 10, 24, 2, 9, 9, 5, 25, 1, 6, 2, 7, 18, 20, 4}
	// sorted(set(a)): [1, 2, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 15, 16, 18, 19, 20, 22, 23, 24, 25]

	tree := New()
	for _, e := range array {
		tree.ReplaceOrInsert(e)
	}
	if tree.root.NDescendants != 21 {
		t.Errorf("tree.root.NDescendants: r: %v, e: %v", tree.root.NDescendants, 21)
	}
	if r, _ := tree.GetRankOf(Int(25)); r != 21 {
		t.Errorf("RankOf(Int(25)): r: %v, e: %v", r, 21)
	}
	if r, _ := tree.GetRankOf(Int(4)); r != 3 {
		t.Errorf("RankOf(Int(4)): r: %v, e: %v", r, 3)
	}
	if r, _ := tree.GetRankOf(Int(16)); r != 14 {
		t.Errorf("RankOf(Int(25)): r: %v, e: %v", r, 14)
	}
	if r, _ := tree.GetRankOf(Int(3)); r != 3 {
		t.Errorf("RankOf(Int(3)): r: %v, e: %v", r, 3)
	}
}
