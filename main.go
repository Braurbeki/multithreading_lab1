package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"runtime"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Path to root folder is expected")
	}
	res_dict_single := benchmark(single_thread_count)
	res_dict_multi := benchmark(multi_thread_count)
	fmt.Println(res_dict_single)
	if !reflect.DeepEqual(res_dict_single, res_dict_multi) {
		fmt.Errorf("Maps are not the same!!!")
	}
}

func benchmark(a func() map[string]int) map[string]int {
	start := time.Now()
	res := a()
	elapsed := time.Since(start).Milliseconds()
	fmt.Printf("Time passed for %s: %dms\n", runtime.FuncForPC(reflect.ValueOf(a).Pointer()).Name(), elapsed)
	return res
}

func single_thread_count() map[string]int {
	res := make(map[string]int)
	list_files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		fmt.Errorf("Could not read folder tree", err)
	}
	for _, file := range list_files {
		path_to_file := path.Join(os.Args[1], file.Name())
		count := count_symbols(path_to_file, nil)
		merge(&res, &count)
		if err != nil {
			fmt.Errorf("Could not count symbols for " + path_to_file, err)
		}
	}
	return res
}

func multi_thread_count() map[string]int {
	symbols := make(chan map[string]int)
	res := make(map[string]int)
	list_files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		fmt.Errorf("Could not read folder tree", err)
	}
	for _, file := range list_files {
		path_to_file := path.Join(os.Args[1], file.Name())
		go count_symbols(path_to_file, symbols)
		if err != nil {
			fmt.Errorf("Could not count symbols for " + path_to_file, err)
		}
	}
	for i := 0; i < len(list_files); i++ {
		file_symbols := <-symbols
		merge(&res, &file_symbols)
	}
	return res
}

func merge(a, b *map[string]int) {
	for key, val := range *b {
		if _, exists := (*a)[key]; exists {
			(*a)[key] += val
		} else {
			(*a)[key] = val
		}
	}
}

func count_symbols(filename string, c chan map[string]int) map[string]int {
	result := make(map[string]int)
	file, err := os.Open(filename)
    if err != nil {
        return nil
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        val := scanner.Text()
		_, has_val := result[val]
		if has_val {
			result[val] += 1
		} else {
			result[val] = 1
		}
    }
	if c != nil {
		c <- result
	}
	return result
}