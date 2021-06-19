package article

import (
	"context"
	"log"
	"net/http"

	"github.com/SergeyParamoshkin/rest/internal/errresponse"
	"github.com/SergeyParamoshkin/rest/internal/model"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// ArticleCtx middleware is used to load an Article object from
// the URL parameters passed through as the request. In case
// the Article could not be found, we stop here and return a 404.
func ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var article *model.Article
		var err error

		if articleID := chi.URLParam(r, "articleID"); articleID != "" {
			article, err = dbGetArticle(articleID)
		} else if articleSlug := chi.URLParam(r, "articleSlug"); articleSlug != "" {
			article, err = dbGetArticleBySlug(articleSlug)
		} else {
			err = render.Render(w, r, errresponse.ErrNotFound)
			if err != nil {
				log.Println(err)
			}

			return
		}
		if err != nil {
			err = render.Render(w, r, errresponse.ErrNotFound)
			if err != nil {
				log.Println(err)
			}

			return
		}

		// nolint
		ctx := context.WithValue(r.Context(), "article", article)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
