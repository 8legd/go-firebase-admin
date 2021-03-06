package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Reference represents a specific location in Database
type Reference struct {
	database *Database
	path     string
}

func (ref *Reference) url(ctx context.Context) (string, error) {
	if ref.database.app.tokenSource == nil {
		return ref.database.app.databaseURL + "/" + ref.path + ".json", nil
	}
	tk, err := ref.database.app.tokenSource.Token()
	if err != nil {
		return "", err
	}
	token := tk.AccessToken
	return ref.database.app.databaseURL + "/" + ref.path + ".json?access_token=" + token, nil
}

// Set writes data to current location
func (ref *Reference) Set(value interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	url, err := ref.url(ctx)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(v))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := ref.database.client.Do(req)
	if err != nil {
		return err
	}
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	return nil
}

// Push pushs data to current location
func (ref *Reference) Push(value interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	url, err := ref.url(ctx)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(v))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := ref.database.client.Do(req)
	if err != nil {
		return err
	}
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	return err
}

// Remove removes data from current location
func (ref *Reference) Remove() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	url, err := ref.url(ctx)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	resp, err := ref.database.client.Do(req)
	if err != nil {
		return err
	}
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	return err
}
