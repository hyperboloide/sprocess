//
// s3bucket.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Feb  8 2015.
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of this source code package.
//

package sprocess

import (
	"errors"
	"github.com/rlmcpherson/s3gof3r"
	"io"
)

type S3Bucket struct {
	Prefix    string
	Suffix    string
	AccessKey string
	SecretKey string
	Domain    string
	Bucket    string
	Name
	bucket *s3gof3r.Bucket
}

func (s *S3Bucket) Start() error {
	if s.Bucket == "" {
		return errors.New("bucket name is undefined")
	}
	var k s3gof3r.Keys
	var err error

	if s.AccessKey == "" || s.SecretKey == "" {
		k, err = s3gof3r.EnvKeys() // get S3 keys from environment
		if err != nil {
			return err
		}
	} else {
		k = s3gof3r.Keys{
			AccessKey: s.AccessKey,
			SecretKey: s.SecretKey,
		}
	}
	s3 := s3gof3r.New(s.Domain, k)
	s.bucket = s3.Bucket(s.Bucket)
	return err
}

func (s *S3Bucket) NewWriter(id string, d *Data) (io.WriteCloser, error) {
	return s.bucket.PutWriter(s.Prefix+id+s.Suffix, nil, nil)
}

func (s *S3Bucket) NewReader(id string, d *Data) (io.ReadCloser, error) {
	r, _, err := s.bucket.GetReader(s.Prefix+id+s.Suffix, nil)
	return r, err
}

func (s *S3Bucket) Delete(id string, d *Data) error {
	return s.bucket.Delete(s.Prefix + id + s.Suffix)
}
