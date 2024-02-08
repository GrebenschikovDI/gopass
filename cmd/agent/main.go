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
	//запуск
	fmt.Println("Greetings")
	//проверка соединения с сервером
	fmt.Println("Service available")
	//запрос логина и пароля
	//first frame
	fmt.Println("Menu: 1 to login; 2 to register; 3 to exit")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a command: ")
	command, _ := reader.ReadString('\n')
	command = strings.TrimSpace(command)
	switch command {
	case "1":
		login(ctx, cfg)
	case "2":
		register()
	case "3":
		exit()
	default:
		fmt.Println("Unknown command")
	}

	//доступ к командам
	//выполнение команд
}

func login(ctx context.Context, cfg *config.ClientConfig) {
	fmt.Print("Enter login: ")
	reader := bufio.NewReader(os.Stdin)
	login, _ := reader.ReadString('\n')
	login = strings.TrimSpace(login)
	fmt.Print("Enter password: ")
	password, _ := terminal.ReadPassword(syscall.Stdin)
	cookie, err := controller.Login(ctx, cfg, login, string(password))
	if err != nil {
		fmt.Printf("%e", err)
	}
	command(ctx, *cfg, cookie)
}

func register() {
	fmt.Print("Enter login: ")
	reader := bufio.NewReader(os.Stdin)
	login, _ := reader.ReadString('\n')
	login = strings.TrimSpace(login)
	// проверить логин
	fmt.Print("Enter password: ")
	password, _ := terminal.ReadPassword(syscall.Stdin)
	// проверить пароль
	fmt.Println("\nPassword is", string(password))
	//отправка данных на сервер и соединение
	//command()
}

func logout() {
	// закрытие соединения
	// выход на первый фрейм
	//startup()

}

func exit() {
	// доабавить graceful shutdown?
	os.Exit(0)
}

func command(ctx context.Context, cfg config.ClientConfig, cookies []*http.Cookie) {
	//clearConsole("unix")
	//screen.MoveTopLeft()
	fmt.Println("Welcome!")
	for {
		fmt.Println("Menu: 1 to open; 2 to create ; 3 to edit; 4 logout;")
		// вывод имен записей
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
			// add new
		case "3":
		// edit
		case "4":
			logout()
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
