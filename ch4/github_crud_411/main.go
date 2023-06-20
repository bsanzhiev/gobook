package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const Template = "This is new issue title, should be just one line.\n" +
	"\nFollow the title is issue body, github support markdown format,\n" +
	"So you can write use markdown grammar. Template is from golang/go: \n\n" +
	"Please answer these questions before submitting your issue. Thanks! \n\n" +
	"### What versions of Go you using (`go version`)? \n\n" +
	"### Does this issue reproduce with the latest release?\n\n" +
	"### What OS and processor arch are you using (`go env`)?\n\n" +
	"### What did you do?\n\n" +
	"If possible, provide a recipe for reproducing the error.\n" +
	"A complete runnable program is good.\n" +
	"A link on play.golang.org is best.\n\n" +
	"### What did you expect to see?\n\n" +
	"### What did you see instead?\n\n"

const basicAuthURL = "https://appi.github.com/user"

// basic info
var userName string = "bsanzhiev"
var password string = "aRgyden@90"

// default repository, cam be set use `repo ` command
var repository string = "bsanzhiev/yetter"
var editor string = "vim"

func main() {
	// basic auth
	req, err := http.NewRequest("HEAD", basicAuthURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bsanzhiev: %v\n", err)
		os.Exit(1)
	}
	req.SetBasicAuth(userName, password)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bsanzhiev: %v\n", err)
		os.Exit(1)
	}
	if res.StatusCode == http.StatusUnauthorized {
		fmt.Fprintln(os.Stdout, "Opps: Your github username or password is not correct.")
		os.Exit(1)
	}

	// if output is 200, then basic auth is ok
	fmt.Fprintln(os.Stdout, "Welcome, Bator!")

	// command loop
	input := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, "=> ")
		line, _ := input.ReadString('\n')

		// delete '\n'
		line = line[:len(line)-1]

		// skip blank lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// cmd and optional arg(s)
		cmd := strings.Split(line, " ")

		switch cmd[0] {
		case "repo":
			if len(cmd) < 2 || cmd[1] == "" {
				fmt.Fprintf(os.Stdout, "Current Repository: %s\n", repository)
				continue
			}
			// repository check
			repository = cmd[1]
			fmt.Fprintf(os.Stdout, "Switch to Repository: %s\n", repository)
		case "editor":
			if len(cmd) < 2 || cmd[1] == "" {
				fmt.Fprintf(os.Stdout, "Current editor: %s\n", editor)
				continue
			}
		case "create":
			create()
		case "list":
			if len(cmd) == 2 {
				list(cmd[1])
			} else {
				list("")
			}
		case "read":
			if len(cmd) < 2 {
				fmt.Fprintf(os.Stderr, "read: lack of issue number")
				continue
			}

			if number, err := strconv.Atoi(cmd[1]); err != nil {
				fmt.Fprintln(os.Stderr, "read: issue id not a number")
			} else {
				read(number)
			}
		case "edit":
			if len(cmd) < 2 {
				fmt.Fprintln(os.Stderr, "edit: lack of issue number")
				continue
			}
			if number, err := strconv.Atoi(cmd[1]); err != nil {
				fmt.Fprintln(os.Stderr, "edit: issue id not a number")
			} else {
				edit(number)
			}
		case "close":
			if len(cmd) < 2 {
				fmt.Fprintf(os.Stderr, "edit: lack of issue number")
				continue
			}
			if number, err := strconv.Atoi(cmd[1]); err != nil {
				fmt.Fprintln(os.Stderr, "close: issue is not a number")
			} else {
				close(number)
			}
		case "quit":
			fallthrough
		case "exit":
			fmt.Fprintln(os.Stdout, "Bye! Have a nice day!")
			os.Exit(0)
		case "help":
			help()
		default:
			fmt.Fprintf(os.Stdout, "Oops... `%s` is invalid.\n", cmd[0])
		}
	}
}

// Issue struct
type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string
}

type NewIssue struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

// help display all supported command and usage massege
func help() {

}

func proxyRequest(method, targetUrl string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, targetUrl, body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "proxyRequest: %V\n", err)
		os.Exit(1)
	}
	// need basic auth in case of POST
	if method == "POST" || method == "PATCH" {
		req.SetBasicAuth(userName, password)
		// must be json format
		req.Header.Set("Content-Type", "application/json")
	}
	return http.DefaultClient.Do(req)
}

