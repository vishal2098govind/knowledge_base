#go

I nice way to set up configuration of our service could be using struct-tags to mention default values, which could be easier to **understand**

a configuration system could be setup in many ways, but it should at least have these things in it:
- default values, that should work a 100% in dev environment
	- anyone could clone and just be able to run the project asap
- be able to override defaults
	- can be done by using CLI variables or flags
	- can be done using conf files like `yaml` or `toml` or `json` files

```go
	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout        time.Duration `conf:"default:5s"`
			WriteTimeout       time.Duration `conf:"default:10s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:20s"`
			APIHost            string        `conf:"default:0.0.0.0:3000"`
			DebugHost          string        `conf:"default:0.0.0.0:3010"`
			CORSAllowedOrigins []string      `conf:"default:*"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "Sales",
		},
	}
```
- using a struct literal for config is easier to understand over easy to do
- using a named struct type for config and pass it around the program, would rather be easy to do but not easier to understand
- having a struct literal for config, at one place in main, if we want to know the config options available, we can simply come to main.go and find all at one place in the struct literal type definition. 
- no one else (no other package) should be using this conf struct literal type
- the conf can be constructed in main, and be passed around the program
- currently the config is set to it's zero value. To assign default values, we can have a parser on the config, which returns an error 
	- if overrides, if passed any, are not correctly passed
	- or `--help` is passed while running the build
```go
	const prefix = "SALES"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if err == conf.ErrHelpWanted {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}
```
- on passing `--help`, we can display something like:
```sh
go run main.go --help
time=2026-01-20T17:48:28.281+05:30 level=INFO file=main.go:46 msg=startup service=SALES build=develop build_date=YYYY-MM-DDTHH:MM:SSZ GOMAXPROCS=12 trace_id=""
Usage: main [options...] [arguments...]

OPTIONS
  -h, --help                                                                     display this help message
  -v, --version                                                                  display version
      --web-api-host              <string>              (default: 0.0.0.0:3000)  
      --web-cors-allowed-origins  <string>,[string...]  (default: *)             
      --web-debug-host            <string>              (default: 0.0.0.0:3010)  
      --web-idle-timeout          <duration>            (default: 120s)          
      --web-read-timeout          <duration>            (default: 5s)            
      --web-shutdown-timeout      <duration>            (default: 20s)           
      --web-write-timeout         <duration>            (default: 10s)           

ENVIRONMENT
  SALES_WEB_API_HOST              <string>              (default: 0.0.0.0:3000)  
  SALES_WEB_CORS_ALLOWED_ORIGINS  <string>,[string...]  (default: *)             
  SALES_WEB_DEBUG_HOST            <string>              (default: 0.0.0.0:3010)  
  SALES_WEB_IDLE_TIMEOUT          <duration>            (default: 120s)          
  SALES_WEB_READ_TIMEOUT          <duration>            (default: 5s)            
  SALES_WEB_SHUTDOWN_TIMEOUT      <duration>            (default: 20s)           
  SALES_WEB_WRITE_TIMEOUT         <duration>            (default: 10s)      
```
- we can have a string method on the config to log the configuration options being used to run the program before starting the service
```go
	s, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Info(ctx, "startup", "config", s)
```
- 