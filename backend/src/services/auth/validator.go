package auth

import (
	"backend/src/repository"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func ValidatePassword(hashedpassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(password))
	if err != nil {
		return err
	}
	
	return nil
}

func ValidateRequester[Trepo any](ctx context.Context, requester, comparer, repoType string, repo Trepo)(bool, error){
	if requester == "" || comparer == "" {
		return false, fmt.Errorf("requester and comparer id can't be empty")
	} 

	registry := map[string]func(ctx context.Context, r Trepo, id string)(string, error) {
		"MERCHANT":func(ctx context.Context, r Trepo, id string)(string, error) {
			storeRepo, ok := any(r).(repository.StoreStoreInterface)
			if !ok {
				return "", fmt.Errorf("provided repository does not implement Store interface")
			}

			return storeRepo.FetchStoreID(ctx, id)
		},
		"PRODUCT":func(ctx context.Context, r Trepo, id string) (string, error) {
			productRepo, ok := any(r).(repository.ProductStoreInterface)

			if !ok {
				return "", fmt.Errorf("provided repository does not implement Product interface")
			}

			return productRepo.FetchStoreID(ctx, id)
		},
	}

	fetcher, exists := registry[repoType]
	if !exists {
		return false, fmt.Errorf("unsupported repository type %s", repoType)
	}

	dbOwnerId, err := fetcher(ctx, repo, comparer)
	if err != nil {
		return false, fmt.Errorf("failed to fetch from db %w", err)
	}

	return requester == dbOwnerId, nil
}