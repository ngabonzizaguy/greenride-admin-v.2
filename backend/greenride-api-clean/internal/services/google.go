package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"greenride/internal/config"
	"greenride/internal/log"
	"greenride/internal/protocol"

	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

var (
	googleServiceInstance *GoogleService
	googleOnce            = sync.Once{}
)

type GoogleService struct {
	config      *config.GoogleConfig
	httpClient  *http.Client
	rateLimiter *rate.Limiter
	logger      *logrus.Logger
}

// GetGoogleService 获取GoogleService单例
func GetGoogleService() *GoogleService {
	if googleServiceInstance == nil {
		SetupGoogleService()
	}
	return googleServiceInstance
}

// SetupGoogleService 初始化GoogleService
func SetupGoogleService() {
	googleOnce.Do(func() {
		cfg := config.Get()
		if cfg == nil || cfg.Google == nil {
			panic("configuration not initialized")
		}

		// 创建HTTP客户端
		httpClient := &http.Client{
			Timeout: time.Duration(cfg.Google.RequestTimeout) * time.Second,
		}

		// 创建限流器
		rateLimiter := rate.NewLimiter(
			rate.Every(time.Duration(cfg.Google.RateLimitWindow)*time.Second/time.Duration(cfg.Google.MaxRequestsPerWindow)),
			cfg.Google.MaxRequestsPerWindow,
		)

		googleServiceInstance = &GoogleService{
			config:      cfg.Google,
			httpClient:  httpClient,
			rateLimiter: rateLimiter,
			logger:      log.GetServiceLogger("google-service"),
		}
	})
}

// CalculateRouteDistance 计算两点间的路线距离
func (g *GoogleService) CalculateRouteDistance(ctx context.Context, req *protocol.RouteRequest) (*protocol.RouteResponse, error) {
	// 开始计时
	startTime := time.Now()

	// 验证配置
	if g.config.MapsAPIKey == "" {
		return nil, fmt.Errorf("google Maps API not configured")
	}

	// 验证请求参数
	if err := req.Validate(); err != nil {
		return &protocol.RouteResponse{
			Success:   false,
			Message:   err.Error(),
			Timestamp: time.Now(),
		}, err
	}

	// 限流控制
	if err := g.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit exceeded: %w", err)
	}

	// 构建请求URL
	apiURL, err := g.buildDirectionsURL(req)
	if err != nil {
		return nil, fmt.Errorf("failed to build request URL: %w", err)
	}

	// 发送API请求
	directionsResp, err := g.makeDirectionsRequest(ctx, apiURL)
	if err != nil {
		duration := time.Since(startTime)
		g.logger.Warnf("API request failed after %v: %v", duration, err)
		return &protocol.RouteResponse{
			Success:   false,
			Message:   "API request failed",
			Error:     &protocol.GoogleAPIError{Code: "REQUEST_FAILED", Message: err.Error()},
			Timestamp: time.Now(),
		}, err
	}

	// 处理响应
	response, err := g.processDirectionsResponse(directionsResp, req)

	// 记录总耗时
	duration := time.Since(startTime)
	if err != nil {
		g.logger.Warnf("Route calculation failed after %v: %v", duration, err)
	} else if response.Success {
		g.logger.Infof("Route calculation completed in %v: distance=%dm, duration=%ds",
			duration, response.Distance.Value, response.Duration.Value)
	}

	return response, err
}

// buildDirectionsURL 构建Directions API请求URL
func (g *GoogleService) buildDirectionsURL(req *protocol.RouteRequest) (string, error) {
	baseURL := "https://maps.googleapis.com/maps/api/directions/json"

	params := url.Values{}
	params.Set("origin", req.Origin)
	params.Set("destination", req.Destination)
	params.Set("key", g.config.MapsAPIKey)
	params.Set("mode", req.Mode)
	params.Set("language", req.Language)
	params.Set("units", req.Units)

	if req.Avoid != "" {
		params.Set("avoid", req.Avoid)
	}
	if req.Region != "" {
		params.Set("region", req.Region)
	}
	if len(req.Waypoints) > 0 {
		waypoints := ""
		for i, wp := range req.Waypoints {
			if i > 0 {
				waypoints += "|"
			}
			waypoints += wp
		}
		params.Set("waypoints", waypoints)
	}

	return baseURL + "?" + params.Encode(), nil
}

