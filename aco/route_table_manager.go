package aco

import (
	"fmt"
	"math"
	// "sync"
	"time"
	"math/rand"
)

type RouteMetric struct {
	delta_rate  float64
	update_time time.Time
	used        int
}

type RouteTableManager struct {
	route_table   map[int]RouteMetric
	average_rate  float64
	full_rate     float64
	update_period int 	// 30s
	to_zero_count int 	// 120
	// rt_lock sync.Mutex
}

func NewRouteTableManager(node_id_list []int, average_rate int, bandwidth int) RouteTableManager{
	route_table := make(map[int]RouteMetric)
	for _, i := range node_id_list{
		route_table[i] = RouteMetric{0.0, time.Unix(0, 0), 0}
	}

	return RouteTableManager{
		route_table,
		float64(average_rate),
		float64(bandwidth),
		5,
		20,
		// sync.Mutex{},
	}
}

func (rt *RouteTableManager) getRouteMetric(route_id int) float64 {
	// rt.rt_lock.Lock()
	// defer rt.rt_lock.Unlock()

	route_item, ok := rt.route_table[route_id]
	// if route_item.used > 0 || ok == false {
	// 	return 0.
	// }
	if !ok {
		return 0.
	}

	if route_item.used > 0{
		return 0.
	}

	if route_item.delta_rate == 0 {
		return rt.average_rate
	}

	elapsed := int(time.Since(route_item.update_time).Seconds())
	if elapsed >= rt.update_period*rt.to_zero_count {
		return rt.average_rate
	}

	rest_delta_rate := (1. - float64(elapsed/rt.update_period)/float64(rt.to_zero_count)) * route_item.delta_rate
	cur_rate := rt.average_rate + rest_delta_rate
	return cur_rate
}

func (rt *RouteTableManager) SelectACONode(node_id_list []int) int {
	return rt.SelectACORoute(node_id_list)
}

func (rt *RouteTableManager) FetchAllRouteMetric() map[int]float64{
	node_metric := make(map[int]float64)

	// rt.rt_lock.Lock()
	// defer rt.rt_lock.Unlock()
	var cur_rate float64
	for k, v := range rt.route_table{
		node_id := k
		route_item := v

		if route_item.delta_rate == 0{
			cur_rate = rt.average_rate
		}else{
			elapsed := int(time.Since(route_item.update_time).Seconds())
			if elapsed >= rt.update_period * rt.to_zero_count{
				cur_rate = rt.average_rate
			}else{
				rest_delta_rate := (1. - float64(elapsed/rt.update_period)/float64(rt.to_zero_count)) * route_item.delta_rate
				cur_rate = rt.average_rate + rest_delta_rate
			}
		}
		cur_metric := cur_rate
		node_metric[node_id] = cur_metric
	}
	return node_metric
}

func (rt *RouteTableManager) SelectACORoute(route_id_list []int) int {
	metric_map := make(map[int]float64)

	metric_sum := 0.
	for _, route_id := range route_id_list {
		cur_metric := rt.getRouteMetric(route_id)
		metric_sum += cur_metric
		metric_map[route_id] = cur_metric
	}

	rand.Seed(int64(time.Now().Nanosecond()))
	select_value := rand.Float64() * metric_sum

	prev_metric_sum := 0.
	for route_id, metric := range metric_map {
		if metric == 0 {
			continue
		}
		cur_metric_sum := prev_metric_sum + metric
		if prev_metric_sum <= select_value && select_value < cur_metric_sum {
			return route_id
		}
		prev_metric_sum = cur_metric_sum
	}

	return route_id_list[len(route_id_list)-1]
}

func (rt *RouteTableManager) OnTaskStart(route_id int) {
	// rt.rt_lock.Lock()
	// defer rt.rt_lock.Unlock()

	prev_item := rt.route_table[route_id]
	rt.route_table[route_id] = RouteMetric{prev_item.delta_rate, prev_item.update_time, prev_item.used + 1}
}

