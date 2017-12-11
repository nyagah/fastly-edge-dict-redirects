////////
//
// This code creates a new Edge Dictionary in a Fastly service.
//
// Ensure the GOPATH environment variable is set to your workspace directory. 
// Place this file in the "src" directory under your workspace directory. See more here: 
// https://github.com/golang/go/wiki/SettingGOPATH#unix-systems
//
// Commands:
//
// To compile: go build -o create-edge-dict create-edge-dict.go
// To run: ./create-edge-dict <SECRET-API-TOKEN> <SERVICE-ID> <EDGE-DICT-NAME>
//
//

package main

import "os"
import "fmt"
import "log"
import "github.com/sethvargo/go-fastly/fastly"

func main() {

   //////
   //
   // Check User Input.
   //
   if len(os.Args) != 4 { 
     log.Fatal("\nWrong arguments count: \n\n\t Usage is: ./create-edge-dict <SECRET-API-TOKEN> <SERVICE-ID> <EDGE-DICT-NAME>\n")
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
   if err == nil {
     log.Fatal("   Edge Dictionary already exists --> Name: ", edgeDict.Name, " ID: ", edgeDict.ID)
   }
   

   ////////
   //
   // Clone the latest version 
   // 
   // So we can make changes without affecting the active configuration.
   //
   version, err := client.CloneVersion(&fastly.CloneVersionInput{
     Service: serviceID,
     Version: latest.Number,
   })
   if err != nil {
     log.Fatal(err)
   }

    
   ////////
   //
   // Make changes to newly cloned version. 
   //
   // In this case, we create a new edge dictionary
   //
   edgeDictNew, err := client.CreateDictionary(&fastly.CreateDictionaryInput{
     Service: serviceID,
     Version: version.Number,
     Name: edgeDictName,
   })
   if err != nil {
     log.Fatal(err)
   }

   // Output: Edge Dictionary Name and ID
   fmt.Println("Edge Dictionary Added --> Name: ", edgeDictNew.Name, "ID: ", edgeDictNew.ID)
   
   
   ////////
   //
   // Validate that version changes are  valid.
   //
   valid, msg, err := client.ValidateVersion(&fastly.ValidateVersionInput{
     Service: serviceID,
     Version: version.Number,
   })
   if err != nil {
     log.Fatal(err)
   }
   if !valid {
     log.Fatal("not valid version", msg)
   }
   

   ////////
   //
   // Finally, activate this new version.
   //
   activeVersion, err := client.ActivateVersion(&fastly.ActivateVersionInput{
     Service: serviceID,
     Version: version.Number,
   })
   if err != nil {
     log.Fatal(err)
   }
   
   // Output: to confirm activation of changes
   if activeVersion.Active == true {
     fmt.Printf("Version Number %d is activated\n", activeVersion.Number)
   } else {
     log.Fatal("Version Number ", activeVersion.Number, " is NOT activated!!!")
   }


}