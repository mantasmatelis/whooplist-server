import (
	"path/filepath"
	"crypto/sha512"
	"encoding/base64"
	"io/ioutil"
	"os"
	"rpc"

type Server struct {
	baseUrl string
	basePath string
}


type AddArgs struct {
	Filename string
	Data     byte[]
}

type AddResponse struct {
	Err error
	Url string
}

func (s *Server) Add(args AddArgs, response  *AddResponse) {
	hasher := sha512.New()
	hasher.Write([]byte(AddArgs.Filename))
	hashFilename := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	err := os.Mkdir filePath.Join(basePath, hashFilename)

	if err != nil {
		log.Print("Failed to create directory: ", e)
		response.Err = err
		return
	}

	err := ioutil.WriteFile(filePath.Join(basePath, hashFilename args.Filename), args.Data, 0777)

	if err != nil {
		log.Print("Failed to write file: ", e)
		response.Err = err
		return
	}

	log.Print("Successfuly wrote file ", filePath.Join(basePath, hashFilename, args.Filename))
	response.Url = filepath.Join(baseUrl, hashFilename, args.Filename)
}

func rpcServe(baseUrl, basePath string) {
	server := new(Server)
	server.baseUrl, server.basePath = baseUrl, basePath
	rpc.Register(server)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":3001")
	if e != nil {
		log.Fatal("RPC listen error: " , e)
	}
	go http.Serve(l, nil)
}