// makeDirectionsRequest 发送Directions API请求
func (g *GoogleService) makeDirectionsRequest(ctx context.Context, apiURL string) (*protocol.DirectionsResponse, error) {
	var lastErr error

	// 重试机制
	for attempt := 0; attempt <= g.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// 指数退避
			backoff := time.Duration(attempt*attempt) * time.Second
			g.logger.Warnf("API请求失败，%v后重试 (第%d次)", backoff, attempt)

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
		if err != nil {
			lastErr = err
			continue
		}

		resp, err := g.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("HTTP error: %d, response: %s", resp.StatusCode, string(body))
			continue
		}

		var directionsResp protocol.DirectionsResponse
		if err := json.Unmarshal(body, &directionsResp); err != nil {
			lastErr = fmt.Errorf("failed to parse response: %w", err)
			continue
		}

		return &directionsResp, nil
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", g.config.MaxRetries, lastErr)
}

// processDirectionsResponse 处理Directions API响应
func (g *GoogleService) processDirectionsResponse(directionsResp *protocol.DirectionsResponse, req *protocol.RouteRequest) (*protocol.RouteResponse, error) {
	response := &protocol.RouteResponse{
		Timestamp: time.Now(),
	}

	// 检查API响应状态
	if !directionsResp.IsSuccess() {
		response.Success = false
		response.Error = &protocol.GoogleAPIError{
			Code:    directionsResp.Status,
			Message: directionsResp.ErrorMessage,
			Status:  directionsResp.Status,
		}

		switch directionsResp.Status {
		case protocol.StatusZeroResults:
			response.Message = "No available routes found"
		case protocol.StatusNotFound:
			response.Message = "Origin or destination address not found"
		case protocol.StatusInvalidRequest:
			response.Message = "Invalid request parameters"
		case protocol.StatusOverQueryLimit:
			response.Message = "API quota exceeded"
		case protocol.StatusRequestDenied:
			response.Message = "API access denied"
		default:
			response.Message = fmt.Sprintf("API call failed: %s", directionsResp.Status)
		}

		return response, fmt.Errorf(response.Message)
	}

	// 检查是否有路线
	if !directionsResp.HasRoutes() {
		response.Success = false
		response.Message = "No available routes found"
		return response, fmt.Errorf("no available routes found")
	}

	// 提取路线信息
	route := directionsResp.GetFirstRoute()
	response.Success = true
	response.Route = route
	response.Message = "Route calculation successful"

	// 计算总距离和时长
	if len(route.Legs) > 0 {
		totalDistance := directionsResp.GetTotalDistance()
		totalDuration := directionsResp.GetTotalDuration()

		response.Distance = &protocol.Distance{
			Value: totalDistance,
			Text:  g.formatDistance(totalDistance),
		}
		response.Duration = &protocol.Duration{
			Value: totalDuration,
			Text:  g.formatDuration(totalDuration),
		}

		g.logger.Infof("Route calculation successful: distance=%dm, duration=%ds", totalDistance, totalDuration)
	}

	return response, nil
}

// formatDistance 格式化距离显示
func (g *GoogleService) formatDistance(meters int) string {
	if meters >= 1000 {
		km := float64(meters) / 1000.0
		return fmt.Sprintf("%.1f km", km)
	}
	return fmt.Sprintf("%d m", meters)
}

