package core

import (
	"time"
	//"net/http"
)

type LruItem struct{
	Key string
	Value string
	UpdateTs int64
	PrevItem *LruItem
	NextItem *LruItem
}

type LruList struct{
	Head *LruItem
	Tail *LruItem
	Length uint32

}

type Lru struct{
	LruDict map[string]*LruItem
	Lst	*LruList
}

func (lst *LruList)IsEmpty()bool{
	if lst.Length==0 {
		return true
	}
	return false
}

func NewLru()*Lru{

}

func NewLruList() *LruList{
	lst:=&LruList{
		Length:		0,
		Head:		nil,
		Tail:		nil,
	}
	return lst
}

func NewLruItem(key string) *LruItem{
	item:=&LruItem{
		Key:		key,
		UpdateTs:	time.Now().Unix(),
		PrevItem:	nil,
		NextItem:	nil,
	}
	return item
}

// push的时候，从尾部追加
func (l *LruList)Push(item *LruItem)uint32{
	if l.Length==0 {
		l.Length=1
		l.Head=item
		l.Tail=item
		item.NextItem=nil
		item.PrevItem=nil
		return l.Length
	}
	l.Tail.NextItem=item
	item.PrevItem=l.Tail
	l.Tail=item
	l.Length=l.Length+1
	return l.Length
}

// pop的时候，从头弹出来
func (l *LruList)Pop() *LruItem{
	if l.Length==0{
		return nil
	}
	if l.Length==1{
		tempItem:=l.Head
		l.Head=nil
		l.Tail=nil
		l.Length=0
		return tempItem
	}
	tempItem:=l.Head
	l.Head=l.Head.NextItem
	l.Length=l.Length-1
	return tempItem
}

























