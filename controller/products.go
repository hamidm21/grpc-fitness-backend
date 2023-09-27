package controller

import (
	"context"
	"fmt"

	"gitlab.com/mefit/mefit-server/entity"
	// "github.com/jalaali/go-jalaali"
	// "gitlab.com/mefit/mefit-server/utils"

	pb "gitlab.com/mefit/mefit-api/proto"
)

func (s *Controller) GetProducts(ctx context.Context, in *pb.Empty) (*pb.ProductRes, error) {
	products := []entity.Product{}
	if err := entity.SimpleCrud(entity.Product{}).List(&products, nil); err != nil {
		return nil, err
	}
	return listProductFrom(products)
}

func listProductFrom(items []entity.Product) (*pb.ProductRes, error) {
	products := &pb.ProductRes{}
	products.Products = []*pb.Product{}
	for _, prd := range items {
		product := &pb.Product{
			ID:          int32(prd.ID),
			Name:        prd.Name,
			Price:       int32(prd.Price),
			Description: prd.Description,
			Off:         int32(prd.Off),
			Label:       prd.Label,
			Recommended: prd.Recommended,
		}
		products.Products = append(products.Products, product)
	}
	fmt.Print(products)
	return products, nil
}
