// Package raw provides a client for interacting with the GitHub raw file API
package raw

import (
	"context"
	"net/http"
	"net/url"

	gogithub "github.com/google/go-github/v72/github"
)

// GetRawClientFn is a function type that returns a RawClient instance.
type GetRawClientFn func(context.Context) (*Client, error)

// Client is a client for interacting with the GitHub raw content API.
type Client struct {
	url    *url.URL
	client *gogithub.Client
}

// NewClient creates a new instance of the raw API Client with the provided GitHub client and provided URL.
func NewClient(client *gogithub.Client, rawURL *url.URL) *Client {
	client = gogithub.NewClient(client.Client())
	client.BaseURL = rawURL
	return &Client{client: client, url: rawURL}
}

func (c *Client) newRequest(method string, urlStr string, body interface{}, opts ...gogithub.RequestOption) (*http.Request, error) {
	req, err := c.client.NewRequest(method, urlStr, body, opts...)
	return req, err
}

func (c *Client) refURL(owner, repo, ref, path string) string {
	if ref == "" {
		return c.url.JoinPath(owner, repo, "HEAD", path).String()
	}
	return c.url.JoinPath(owner, repo, ref, path).String()
}

func (c *Client) URLFromOpts(opts *RawContentOpts, owner, repo, path string) string {
	if opts == nil {
		opts = &RawContentOpts{}
	}
	if opts.SHA != "" {
		return c.commitURL(owner, repo, opts.SHA, path)
	}
	return c.refURL(owner, repo, opts.Ref, path)
}

// BlobURL returns the URL for a blob in the raw content API.
func (c *Client) commitURL(owner, repo, sha, path string) string {
	return c.url.JoinPath(owner, repo, sha, path).String()
}

type RawContentOpts struct {
	Ref string
	SHA string
}

// GetRawContent fetches the raw content of a file from a GitHub repository.
func (c *Client) GetRawContent(ctx context.Context, owner, repo, path string, opts *RawContentOpts) (*http.Response, error) {
	url := c.URLFromOpts(opts, owner, repo, path)
	req, err := c.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.client.Client().Do(req)
}
