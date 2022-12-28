package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/sys/windows"
)

func main() {

	hostname := hostname()
	ipAddresses := ipAddress()
	ipAddress := ipAddresses[0]
	memoryUtilization := memoryUtilization()
	diskUtilization := CDiskUtilization()
	localUsers := localUsers()
	runningProcesses := runningProcesses()

	fmt.Println("Hostname:", hostname)
	fmt.Println("IP address:", ipAddress)
	fmt.Printf("Memory utilization: %.2f%%\n", memoryUtilization)
	fmt.Printf("C disk utilization: %.2f%%\n", diskUtilization)
	fmt.Print(localUsers)
	fmt.Print(runningProcesses)

	sqlLite(hostname, ipAddress, memoryUtilization, diskUtilization, localUsers, runningProcesses)
}

func hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
	}
	return hostname
}

func ipAddress() []string {
	var ips []string
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println(err)
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Println(err)
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			ips = append(ips, ip.String())
		}
	}
	return ips
}

func memoryUtilization() float64 {
	var mem runtime.MemStats

	runtime.ReadMemStats(&mem)

	memoryUtilization := float64(mem.TotalAlloc) / float64(mem.Sys) * 100
	return memoryUtilization
}

func CDiskUtilization() float64 {
	var freeBytes, totalBytes, totalFreeBytes uint64
	dir := "C:\\"
	err := windows.GetDiskFreeSpaceEx(syscall.StringToUTF16Ptr(dir), &freeBytes, &totalBytes, &totalFreeBytes)
	if err != nil {
		log.Println(err)
	}

	diskUtilization := float64(totalBytes-freeBytes) / float64(totalBytes)
	return diskUtilization * 100
}

func localUsers() string {
	cmd := exec.Command("net", "user")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	return out.String()
}

func runningProcesses() string {
	cmd := exec.Command("tasklist")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	return out.String()
}

func sqlLite(hostname string, ipAddress string, memoryUtilization float64, diskUtilization float64, localUsers string, runningProcesses string) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS data (
		id INTEGER PRIMARY KEY,
		hostname TEXT,
		ip_address TEXT,
		memory_utilization REAL,
		c_disk_utilization REAL,
		local_users TEXT,
		running_processes TEXT
	)`)
	if err != nil {
		log.Println(err)
	}

	_, err = db.Exec(`INSERT INTO data (hostname, ip_address, memory_utilization, c_disk_utilization, local_users, running_processes)
    VALUES (?, ?, ?, ?, ?, ?)`, hostname, ipAddress, memoryUtilization, diskUtilization, localUsers, runningProcesses)
	if err != nil {
		log.Println(err)
	}
}
