# Cake
Running shell scripts on remote machine just by invoking urls.

### Installing
```bash
	go get -u github.com/Hu13er/cakeserver
```
### Usage
write `Cakefile.yaml` with these rules:

```yaml
# Cakefile.yaml

# Your listening port:
addr:   ':2128'

# Your secret that needs for authentication:
secret: 'mysecret'

# List your commands in here:
commands:
  - name: 'foobar'      # Your command name
    trusted: ['.*']     # trusted remote machine with regular expression

    # Your script: 
    script: |
     	FILENAME="foobar.txt"
     	touch $FILENAME
     	echo "someone invoked foobar" >> $FILENAME
```



