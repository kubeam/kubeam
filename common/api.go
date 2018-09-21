package common

import (
	"github.com/go-chi/render"
	"net/http"
)

// ResponseData contains the actual data for the perticular API.
// e.g. for iksm operation it has all the namespace info
type ResponseData struct {
	//Deploymentgroup *ArgoCdAppRepoResponse `json:"deploymentgroup,omitempty"`
}

// AsyncStatus represents the onboarding status with message
type AsyncStatus struct {
	Code int    `json:"code"`
	Text string `json:"text,omitempty"`
}

// SyncStatus represents the onboarding status with message
type SyncStatus struct {
	Code int    `json:"code"`
	Text string `json:"text,omitempty"`
}

// AsyncResponsePayload is the payload returned when initiated the asynchronous API call
type AsyncResponsePayload struct {
	ID       string `json:"id,omitempty"`
	Location string `json:"location,omitempty"`
}

// StatusResponsePayload is for asynchronous API status
type StatusResponsePayload struct {
	Message string `json:"message,omitempty"`

	Status AsyncStatus `json:"status,omitempty"`
	AsyncResponsePayload
}

// Render implements the chi Renderer interface, it makes sure to set the correct Content-Type
func (res *StatusResponsePayload) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	render.SetContentType(render.ContentTypeJSON)

	return nil
}

// Error - Error response on rest calls
type Error struct {
	Message string
	Code    int
}

// EmptyResponse - Is a empty struct used when Search/List returns no rows.
type EmptyResponse struct {
}

// Render implements the chi Renderer interface, for EmptyResponse
func (res *EmptyResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	render.SetContentType(render.ContentTypeJSON)

	return nil
}

// ErrResponse is the error schema based on intuit standard
type ErrResponse struct {
	Code     string `json:"code"`
	Message  string `json:"message,omitempty"`
	Detail   string `json:"detail,omitempty"`
	MoreInfo string `json:"moreInfo,omitempty"`
}

// ErrResponsePayload is the generic API response structure for error return
type ErrResponsePayload struct {
	HTTPStatusCode int `json:"-"` // http response status code

	Error []ErrResponse `json:"error"`
}

// Render implements the chi Renderer interface, it makes sure to set the correct Content-Type
func (res *ErrResponsePayload) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, res.HTTPStatusCode)
	render.SetContentType(render.ContentTypeJSON)

	return nil
}
