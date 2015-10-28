// Roster is a tool for working transforming Terraform state
// information into other forms (it's original calling was as an
// Ansible dynamic inventory provider).
//
// It uses terraform routines to work with the state information,
// which seems better than reverse engineering the state files
// themselves.  But, these routines are not mainstream, so we'll have
// to see what happens going forward.  More info in
// [terraform issue 3405](https://github.com/hashicorp/terraform/issues/3405).
//
// It uses Go templates to generate it's output.  One can dump out
// these templates, e.g.:
//
//    $ ./roster dump-template /templates/dynamicInventoryTemplate
//
// and see what it's doing.  One can also supply one's own templates
//
//    $ ./roster execute-template myTemplate.txt
//
// The input to the template is a slice of *instanceInfo.  There are
// two helper functions, groups and hostvars that make templates less
// nasty by transforming the instanceInfo into slices of Group and
// HostVar, respectively.  Slices are used so that it is easier to
// avoid trailing commas when generating JSON, using
// [this technique]()http://stackoverflow.com/questions/21305865/golang-separating-items-with-comma-in-template).
//
// Roster parses state data for instances from:
//
//    OpenStack
//    - the instance name is used as the Name
//    - the access_ip_v4 is used as the Address
//    - host_vars and groups can be specified in the resource configuration
//      metadata section like so:
//
//        metadata {
//            ansible_groups = "foo,  bar"
//       			ansible_host_vars = "color=red, importance=high"
//        }
//
//    parsing those strings is very unsophisticated, I'm waiting for better
//    use cases before I get fancy.
//
//    DigitalOcean
//    - the name is used as the Name
//    - the ipv4address is used as the Address
//    - there is no support for groups or hostvars
//
//    Cloudstack
//    - the name is used as the Name
//    - the ipaddress is used as the address
//    - there is no support for groups or hostvars
//
//    AWS
//    - I have preliminary support for AWS instances, BUT
//    - Every instance is named "moose", I can't figure out what info to use
//      as the Name.
//    - If an instance has a public_ip, that is used as the Address, if not
//      and it has a private_ip then that is used.
//
//  Usage:
//
//    These two invocations support the standard Ansible dynamic
//    inventory calling conventions, returning a dynamic inventory and
//    an (empty) host specific inventory as required by the
//    specification.
//
//    $ ./roster --list
//    $ ./roster --host MyHostName
//
//    Beyond that:
//
//    $ ./roster --help
//    usage: roster [--version] [--help] <command> [<args>]
//
//    Available commands are:
//        dump-template       Dump one of roster's built in templates.
//        execute-template    Execute a user supplied template.
//        hosts               Generate an /etc/hosts fragment for the Terraform instances
//        inventory           Generate an Ansible dynamic inventory
//
//    $ ./roster inventory --help
//    Usage of inventory:
//      -dir string
//        	The path to the terraform directory (default ".")
//      -host string
//        	Generate a host-specific inventory for this host.
//      -list
//        	Generate a full inventory (the default behavior).
//    $ ./roster hosts --help
//    Usage of inventory:
//      -dir string
//        	The path to the terraform directory (default ".")
//
//    $ ./roster dump-template --help
//    Usage of inventory:
//      -dir string
//        	The path to the terraform directory (default ".")
//      -template string
//        	The name of the template to dump.
//
//    $ ./roster dump-template --help
//    Usage of inventory:
//      -dir string
//        	The path to the terraform directory (default ".")
//      -template string
//        	The name of the template to dump: /templates/{dynamicInventoryTemplate,etcHostsTemplate}.
//
//    $ ./roster execute-template --help
//    Usage of inventory:
//      -dir string
//        	The path to the terraform directory (default ".")
//      -template string
//        	The filename of the template to dump.
//
package main
