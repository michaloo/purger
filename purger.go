package main

import (
    "github.com/mitchellh/goamz/aws"
    "github.com/mitchellh/goamz/s3"
    "os"
    "sort"
    "strconv"
)

type byLastModified []s3.Key
func (v byLastModified) Len() int { return len(v) }
func (v byLastModified) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v byLastModified) Less(i, j int) bool { return v[i].LastModified < v[j].LastModified }

func main() {

    count, err := strconv.ParseInt(os.Args[1], 10, 0)

    if count == 0 {
        count = 10
    }

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
        println("Nothing to purge")
        os.Exit(0)
    }

    sort.Sort(byLastModified(arr.Contents))

    toDelete := arr.Contents[0:len(arr.Contents) - int(count)]

    println("Deleting", len(toDelete), "items")

    for n := 0; n < len(toDelete); n ++ {

        delErr := bucket.Del(toDelete[n].Key)

        if delErr != nil {
            panic("Error in deleting")
        }

        println("Deleted", n + 1, toDelete[n].Key)
    }
}
