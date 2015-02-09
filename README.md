# sprocess

sprocess is a library designed to do transformations on golang streams and http requests.

## installation

```shell
go get github.com/hyperboloide/sprocess
```

## How it works

**sprocess** is use to transform streams and save them somewhere. It can also do the inverse operation: get a stream and untransform it.

### Encoders and Decoders

Transformations are done with **encoders** and **decoders**. These are objects that transform and untransform a stream, like compress/decompress or encrypt/uncrypt.


### Outputs, Inputs

**outputs** can save a stream (to a file or an s3 bucket) and **inputs** read this stream back.

These may also allow for deletetion.

### Http Handlers

**sprocess** also contains HTTP handlers to save, get and delete files.






