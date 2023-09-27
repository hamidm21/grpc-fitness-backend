package controller

import (
	"context"

	"gitlab.com/mefit/mefit-server/services/storage/mongodb"

	pb "gitlab.com/mefit/mefit-api/proto"
	"gitlab.com/mefit/mefit-server/utils"
	"gitlab.com/mefit/mefit-server/utils/log"
)

func (s *Controller) GetArticles(ctx context.Context, in *pb.ListReq) (*pb.ListArticleRes, error) {
	// usrID, _ := ctx.Value(utils.KeyEmail).(uint)
	articles, err := mongodb.GetArticles(int(in.Page))
	if err != nil {
		return nil, utils.ErrNotFound
	}
	return listArticlesFrom(articles)
}

func (s *Controller) ArticleInfo(ctx context.Context, in *pb.FindByMongoIdReq) (*pb.Article, error) {
	article, err := mongodb.GetArticle(in.Id)
	if err != nil {
		return nil, utils.ErrNotFound
	}
	return articleInfoFrom(article)
}

func articleInfoFrom(item *mongodb.Article) (*pb.Article, error) {
	author := "فیتکس"
	if item.Author != "" {
		author = item.Author
	}
	return &pb.Article{
		ThumbnailUrl: "https://mag.fitexapp.ir" + item.Thumb,
		CoverUrl:     "https://mag.fitexapp.ir" + item.Images.Small,
		Title:        item.Title,
		Body:         item.Body,
		Id:           item.ID.Hex(),
		Author:       author,
		ShareUrl:     "https://mag.fitexapp.ir/post/" + item.Slug,
		CreatedAt:    item.Time,
	}, nil
}

func listArticlesFrom(items []mongodb.Article) (*pb.ListArticleRes, error) {
	count, err := mongodb.ArticleCount()
	if err != nil {
		log.Logger().Errorf("While get article counts: %v", err)
		return nil, utils.ErrInternal
	}
	articles := &pb.ListArticleRes{}
	articles.Articles = []*pb.Article{}
	for _, art := range items {
		author := "فیتکس"
		if art.Author != "" {
			author = art.Author
		}
		article := &pb.Article{
			ThumbnailUrl: "https://mag.fitexapp.ir" + art.Thumb,
			CoverUrl:     "https://mag.fitexapp.ir" + art.Images.Small,
			Title:        art.Title,
			Body:         art.Body,
			Id:           art.ID.Hex(),
			Author:       author,
			Sum:          art.Sum,
			CreatedAt:    art.Time,
			ShareUrl:     "https://mag.fitexapp.ir/post/" + art.Slug,
		}
		articles.Articles = append(articles.Articles, article)
	}
	articles.TotalCount = int32(count)
	return articles, nil
}
