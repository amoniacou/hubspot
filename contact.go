package hubspot

import (
	"encoding/json"
)

type Contact struct {
	Vid        int `json:"vid"`
	AddedAt    int64 `json:"addedAt"`
	Email      string `json:"email"`
	PortalID   int `json:"portal-id"`
	IsContact  bool `json:"is-contact"`
	ProfileURL string `json:"profile-url"`
	Properties map[string]Property `json:"properties"`
}

type Contacts struct {
	Client         *Client `json:"-"`
	Contacts       []Contact `json:"contacts"`
	Offset         int `json:"-"`
	SubmissionMode string `json:"-"`
	Count          int `json:"-"`
	Properties     []string `json:"-"`
	HasMore        bool `json:"has-more"`
	VidOffset      int `json:"vid-offset"`
}

func (c *Contacts) Next() bool {
	return c.HasMore
}

func (c *Contacts) GetNext() error {
	c.Offset = c.VidOffset
	return c.getAll()
}

func (c *Contacts) getRequestParams() map[string]interface{} {
	return map[string]interface{}{
		"vidOffset":          c.Offset,
		"count":              c.Count,
		"formSubmissionMode": c.SubmissionMode,
		"property":           c.Properties,
	}
}

func (c *Contacts) getAll() error {
	var data []byte
	body, err := c.Client.doRequest("contacts/v1/lists/all/contacts/all", "GET", data, c.getRequestParams())
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, c)
	return err
}

func (c *Client) AddContact(contact *Contact) error {
	return nil
}

func (c *Client) GetContacts(arguments ...interface{}) (*Contacts, error) {
	var offset int = 0
	var count int = 100
	var submissionMode string = "all"
	var properties []string = []string{
		"firstname", "lastname", "company", "email",
	}
	for i, p := range arguments {
		switch i {
		case 0:
			param, ok := p.(int)
			if !ok {
				panic("1st parameter not type int.")
			}
			offset = param

		case 1:
			param, ok := p.(int)
			if !ok {
				panic("2nd parameter not type int.")
			}
			count = param
		case 2:
			param, ok := p.(string)
			if !ok {
				panic("3rd parameter not type string.")
			}
			submissionMode = param

		case 3:
			param, ok := p.([]string)
			if !ok {
				panic("4th parameter not type []string.")
			}
			properties = param
		}
	}

	contacts := &Contacts{
		Offset:         offset,
		Count:          count,
		Properties:     properties,
		SubmissionMode: submissionMode,
		Client:         c,
	}
	err := contacts.getAll()

	return contacts, err
}
