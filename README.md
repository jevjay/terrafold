# Terrafold

Terrafold is a tool, used to configure terraform based repositories and provide unified YAMl defined interface to Terraform defined objects

## Dependencies

Assuming you would want to compile your appkication following Go package needs to be installed by running:

    go get gopkg.in/yaml.v2

More information can be found here:

https://github.com/go-yaml/yaml


## Installing

### Source

In order to compile from source please run

    go build src/main.go src/ec2.go src/bucket.go src/rds.go src/elb.go 

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
