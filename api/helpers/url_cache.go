// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package helpers

import (
	"github.com/kiprotect/kodex"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Content struct {
	Type      string
	FetchedAt time.Time
	Data      []byte
}

type Cache struct {
	mutex sync.Mutex
	urls  map[string]*Content
}

func (c *Cache) update(url string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if existingContent, ok := c.urls[url]; ok {
		// we only update once per minute
		if time.Now().UTC().Sub(existingContent.FetchedAt) < time.Second*60 {
			return nil
		}
	}

	if content, err := c.fetch(url); err != nil {
		kodex.Log.Error(err)
		return err
	} else {
		c.urls[url] = content
		return nil
	}
}

func (c *Cache) fetch(url string) (*Content, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return &Content{
		Type:      resp.Header.Get("content-type"),
		FetchedAt: time.Now().UTC(),
		Data:      body,
	}, nil
}

func (c *Cache) Get(url string) (*Content, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if content, ok := c.urls[url]; !ok {
		// we need to fetch this url
		if content, err := c.fetch(url); err != nil {
			return nil, err
		} else {
			c.urls[url] = content
			return content, nil
		}
	} else {
		// we asynchronously update the URL and return the content
		go c.update(url)
		return content, nil
	}
}

func MakeURLCache() *Cache {
	return &Cache{
		urls: make(map[string]*Content),
	}
}
