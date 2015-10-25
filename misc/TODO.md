- Multi-word commands, mitchellh/cli default command, etc...

  I need to figure out whether one can use "multi-word" commands after
  `-i` in ansible commands.  When I use the `ansible` command with `-i
  .../roster inventory` it fails, but if I create a little shell
  script that just calls `roster -inventory` then everything works.

        ansible-playbook -i 'ec2.py --profile prod' myplaybook.yml

- I might *have* to generate a group for each individual host in order
  for things to work, in which case I can clean up that silly boolean.

- parser tests.

- AWS name?  "moose" isn't going to cut it.

- Play with ways to integrate dynamic inventory w/ static files.  In
  particular, are groups recursive?
