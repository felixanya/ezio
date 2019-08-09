module github.com/valeamoris/ezio

go 1.12

require (
	cloud.google.com/go v0.43.0 // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20190717042225-c3de453c63f4 // indirect
	github.com/denisenkom/go-mssqldb v0.0.0-20190806190131-db2462fef53b // indirect
	github.com/fatih/color v1.7.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/xorm v0.7.6
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/jackc/pgx v3.5.0+incompatible // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/lib/pq v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/mattn/go-sqlite3 v1.11.0 // indirect
	github.com/prometheus/common v0.6.0
	github.com/sirupsen/logrus v1.4.2
	github.com/urfave/cli v1.21.0
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4 // indirect
	golang.org/x/sys v0.0.0-20190804053845-51ab0e2deafa
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)

replace github.com/golang/lint v0.0.0-20190313153728-d0100b6bd8b3 => golang.org/x/lint v0.0.0-20190313153728-d0100b6bd8b3

replace github.com/testcontainers/testcontainer-go => github.com/testcontainers/testcontainers-go v0.0.0-20190108154635-47c0da630f72

replace github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1

replace gopkg.in/jcmturner/rpc.v1 v1.1.1 => github.com/jcmturner/rpc v0.0.0-20190727145011-a5898cb6c474
