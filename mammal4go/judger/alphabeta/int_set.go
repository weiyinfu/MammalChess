package alphabeta

import (
	"github.com/weiyinfu/MammalChess/mammal4go/util"
	"sort"
)

type IntSet struct {
	//无序集合
	a        []int
	ma       []int //value=>下标的映射
	end      int
	maxSize  int
	maxValue int
}

func NewIntSet(maxSize int, maxValue int) *IntSet {
	x := IntSet{
		a:        make([]int, maxSize),
		ma:       make([]int, maxSize),
		end:      0,
		maxSize:  maxSize,
		maxValue: maxValue,
	}
	for i, _ := range x.ma {
		x.ma[i] = -1
	}
	return &x
}
func (self *IntSet) Get() []int {
	return self.a[:self.end]
}
func (self *IntSet) Add(v int) {
	if self.ma[v] != -1 {
		return
	}
	self.ma[v] = self.end
	self.a[self.end] = v
	self.end++
}
func (self *IntSet) MultiAdd(v []int) {
	for _, i := range v {
		self.Add(i)
	}
}
func (self *IntSet) Update(k int, v int) {
	//把旧的k更新为v
	ind := self.ma[k]
	self.a[ind] = v
	self.ma[k] = -1
	self.ma[v] = ind
}
func (self *IntSet) Remove(v int) {
	pos := self.ma[v]
	if pos == -1 {
		return
	}
	self.ma[v] = -1
	//把末尾元素移动到当前位置
	self.a[pos] = self.a[self.end-1]
	if pos != self.end-1 {
		//如果是不同元素才更新ma
		self.ma[self.a[pos]] = pos
	}
	self.end--
}

func (self *IntSet) Size() int {
	return self.end
}

func (self *IntSet) Eq(x *IntSet) bool {
	mine := util.ArrayCopy(self.Get())
	his := util.ArrayCopy(x.Get())
	if len(mine) != len(his) {
		return false
	}
	sort.Slice(mine, func(i, j int) bool {
		return mine[i] < mine[j]
	})
	sort.Slice(his, func(i, j int) bool {
		return his[i] < his[j]
	})
	for ind, _ := range mine {
		if mine[ind] != his[ind] {
			return false
		}
	}
	return true
}
