package controllers

/*
The file based controller uses the in-memory controller and syncs all changes
done by it to disk, to a collection of YAML or JSON files. The folder structure
is so that the resulting data can be easily versioned in a version control
system like Git.

actions/
configs/
streams/
projects/
  ac4acdcef.yml
*/
