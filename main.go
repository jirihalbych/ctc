package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

type Config struct {
	Cars      Cars      `json:"cars"`
	Stations  []Station `json:"stations"`
	Registers Registers `json:"registers"`
}

type Registers struct {
	Count     int `json:"count"`
	HandleMin int `json:"handle_time_min"`
	HandleMax int `json:"handle_time_max"`
	Queue     chan Car
	Times     []time.Duration
	TotalTime time.Duration
	Cars      int
}

type Cars struct {
	Count      int `json:"count"`
	ArrivalMin int `json:"arrival_time_min"`
	ArrivalMax int `json:"arrival_time_max"`
}

type Car struct {
	Type string
}

type Station struct {
	Type      string `json:"type"`
	Count     int    `json:"count"`
	ServeMin  int    `json:"serve_time_min"`
	ServeMax  int    `json:"serve_time_max"`
	Queue     chan Car
	Times     []time.Duration
	TotalTime time.Duration
	Cars      int
}

func GetMaxDuration(times []time.Duration) time.Duration {
	maximum := time.Duration(0)
	for i := 0; i < len(times)-1; i++ {
		if times[i] > maximum {
			maximum = times[i]
		}
	}
	return maximum
}

func GetStation(stations []Station, station string) *Station {
	for v := range stations {
		if stations[v].Type == station {
			return &stations[v]
		}
	}
	return &Station{}
}

func UseStation(car Car, stations *[]Station) {
	start := time.Now()
	defer wg.Done()
	station := GetStation(*stations, car.Type)
	cas := time.Duration(rand.Intn(station.ServeMax-station.ServeMin+1)+station.ServeMin) * time.Millisecond
	station.Cars = station.Cars + 1
	station.Queue <- car
	done := time.Since(start)
	station.Times = append(station.Times, done)
	station.TotalTime = station.TotalTime + done
	time.Sleep(cas)
	<-station.Queue
}

func UseRegister(car Car, registers *Registers) {
	start := time.Now()
	defer wg.Done()
	cas := time.Duration(rand.Intn(registers.HandleMax-registers.HandleMin+1)+registers.HandleMin) * time.Millisecond
	registers.Cars = registers.Cars + 1
	registers.Queue <- car
	done := time.Since(start)
	registers.Times = append(registers.Times, done)
	registers.TotalTime = registers.TotalTime + done
	time.Sleep(cas)
	<-registers.Queue
}

func main() {
	var types []string
	jsonFile, err := os.ReadFile("config.conf")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened config")
	var config Config
	var stations []Station
	var cars Cars
	var registers Registers
	err = json.Unmarshal(jsonFile, &config)
	stations = config.Stations
	cars = config.Cars
	registers = config.Registers
	if err != nil {
		fmt.Println("error:", err)
	}
	for i := 0; i < len(stations); i++ {
		stations[i].Queue = make(chan Car, stations[i].Count)
		types = append(types, stations[i].Type)
	}
	for i := 0; i < registers.Count; i++ {
		registers.Queue = make(chan Car, registers.Count)
	}
	for i := 0; i < cars.Count; i++ {
		wg.Add(1)
		cas := time.Duration(rand.Intn(cars.ArrivalMax-cars.ArrivalMin+1)+cars.ArrivalMin) * time.Millisecond
		time.Sleep(cas)
		car := Car{types[rand.Intn(len(types))]}
		go UseStation(car, &stations)
		wg.Add(1)
		go UseRegister(car, &registers)
	}
	wg.Wait()
	fmt.Println("stations:")
	for i := 0; i < len(stations); i++ {
		fmt.Println("\t", stations[i].Type, ":")
		fmt.Println("\t\ttotal_cars: ", stations[i].Cars)
		fmt.Println("\t\ttotal_queue_time: ", stations[i].TotalTime)
		fmt.Println("\t\tmax_queue_time: ", GetMaxDuration(stations[i].Times))
		fmt.Println("\t\tavg_queue_time: ", time.Duration(float64(float64(stations[i].TotalTime.Milliseconds()))/float64(len(stations[i].Times))*float64(time.Millisecond)))

	}
	fmt.Println("registers:")
	fmt.Println("\ttotal_cars: ", registers.Cars)
	fmt.Println("\ttotal_queue_time: ", registers.TotalTime)
	fmt.Println("\tmax_queue_time: ", GetMaxDuration(registers.Times))
	fmt.Println("\tavg_queue_time: ", time.Duration(float64(float64(registers.TotalTime.Milliseconds()))/float64(len(registers.Times))*float64(time.Millisecond)))
}