// formatDuration 格式化时长显示
func (g *GoogleService) formatDuration(seconds int) string {
	if seconds >= 3600 {
		hours := seconds / 3600
		minutes := (seconds % 3600) / 60
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else if seconds >= 60 {
		minutes := seconds / 60
		return fmt.Sprintf("%dm", minutes)
	}
	return fmt.Sprintf("%ds", seconds)
}

// GetDistanceBetweenPoints 简化的距离计算接口
func (g *GoogleService) GetDistanceBetweenPoints(ctx context.Context, origin, destination string) (int, error) {
	startTime := time.Now()

	req := &protocol.RouteRequest{
		Origin:      origin,
		Destination: destination,
		Mode:        "driving",
		Language:    "en",
		Units:       "metric",
	}

	resp, err := g.CalculateRouteDistance(ctx, req)
	duration := time.Since(startTime)

	if err != nil {
		g.logger.Warnf("GetDistanceBetweenPoints failed after %v: %v", duration, err)
		return 0, err
	}

	if resp.Distance != nil {
		g.logger.Infof("GetDistanceBetweenPoints completed in %v: %dm", duration, resp.Distance.Value)
		return resp.Distance.Value, nil
	}

	g.logger.Warnf("GetDistanceBetweenPoints completed in %v but no distance data", duration)
	return 0, fmt.Errorf("unable to get distance information")
}

// CalculateRidehailingRoute 网约车路线计算专用函数
// 针对网约车场景优化，提供准确的驾车路线、距离和预估时间
func (g *GoogleService) CalculateRidehailingRoute(ctx context.Context, originLat, originLng, destLat, destLng float64, avoidTolls bool) (*protocol.RouteResponse, error) {
	// 开始计时
	startTime := time.Now()

	// 验证配置
	if g.config.MapsAPIKey == "" {
		return nil, fmt.Errorf("google Maps API not configured")
	}

	// 验证坐标有效性
	if originLat == 0 || originLng == 0 || destLat == 0 || destLng == 0 {
		return &protocol.RouteResponse{
			Success:   false,
			Message:   "Valid origin and destination coordinates required",
			Timestamp: time.Now(),
		}, fmt.Errorf("valid origin and destination coordinates required")
	}

	// 构建坐标字符串
	origin := fmt.Sprintf("%f,%f", originLat, originLng)
	destination := fmt.Sprintf("%f,%f", destLat, destLng)

	// 构建网约车专用请求
	req := &protocol.RouteRequest{
		Origin:      origin,
		Destination: destination,
		Mode:        "driving",
		Language:    "en",
		Units:       "metric",
		Region:      "RW", // 卢旺达地区代码
	}

	// 根据参数设置是否避开收费站
	if avoidTolls {
		req.Avoid = "tolls"
	}

	// 限流控制
	if err := g.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit exceeded: %w", err)
	}

	// 构建请求URL
	apiURL, err := g.buildDirectionsURL(req)
	if err != nil {
		return nil, fmt.Errorf("failed to build request URL: %w", err)
	}

	// 发送API请求
	directionsResp, err := g.makeDirectionsRequest(ctx, apiURL)
	if err != nil {
		duration := time.Since(startTime)
		g.logger.Warnf("Ridehailing API request failed after %v: %v", duration, err)
		return &protocol.RouteResponse{
			Success:   false,
			Message:   "API request failed",
			Error:     &protocol.GoogleAPIError{Code: "REQUEST_FAILED", Message: err.Error()},
			Timestamp: time.Now(),
		}, err
	}

	// 处理响应
	response, err := g.processDirectionsResponse(directionsResp, req)

	// 记录总耗时
	duration := time.Since(startTime)
	if err != nil {
		g.logger.Warnf("Ridehailing route calculation failed after %v: %v", duration, err)
	} else if response.Success && response.Route != nil {
		g.logger.Infof("Ridehailing route calculated in %v: origin=(%.6f,%.6f), dest=(%.6f,%.6f), distance=%dm, duration=%ds, avoid_tolls=%v",
			duration, originLat, originLng, destLat, destLng, response.Distance.Value, response.Duration.Value, avoidTolls)
	}

	return response, nil
} // CalculateMultipleRoutes 批量计算多个路线
func (g *GoogleService) CalculateMultipleRoutes(ctx context.Context, requests []*protocol.RouteRequest) ([]*protocol.RouteResponse, error) {
	if len(requests) == 0 {
		return nil, fmt.Errorf("request list is empty")
	}

	responses := make([]*protocol.RouteResponse, len(requests))

	// 并发处理多个请求，但受限流器控制
	sem := make(chan struct{}, 5) // 最多并发5个请求
	var wg sync.WaitGroup

	for i, req := range requests {
		wg.Add(1)
		go func(index int, request *protocol.RouteRequest) {
			defer wg.Done()

			sem <- struct{}{}        // 获取信号量
			defer func() { <-sem }() // 释放信号量

			resp, err := g.CalculateRouteDistance(ctx, request)
			if err != nil {
				resp = &protocol.RouteResponse{
					Success:   false,
					Message:   err.Error(),
					Timestamp: time.Now(),
				}
			}
			responses[index] = resp
		}(i, req)
	}

	wg.Wait()
	return responses, nil
}

// =============================================================================
// Google Places API 功能
// =============================================================================

