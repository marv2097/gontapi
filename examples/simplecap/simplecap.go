/* Simple Capture 
* 
* Opens a Napatech interface and prints raw packet data to the screen. 
* Should only be used for testing or the basis of a more complete program.
* 
* Note: napatech api and suitable hardware should already be installed on your system.
* Tested with a NT20E3-2 (PCIe Gen3 2x10Gb SFP+) runing v3.10.
*
* To get - go get github.com/marv2097/gontapi
* To build - go build github.com/marv2097/gontapi/examples/simplecap
* To run - sudo ./simplecap -c 1
*
*/

package main

import (
    "flag"
    "fmt"
    "log"
    "github.com/marv2097/gontapi"
)


// Command Line Args
var port = flag.String("p", "0", "Port to read packets from")
var maxcount = flag.Int("c", -1, "Only grab this many packets, then exit")

func main() {
    flag.Parse()

    // Initialise NT API
    if err := ntapi.NtInit(); err != nil {
        log.Fatalln(err)
    }

    if err := ntapi.NtConfigOpen("config"); err != nil {
        log.Fatalln(err)
    }

    // Clear out existing filters
    if err, _ := ntapi.NtNtpl("Delete = ALL"); err != nil {
        log.Fatalln(err)
    }

    // Define filter to get packets from a specific port
    if err, _ := ntapi.NtNtpl("Define FilterPort = Filter(Port=="+*port+")"); err != nil {
        log.Fatalln(err)
    }

    // Assign the filter
    if err, ntplInfo := ntapi.NtNtpl("Assign[streamid=1] = FilterPort"); err != nil {
        // Unable to Assign, print the details
        log.Println(err, "\n\n    " + ntplInfo.ErrDesc[0] + "    " +ntplInfo.ErrDesc[1] + "\n    " + ntplInfo.ErrDesc[2] + "\n")
        log.Fatalln("Quitting")
    } else {
        // Get the ID
        log.Println("Filter ID:", ntplInfo.NtplId)
    }

    // Open config
    if err := ntapi.NtNetRxOpen("config"); err != nil {
        log.Fatalln(err)
    }

    // Get packets
    data := make([]byte, 1522)
    count := int(0)

    for true {
        if ci, err := ntapi.NtNetRxGetTo(data); err != nil{
            log.Fatalln(err)
        } else {
            count++

            // Print the Capture Info and Raw data
            fmt.Println("-----------------------------------------------------------------")
            fmt.Println(ci)
            fmt.Println(data[:ci.CaptureLength])

            // Check if over the packet limit
            limit := *maxcount > 0 && count >= *maxcount
            if limit {
                break
            }
        }
    }
    log.Println("Done...")
}
