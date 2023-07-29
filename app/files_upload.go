package app

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/andrewarrow/feedback/filestorage"
	"github.com/andrewarrow/feedback/router"
	"github.com/andrewarrow/feedback/util"
	"google.golang.org/api/option"
)

func handleFilePost(c *router.Context) {
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		router.SetFlash(c, err.Error())
		http.Redirect(c.Writer, c.Request, "/", 302)
		return
	}
	keyPath := ""
	client, err := filestorage.NewClient(context.Background(),
		option.WithCredentialsFile(keyPath))
	client.BucketPath = c.Router.BucketPath

	files := c.Request.MultipartForm.File["file"]

	for _, fileHeader := range files {
		name := fileHeader.Filename
		file, _ := fileHeader.Open()
		asBytes, _ := io.ReadAll(file)
		filename := putFile(client, name, asBytes)
		fmt.Println(filename)
	}
	http.Redirect(c.Writer, c.Request, "/", 302)
}

func putFile(client *filestorage.Client, name string, data []byte) string {
	if !strings.Contains(name, ".") {
		name = name + ".bin"
	}
	tokens := strings.Split(name, ".")
	ext := tokens[len(tokens)-1]
	guid := util.PseudoUuid()
	filename := guid + "." + ext

	bucket := ""
	w := client.Bucket(bucket).Object(filename).NewWriter(context.Background())
	w.ContentType = "application/octet-stream"
	_, err := w.Write(data)
	fmt.Println("write", err)
	w.Close()
	return filename
}
