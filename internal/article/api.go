package article

import (
	"log"
	"net/http"

	"github.com/SergeyParamoshkin/rest/internal/articlerequest"
	"github.com/SergeyParamoshkin/rest/internal/articleresponse"
	"github.com/SergeyParamoshkin/rest/internal/errresponse"
	"github.com/SergeyParamoshkin/rest/internal/model"
	"github.com/SergeyParamoshkin/rest/internal/user"
	"github.com/SergeyParamoshkin/rest/internal/userpayload"
	"github.com/go-chi/render"
)

func ListArticles(w http.ResponseWriter, r *http.Request) {
	if err := render.RenderList(w, r, articleresponse.NewArticleListResponse(articles)); err != nil {
		err = render.Render(w, r, errresponse.ErrRender(err))
		if err != nil {
			// sugarLogger.Errorw(err.Error())
		}

		return
	}
}

// CreateArticle persists the posted Article and returns it
// back to the client as an acknowledgement.
func CreateArticle(w http.ResponseWriter, r *http.Request) {
	data := &articlerequest.ArticleRequest{}
	if err := render.Bind(r, data); err != nil {
		err = render.Render(w, r, errresponse.ErrInvalidRequest(err))
		if err != nil {
			log.Println(err)
		}

		return
	}

	article := data.Article
	_, err := dbNewArticle(article)
	if err != nil {
		log.Println(err)
	}

	render.Status(r, http.StatusCreated)
	err = render.Render(w, r, articleresponse.NewArticleResponse(article))
	if err != nil {
		log.Println(err)
	}
}

// GetArticle returns the specific Article. You'll notice it just
// fetches the Article right off the context, as its understood that
// if we made it this far, the Article must be on the context. In case
// its not due to a bug, then it will panic, and our Recoverer will save us.
func GetArticle(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the article
	// context because this handler is a child of the ArticleCtx
	// middleware. The worst case, the recoverer middleware will save us.
	// nolint
	article := r.Context().Value("article").(*model.Article)

	if err := render.Render(w, r, articleresponse.NewArticleResponse(article)); err != nil {
		err = render.Render(w, r, errresponse.ErrRender(err))
		if err != nil {
			log.Println(err)
		}

		return
	}
}

// UpdateArticle updates an existing Article in our persistent store.
func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	// nolint
	article := r.Context().Value("article").(*model.Article)

	data := &articlerequest.ArticleRequest{
		Article: article,
		User: &userpayload.UserPayload{
			User: &user.User{
				ID:   0,
				Name: "",
			},
			Role: "",
		},
		ProtectedID: "",
	}
	if err := render.Bind(r, data); err != nil {
		err = render.Render(w, r, errresponse.ErrInvalidRequest(err))
		if err != nil {
			log.Println(err)
		}

		return
	}

	article = data.Article
	_, err := dbUpdateArticle(article.ID, article)
	if err != nil {
		log.Println(err)
	}

	err = render.Render(w, r, articleresponse.NewArticleResponse(article))
	if err != nil {
		log.Println(err)
	}
}

// DeleteArticle removes an existing Article from our persistent store.
func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	var err error

	// Assume if we've reach this far, we can access the article
	// context because this handler is a child of the ArticleCtx
	// middleware. The worst case, the recoverer middleware will save us.
	// nolint
	article := r.Context().Value("article").(*model.Article)

	article, err = dbRemoveArticle(article.ID)
	if err != nil {
		err = render.Render(w, r, errresponse.ErrInvalidRequest(err))
		if err != nil {
			log.Println(err)
		}

		return
	}

	err = render.Render(w, r, articleresponse.NewArticleResponse(article))
	if err != nil {
		log.Println(err)
	}
}

// SearchArticles searches the Articles data for a matching article.
// It's just a stub, but you get the idea.
func SearchArticles(w http.ResponseWriter, r *http.Request) {
	if err := render.RenderList(w, r, articleresponse.NewArticleListResponse(articles)); err != nil {
		err = render.Render(w, r, errresponse.ErrRender(err))
		if err != nil {
			log.Println(err)
		}

		return
	}
}
