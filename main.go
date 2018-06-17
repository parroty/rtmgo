package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/parroty/rtmgo/rtm"
)

var (
	mode         = flag.String("mode", "getTasks", "")
	taskID       = flag.String("task-id", "", "")
	listID       = flag.String("list-id", "", "")
	taskseriesID = flag.String("taskseries-id", "", "")
	due          = flag.String("due", "", "")
	filter       = flag.String("filter", "", "")
	name         = flag.String("name", "", "")
	tags         = flag.String("tags", "", "")
	test         = flag.Bool("test", false, "")
	errorMode    = flag.Bool("error", false, "")
)

func init() {
	flag.Parse()
}

func main() {
	fmt.Println("rtmgo api examples")
	token := os.Getenv("RTM_TOKEN")
	apiKey := os.Getenv("RTM_API_KEY")
	client := rtm.NewClient(token, apiKey)

	if *test {
		mux := http.NewServeMux()
		mux.HandleFunc(
			"/services/rest",
			func(w http.ResponseWriter, r *http.Request) {
				q := r.URL.Query()
				method := q["method"]
				fmt.Println(r.URL)
				fmt.Println()
				data, err := ioutil.ReadFile("testdata/" + method[0] + ".json")
				if err != nil {
					panic(err)
				}
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, string(data))
			},
		)
		ts := httptest.NewServer(mux)
		defer ts.Close()
		fmt.Println(ts.URL)
		client.HTTP.BaseURL = ts.URL
	}

	timeline, err := client.Timelines.Create()
	if err != nil {
		panic(err)
	}
	fmt.Println(timeline)

	if *mode == "getTasks" {
		tasks, _ := client.Tasks.GetList("", *filter, "")
		for _, task := range tasks {
			fmt.Println(task.Name)
		}
	} else if *mode == "addTask" {
		if *errorMode {
			*name = ""
		}
		tasks, transaction, err := client.Tasks.Add(timeline, *name)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, task := range tasks {
			fmt.Println(task.Name)
		}
		fmt.Println(transaction)
	} else if *mode == "deleteTask" {
		tasks, _ := client.Tasks.GetList("", *filter, "")
		if len(tasks) > 0 {
			task := tasks[0]
			fmt.Println(task)
			if *errorMode {
				task.Chunks[0].ID = ""
			}
			client.Tasks.Delete(task, timeline)
		}
	} else if *mode == "modifyTasks" {
		tasks, _ := client.Tasks.GetList("", *filter, "")
		if len(tasks) > 0 {
			task := tasks[0]
			fmt.Println(task.Name)
			client.Tasks.SetDueDate(tasks[0], timeline, *due)
			client.Tasks.Complete(tasks[0], timeline)
			client.Tasks.Uncomplete(tasks[0], timeline)
		}
	} else if *mode == "addTags" {
		tasks, _ := client.Tasks.GetList("", *filter, "")
		task := tasks[0]
		transactions, _ := client.Tasks.AddTags(task, timeline, *tags)
		fmt.Println(transactions)
	} else if *mode == "getLists" {
		lists, _ := client.Lists.GetList()
		for _, list := range lists {
			fmt.Println(list.Name)
		}
	} else if *mode == "deleteList" {
		lists, _ := client.Lists.GetList()
		for _, list := range lists {
			if list.Name == *name {
				fmt.Println(client.Lists.Delete(list, timeline))
			}
		}
	} else if *mode == "addList" {
		list, _, err := client.Lists.Add(*name, timeline)
		if err != nil {
			fmt.Println("Add List error")
			fmt.Println(err)
		} else {
			fmt.Println(list.Name)
		}
	} else if *mode == "createTimeline" {
		timeline, _ := client.Timelines.Create()
		fmt.Println(timeline)
	} else if *mode == "getFrob" {
		frob, _ := client.Auth.GetFrob()
		fmt.Println(frob)
		perms := []string{"delete"}
		url := client.GetAuthURL(frob, perms)
		fmt.Println(url)
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Access the url above and then input enter to get token: ")
		text, _ := reader.ReadString('\n')
		fmt.Println(text)
		result, _ := client.Auth.GetToken(frob)
		fmt.Println(result.Token)
		fmt.Println(result.User)
	}
}
