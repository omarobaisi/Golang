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
	"strings"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/sys/windows"
)

func main() {

	hostname := hostname()
	ipAddress := ipAddress()
	memoryUtilization := memoryUtilization()
	diskUtilization := CDiskUtilization()
	localUsers := localUsers()
	runningProcesses := runningProcesses()
	installedApplications := installedApplications()

	fmt.Println("Hostname:", hostname)
	fmt.Println("IP address:", ipAddress)
	fmt.Printf("Memory utilization: %.2f%%\n", memoryUtilization)
	fmt.Printf("C disk utilization: %.2f%%\n", diskUtilization)
	fmt.Print(localUsers)
	fmt.Print(runningProcesses)
	fmt.Print(installedApplications)

	sqlLite(hostname, ipAddress, memoryUtilization, diskUtilization, localUsers, runningProcesses, installedApplications)
}

func hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
	}
	return hostname
}

func ipAddress() string {
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
	return fmt.Sprintf("%v", ips)
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
		log.Println(err)
	}

	return out.String()
}

func runningProcesses() string {
	cmd := exec.Command("tasklist")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}

	return out.String()
}

func installedApplications() string {
	cmd := exec.Command("wmic", "product", "get", "name")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}

	// Split the output by newline characters to get a list of installed applications
	apps := strings.Split(out.String(), "\n")
	// Remove the first and last elements, which are the column headers and an empty string
	apps = apps[1 : len(apps)-1]
	return fmt.Sprintf("%v", apps)
}

func sqlLite(hostname string, ipAddress string, memoryUtilization float64, diskUtilization float64, localUsers string, runningProcesses string, installedApplications string) {
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
		running_processes TEXT,
		installed_applications TEXT
	)`)
	if err != nil {
		log.Println(err)
	}

	_, err = db.Exec(`INSERT INTO data (hostname, ip_address, memory_utilization, c_disk_utilization, local_users, running_processes, installed_applications)
    VALUES (?, ?, ?, ?, ?, ?)`, hostname, ipAddress, memoryUtilization, diskUtilization, localUsers, runningProcesses, installedApplications)
	if err != nil {
		log.Println(err)
	}
}
