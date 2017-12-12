  # Overview
 
  This workflow describes how to deploy URL redirects using an Edge Dictionary on Fastly. 
  The source and destination URLs are defined in a CSV file and Go is used to read in the
  file and to programmatically (via API) write the URLs into an Edge Dictionary
  on Fastly. Note that the redirects will immediately take effect after the API call to add a URL pair
  is complete.
  
  There are 2 main Go files:
  
  1.) *create-edge-dict.go -> To create a new Edge Dictionary (shouldn't be used that often)*  
  2.) *add-csv-to-edge-dict.go -> To upload new redirects from CVS file*

  # Installation  
  

  **Go Code**
  
  Create a workspace directory and set the `$GOPATH` environment variable to point to it. See more [here]( https://github.com/golang/go/wiki/SettingGOPATH#unix-systems). In this workspace directory create a directory called `src`. Then do the following:
  
  1.) Download the Go source code into the `$GOPATH/src` folder:

```
  $ git clone git@github.com:nyagah/fastly-edge-dict-redirects.git
```
  2.) Download the Fastly Go client library into the `$GOPATH/src` folder:

```
  $ go get github.com/sethvargo/go-fastly/fastly
```  

  3.) Compile the Go code `$GOPATH/src/fastly-edge-dict-redirects` folder:
 
 ```
  $ go build -o create-edge-dict create-edge-dict.go
  $ go build -o add-csv-to-edge-dict add-csv-to-edge-dict.go
```  


  **VCL Code**

  1.) Add the code below at the top of `vcl_recv`:

  ```vcl
  # NOTE: Make sure to replace <EDGE-DICT-NAME> with actual name of Edge Dictionary
  #
  set req.http.redir_location = table.lookup(<EDGE-DICT-NAME>, req.url, "")
 
  if (req.http.redir_location != "" ) {
       error 801 "Permanent Redirect";
   }
 ```

  2.) Add the code below at the top of `vcl_error`:

  ```vcl
  # Permanent Redirects
  #
  if (obj.status == 801) {
     set obj.status = 301;
     set obj.response = "Moved Permanently";
     set obj.http.Location = req.http.redir_location;
     synthetic {""};
     return (deliver);
  }  

  # Temporary Redirects
  #
  if (obj.status == 802) {
     set obj.status = 302;
     set obj.response = "Found";
     set obj.http.Location = req.http.redir_location;
     synthetic {""};
     return (deliver);
  }
```
  
  # Usage
  
  1.) To create a new edge dictionary run the command below (shouldn't be used that often):
  
  ```
  $ ./create-edge-dict <SECRET-API-TOKEN> <SERVICE-ID> <EDGE-DICT-NAME>
  ```
  
  2.) To add new redirects to an edge dictionary run the command below:
  
```
$ ./add-csv-to-edge-dict <SECRET-API-TOKEN> <SERVICE-ID> <EDGE-DICT-NAME> <CVS-FILE-NAME>
```


  # NOTES
  1.) Edge Dictionaries have a limit of 1000 entires. Please reach out to Fastly to have that limit increased  
  2.) Dictionary item keys are limited 256 characters and values are limited to 8000 characters  
  3.) Dictionary item keys are case sensitive  
  4.) Event logs don't exist for Edge Dictionary changes  

