# sprocess

sprocess is a library designed to do transformations on golang streams and http requests.

## installation

```shell
go get github.com/hyperboloide/sprocess
```

## How it works

**sprocess** is use to transform streams and save them somewhere. It can also do the inverse operation: get a stream and untransform it.

### Encoders and Decoders

Transformations are done with **encoders** and **decoders**. These are structs that transform and untransform a stream, like compress/decompress or encrypt/uncrypt.

#### AES

AES encrypts and decrypts the stream.

```golang
type AES struct {
    Key          []byte // Encryption key
    Base64String string // or a base64 encoded string
    Name         string
}
```
You must provide either a `Key` or a `Base64String`. To encrypt in AES 256, provide a 32 bytes long key. To encrypt in AES 128, provide a 16 bytes long key

#### Bash

Run a command in a Bash shell. The stream is piped to the command.

```golang
type Bash struct {
    Cmd  string // Command to run.
    Name string
}
```

#### Gzip

Compress a stream.

```golang
type Gzip struct {
    Algo string // compression level
    Name string
}
```
`Algo` can have one of 3 values that correspond to the compression algorithm used :

*  `"best"` : `gzip.BestCompression`
*  `"speed"` : `gzip.BestSpeed`
*  `"default"`: `gzip.DefaultCompression`


### Outputs, Inputs

**outputs** can save a stream (to a file or an s3 bucket) and **inputs** read this stream back.

These may also allow for deletetion.

#### File

Save to the local filesystem.

```golang
type File struct {
    Dir  string // directory
    Name string
}
```

Streams are saved in directory. If doesn't exists, directory will be created.

### Http Handlers

**sprocess** also contains HTTP handlers to save, get and delete files.