// PlaceSearchResult Places API搜索结果
type PlaceSearchResult struct {
	PlaceID          string   `json:"place_id"`
	Name             string   `json:"name"`
	FormattedAddress string   `json:"formatted_address"`
	Rating           float64  `json:"rating,omitempty"`
	UserRatingsTotal int      `json:"user_ratings_total,omitempty"`
	Types            []string `json:"types,omitempty"`
	Geometry         struct {
		Location struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
	} `json:"geometry"`
}

// PlaceDetailsResult Places API详情结果
type PlaceDetailsResult struct {
	PlaceID                  string   `json:"place_id"`
	Name                     string   `json:"name"`
	FormattedAddress         string   `json:"formatted_address"`
	InternationalPhoneNumber string   `json:"international_phone_number,omitempty"`
	Website                  string   `json:"website,omitempty"`
	Rating                   float64  `json:"rating,omitempty"`
	UserRatingsTotal         int      `json:"user_ratings_total,omitempty"`
	PriceLevel               int      `json:"price_level,omitempty"`
	Types                    []string `json:"types,omitempty"`
	OpeningHours             struct {
		OpenNow     bool     `json:"open_now,omitempty"`
		WeekdayText []string `json:"weekday_text,omitempty"`
	} `json:"opening_hours,omitempty"`
	Geometry struct {
		Location struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
	} `json:"geometry"`
	Photos []struct {
		Height           int      `json:"height"`
		Width            int      `json:"width"`
		PhotoReference   string   `json:"photo_reference"`
		HTMLAttributions []string `json:"html_attributions"`
	} `json:"photos,omitempty"`
}

// PlacesSearchResponse Places API搜索响应
type PlacesSearchResponse struct {
	Results []PlaceSearchResult `json:"results"`
	Status  string              `json:"status"`
	Error   string              `json:"error_message,omitempty"`
}

// PlaceDetailsResponse Places API详情响应
type PlaceDetailsResponse struct {
	Result PlaceDetailsResult `json:"result"`
	Status string             `json:"status"`
	Error  string             `json:"error_message,omitempty"`
}

// SearchPlacesByName 根据名称搜索地点，返回Place ID
func (g *GoogleService) SearchPlacesByName(ctx context.Context, name, city, country string) (*PlaceSearchResult, error) {
	// 验证配置
	if g.config.MapsAPIKey == "" {
		return nil, fmt.Errorf("google Maps API not configured")
	}

	// 限流控制
	if err := g.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit exceeded: %w", err)
	}

	// 构建搜索查询
	query := name
	if city != "" && country != "" {
		query = fmt.Sprintf("%s, %s, %s", name, city, country)
	} else if city != "" {
		query = fmt.Sprintf("%s, %s", name, city)
	}

	// 构建请求URL
	baseURL := "https://maps.googleapis.com/maps/api/place/textsearch/json"
	params := url.Values{}
	params.Set("query", query)
	params.Set("key", g.config.MapsAPIKey)
	params.Set("language", "en")

	// 如果指定了国家，添加地区偏好
	if country != "" {
		params.Set("region", country)
	}

	apiURL := baseURL + "?" + params.Encode()
	g.logger.Infof("Searching place: query=%s, url=%s", query, apiURL)

	// 发送请求
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d, response: %s", resp.StatusCode, string(body))
	}

	var searchResp PlacesSearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 检查API响应状态
	if searchResp.Status != "OK" {
		if searchResp.Status == "ZERO_RESULTS" {
			return nil, fmt.Errorf("no places found for query: %s", query)
		}
		return nil, fmt.Errorf("places API error: %s - %s", searchResp.Status, searchResp.Error)
	}

	// 返回第一个结果
	if len(searchResp.Results) == 0 {
		return nil, fmt.Errorf("no places found for query: %s", query)
	}

	result := &searchResp.Results[0]
	g.logger.Infof("Found place: name=%s, place_id=%s, rating=%.1f, address=%s",
		result.Name, result.PlaceID, result.Rating, result.FormattedAddress)

	return result, nil
}

