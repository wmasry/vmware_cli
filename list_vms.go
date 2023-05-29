package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"crypto/tls"
	"fmt"
	"encoding/json"
	"flag"
	"os"
)

func getvms(authtoken string, address string, vmname string) {
        var result []map[string]interface{}
	var singleresult map[string]interface{}

	url := "https://"+address+"/api/vcenter/vm"
	singlevmurl := "https://"+address+"/api/vcenter/vm/"+vmname
	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
              }
        client := &http.Client{Transport: tr}
	authtoken = authtoken[1 : len(authtoken)-1]

	if (vmname == "all") || (vmname == "All") {
	  req, err := http.NewRequest("GET", url, nil)
	  req.Header.Set("vmware-api-session-id", authtoken)
	  resp, err := client.Do(req)
	  defer resp.Body.Close()
	  if err != nil{
		fmt.Println(err)
	  }
	  bodyText, err := ioutil.ReadAll(resp.Body)
	  fmt.Println(string(bodyText))
	  json.Unmarshal([]byte(bodyText), &result)
	  for _,value  := range result {
		fmt.Println("##################")
		for x, y := range(value){
			fmt.Println(x, ":", y)
	  }}
	}else if (vmname != "All") && (vmname != "all") && (vmname != "") {
	  fmt.Println("Single VM Section")
	  req, err := http.NewRequest("GET", singlevmurl, nil)
          req.Header.Set("vmware-api-session-id", authtoken)
          resp, err := client.Do(req)
	  defer resp.Body.Close()
          if err != nil{
                fmt.Println(err)
            }
          bodyText, err := ioutil.ReadAll(resp.Body)
	  json.Unmarshal([]byte(bodyText), &singleresult)
	  //Dont remove the following, for future use
	  //fmt.Println(singleresult["name"])
	  for key, value := range singleresult{
		  fmt.Println(key, ":", value)
		   fmt.Println("##################")
	  }
	  }else{
		fmt.Println("Usage Error, cannot list VM(s)")
  }
}
func basicAuth(username string, passwd string, address string) string {
    vcenterurl := "https://"+address+"/api/session"
    // Warning, the following line ignore the verification if client certificate
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}
    req, err := http.NewRequest("POST", vcenterurl, nil)
    req.SetBasicAuth(username, passwd)
    resp, err := client.Do(req)
    if err != nil{
        log.Fatal(err)
    }
    bodyText, err := ioutil.ReadAll(resp.Body)
    s := string(bodyText)
    return s
}

func main() {
	ListVMsCmd := flag.NewFlagSet("listvms", flag.ExitOnError)
	vmname := ListVMsCmd.String("vmname", "all", "Please Specify VM Name")
	username := ListVMsCmd.String("u", "", "Please specify username")
	password := ListVMsCmd.String("p", "", "Please specify password")
	url := ListVMsCmd.String("url", "", "Please specify IP/FQDN")


	switch os.Args[1]{
	case "listvms":
	    ListVMsCmd.Parse(os.Args[2:])
            authtoken := basicAuth(*username, *password,  *url )
            getvms(authtoken, *url, *vmname)
	default:
            fmt.Println("expected  subcommands")
            os.Exit(1)
    }

}

