package catalog

import (
	"context"

	"github.com/leminkhoa/go-grpc-graphql-microservice/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.CatalogServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := pb.NewCatalogServiceClient(conn)

	return &Client{
		conn:    conn,
		service: c,
	}, nil
}

func (c *Client) PostProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	r, err := c.service.PostProduct(ctx, &pb.PostProductRequest{
		Name:        name,
		Description: description,
		Price:       price,
	})

	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          r.Product.Id,
		Name:        r.Product.Name,
		Description: r.Product.Description,
		Price:       r.Product.Price,
	}, nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	r, err := c.service.GetProduct(ctx, &pb.GetProductRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          r.Product.Id,
		Name:        r.Product.Name,
		Description: r.Product.Description,
		Price:       r.Product.Price,
	}, nil
}

func (c *Client) GetProducts(ctx context.Context, skip uint64, take uint64, ids []string, query string) ([]Product, error) {
	r, err := c.service.GetProducts(ctx, &pb.GetProductsRequest{
		Ids:   ids,
		Query: query,
		Skip:  skip,
		Take:  take,
	})

	if err != nil {
		return nil, err
	}

	var products []Product

	for _, r := range r.Products {
		products = append(products, Product{
			ID:          r.Id,
			Name:        r.Name,
			Description: r.Description,
			Price:       r.Price,
		})
	}

	return products, nil
}
