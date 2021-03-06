# sprocess 

this library is deprecated, use [Pipe](https://github.com/hyperboloide/pipe) instead.

[![Build Status](https://travis-ci.org/hyperboloide/sprocess.svg?branch=master)](https://travis-ci.org/hyperboloide/sprocess)

sprocess is a library designed to do transformations on golang streams and http requests.


## installation

```bash
go get github.com/hyperboloide/sprocess
```

## How it works

**sprocess** is use to transform streams and save them somewhere. It can also do the inverse operation: get a stream and untransform it.

If you want to jump straight to the code [see this example](https://gist.github.com/fdelbos/0c1a0b47ae2cab0e971f#file-example-go)

### Encoders and Decoders

Transformations are done with **encoders** and **decoders**. These are structs that transform and untransform a stream, like compress/decompress or encrypt/uncrypt.

Note that all Encoders and Decoders defines a field `Name` that is used for error reporting.

#### AES

AES encrypts and decrypts the stream.

```go
type AES struct {
    Key          []byte // Encryption key
    Base64String string // or a base64 encoded string
    Name         string
}
```
You must provide either a `Key` or a `Base64String`. To encrypt in AES 256, provide a 32 bytes long key. To encrypt in AES 128, provide a 16 bytes long key

#### Bash

Run a command in a Bash shell. The stream is piped directly to the shell process.

```go
type Bash struct {
    Cmd  string // Command to run ,for example "gzip" or "gzip -d" to compress and decompress.
    Name string
}
```

#### CheckSum

Make a checksum of the stream during encoding and return an error if the checksum dosen't match during decoding. The checksum is save in data with the key corresponding to the Encoder name

```go
type CheckSum struct {
    Hash string
    Name string
}
```

`Hash` can be either "md5" or "sha1" (default if not set)

`Name` is mandatory

#### Gzip

Compress a stream.

```go
type Gzip struct {
    Algo string // compression level
    Name string
}
```
`Algo` can have one of 3 values that correspond to the compression algorithm used :

*  `"best"` : `gzip.BestCompression`
*  `"speed"` : `gzip.BestSpeed`
*  `"default"`: `gzip.DefaultCompression` (default if not set)

#### Image
Transforms an image (for resize and thumbnails). Note that **Image is Encoder only**.

```go
type Image struct {
    Operation     ImageOperation
    Height        uint
    Width         uint
    Interpolation string
    Output        string
    Name          string
}
```

`ImageOperation` defines the type of operation you want, it can be either:

* `ImageThumbnail` : to downscale an image preserving its aspect ratio to the maximum dimensions (`Width`, `Height`)
* `ImageResize` :  to create a scaled image with new dimensions. If either `Width` or `Height` is set to 0, it will be set to an aspect ratio preserving value.

`Interpolation` defines the interpolation function to use (from fast to slow execution time):

* `NearestNeighbor`: [Nearest-neighbor interpolation](http://en.wikipedia.org/wiki/Nearest-neighbor_interpolation) (default if not set)
* `Bilinear`: [Bilinear interpolation](http://en.wikipedia.org/wiki/Bilinear_interpolation)
* `Bicubic`: [Bicubic interpolation](http://en.wikipedia.org/wiki/Bicubic_interpolation)
* `MitchellNetravali`: [Mitchell-Netravali interpolation](http://dl.acm.org/citation.cfm?id=378514)
* `Lanczos2`: [Lanczos resampling](http://en.wikipedia.org/wiki/Lanczos_resampling) with a=2
* `Lanczos3`: [Lanczos resampling](http://en.wikipedia.org/wiki/Lanczos_resampling) with a=3

`Output` defines the output format:

* `jpg` (default if not set)
* `png`
* `gif`

#### Size

Count the number of bytes

```go
type Size struct {
    Name string
}
```

The sum of bytes will be recored in data with the key `Name`.

#### Tee

Divides an input stream, this is usefull if you want to save multiple files at once. For example you may want to save a picture and at the same time create a thumbnail. Note that **Tee is Encoder only**.

```go
type Tee struct {
    Encoders []Encoder
    Output   Outputer
    Name     string
}
```

`Encoders` defines the operations to be applied on the new stream. This can be empty.


`Output` defines where to save the output. Note that a unique `Outputer` should be used to avoid conflicts.


### Outputs, Inputs

**outputs** can save a stream (to a file or an s3 bucket) and **inputs** read this stream back.

These may also allow for deletetion.

#### File

Save to the local filesystem.

```go
type File struct {
    Prefix string // add a prefix before id to every file
	Suffix string // add a suffix after id to every file
    Dir  string // directory
    Name string
}
```

Streams are saved in directory. If doesn't exists, directory will be created.

#### S3Bucket

Save to an S3 Bucket

```go
type S3Bucket struct {
    Prefix string // add a prefix before id to every file
	Suffix string // add a suffix after id to every file
    AccessKey string // AWS Acces Key
    SecretKey string // AWS Secret Key
    Domain    string // AWS Domain
    Bucket    string // Bucket name
    Name      string
}
```

`AccessKey` and `SecretKey` will be read directly from the environment if not set. To set your environment do something like this :

```bash
export AWS_ACCESS_KEY_ID="my_access_key_id"
export AWS_SECRET_ACCESS_KEY="my_secret_access_key"
```

`Domain` correspond to S3 endpoint where you created the bucket (see : [http://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region](http://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region))



### Http Handlers

**sprocess** also contains HTTP handlers to save, get and delete files.

An example can be found here : [https://gist.github.com/fdelbos/0c1a0b47ae2cab0e971f#file-example-go](https://gist.github.com/fdelbos/0c1a0b47ae2cab0e971f#file-example-go)
