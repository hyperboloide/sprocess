package sprocess

import (
	"errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
	"io"
	"io/ioutil"
	"net/http"
)

const GoogleCloudScope = "https://www.googleapis.com/auth/devstorage.read_write"

type GoogleCloud struct {
	Prefix      string
	Suffix      string
	ProjectId   string
	Bucket      string
	Name        string
	JsonKeyPath string
	Client      *http.Client
	context     context.Context
}

func (g *GoogleCloud) GetName() string {
	return g.Name
}

func (g *GoogleCloud) Start() error {
	if g.Bucket == "" {
		return errors.New("bucket name is undefined")
	} else if g.ProjectId == "" {
		return errors.New("project id is undefined")
	}

	if g.Client != nil {
		g.context = cloud.NewContext(g.ProjectId, g.Client)
		return nil
	}

	if g.JsonKeyPath == "" {
		return errors.New("no mean of identification provided")
	}
	data, err := ioutil.ReadFile(g.JsonKeyPath)
	if err != nil {
		return errors.New("cannot read json key")
	}
	conf, err := google.JWTConfigFromJSON(data, GoogleCloudScope)
	if err != nil {
		return err
	}
	g.context = cloud.NewContext(g.ProjectId, conf.Client(oauth2.NoContext))

	return nil
}

func (g *GoogleCloud) NewWriter(id string, d *Data) (io.WriteCloser, error) {
	return storage.NewWriter(g.context, g.Bucket, g.Prefix+id+g.Suffix), nil
}

func (g *GoogleCloud) NewReader(id string, d *Data) (io.ReadCloser, error) {
	return storage.NewReader(g.context, g.Bucket, g.Prefix+id+g.Suffix)
}

func (g *GoogleCloud) Delete(id string, d *Data) error {
	return storage.DeleteObject(g.context, g.Bucket, g.Prefix+id+g.Suffix)
}
