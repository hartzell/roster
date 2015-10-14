package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hashicorp/terraform/builtin/providers/openstack"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func main() {

	exitStatus, err := doIt()
	if exitStatus != 0 {
		log.Println(err)
	}

	os.Exit(exitStatus)
}

func doIt() (exitStatus int, errorMessage string) {
	var list bool
	var host string

	flag.BoolVar(&list, "list", false, "Run the list command")
	flag.StringVar(&host, "host", "", "specify a host")
	flag.Parse()

	if list && host != "" {
		return 1,
			"Provide either \"--list\" or \"--host\", not both (\"--help\" for help)."
	}

	if host != "" {
		return doHost(host)
	}

	file := "terraform.tfstate"
	if flag.Arg(0) != "" {
		file = flag.Arg(0)
	}
	f, err := os.Open(file)
	if err != nil {
		return 1, "Unable to open state file "
	}

	return doList(f)
}

func doHost(host string) (exitStatus int, errorMessage string) {
	fmt.Println("{}")
	return 0, ""
}

func doList(src io.Reader) (exitStatus int, errorMessage string) {
	state, err := terraform.ReadState(src)
	if err != nil {
		return 1, "Unable to read state file"
	}

	i := inventory{}
	for _, m := range state.Modules {
		for _, rs := range m.Resources {
			switch rs.Type {
			case "openstack_compute_instance_v2":
				// Thanks to @apparentlymart for this bit of code that pulls
				// values out of the state file info.  It's a bit underhanded
				// and behind terraform's public interface, but until there's
				// a better way....  See:
				// https://github.com/hashicorp/terraform/issues/3405
				provider := openstack.Provider().(*schema.Provider)
				instanceSchema := provider.ResourcesMap["openstack_compute_instance_v2"].Schema
				stateReader := &schema.MapFieldReader{
					Schema: instanceSchema,
					Map:    schema.BasicMapReader(rs.Primary.Attributes),
				}
				metadataResult, err := stateReader.ReadField([]string{"metadata"})
				if err != nil {
					return 1, "dammit"
				}

				for a, b := range metadataResult.ValueOrZero(instanceSchema["metadata"]).(map[string]interface{}) {
					fmt.Println(a, "=", b)
				}
				nameResult, err := stateReader.ReadField([]string{"name"})
				if err != nil {
					return 1, "dammit #2"
				}
				fmt.Println("name = ", nameResult.ValueOrZero(instanceSchema["name"]).(string))
				accessResult, err := stateReader.ReadField([]string{"access_ip_v4"})
				if err != nil {
					return 1, "dammit #3"
				}
				fmt.Println("access_ip_v4 = ",
					accessResult.ValueOrZero(instanceSchema["access_ip_v4"]).(string))
				instanceName := rs.Primary.Attributes["name"]
				addr := rs.Primary.Attributes["access_ip_v4"]
				i.AddHostToGroup(addr, instanceName)
			}
		}
	}
	b, err := json.Marshal(i)
	if err != nil {
		return 1, "unable to json.Marshal inventory"
	}
	_, err = os.Stdout.Write(b)
	if err != nil {
		return 1, "unable to write json to stdout"
	}

	return 0, ""
}
