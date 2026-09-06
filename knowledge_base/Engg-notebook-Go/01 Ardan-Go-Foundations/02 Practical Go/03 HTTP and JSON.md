#engineers-notebook #golang #http #json

- by default numbers in JSON converted to `float64`
- `encoding/json` api in Go
	- JSON to []byte to Go struct -> use `json.Unmarshall` -> `Serialize`
	- Go struct to []byte to JSON -> use `json.Marshall` -> `Deserialize`
	- JSON to `io.Reader` (e.g. http.Response.Body) to Go struct -> use `json.Decoder`
	- Go struct to `io.Writer` to JSON -> use `json.Encoder` -> use `json.Encoder`