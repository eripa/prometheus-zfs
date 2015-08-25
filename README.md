# prometheus-zfs

----
**Note:** This project was mainly built in order to learn more Go and Prometheus. Given the nature of ZFS and the zpool command, this service requires root to be work. This generally a bad idea for a network exposed service.

I've come to realize that Prometheus node_exporter does support a text-file parsing feature. My intention is to rewrite this program to instead of being a network service for Prometheus it will be simplified to produce text output in the [Prometheus Exposition Format](http://prometheus.io/docs/instrumenting/exposition_formats/) that can then be used by node_exporter.
----

Prometheus metric endpoint to get ZFS pool stats, written in Go.

Using Go gives the nice benefit of static binaries on different platforms. The only external dependency is 'zpool', which you probably have where you want to use this.

Heavily borrowed from my [nagios-zfs-go](https://github.com/eripa/nagios-zfs-go) utility which is used to do Nagios status checks of zpools.

## Usage

prometheus-zfs runs in the foreground, providing a HTTP endpoint for Prometheus collection.

Listen port and endpoint name can be configured using command lines, as shown in the help text.

    Usage of ./prometheus-zfs:
      -endpoint string
            HTTP endpoint to export data on (default "metrics")
      -p string
            what ZFS pool to monitor (shorthand) (default "tank")
      -pool string
            what ZFS pool to monitor (default "tank")
      -port string
            Port to listen on (default "8080")
      -version
            display current tool version

## Example run

Launch exporter:

    $ ./prometheus-zfs -p zones -port 8090 -endpoint zonesmetrics
    Starting zpool metrics exporter on :8090/zonesmetrics

And collect using curl:

    $ curl http://localhost:8090/zonesmetrics 2> /dev/null | grep "^zpool"
    zpool_capacity_percentage 53
    zpool_faulted_providers_count 0
    zpool_online_providers_count 6

## Build

I recommend to use Go 1.5, to make cross-compilation a lot easier.

SmartOS (x86_64):

    env GOOS=solaris GOARCH=amd64 go build -o bin/prometheus-zfs-solaris

Linux (x86_64):

    env GOOS=linux GOARCH=amd64 go build -o bin/prometheus-zfs-linux

Mac OS X:

    env GOOS=darwin GOARCH=amd64 go build -o bin/prometheus-zfs-mac

## Tests

There are some simple test cases to make sure that no insane results occur. All test cases are based on a raidz2 setup with 6 disks. So perhaps more variants of pool configurations would be good to add.. also one could create different, real, pool using disk images. Contributions are welcome!

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
