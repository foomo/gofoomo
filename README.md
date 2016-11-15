[![Travis CI](https://travis-ci.org/foomo/gofoomo.svg?branch=master)](https://travis-ci.org/foomo/gofoomo)

# gofoomo

Gofoomo lets you use Go in your foomo project. It also lets you use php in your Go project.

We want to use Go, but it is not the right language for everyone, who is using php.

## Complementing your LAMP stack

Go is a much younger and cleaner stack than LAMP.

* Serve static files without bugging your prefork apache
* Keep slow connections away from your php processes (not implemented yet)
* Hijack foomo json rpc services methods
  * Your code is also running the server, this puts you in a place, whereyou can solve problems, that you can not solve in php
* Go´s runtime model is pretty much the opposite of the php runtime model
  * all requests vs one request per lifetime
  * shared memory vs process and memory isolation
  * one bug to kill them all vs one bug kills one request
  * hard, but fast vs easy but slow

## Sitting in front of your foomo LAMP app with Go

Go or php? It is up to you, to decide which tool provides better solutions for your problem and who on your team will be more productive with php or Go.

## Hijacking json rpc calls

Gofoomo lets you intercept and implement calls to foomo json rpc services. In addition [Foomo.Go](https://github.com/foomo/Foomo.Go) gives you an infrastructure to generate golang structs for php value objects.

## Access foomo configurations

Gofoomo gives you access to foomo configurations from Go. Hint: if your php configuration objects are well annotated they are essentially value objects and corresponding structs can easily be generated with Foomo.Go.

## foomo-bert

Is a command line utility, that helps you with the setup of foomo installations.

```bash
go install github.com/foomo/gofoomo/foomo-bert
foomo-bert -help
usage: foomo-bert <command>
foomo-bert prepare :
  -dir string
    	path/to/your/foomo/root
  -run-mode string
    	foomo run mode test | development | production
foomo-bert reset :
  -addr string
    	address of the foomo server
  -dir string
    	path/to/your/foomo/root
  -main-module string
    	name of main module (default "Foomo")
  -run-mode string
    	foomo run mode test | development | production

```

## More to come, but not much more

We are going to add features, as we are going to need them. The focus is to have a simple interface between foomo and Go.

## A little more

This example shows how to access a remote server to read configs.

```Go
rc, err := core.NewRemoteClient("https://user:password@host.com")
if err != nil {
	log.Fatal("could not start remote client", err)
}
c := &MyConf{}
configErr := rc.GetConfig(c, "My.Module", "My.conf", "")
fmt.Println(configErr, c)
```
