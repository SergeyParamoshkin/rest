package articleresponse

import (
	"net/http"

	"github.com/SergeyParamoshkin/rest/internal/model"
	"github.com/SergeyParamoshkin/rest/internal/user"
	"github.com/SergeyParamoshkin/rest/internal/userpayload"
	"github.com/go-chi/render"
)

// ArticleResponse is the response payload for the Article data model.
// See NOTE above in ArticleRequest as well.
//
// In the ArticleResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
type ArticleResponse struct {
	*model.Article

	User *userpayload.UserPayload `json:"user,omitempty"`

	// We add an additional field to the response here.. such as this
	// elapsed computed property
	Elapsed int64 `json:"elapsed"`
}

func NewArticleListResponse(articles []*model.Article) []render.Renderer {
	list := []render.Renderer{}
	for _, article := range articles {
		list = append(list, NewArticleResponse(article))
	}

	return list
}

func NewArticleResponse(article *model.Article) *ArticleResponse {
	resp := &ArticleResponse{
		Elapsed: 0,
		Article: article,
	}

	if resp.User == nil {
		if user, _ := user.DBGetUser(resp.User.ID); user != nil {
			resp.User = userpayload.NewUserPayloadResponse(user)
		}
	}

	return resp
}

func (rd *ArticleResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	rd.Elapsed = 10

	return nil
}
