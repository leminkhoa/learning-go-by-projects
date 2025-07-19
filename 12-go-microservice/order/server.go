package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/leminkhoa/go-grpc-graphql-microservice/account"
	"github.com/leminkhoa/go-grpc-graphql-microservice/catalog"
	"github.com/leminkhoa/go-grpc-graphql-microservice/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pb.UnimplementedOrderServiceServer
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
}

func ListenGRPC(s Service, accountURL, catalogURL string, port int) error {

	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return err
	}

	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{
		service:       s,
		accountClient: accountClient,
		catalogClient: catalogClient,
	})
	reflection.Register(serv)
	return serv.Serve(lis)

}

func (s *grpcServer) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, r.AccountId)
	if err != nil {
		log.Println("Error getting account:", err)
		return nil, errors.New("account not found")
	}

	productIDs := []string{}
	for _, p := range r.Products {
		productIDs = append(productIDs, p.ProductId)
	}

	// Retrieve product information from catalog client
	products, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("Error getting products")
		return nil, errors.New("products not found")
	}

	// Construct ordered products
	orderedProducts := []OrderedProduct{}
	for _, p := range products {
		orderedProduct := OrderedProduct{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    0, // Place holder
		}

		for _, rp := range r.Products {
			if rp.ProductId == p.ID {
				orderedProduct.Quantity = rp.Quantity
				break
			}
		}

		if orderedProduct.Quantity != 0 {
			orderedProducts = append(orderedProducts, orderedProduct)
		}
	}

	// Call order service implementation
	order, err := s.service.PostOrder(ctx, r.AccountId, orderedProducts)
	if err != nil {
		log.Println("Error posting order:", err)
		return nil, errors.New("could not post order")
	}

	orderProto := &pb.Order{
		Id:         order.ID,
		AccountId:  order.AccountID,
		TotalPrice: order.TotalPrice,
		Products:   []*pb.Order_OrderProduct{},
	}
	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()
	for _, p := range order.Products {
		orderProto.Products = append(orderProto.Products, &pb.Order_OrderProduct{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
		})
	}

	return &pb.PostOrderResponse{
		Order: orderProto,
	}, nil

}

func (s *grpcServer) GetOrdersForAccount(ctx context.Context, r *pb.GetOrdersForAccountRequest) (*pb.GetOrdersForAccountResponse, error) {
	accountOrders, err := s.service.GetOrdersForAccount(ctx, r.AccountId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Construct a map to find unique product id only
	productIDMap := map[string]bool{}
	for _, o := range accountOrders {
		for _, p := range o.Products {
			productIDMap[p.ID] = true
		}
	}

	// Loop through the map to extract the productIDs
	productIDs := []string{}
	for id := range productIDMap {
		productIDs = append(productIDs, id)
	}

	// Retrieve product information from the product client
	products, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("Error getting account products ", err)
		return nil, err
	}

	order := []*pb.Order{}
	for _, o := range accountOrders {
		op := &pb.Order{
			Id:         o.ID,
			TotalPrice: o.TotalPrice,
			Products:   []*pb.Order_OrderProduct{},
		}
		op.CreatedAt, _ = o.CreatedAt.MarshalBinary()

		// Mapping
		for _, product := range o.Products {
			for _, p := range products {
				if product.ID == p.ID {
					product.Name = p.Name
					product.Description = p.Description
					product.Price = p.Price
				}
			}
			op.Products = append(op.Products, &pb.Order_OrderProduct{
				Id:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				Quantity:    product.Quantity,
			})
		}

		// Append each order protobuf to the list
		order = append(order, op)

	}

	return &pb.GetOrdersForAccountResponse{
		Orders: order,
	}, nil

}
