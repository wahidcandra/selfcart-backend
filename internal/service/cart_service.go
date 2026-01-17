package service

import (
	"context"
	"errors"
	"fmt"
	"selfcart/internal/repository"

	"github.com/google/uuid"
)

type CartService struct {
	cartRepo           *repository.CartRepo
	productRepo        *repository.ProductRepo
	transactionService *TransactionService
}

func NewCartService(c *repository.CartRepo, p *repository.ProductRepo, t *TransactionService) *CartService {
	return &CartService{
		cartRepo:           c,
		productRepo:        p,
		transactionService: t,
	}
}
func (s *CartService) CreateCart(ctx context.Context, customerID int64) (*repository.Cart, error) {
	cartCode := "CART-" + uuid.New().String()
	return s.cartRepo.Create(ctx, customerID, cartCode)
}

func (s *CartService) AddItem(ctx context.Context, cartID int64, barcode string, act string) error {
	product, err := s.productRepo.GetByBarcode(ctx, barcode)
	fmt.Println(product)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}
	//Get Cart By Item Already Exist
	item, err := s.cartRepo.GetItem(ctx, cartID, product.ID)
	fmt.Println(item)
	if err != nil {
		return err
	}
	if item.ProductID == 0 && act == "add" {
		fmt.Println("Item New")
		return s.cartRepo.AddItem(ctx, cartID, repository.CartItem{
			ProductID: product.ID,
			Qty:       1,
			Price:     int(product.Price),
			Discount:  0,
			SubTotal:  int(product.Price),
		})
	} else {
		fmt.Println("Item Already Exist")
		if act == "add" {
			return s.cartRepo.UpdateItem(ctx, item.ID, repository.CartItem{
				ProductID: product.ID,
				Qty:       item.Qty + 1,
				Price:     int(product.Price),
				Discount:  0,
				SubTotal:  (item.Qty + 1) * int(product.Price),
			})
		} else if act == "reduce" {

			if item.Qty == 1 {
				return s.cartRepo.DeleteItem(ctx, item.ID)
			} else {
				return s.cartRepo.UpdateItem(ctx, item.ID, repository.CartItem{
					ProductID: product.ID,
					Qty:       item.Qty - 1,
					Price:     int(product.Price),
					Discount:  0,
					SubTotal:  (item.Qty - 1) * int(product.Price),
				})
			}
		}
		return nil
	}
}
func (s *CartService) RemoveItem(ctx context.Context, cartID int64, itemID int64) error {
	return s.cartRepo.DeleteItem(ctx, itemID)
}
func (s *CartService) GetCart(ctx context.Context, cartID int64) (*repository.Cart, error) {
	return s.cartRepo.GetCart(ctx, cartID)
}
func (s *CartService) GetCartItems(ctx context.Context, cartID int64) ([]repository.CartItem, error) {
	return s.cartRepo.GetItems(ctx, cartID)
}
func (s *CartService) CheckOut(ctx context.Context, cartID int64) (*repository.Transaction, error) {
	cart, err := s.cartRepo.GetCart(ctx, cartID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, errors.New("cart not found")
	}
	if cart.Status != "ACTIVE" {
		return nil, errors.New("cart is not active")
	}

	err = s.cartRepo.UpdateStatus(ctx, cartID, "CHECKOUT")

	if err != nil {
		return nil, err
	}
	cartItems, err := s.cartRepo.GetItems(ctx, cartID)
	if err != nil {
		return nil, err
	}
	for _, item := range cartItems {
		item.SubTotal = item.Qty * item.Price
		item.Discount = 0
		cart.Total = cart.Total + item.SubTotal
	}

	transaction, err := s.transactionService.Create(ctx, cartID, cart.Total)
	if err != nil {
		return nil, err
	}
	return transaction, nil

	// cart.Status = "CHECKOUT"
	// cart.Total = cart.Total - cart.Discount
	// cartItems, err := s.cartRepo.GetItems(ctx, cartID)
	// if err != nil {
	// 	return err
	// }
	// for _, item := range cartItems {
	// 	item.SubTotal = item.Qty * item.Price
	// 	item.Discount = 0
	// 	cart.Total = cart.Total + item.SubTotal
	// }
	// cart.Total = cart.Total - cart.Discount
}
