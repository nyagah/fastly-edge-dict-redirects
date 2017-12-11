  # Overview
 
  This workflow describes how to create URL redirects using an Edge Dictionary on Fastly. 
  The source and destination URLs are defined in a CSV file. Go is used to read in the
  CSV file and to programmatically (via API) write the URLs into an Edge Dictionary
  on Fastly. The redirects will immediately take effect on Fastly after a source and 
  destination URL has been uploaded.
  
  There are 2 main Go files:
  
  1.) *create-edge-dict.go -> To create a new Edge Dictionary*  
  2.) *add-csv-to-edge-dict.go -> To upload new redirects from CVS file*

  # Installation  
  

  **Go Code**
  
  Download the Go source code into the $GOPATH/src folder:

```
  $ git clone git@github.com:nyagah/fastly-edge-dict-redirects.git
```
  Download the Fastly Go client library into the $GOPATH/src folder:

```
  $ go get github.com/sethvargo/go-fastly/fastly
```

  Ensure the GOPATH environment variable is set to your workspace directory.
  Place this file in the "src" directory under your workspace directory. See more here:
  https://github.com/golang/go/wiki/SettingGOPATH#unix-systems

  Compile the Go code:
 
 ```
  $ go build -o create-edge-dict create-edge-dict.go
  $ go build -o add-csv-to-edge-dict add-csv-to-edge-dict.go
```  


  **VCL Code**

  Add the code below at the top of `vcl_recv`:

  ```vcl
  set req.http.redir_location = table.lookup(<EDGE-DICT-NAME>, req.url, "")
 
  if (req.http.redir_location != "" ) {
       error 801 "Permanent Redirect";
   }
 ```

  Add the code below at the top of `vcl_error`:

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
  
  To create a new edge dictionary run the command below:
  
  ```
  $ ./create-edge-dict <SECRET-API-TOKEN> <SERVICE-ID> <EDGE-DICT-NAME>
  ```
  
  To add new redirects to an edge dictionary run the command below:

```
  $ ./add-csv-to-edge-dict <SECRET-API-TOKEN> <SERVICE-ID> <EDGE-DICT-NAME> <CVS-FILE-NAME>
```


  # NOTES
  1.) Edge Dictionaries have a limit of 1000 entires. Please reach out to Fastly to have that limit increased  
  2.) Dictionary item keys are limited 256 characters and values are limited to 8000 characters  
  3.) Dictionary item keys are case sensitive  
  4.) Event logs don't exist for Edge Dictionary changes  

