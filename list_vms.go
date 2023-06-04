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
	"bytes"
	"text/tabwriter"
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
          json.Unmarshal([]byte(bodyText), &result)

	  w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, '.', tabwriter.AlignRight|tabwriter.Debug)
	  for _,value  := range result {
	     fmt.Fprintln(w, "Name" , "\t", "VM-UniqueID" ,"\t",  "Power_State", "\t" ,"Memory_Size_MiB", "\t" , "CPU_Count" )
	     w.Flush()
             fmt.Fprintln(w,  value["name"], "\t", value["vm"], "\t", value["power_state"],
	     "\t", value["memory_size_MiB"], "\t", value["cpu_count"] )
	     w.Flush()
        }

	}else if (vmname != "All") && (vmname != "all") && (vmname != "") {
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
		   fmt.Println("--------------")
	  }
	  }else{
		fmt.Println("Usage Error, cannot list VM(s)")
  }
}
func basicAuth(username string, passwd string, address string) string {
    vcenterurl := "https://"+address+"/api/session"
    // Warning, the following line ignore the verification of client certificate
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

func createvm(authtoken string, address string) string{

	url := "https://"+address+"/api/vcenter/vm"

	var createvm = []byte(`{
		"name": "golangtest4",
		"guest_OS": "DOS",
		"placement": {
			"datastore": "datastore-3024",
			"folder": "group-v4",
			"resource_pool": "resgroup-3022"			
			}
		}`)

        tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
              }
        client := &http.Client{Transport: tr}
        authtoken = authtoken[1 : len(authtoken)-1]

        req, err := http.NewRequest("POST", url, bytes.NewBuffer(createvm))
        req.Header.Set("vmware-api-session-id", authtoken)
        req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	resp, err := client.Do(req)
        defer resp.Body.Close()
        if err != nil{
              fmt.Println(err)
              }
        fmt.Println(resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return  (string(body))
	}

func main() {
	listVMsCmd := flag.NewFlagSet("listvms", flag.ExitOnError)
	listVMsVMName := listVMsCmd.String("vmname", "all", "Please specify VM Name")
	listVMsUsername := listVMsCmd.String("u", "", "Please specify username")
	listVMsPassword := listVMsCmd.String("p", "", "Please specify password")
	listVMsURL := listVMsCmd.String("url", "", "Please specify IP/FQDN")

	createVMCmd := flag.NewFlagSet("createvm", flag.ExitOnError)
	createVMUsername := createVMCmd.String("u", "", "Please specify username")
	createVMPassword := createVMCmd.String("p", "", "Please specify password")
	createVMURL := createVMCmd.String("url", "", "Please specify IP/FQDN")

	createTokenCmd := flag.NewFlagSet("createtoken", flag.ExitOnError)
        createTokenUsername := createTokenCmd.String("u", "", "Please specify username")
        createTokenPassword := createTokenCmd.String("p", "", "Please specify password")
        createTokenURL := createTokenCmd.String("url", "", "Please specify IP/FQDN")


	if len(os.Args) < 2 {
		fmt.Println("expected subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	//That should list all VMs, or a specific VM
	case "listvms":
		listVMsCmd.Parse(os.Args[2:])
		authtoken := basicAuth(*listVMsUsername, *listVMsPassword, *listVMsURL)
		getvms(authtoken, *listVMsURL, *listVMsVMName)
		//That should create a virtual machine. Note:  still missing alot of configuration.
	case "createvm":
		createVMCmd.Parse(os.Args[2:])
		authtoken := basicAuth(*createVMUsername, *createVMPassword, *createVMURL)
		createvm(authtoken, *createVMURL)
	//That will just print a new token, in case you want to use it with other tools (e.g. curl)
	case "createtoken":
		createTokenCmd.Parse(os.Args[2:])
		authtoken := basicAuth(*createTokenUsername, *createTokenPassword, *createTokenURL)
		fmt.Println(authtoken)

	default:
		fmt.Println("expected subcommands")
		os.Exit(1)
	}
}