func (rt *RouteTableManager) OnTaskFinish(route_id int, download_rate int) {

	var delta_rate float64
	dl_rate := float64(download_rate)
	if dl_rate > rt.average_rate {
		delta_rate = math.Min(dl_rate, rt.full_rate) - rt.average_rate
	} else {
		delta_rate = math.Max(0, dl_rate) - rt.average_rate
	}

	cur_time := time.Now()
	new_item := RouteMetric{delta_rate, cur_time, 0}

	// rt.rt_lock.Lock()
	// defer rt.rt_lock.Unlock()
	rt.route_table[route_id] = new_item
}

func (rt *RouteTableManager) Test(){
	// var route_id int
	// var metric float64

	// route_id = 0
	// metric = rt.getRouteMetric(route_id)
	// fmt.Printf("%v: %.0f\n", route_id, metric)

	// route_id = 1
	// metric = rt.getRouteMetric(route_id)
	// fmt.Printf("%v: %.0f\n", route_id, metric)

	// route_id = 2
	// metric = rt.getRouteMetric(route_id)
	// fmt.Printf("%v: %.0f\n", route_id, metric)

	// var metric float64
	// for route_id := range rt.route_table{
	// 	metric = rt.getRouteMetric(route_id)
	// 	fmt.Printf("%v: %.0f\n", route_id, metric)
	// }

	// keys := make([]int, 0, len(rt.route_table))
	// for k := range rt.route_table {
	// 	keys = append(keys, k)
	// }
	
	// println()
	// route_id := len(rt.route_table)+1
	// metric = rt.getRouteMetric(route_id)
	// fmt.Printf("%v: %.0f\n", route_id, metric)

	ret_table := make(map[int]int)
	for i:=0; i<100000; i++{
		route_id_list := []int{0,1,2,3}
		ret_id := rt.SelectACORoute(route_id_list)
		// fmt.Printf("select route: %v\n", ret_id)

		prev_count := ret_table[ret_id]
		ret_table[ret_id] = prev_count+1
	}

	for k,v := range ret_table{
		route_id, count := k, v
		fmt.Printf("%v: %v\n", route_id, count)
	}
}

func test01(){
	// route_table := make(map[int]RouteItem)
	route_table := map[int]RouteMetric{
		0: {3_000_000, time.Now().Add(time.Duration(-50)*time.Minute), 0},
		1: {3_000_000, time.Now().Add(time.Duration(-40)*time.Minute), 0},
		2: {3_000_000, time.Now().Add(time.Duration(-30)*time.Minute), 0},
		3: {3_000_000, time.Now().Add(time.Duration(-20)*time.Minute), 0},
		4: {3_000_000, time.Now().Add(time.Duration(-10)*time.Minute), 0},
		5: {3_000_000, time.Now().Add(time.Duration(-1)*time.Minute), 0},
		// 2: {3_000_000, time.Now().Add(time.Duration(-10)*time.Minute), 1},
		// 3: {3_000_000, time.Now().Add(time.Duration(-100)*time.Minute), 0},
		// 4: {-200_000, time.Now().Add(time.Duration(-100)*time.Minute), 0},
		// 5: {-200_000, time.Now().Add(time.Duration(-1)*time.Minute), 0},
		// 6: {-2_000_000, time.Now().Add(time.Duration(-10)*time.Minute), 0},
		// 7: {-2_000_000, time.Now().Add(time.Duration(-20)*time.Minute), 0},
		// 8: {-2_000_000, time.Now().Add(time.Duration(-40)*time.Minute), 0},
		// 9: {-2_000_000, time.Now().Add(time.Duration(-60)*time.Minute), 0},
		// 10: {-2_000_000, time.Now().Add(time.Duration(-80)*time.Minute), 0},
		// 1: RouteItem{0, 1, time.Unix(0, 0)},
	}

	rt := RouteTableManager{
		route_table: route_table,
		average_rate: 5_000_000, // 3Mbps
		full_rate: 10_000_000,	 // 10Mbps
		to_zero_count: 120,		 // 120*30s = 1h
		update_period: 30,		 // 30s
	}

	rt.Test()
}

func test02(){
	now := time.Now()
	fmt.Printf("%v\n", now)
	next := now.Add(time.Duration(-100)*time.Minute)
	fmt.Printf("%v\n", next)

	// t := time.Parse("2016-01-03", "2021-08-01")
	t := time.Date(2020,2,2,0,0,0,0,time.Local)
	fmt.Printf("%v\n", t)
}

// func main(){
// 	test01()
// }