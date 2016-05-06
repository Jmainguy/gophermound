package main

import (
    "fmt"
    "net"
    "bufio"
    "os"
    "github.com/gin-gonic/gin"
    "io/ioutil"
    "time"
    "strings"
    "strconv"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func zk2file(url, command, file string) {
    // Truncate file
    f, err := os.Create(file)
    check(err)
    defer f.Close()
    // Connect to ZK
    //conn, err := net.Dial("tcp", "localhost:32770")
    conn, err := net.Dial("tcp", url)
    check(err)
    // stat
    fmt.Fprintf(conn, command)
    // Write lines to file
    w := bufio.NewWriter(f)
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        line := scanner.Text() + "\n"
        w.WriteString(line)
    }
    w.Flush()
    f.Sync()
    conn.Close()
}

func zkconnections(url string) (connections string) {
    // Connect to ZK
    conn, err := net.Dial("tcp", url)
    check(err)
    // stat
    fmt.Fprintf(conn, "stat")
    // Parse out just connections line
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        if strings.Contains(scanner.Text(), "Connections:") {
            connections := strings.Split(scanner.Text(), " ")[1]
            return connections
        }
    }
    return
}


var zkurl = "localhost:2181"
//var zkurl = "localhost:32770"

func main() {
    go func () {
        for {
              // stat
              zk2file(zkurl, "stat\n", "/opt/gomound/stat")
              // mntr
              zk2file(zkurl, "mntr\n", "/opt/gomound/mntr")
              // ruok
              zk2file(zkurl, "ruok\n", "/opt/gomound/ruok")
              time.Sleep(5 * time.Second)
        }
    }()
    validurl := map[string]bool {
        "stat": true,
        "mntr": true,
        "ruok": true,
    }
    r := gin.Default()
    r.GET("/:url", func(c *gin.Context) {
        url := c.Param("url")
        if url == "connections" {
            connections := zkconnections(zkurl)
            i, err := strconv.Atoi(connections)
            check(err)
            // Change this to a reasonable value if we want to serioussly start using this
            if i > 20000 {
                var msg = "Warning, too many connections, Connections = " + connections
                c.String(406, msg)
            } else {
                var msg = "Connections = " + connections
                c.String(200, msg)
            }
        } else if validurl[url] {
            f := "/opt/gomound/" + url
            file, err := ioutil.ReadFile(f)
            check(err)
            message := string(file)
            c.String(200, message)
        } else {
            c.String(404, "Not Found")
        }
    })
    r.Run(":8080") // listen and server on 0.0.0.0:8080
}

