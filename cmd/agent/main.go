package main

import (
	"GoPass/internal/agent/config"
	"GoPass/internal/agent/controller"
	"GoPass/internal/agent/records"
	"GoPass/internal/agent/transport"
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	cfg, _ := config.LoadConfig()
	ctx := context.Background()
	//p := fmt.Sprintf("%s/ping", cfg.Server)
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//go func() {
	//	err := transport.Ping(ctx, p, cfg.Ping)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}()
	//запуск
	fmt.Println(cfg.OS)
	fmt.Println("Greetings")
	//проверка соединения с сервером
	fmt.Println("Service available")
	//clearConsole(cfg.OS)
	startup(ctx, cfg)
}

func startup(ctx context.Context, cfg *config.ClientConfig) {
	fmt.Println("Menu: 1 to login; 2 to register; 3 to exit")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a command: ")
	command, _ := reader.ReadString('\n')
	command = strings.TrimSpace(command)
	switch command {
	case "1":
		login(ctx, cfg)
	case "2":
		register(ctx, cfg)
	case "3":
		exit()
	default:
		fmt.Println("Unknown command")
	}
}

func login(ctx context.Context, cfg *config.ClientConfig) {
	fmt.Println("Login")
	fmt.Print("Enter login: ")
	reader := bufio.NewReader(os.Stdin)
	login, _ := reader.ReadString('\n')
	login = strings.TrimSpace(login)
	fmt.Print("Enter password: ")
	password, _ := terminal.ReadPassword(syscall.Stdin)
	cookie, err := controller.Login(ctx, cfg, login, string(password))
	if err != nil {
		fmt.Printf("%e", err)
		startup(ctx, cfg)
	}
	command(ctx, *cfg, cookie)
}

func register(ctx context.Context, cfg *config.ClientConfig) {
	fmt.Print("Enter login: ")
	reader := bufio.NewReader(os.Stdin)
	login, _ := reader.ReadString('\n')
	login = strings.TrimSpace(login)
	fmt.Print("Enter password: ")
	password, _ := terminal.ReadPassword(syscall.Stdin)
	cookie, err := controller.Register(ctx, cfg, login, string(password))
	if err != nil {
		fmt.Printf("%e", err)
		startup(ctx, cfg)
	}
	command(ctx, *cfg, cookie)
}

func logout(ctx context.Context, cfg *config.ClientConfig) {
	startup(ctx, cfg)
}

func exit() {
	os.Exit(0)
}

func command(ctx context.Context, cfg config.ClientConfig, cookies []*http.Cookie) {
	fmt.Println("Welcome!")
	run := true
	for run {
		fmt.Println("Menu: 1 to open; 2 to create ; 3 to edit; 4 to delete; 5 to exit;")
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter a command: ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		switch command {
		case "1":
			// вывод  записи полностью
			fmt.Println("List of records")
			text, err := transport.GetList(ctx, cfg, cookies)
			if err != nil {
				fmt.Printf("%e", err)
			}
			fmt.Println(string(text))
		case "2":
			record := records.Record{}
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Enter Name: ")
			name, _ := reader.ReadString('\n')
			record.Name = strings.TrimSpace(name)

			fmt.Print("Enter Site: ")
			site, _ := reader.ReadString('\n')
			record.Site = strings.TrimSpace(site)

			fmt.Print("Enter Login: ")
			login, _ := reader.ReadString('\n')
			record.Login = strings.TrimSpace(login)

			fmt.Print("Enter Password: ")
			password, _ := reader.ReadString('\n')
			record.Password = strings.TrimSpace(password)

			fmt.Print("Enter Info: ")
			info, _ := reader.ReadString('\n')
			record.Info = strings.TrimSpace(info)
			err := transport.CreateRecord(ctx, cfg, record, cookies)
			if err != nil {
				fmt.Printf("%e", err)
			}
		case "3":
			record := records.Record{}
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Enter ID of a record: ")
			idStr, _ := reader.ReadString('\n')
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				fmt.Println("Invalid ID format:", err)
				return
			}
			record.ID = id

			fmt.Print("Enter Name: ")
			name, _ := reader.ReadString('\n')
			record.Name = strings.TrimSpace(name)

			fmt.Print("Enter Site: ")
			site, _ := reader.ReadString('\n')
			record.Site = strings.TrimSpace(site)

			fmt.Print("Enter Login: ")
			login, _ := reader.ReadString('\n')
			record.Login = strings.TrimSpace(login)

			fmt.Print("Enter Password: ")
			password, _ := reader.ReadString('\n')
			record.Password = strings.TrimSpace(password)

			fmt.Print("Enter Info: ")
			info, _ := reader.ReadString('\n')
			record.Info = strings.TrimSpace(info)
			err = transport.UpdateRecord(ctx, cfg, record, cookies)
			if err != nil {
				fmt.Printf("%e", err)
			}
		case "4":
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter ID: ")
			idStr, _ := reader.ReadString('\n')
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				fmt.Println("Invalid ID format:", err)
				return
			}
			err = transport.DeleteRecord(ctx, &cfg, cookies, id)
			if err != nil {
				fmt.Printf("%e", err)
			}
		case "5":
			run = false
			logout(ctx, &cfg)
		default:
			fmt.Println("Unknown command")
		}
	}
}

func clearConsole(platform string) {
	var cmd *exec.Cmd

	switch platform {
	case "unix":
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}