// GetPlaceDetails 根据Place ID获取地点详情
func (g *GoogleService) GetPlaceDetails(ctx context.Context, placeID string) (*PlaceDetailsResult, error) {
	// 验证配置
	if g.config.MapsAPIKey == "" {
		return nil, fmt.Errorf("google Maps API not configured")
	}

	if placeID == "" {
		return nil, fmt.Errorf("place ID is required")
	}

	// 限流控制
	if err := g.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit exceeded: %w", err)
	}

	// 构建请求URL
	baseURL := "https://maps.googleapis.com/maps/api/place/details/json"
	params := url.Values{}
	params.Set("place_id", placeID)
	params.Set("key", g.config.MapsAPIKey)
	params.Set("language", "en")

	// 指定需要的字段
	fields := []string{
		"place_id", "name", "formatted_address", "international_phone_number",
		"website", "rating", "user_ratings_total", "price_level", "types",
		"opening_hours", "geometry", "photos",
	}
	// Google Places API需要使用逗号分隔的字段列表
	params.Set("fields", strings.Join(fields, ","))

	apiURL := baseURL + "?" + params.Encode()
	g.logger.Infof("Getting place details: place_id=%s", placeID)

	// 发送请求
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		g.logger.Errorf("HTTP error: %d, response: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("HTTP error: %d, response: %s", resp.StatusCode, string(body))
	}

	var detailsResp PlaceDetailsResponse
	if err := json.Unmarshal(body, &detailsResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 检查API响应状态
	if detailsResp.Status != "OK" {
		if detailsResp.Status == "NOT_FOUND" {
			return nil, fmt.Errorf("place not found: %s", placeID)
		}
		return nil, fmt.Errorf("place details API error: %s - %s", detailsResp.Status, detailsResp.Error)
	}

	g.logger.Infof("Got place details: name=%s, rating=%.1f, phone=%s, website=%s",
		detailsResp.Result.Name, detailsResp.Result.Rating,
		detailsResp.Result.InternationalPhoneNumber, detailsResp.Result.Website)

	return &detailsResp.Result, nil
}

// GetPlacePhotoURL 获取Google Places照片URL
func (g *GoogleService) GetPlacePhotoURL(photoReference string, maxWidth, maxHeight int) string {
	if photoReference == "" || g.config.MapsAPIKey == "" {
		return ""
	}

	if maxWidth <= 0 {
		maxWidth = 400
	}
	if maxHeight <= 0 {
		maxHeight = 400
	}

	return fmt.Sprintf("https://maps.googleapis.com/maps/api/place/photo?maxwidth=%d&maxheight=%d&photoreference=%s&key=%s",
		maxWidth, maxHeight, photoReference, g.config.MapsAPIKey)
}

// ExtractCompleteGoogleData 提取完整的Google地点数据
func (g *GoogleService) ExtractCompleteGoogleData(details *PlaceDetailsResult) map[string]interface{} {
	data := make(map[string]interface{})

	// 基本信息
	if details.PlaceID != "" {
		data["place_id"] = details.PlaceID
	}
	if details.Name != "" {
		data["name"] = details.Name
	}
	if details.FormattedAddress != "" {
		data["formatted_address"] = details.FormattedAddress
	}

	// 联系信息
	if details.InternationalPhoneNumber != "" {
		data["international_phone_number"] = details.InternationalPhoneNumber
	}
	if details.Website != "" {
		data["website"] = details.Website
	}

	// 评分信息
	if details.Rating > 0 {
		data["rating"] = details.Rating
	}
	if details.UserRatingsTotal > 0 {
		data["user_ratings_total"] = details.UserRatingsTotal
	}

	// 位置信息
	if details.Geometry.Location.Lat != 0 {
		data["latitude"] = details.Geometry.Location.Lat
	}
	if details.Geometry.Location.Lng != 0 {
		data["longitude"] = details.Geometry.Location.Lng
	}

	// 价格等级
	if details.PriceLevel >= 0 && details.PriceLevel <= 4 {
		data["price_level"] = details.PriceLevel
	}

	// 商家类型
	if len(details.Types) > 0 {
		data["types"] = details.Types
	}

	// 营业时间
	if len(details.OpeningHours.WeekdayText) > 0 {
		// 将营业时间数组转换为JSON字符串
		hoursJSON, _ := json.Marshal(map[string]interface{}{
			"weekday_text": details.OpeningHours.WeekdayText,
			"open_now":     details.OpeningHours.OpenNow,
		})
		data["opening_hours"] = string(hoursJSON)
	}

	// 图片
	if len(details.Photos) > 0 {
		// 获取第一张照片的URL
		photoURL := g.GetPlacePhotoURL(details.Photos[0].PhotoReference, 800, 600)
		if photoURL != "" {
			data["image_url"] = photoURL
		}
	}

	return data
}
