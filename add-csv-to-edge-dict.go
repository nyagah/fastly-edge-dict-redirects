////////
//
// This code reads in a CSV file for source/desination URL pairs and adds them to
// to an Edge Dictionary.
//
// Ensure the GOPATH environment variable is set to your workspace directory. 
// Place this file in the "src" directory under your workspace directory. See more here: 
// https://github.com/golang/go/wiki/SettingGOPATH#unix-systems
//
// Commands:
//
// To compile: go build -o add-csv-to-edge-dict add-csv-to-edge-dict.go
// To run: ./add-csv-to-edge-dict <SECRET-API-TOKEN> <SERVICE-ID> <EDGE-DICT-NAME> <CVS-FILE-NAME>
//
//

package main

import "os"
import "io"
import "fmt"
import "log"
import "bufio"
import "encoding/csv"
import "github.com/sethvargo/go-fastly/fastly"

func main() {

   //////
   //
   // Check User Input.
   //
   if len(os.Args) != 5 { 
     log.Fatal("\nWrong arguments count: \n\n\t Usage is: ./add-csv-to-edge-dict <SECRET-API-TOKEN> <SERVICE-ID> <EDGE-DICT-NAME> <CVS-FILE-NAME>\n")
   } 

  
   ////////
   //
   // Create a client object. 
   //
   // The client has no state, so it can be persisted
   // and re-used. It is also safe to use concurrently due to its lack of state.
   // There is also a DefaultClient() method that reads an environment variable.
   // Please see the documentation for more information and details.
   //
   client, err := fastly.NewClient(os.Args[1])
   if err != nil {
     log.Fatal(err)
   }

   
   ////////
   //
   // Retrieve latest active version of service
   //
   var serviceID = os.Args[2]
   
   latest, err := client.LatestVersion(&fastly.LatestVersionInput{
     Service: serviceID,
   })
   if err != nil {
        log.Fatal(err)
   }


   ////////
   //
   // Check dictionary exists
   //
   var edgeDictName = os.Args[3]
   
   edgeDict, err := client.GetDictionary(&fastly.GetDictionaryInput{
     Service: serviceID,
     Version: latest.Number,
     Name: edgeDictName,
   })
   if err != nil {
     log.Fatal("Edge Dictionary called ", edgeDictName, " does not exist on this service")
   }
   

   ////////
   //
   // Add entries from CSV file to Edge Dictionary
   //
   csvFile, _ := os.Open(os.Args[4])
   reader := csv.NewReader(bufio.NewReader(csvFile))
 
   for {
       entry, err := reader.Read()
       if err == io.EOF {
           break
       } else if err != nil {
           log.Fatal(err)
       } else if len(entry[0]) == 0 {  // 1st column is source URL
         fmt.Println("Key is empty")
         continue
       } else if len(entry[2]) == 0 { // 3rd column is destination URL
         fmt.Println("Value is empty")
         continue
       }

       _, err_:= client.CreateDictionaryItem(&fastly.CreateDictionaryItemInput{
         Service: serviceID,
         Dictionary: edgeDict.ID,
         ItemKey: entry[0],
         ItemValue: entry[2],
       })
       if err_ != nil {
         fmt.Println("Something went wrong - maybe the entry key already exists.")
       } else {
         fmt.Println("Added: ", entry[0], "-->", entry[2])
       }
   }
   
}