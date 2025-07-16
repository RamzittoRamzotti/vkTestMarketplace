package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"vkTestMarketplace/internal/models"
	"vkTestMarketplace/internal/storage"
	utilsads "vkTestMarketplace/internal/utils/ads"
	utilsauth "vkTestMarketplace/internal/utils/auth"
	"vkTestMarketplace/internal/utils/logger"
)

type AdStorage interface {
	CreateAd(ad *models.Ad) error
	ListAds(filter storage.AdFilter) ([]models.Ad, error)
}

type AdHandlers struct {
	Ads       AdStorage
	Users     UserStorage
	JWTSecret string
}

func (h *AdHandlers) CreateAdHandler(w http.ResponseWriter, r *http.Request) {
	userID := utilsauth.UserIDFromContext(r.Context())
	if userID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		ImageURL    string  `json:"image_url"`
		Price       float64 `json:"price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := utilsads.ValidateAdTitle(req.Title); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := utilsads.ValidateAdDescription(req.Description); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := utilsads.ValidateAdPrice(req.Price); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := utilsads.ValidateAdImageURL(req.ImageURL); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ad := &models.Ad{
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Price:       req.Price,
		AuthorID:    userID,
	}
	if err := h.Ads.CreateAd(ad); err != nil {
		logger.Error("CreateAd failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	user, err := h.Users.GetUserByID(userID)
	if err == nil {
		ad.AuthorUsername = user.Login
	}
	logger.Info("Ad created: id=%d by user=%s", ad.ID, ad.AuthorUsername)
	resp := map[string]interface{}{
		"id":          ad.ID,
		"title":       ad.Title,
		"description": ad.Description,
		"image_url":   ad.ImageURL,
		"price":       ad.Price,
		"author":      ad.AuthorUsername,
		"is_mine":     true,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AdHandlers) ListAdsHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := storage.AdFilter{}
	filter.Page, _ = strconv.Atoi(q.Get("page"))
	if filter.Page < 1 {
		filter.Page = 1
	}
	filter.Limit, _ = strconv.Atoi(q.Get("limit"))
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	filter.SortBy = q.Get("sort_by")
	if filter.SortBy != "price" {
		filter.SortBy = "created_at"
	}
	filter.SortOrder = q.Get("sort_order")
	if filter.SortOrder != "asc" {
		filter.SortOrder = "desc"
	}
	filter.MinPrice, _ = strconv.ParseFloat(q.Get("min_price"), 64)
	filter.MaxPrice, _ = strconv.ParseFloat(q.Get("max_price"), 64)

	ads, err := h.Ads.ListAds(filter)
	if err != nil {
		logger.Error("ListAds failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	userID := utilsauth.UserIDFromContext(r.Context())
	resp := make([]map[string]interface{}, 0, len(ads))
	for _, ad := range ads {
		item := map[string]interface{}{
			"title":       ad.Title,
			"description": ad.Description,
			"image_url":   ad.ImageURL,
			"price":       ad.Price,
			"author":      ad.AuthorUsername,
		}
		if userID != 0 && ad.AuthorID == userID {
			item["is_mine"] = true
		}
		resp = append(resp, item)
	}
	logger.Info("Ads listed: count=%d", len(resp))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
