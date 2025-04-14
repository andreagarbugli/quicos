# NOTES

To create the `tar.gz` file of the job directory we can use the following command (to be executed from inside the job directory):

```sh
$ tar -czvf $(basename "$PWD").tar.gz .
```