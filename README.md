# Cystack Task

This project consists of a client and server written in Go that gather and process information about the operating system they are installed on.

## Project Components

The client is an executable that runs as a service in Windows, gathering the following information:
- Windows hostname
- IP address
- Memory utilization
- C disk utilization
- All local users in Windows
- All running processes in Windows
- All installed applications in the control panel

This information is stored in a SQLite database and sent to the server for processing.

## The following features have been implemented:

- Create a client written in Go
- Gather information about the operating system and store it in a SQLite database

## The following tasks remain to be completed:

- Send the data from the client to the server for processing
- Build and deploy the client and server executables
- Install the client as a service in Windows

## Dependencies

The following dependencies are required to build and run this project:
- Go
- sqlite3
