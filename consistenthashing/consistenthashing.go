package consistenthashing

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type Hash func(data []byte) uint32

type Map struct {
	// hash算法函数
	hash Hash

	// 虚拟节点倍数
	replicas int

	// 排序的hash虚拟节点
	hashSortedNodes []uint32

	//已绑定的节点
	nodes map[string]bool

	// 虚拟节点对应节点信息
	hashMap map[uint32]string

	rwlock sync.RWMutex
}

func NewMap(replicas int, hashFunc Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     hashFunc,
		hashMap:  make(map[uint32]string),
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(node string) error {
	if node == "" {
		return nil
	}

	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	if m.nodes == nil {
		m.nodes = map[string]bool{}
	}

	if m.hashMap == nil {
		m.hashMap = map[uint32]string{}
	}

	if _, ok := m.nodes[node]; ok {
		return fmt.Errorf("node:%v, is already existed", node)
	}

	m.nodes[node] = true

	for i := 0; i < m.replicas; i++ {
		hash := m.hash([]byte(node + strconv.Itoa(i)))
		m.hashMap[hash] = node
		m.hashSortedNodes = append(m.hashSortedNodes, hash)
	}

	sort.Slice(m.hashSortedNodes, func(i, j int) bool {
		return m.hashSortedNodes[i] < m.hashSortedNodes[j]
	})

	return nil
}

// 根据key字段, 返回匹配中的serverIP
func (m *Map) Get(key string) string {
	m.rwlock.RLock()
	defer m.rwlock.RUnlock()

	hash := m.hash([]byte(key))

	i := sort.Search(len(m.hashSortedNodes), func(i int) bool {
		return m.hashSortedNodes[i] >= hash
	})

	return m.hashMap[m.hashSortedNodes[i%len(m.hashSortedNodes)]]
}
