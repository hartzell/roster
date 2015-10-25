# Feedback Wanted!

I'm looking for all kinds of feedback; on the code, real world use
cases, feature requests, etc....

Oh yeah, does it work for you?

## From the Go community.

This is my first piece of go code.  I'm sure there's lots of room for
improvement.  Things I know I want feedback on:

- variable and function naming

- structs vs. pointers to structs, consistency and sensibility

- Could I have made better use of templates?

- documentation, go doc

## From the Terraform community

- Thoughts about parsing state files and grubbing around in resource
  data.

- What should I use in AWS for a name, and for group/hostvar info.

- What should I use in DigitalOcean and Cloudstack providers for
  groups/hostvar info.

## From the Ansible communtiy

- Should `-i 'roster inventory'` work?  There's an `ec2.py` example in
  the docs that suggests it should, but it does not for me.

- is it worth making `roster inventory --host=moose.example.org` do
  something more than spit out `{}`?

- Could I have made better use of Go templates?
