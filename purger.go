package main

import (
    "github.com/mitchellh/goamz/aws"
    "github.com/mitchellh/goamz/s3"
    "os"
    "sort"
    // "strconv"
    "flag"
    "fmt"
)

type byLastModified []s3.Key
func (v byLastModified) Len() int { return len(v) }
func (v byLastModified) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v byLastModified) Less(i, j int) bool { return v[i].LastModified < v[j].LastModified }

func main() {

    var count int
    var view bool
    var sort_key string

    flag.IntVar(&count, "count", 10, "number of files to left")
    flag.BoolVar(&view, "view", false, "just show the files")
    flag.StringVar(&sort_key, "sort", "asc", "define sort - asc or desc")

    flag.Parse()

    auth := aws.Auth{
        os.Getenv("S3_KEY"),
        os.Getenv("S3_SECRET"),
        "",
    }

    region := aws.Region{
        "dreamobjects",
        "",
        "https://" + os.Getenv("S3_ENDPOINT"),
        "",
        true,
        true,
        "",
        "",
        "",
        "",
        "",
        "",
        "",
        "",
    }

    new_s3 := s3.New(auth, region)

    bucket := s3.Bucket{
        new_s3,
        os.Getenv("S3_BUCKET"),
    }

    arr, err := bucket.List(os.Getenv("S3_PATH"), "", "", 0)

    if err != nil {
        panic("Error in listing")
    }

    if len(arr.Contents) < int(count) {
        fmt.Println("Nothing to purge")
        os.Exit(0)
    }

    if sort_key == "asc" {
        sort.Sort(byLastModified(arr.Contents))
    }

    if sort_key == "desc" {
        sort.Sort(sort.Reverse(byLastModified(arr.Contents)))
    }

    toDelete := arr.Contents[0:len(arr.Contents) - int(count)]

    if view == true {
        for n := 0; n < len(toDelete); n ++ {
            fmt.Println(toDelete[n].Key, toDelete[n].LastModified)
        }
        os.Exit(0)
    }

    fmt.Println("Deleting", len(toDelete), "items")

    for n := 0; n < len(toDelete); n ++ {

        delErr := bucket.Del(toDelete[n].Key)

        if delErr != nil {
            panic("Error in deleting")
        }

        fmt.Println("Deleted", n + 1, toDelete[n].Key)
    }
}