// list repo/issues
func list(state string) {
	issuesURL := "https://api.github.com/repos/" + repository + "/issues"

	if state == "all" || state == "open" || state == "closed" {
		issuesURL += "?state=" + state
	}

	res, err := proxyRequest("GET", issuesURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "list: %v\n", err)
		return
	}
	// most close res.Body on all exec path.
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Oops... %s\n", res.Status)
		fmt.Fprintf(os.Stderr, "Is %s a valid repository? Or we blocked by Github.\n", repository)
		return
	}

	var result []Issue
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		fmt.Fprintf(os.Stderr, "list: %v\n", err)
		return
	}

	// show issues (the first page or the current page)
	for _, item := range result {
		fmt.Printf("#%-5d %10.10s %s %.55s\n",
			item.Number, item.User.Login, item.CreatedAt, item.Title)
	}

}

// read/check secific issue
func read(number int) {
	issuesURL := "https://api.github.com/repos/" + repository + "/issues/" + strconv.Itoa(number)
	res, err := proxyRequest("GET", issuesURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read: %v\n", err)
		return
	}

	if res.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Is %s a valid reporitory, %d a valid issue?\n", repository, number)
		return
	}

	var result Issue
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		fmt.Fprintf(os.Stderr, "read: %v\n", err)
		return
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Title: %s\n", result.Title)
	fmt.Fprintf(os.Stdout, "From: %s %s\n", result.User.Login, result.CreatedAt)
	fmt.Fprintln(os.Stdout)
	// MarkDown format
	fmt.Fprintln(os.Stdout, result.Body)
}

func create() {
	postURL := "https://api.github.com/repos/" + repository + "/ussies"

	// create template
	tmpfile, err := ioutil.TempFile(os.TempDir(), "bsanzhiev-")
	if err != nil {
		log.Fatal(err)
	}
	tmpfile.WriteString(Template)
	tmpfile.Close()
	cmd := exec.Command(editor, tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	//parse issue file
	infile, err := os.Open(tmpfile.Name())
	if err != nil {
		log.Fatal(err)
	}

	// construct request body
	body := bytes.NewBuffer(make([]byte, 100))

	var issue NewIssue

	// set title
	reader := bufio.NewReader(infile)
	title, _ := reader.ReadString('\n')
	issue.Title = strings.TrimSpace(title)

	// set body
	coll := []byte{}
	b, e := reader.ReadByte()
	for e == nil {
		coll = append(coll, b)
		b, e = reader.ReadByte()
	}
	issue.Body = string(coll)

	data, err := json.Marshal(issue)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	body = bytes.NewBuffer(data)
	res, err := proxyRequest("POST", postURL, body)
	if err != nil {
		log.Fatalf("Post new Issue failed: %v\n", res.Status)
	}
}

// edit/update issues
func edit(number int) {
	// get original issue first
	// ---------------------------------------------------------------------------------------------------
	// same as read()
	issuesURL := "https://api.github.com/repos/" + repository + "/issues/" + strconv.Itoa(number)
	resp, err := proxyRequest("GET", issuesURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read: %v\n", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Is %s a valid reporitory, %d a valid issue?\n", repository, number)
		return
	}

	var result Issue
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Fprintf(os.Stderr, "read: %v\n", err)
		return
	}

	// edit/save, then push
	// same as create()
	// ---------------------------------------------------------------------------------------------------
	// create a template of issue, markdown format is better
	// create a temp file
	tmpfile, err := ioutil.TempFile(os.TempDir(), "trailsman-")
	if err != nil {
		log.Fatal(err)
	}

	tmpfile.WriteString(result.Title)
	tmpfile.Write([]byte("\n"))
	tmpfile.WriteString(result.Body)
	tmpfile.Close()

	// there should be a template for new issue
	// method is put template in the temp first
	// invoke user's favourite editor: e.g. vim or code
	// template file
	cmd := exec.Command(editor, tmpfile.Name())

	// support vim first, then code
	if editor == "vim" {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
	}

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	// parse issue file
	infile, err := os.Open(tmpfile.Name())
	if err != nil {
		log.Fatal(err)
	}

	// construct request body
	// 100 bytes
	body := bytes.NewBuffer(make([]byte, 100))

	//TODO: the best solution may just use json format file
	var issue NewIssue
	//issue.Title = "Another Message from trailsman!"
	//issue.Body = "你好，世界"

	// get title
	reader := bufio.NewReader(infile)
	title, _ := reader.ReadString('\n')
	issue.Title = strings.TrimSpace(title)
	//fmt.Fprintln(os.Stdout, issue.Title)

	// get body
	coll := []byte{}
	b, e := reader.ReadByte();
	for e == nil {
		coll = append(coll, b)
		b, e = reader.ReadByte()
	}
	issue.Body = string(coll)
	//fmt.Fprintln(os.Stdout, issue.Body)

	data, err :=
		json.Marshal(issue)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}

	body = bytes.NewBuffer(data)
	resp, err = proxyRequest("PATCH", issuesURL, body);
	if err != nil {
		log.Fatalf("Edit Issue failed: %v\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Edit Issue failed: %v\n", resp.Status)
	}
}

// close issue
// close is a kind of comment
func close(number int) {

}