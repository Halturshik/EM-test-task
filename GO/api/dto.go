package api

type SubResponse struct {
	ServiceName string  `json:"service_name" example:"Yandex Plus"`
	Price       int     `json:"price" example:"100"`
	StartDate   string  `json:"start_date" example:"07-2025"`
	EndDate     *string `json:"end_date,omitempty" example:"12-2025"`
}

type TotalCostResponse struct {
	Message string `json:"message" example:"Общая стоимость подписки Yandex Plus за указанный период составила: 300"`
}

type TotalCostRequest struct {
	TotalFrom string `json:"total_from" example:"07-2025"`
	TotalTo   string `json:"total_to" example:"09-2025"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Некорректный формат идентификатора пользователя"`
}

type CreateSubRequest struct {
	UserID      string  `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ServiceName string  `json:"service_name" example:"Yandex Plus"`
	Price       int     `json:"price" example:"100"`
	StartDate   string  `json:"start_date" example:"07-2025"`
	EndDate     *string `json:"end_date,omitempty" example:"12-2025"`
}

type CreateSubResponse struct {
	Message string `json:"message" example:"Подписка успешно создана"`
}

type UpdateSubRequest struct {
	NewPrice   *int    `json:"new_price,omitempty" example:"100"`
	NewEndDate *string `json:"new_end_date,omitempty" example:"12-2025"`
}

type UpdateSubResponse struct {
	Message string `json:"message" example:"Уровень подписки повышен и уже действует. Дата окончания подписки изменена"`
}

type DeleteSubRequest struct {
	StartDate string `json:"start_date" example:"07-2025"`
}

type DeleteSubResponse struct {
	Message string `json:"message" example:"Подписка успешно удалена"`
}
