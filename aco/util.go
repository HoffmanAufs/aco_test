package aco

import (
	"fmt"
	"math/rand"
	"time"
)

func StartDownloadTask(chunck_id int, node_id int, download_rate int, dl_delay int, ch chan TaskInfo){
	chunck_size := CHUNCK_SIZE
	elapsed := float64(chunck_size)*(8.)/float64(download_rate)
	// elapsed := 5
	// fmt.Printf("sleep: %v s\n", elapsed)

	fmt.Printf("-->: (c: %v, n: %v) at %vKB/s %.2fs...\n", chunck_id, node_id, download_rate, elapsed)
	// if dl_delay > 0{
	// 	time.Sleep(time.Duration(dl_delay)*time.Millisecond)
	// }
	ch <- TaskInfo{chunck_id, node_id, "start", download_rate}
	time.Sleep(time.Duration(100)*time.Millisecond)

	time.Sleep(time.Duration(elapsed)*time.Second)
	ch <- TaskInfo{chunck_id, node_id, "finish", download_rate}
	// fmt.Printf("-->: chunck_id: %v from node_id: %v at %v KB/s\n\n", chunck_id, node_id, download_rate)
	fmt.Printf("<<<<<<<<<<<<<<<<<<<<<< (c: %v, n: %v) \n", chunck_id, node_id)
}

func GenerateRandomChunckTable(node_id_list []int, chunck_count int) map[int]NodeAttribute{
	chunck_table := make(map[int]NodeAttribute)

	// init bandwidth
	// bandwidth_list := []int{512, 1_000, 2_000, 5_000, 10_000, 20_000, 20_000, 50_000, 50_000, 100_000}
	// bandwidth_list := []int{512, 2_500, 2_500, 5_000, 10_000, 10_000}
	bandwidth_list := []int{1500, 1600, 1700, 1800, 1900, 2000, 2100, 2200, 2300, 2400, 2500}
	// bandwidth_list := []int{2_000, 2_000, 2_000}
	bandwidth_list_len := len(bandwidth_list)

	// init bandwidth
	node_count := len(node_id_list)
	for i:=0; i<node_count;i++{
		chunck_bitvec := make([]bool, chunck_count)

		rand_index := rand.Int()%bandwidth_list_len
		band_width := bandwidth_list[rand_index]
		new_node := NodeAttribute{chunck_bitvec, band_width, 0}
		node_id := node_id_list[i]
		chunck_table[node_id] = new_node
	}

	// init chunck_bitvec
	half_true_count := node_count/2
	for chunck_index:=0; chunck_index<chunck_count;chunck_index++{
		true_count := 0
		for{
			rand_node_index := rand.Int()%node_count
			if !chunck_table[rand_node_index].chunck_bitvec[chunck_index] {
				chunck_table[rand_node_index].chunck_bitvec[chunck_index] = true
				true_count += 1
			}
			if true_count >= half_true_count{
				break
			}
		}
	}

	println("节点带宽信息： ")
	for i, d := range chunck_table{
		node_id := i
		bandwidth := d.bandwidth
		fmt.Printf("%d: %v\n", node_id, bandwidth)
	}


	return chunck_table
}


// func GenerateRandomChunckTable_b2(node_count int, chunck_count int) map[int]NodeAttribute{
// 	chunck_table := make(map[int]NodeAttribute)

// 	// init bandwidth
// 	bandwidth_list := []int{64_000, 256_000, 512_000, 1_000_000, 2_000_000, 5_000_000, 
// 		10_000_000, 20_000_000, 20_000_000, 50_000_000, 50_000_000, 100_000_000}
// 	bandwidth_list_len := len(bandwidth_list)
// 	for i:=0; i<node_count;i++{
// 		chunck_bitvec := make([]bool, chunck_count)

// 		rand_index := rand.Int()%bandwidth_list_len
// 		band_width := bandwidth_list[rand_index]
// 		new_node := NodeAttribute{chunck_bitvec, band_width, 0}
// 		chunck_table[i] = new_node
// 	}

// 	// init chunck_bitvec
// 	half_true_count := node_count/2
// 	for chunck_index:=0; chunck_index<chunck_count;chunck_index++{
// 		true_count := 0
// 		for{
// 			rand_node_index := rand.Int()%node_count
// 			if chunck_table[rand_node_index].chunck_bitvec[chunck_index] == false{
// 				chunck_table[rand_node_index].chunck_bitvec[chunck_index] = true
// 				true_count += 1
// 			}
// 			if true_count >= half_true_count{
// 				break
// 			}
// 		}
// 	}

// 	return chunck_table
// }

// func GenerateRandomChunckTable_b1(node_count int, chunck_count int) map[int][]bool{
// 	chunck_table := make(map[int][]bool)

// 	for i:=0; i<node_count;i++{
// 		chunck_table[i] = make([]bool, chunck_count)
// 	}

// 	half_node_count := node_count/2
// 	for chunck_index:=0; chunck_index<chunck_count;chunck_index++{
// 		true_count := 0
// 		for {
// 			rand_node_index := rand.Int()%node_count
// 			if chunck_table[rand_node_index][chunck_index] == false{
// 				chunck_table[rand_node_index][chunck_index] = true
// 				true_count += 1
// 			}
// 			if true_count >= half_node_count{
// 				break
// 			}
// 		}
// 	}
// 	return chunck_table
// }