package main

//Made by chkrr00k
//GPLv3 license

import (
	"fmt"
	"net/textproto"
	"os"
	"bufio"
	"strings"
)

func connectTo(name string, port int) (bool, *textproto.Conn){
	sock, err := textproto.Dial("tcp", fmt.Sprintf("%s:%d", name, port))
	if err != nil {
		return false, sock
	}else{
		return true, sock
	}
}

func pingHandler(input string) string {
	return fmt.Sprintf("PONG %s", input[5:])
}
func joinString(chanName string) string {
	return fmt.Sprintf("JOIN %s\r\n", chanName)
}
func connectString(nick string, ident string, realname string) string{
	return fmt.Sprintf("NICK %s\r\nUSER %s %s ayy :%s", nick, ident, realname)
}

func writeSock(sock *textproto.Conn, msg string){
	sock.Writer.PrintfLine("%s", msg)
}

type parser func(string) (bool, string)

func privmsgHandler(input string) (bool, string){
	res := strings.Split(input, " ")
	if res[1] != "PRIVMSG"{
		return false, ""
	}
	return true, fmt.Sprintf("%s <%s> %s", res[2], strings.Split(res[0], "!")[0][1:], strings.Join(res[3:], " "))

}
func noticeHandler(input string) (bool, string){
	res := strings.Split(input, " ")
	if res[1] != "NOTICE"{
		return false, ""
	}
	return true, fmt.Sprintf("%s <%s> %s", res[2], strings.Split(res[0], "!")[0][1:], strings.Join(res[3:], " "))
}
func joinHandler(input string) (bool, string){
	res := strings.Split(input, " ")
	if res[1] != "JOIN"{
		return false, ""
	}
	return true, fmt.Sprintf("%s joined %s", strings.Split(res[0], "!")[0][1:], res[2])
}
func quitHandler(input string) (bool, string){
	res := strings.Split(input, " ")
	if res[1] != "QUIT"{
		return false, ""
	}
	return true, fmt.Sprintf("%s has quttied", strings.Split(res[0], "!")[0][1:])
}
func leftHandler(input string) (bool, string){
	res := strings.Split(input, " ")
	if res[1] != "LEFT"{
		return false, ""
	}
	return true, fmt.Sprintf("%s has left %s", strings.Split(res[0], "!")[0][1:], res[2])
}
func modeHandler(input string) (bool, string){
	res := strings.Split(input, " ")
	if res[1] != "MODE"{
		return false, ""
	}
	return true, fmt.Sprintf("%s has set modes %s %s", strings.Split(res[0], "!")[0][1:], res[2], res[3])
}
func handle(funcs []parser, line string) string{
	var b bool
	var result string
	for _, f := range funcs{
		b, result = f(line)
		if b {
			return result
		}
	}
	return line
}

func listener(conn *textproto.Conn) {
	var num int = 0
	for keep {
		if num < 6 {
			num++
		}else if num == 6 {
			writeSock(conn, joinString("#nietzsche"))
			num++
		}

		line, _ := conn.Reader.ReadLine()
		if line == "" {
			break
		}
		if strings.HasPrefix(line, "PING") {
			writeSock(conn, pingHandler(line))
		}
		fmt.Println(handle([]parser{privmsgHandler,joinHandler, quitHandler, noticeHandler,leftHandler, modeHandler}, line))
	}
}
var keep bool
var channel string

func sendPM(input string) (bool, string){
	if strings.HasPrefix(input, "message "){
		return true, fmt.Sprintf("PRIVMSG %s :%s", channel, input[8:])
	}
	return false, ""
}
func selectChan(input string) (bool, string){
	if strings.HasPrefix(input, "chan "){
		channel = input[5:]
	}
	return false, ""
}
func joinChan(input string) (bool, string){
	if strings.HasPrefix(input, "join "){
		return true, joinString(input[5:])
	}
	return false, ""
}

func consoleHandler(sock *textproto.Conn, funcs []parser, input string){
	var b bool
	var msg string
	for _, f := range funcs{
		b, msg = f(input)
		if b {
			writeSock(sock, msg)
			break
		}
	}
}

func main(){
	fmt.Println("Program started")
	_, conn := connectTo("irc.rizon.net", 6667)

	console := bufio.NewScanner(os.Stdin)
	writeSock(conn, connectString("chkrB0t", "ident", "realname"))
	keep = true

	go listener(conn)
	
	var input string
	for console.Scan() {
		input = console.Text()
		if strings.HasPrefix(input, "quit"){
			break
		}else{
			consoleHandler(conn, []parser{sendPM, selectChan, joinChan}, input)
		}
	}
	if console.Err() != nil{
		fmt.Println("Everyone died and it's your fault")
	}
	keep = false
	conn.Close()
}
