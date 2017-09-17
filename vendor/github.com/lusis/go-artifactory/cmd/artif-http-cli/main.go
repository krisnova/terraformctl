package main

import (
	"bytes"
	"fmt"

	"os"

	artifactory "github.com/lusis/go-artifactory/artifactory.v51"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app         = kingpin.New("http-cli", "Make an http request to the artifactory server.")
	contentType = app.Flag("content_type", "set the content-type for the request").Short('t').Default("application/json").String()

	get       = app.Command("get", "make an http get request")
	getPath   = get.Arg("path", "path to get").Required().String()
	getParams = get.Flag("query_params", "key=value query parameter. specify multiple times if neccessary").Short('p').PlaceHolder("KEY=VALUE").StringMap()

	post       = app.Command("post", "make an http post request")
	postPath   = post.Arg("path", "path to post").Required().String()
	postParams = post.Flag("query_params", "key=value query parameter. specify multiple times if neccessary").Short('p').PlaceHolder("KEY=VALUE").StringMap()
	postFile   = post.Flag("file", "Full /path/to/data/to/post").Short('f').File()
	postBody   = post.Flag("body", "contents of data to post").Short('b').String()

	put       = app.Command("put", "make an http put request")
	putPath   = put.Arg("path", "path to put").Required().String()
	putParams = put.Flag("query_params", "key=value query parameter. specify multiple times if neccessary").Short('p').PlaceHolder("KEY=VALUE").StringMap()
	putFile   = put.Flag("file", "Full /path/to/data/to/post").Short('f').File()
	putBody   = put.Flag("body", "contents of data to post").Short('b').String()

	delete     = app.Command("delete", "make an http delete request")
	deletePath = delete.Arg("path", "path to delete").Required().String()
)

func main() {
	client := artifactory.NewClientFromEnv()
	var request artifactory.Request

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case get.FullCommand():
		request.Verb = "GET"
		request.Path = *getPath
		request.QueryParams = *getParams
		request.ContentType = *contentType
		data, err := client.HTTPRequest(request)
		if err != nil {
			app.FatalIfError(err, "", nil)
		}
		fmt.Printf("%s\n", string(data))
	case post.FullCommand():
		if *postBody != "" && *postFile != nil {
			app.Fatalf("Cannot specify both file and body at the same time")
		} else {
			if *postBody != "" {
				request.Body = bytes.NewBufferString(*postBody)
			}
			if *postFile != nil {
				request.Body = *postFile
			}
			request.Verb = "POST"
			request.Path = *postPath
			request.QueryParams = *postParams
			request.ContentType = *contentType
			data, err := client.HTTPRequest(request)
			//data, err := client.Post(*postPath, postData, *postParams)
			if err != nil {
				app.FatalIfError(err, "", nil)
			} else {
				fmt.Printf("%s\n", string(data))
			}
		}
	case put.FullCommand():
		request.Verb = "PUT"
		request.Path = *putPath
		request.QueryParams = *putParams
		request.ContentType = *contentType
		if *putBody != "" && *putFile != nil {
			app.Fatalf("Cannot specify both file and body at the same time")
		} else {
			if *putBody != "" {
				request.Body = bytes.NewBufferString(*putBody)
			}
			if *putFile != nil {
				request.Body = *putFile
			}
			data, err := client.HTTPRequest(request)
			//data, err := client.Put(*putPath, putData, *putParams)
			if err != nil {
				app.FatalIfError(err, "", nil)
			} else {
				fmt.Printf("%s\n", string(data))
			}
		}
	case delete.FullCommand():
		err := client.Delete(*deletePath)
		if err != nil {
			app.FatalIfError(err, "", nil)
		}
		fmt.Println("OK")
	}
}
