package consistenthashing

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"testing"
)

func TestNewMap(t *testing.T) {
	// 虚拟节点数列表
	replicasList := []int{100, 150, 200}

	// 测试真实服务器器10台
	nodeNum := 10

	// 测试100w数据分布
	testCount := 1000000

	for _, v := range replicasList {
		m := NewMap(100, nil)
		m.replicas = v
		// 每台服务器上，数据分布
		resultMap := make(map[string]int64)

		for i := 1; i <= nodeNum; i++ {
			serverIP := "172.17.0." + strconv.Itoa(i)
			m.Add(serverIP)
			resultMap[serverIP] = 0
		}

		for i := 0; i < testCount; i++ {
			testName := "testName"
			serverName := m.Get(testName + strconv.Itoa(i))
			resultMap[serverName] ++
		}

		fmt.Printf("###测试%d个节点, 一个节点有%d个虚拟节点, %d条测试数据###\n", nodeNum, m.replicas, testCount)

		var keys []string
		var value []float64
		for k, v := range resultMap {
			keys = append(keys, k)
			value = append(value, float64(v))
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("服务器地址:%v, 分布数据:%d\n", k, resultMap[k])
		}
		fmt.Printf("标准差:%f\n", getStandardDeviation(value))
	}

}

//获取标准差
func getStandardDeviation(list []float64) float64 {
	var total float64
	for _, item := range list {
		total += item
	}
	//平均值
	avg := total / float64(len(list))

	var dTotal float64
	for _, value := range list {
		dValue := value - avg
		dTotal += dValue * dValue
	}

	return math.Sqrt(dTotal / avg)
}
