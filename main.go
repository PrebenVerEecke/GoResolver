package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "github.com/miekg/dns"
    "net"
    "time"
    "sync"
)

func main() {


    jobs := make(chan string)

    var wg sync.WaitGroup

    // Read entire file content, giving us little control but
    // making it very simple. No need to close the file.
    
    file, err := os.Open("../resolvers2.txt")
    if err != nil {
        log.Fatal(err)
    }
    
    for i := 0; i < 10; i++ {
		// tell the waitgroup about the new worker
		wg.Add(1)
		go func() {
			for nameserver := range jobs {
                //fmt.Println(nameserver)
				queryDnsServer("da.dwadawd-23123213.sc-corp.net", nameserver)
            }
            
            

			// when the jobs channel is closed the loop
			// above will stop; then we need to tell the
			// waitgroup that the worker is done
			wg.Done()
		}()
    }

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {

        //fmt.Println(scanner.Text())

        nameserver := scanner.Text()
        jobs <- nameserver
        //queryDnsServer("daidjwbdidbiawbdw.google.com", nameserver)
    }

    close(jobs)

    if err := scanner.Err(); err != nil {
        fmt.Println(err)
    }

    wg.Wait()
}

func queryDnsServer(target, nameserver string) {

    c := new(dns.Client)

    c.Dialer = &net.Dialer {
        Timeout: 300 * time.Millisecond,
    }

    m := dns.Msg{}

    m.SetQuestion(target+".", dns.TypeA)
    r, _, err := c.Exchange(&m, nameserver+":53")

    if err != nil {
        
        fmt.Println(err)
        //fmt.Println("Bad nameserver: ", nameserver)
        return
    }

    //log.Printf("Took %v", t)

    if len(r.Answer) == 0 {
        fmt.Println("Good NS: ", nameserver, )
    } else {
        fmt.Println("Bad NS: ", nameserver)
        fmt.Println("annoying server:", nameserver)
    }
   

    // for _, ans := range r.Answer {
    //     Arecord := ans.(*dns.A)
    //     log.Printf("%s", Arecord.A)
    // }
}
