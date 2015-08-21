# prometheus-zfs

Prometheus metric endpoint to get ZFS pool stats, written in Go.

Using Go gives the nice benefit of static binaries on different platforms. The only external dependency is 'zpool', which you probably have where you want to use this.

Heavily borrowed from my [nagios-zfs-go](https://github.com/eripa/nagios-zfs-go) utility which is used to do Nagios status checks of zpools.

## Usage

TBW

## Example run

TBW

## Build

I recommend to use Go 1.5, to make cross-compilation a lot easier.

SmartOS (x86_64):

    env GOOS=solaris GOARCH=amd64 go build -o bin/check_zfs-solaris

Linux (x86_64):

    env GOOS=linux GOARCH=amd64 go build -o bin/check_zfs-linux

Mac OS X:

    env GOOS=darwin GOARCH=amd64 go build -o bin/check_zfs-mac

## Tests

There are some simple test cases to make sure that no insane results occur. All test cases are based on a raidz2 setup with 6 disks. So perhaps more variants of pool configurations would be good to add.. Contributions are welcome!

Run `go test -v` to run the tests with some verbosity.

## bin/zpool

`bin/zpool` is a shell-script that can be used to fake a 'zpool' command on your local development machine where you might not have ZFS installed. It will simply run zpool over SSH on a remote host. Set environment variable ZFSHOST to whatever host you want to remote to.

The script also has some simple sed statements prepared (you will have to remove the hash signs manually) to fake different pool statuses for testing purposes.

## License

The MIT License, see separate LICENSE file for full text.

## Contributing

  * Fork it
  * Create your feature branch (git checkout -b my-new-feature)
  * Commit your changes (git commit -am 'Add some feature')
  * Push to the branch (git push origin my-new-feature)
  * Create new Pull Request
